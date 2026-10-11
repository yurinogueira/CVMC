---
type: domain
title: "Subdomínio: Anexos e Comprovantes Digitais (Attachment)"
description: "Higienização de notas fiscais, fotos de odômetro, tipos MIME válidos e proteção contra upload malicioso."
tags:
  - domain
  - attachment
  - upload
  - security
resource: backend/internal/domain/attachment
timestamp: 2026-10-10
---

# 📎 Subdomínio: Anexos e Comprovantes (Attachment)

O subdomínio **Attachment** gerencia o ciclo de vida e a integridade de notas fiscais, fotos e documentos vinculados a manutenções e despesas.

Para navegação geral, retorne ao [Catálogo Canônico](../index.md).

---

## 🧩 Modelo de Domínio (`Attachment`)

Estruturado em `backend/internal/domain/attachment/attachment.go`:

```go
type Attachment struct {
    ID          string    `json:"id"`
    Filename    string    `json:"filename"`
    ContentType string    `json:"contentType"`
    Size        int64     `json:"size"`
    StoragePath string    `json:"storagePath"`
    UploadedAt  time.Time `json:"uploadedAt"`
}
```

---

## 📋 Regras de Segurança e Whitelist

1. **Tipos MIME Permitidos**:
   - `image/jpeg` (`.jpg`, `.jpeg`)
   - `image/png` (`.png`)
   - `application/pdf` (`.pdf`)
2. **Defesa Contra Stored XSS**:
   - Arquivos `.svg`, `.html` e `.js` são sumariamente rejeitados.
   - Cabeçalho `X-Content-Type-Options: nosniff` é forçado na entrega do binário.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Armazenamento e Gestão de Mídia](../architecture/storage-and-media.md)
- [Subdomínio de Manutenções](maintenance.md)
