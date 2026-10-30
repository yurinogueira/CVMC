---
type: domain
title: "Subdomínio: Relatórios & Exportações (Report)"
description: "Processamento assíncrono de relatórios de manutenção em PDF com fila sequencial, paginação de comprovantes e retenção de 30 dias."
tags:
  - domain
  - report
  - pdf
  - queue
  - storage
resource: backend/internal/domain/report
timestamp: 2026-10-29
---

# 📄 Subdomínio: Relatórios & Exportações (Report)

O subdomínio **Report** gerencia o processamento assíncrono, a geração otimizada de relatórios consolidados em formato PDF, a fila de processamento com controle estrito de concorrência e a política de retenção temporária de arquivos exportados.

Para navegação geral, retorne ao [Catálogo Canônico](../index.md).

---

## 🧩 Modelo de Domínio (`ReportJob`)

Estruturado em `backend/internal/domain/report/report.go`:

```go
type Status string

const (
    StatusPending    Status = "pending"
    StatusProcessing Status = "processing"
    StatusCompleted  Status = "completed"
    StatusFailed     Status = "failed"
    StatusExpired    Status = "expired"
)

type ReportType string

const (
    TypeMaintenancePDF ReportType = "maintenance_pdf"
)

type ReportJob struct {
    ID          string     `json:"id" bson:"_id,omitempty"`
    UserID      string     `json:"userId" bson:"userId"`
    CarID       string     `json:"carId" bson:"carId"`
    CarName     string     `json:"carName" bson:"carName"`
    Type        ReportType `json:"type" bson:"type"`
    Status      Status     `json:"status" bson:"status"`
    FilePath    string     `json:"filePath,omitempty" bson:"filePath,omitempty"`
    Error       string     `json:"error,omitempty" bson:"error,omitempty"`
    DurationMS  int64      `json:"durationMs,omitempty" bson:"durationMs,omitempty"`
    CreatedAt   time.Time  `json:"createdAt" bson:"createdAt"`
    CompletedAt *time.Time `json:"completedAt,omitempty" bson:"completedAt,omitempty"`
    ExpiresAt   *time.Time `json:"expiresAt,omitempty" bson:"expiresAt,omitempty"`
    ExpiredAt   *time.Time `json:"expiredAt,omitempty" bson:"expiredAt,omitempty"`
}
```

---

## 📋 Regras de Negócio e Arquitetura

1. **Fila Sequencial Otimizada (Concorrência = 1)**:
   - Para evitar estouro de memória RAM com processamento de imagens e compilação de PDFs, o backend processa estritamente **1 relatório por vez** através de um worker em Go (`WorkerQueue`) com canal bufferizado.
   - Chamada a `runtime.GC()` ao final de cada execução para liberação ativa de memória.
   - Reconciliação automática no startup do servidor: jobs em estado `pending` ou `processing` são re-enfileirados.

2. **Diagramação Inteligente de Comprovantes e Resiliência Visual**:
   - **Card de Informações do Veículo**: Colunas com distribuição balanceada (40mm, 64mm, 36mm, 34mm) e truncamento inteligente com reticências (`fitText`) garantindo que nomes de modelos extensos vindos da tabela FIPE não sobreponham o ano de fabricação ou outras colunas adjacentes.
   - **Tabela Cronológica & Quebra de Página**: Cálculo preventivo de altura de linha via `pdf.SplitLines` antes de iniciar cada registro e restauração mandatória da cor de preenchimento zebra (`SetFillColor`) na quebra de página (evitando que a primeira linha de páginas subsequentes herde o fundo Slate 900 do cabeçalho).
   - **Card de Comprovantes de Manutenção**: Altura adaptativa (28mm a 34mm) com quebra controlada em até 2 linhas para títulos longos de múltiplos serviços, sanitização de nomes de arquivos UUID internos (ex.: `7dd77e08...jpeg`) limitando a 25 caracteres, e cálculo dinâmico de enquadramento da imagem (`posY = headerBottomY + margin` e `availH = 275 - topY`) garantindo que a imagem nunca colida com o cabeçalho nem com o rodapé.
   - Cada comprovante de imagem (JPG, PNG, WebP) ganha página própria A4 dedicada ao final do dossiê, e comprovantes em formato PDF (ex: notas fiscais eletrônicas) têm suas páginas vetoriais mescladas diretamente ao final do documento gerado via `pdfcpu`.

3. **Notificações Push (FCM)**:
   - Ao concluir a geração do relatório (`StatusCompleted`), uma notificação push é enviada aos dispositivos registrados do usuário com metadados para redirecionamento.

4. **Retenção Estrita de 30 Dias**:
   - Todo relatório recebe `expiresAt = completedAt + 30 dias`.
   - O cronjob diário `ReportCleanupCron` remove os arquivos físicos do `storage.Provider` e marca o job como `StatusExpired`.
   - Requisições de download de relatórios expirados retornam `410 Gone`.
