---
type: architecture
title: "Armazenamento e Gestão de Mídia — CVMC"
description: "Padrões de upload de notas fiscais e comprovantes, abstração de storage local/OCI e defesas contra Stored XSS e Path Traversal."
tags:
  - architecture
  - storage
  - oci
  - security
  - attachments
timestamp: 2026-10-10
---

# 📦 Armazenamento e Gestão de Mídia — CVMC

O subsistema de armazenamento gerencia comprovantes de abastecimento, ordens de serviço e notas fiscais de manutenções no **CVMC**.

Para navegação geral, retorne ao [Catálogo Canônico](../index.md).

---

## 🔌 1. Abstração de Provedores de Armazenamento

A interface de storage reside em `backend/internal/application/ports/storage/storage.go`:

```go
type Provider interface {
    Save(ctx context.Context, filename string, r io.Reader) (string, error)
    Get(ctx context.Context, path string) (io.ReadCloser, error)
    Delete(ctx context.Context, path string) error
}
```

- **Local Storage (`localstorage`)**: Utilizado em ambiente de desenvolvimento local (`UPLOAD_PATH=./data/uploads`).
- **OCI Object Storage (`ocistorage`)**: Provedor Always Free da Oracle Cloud Infrastructure utilizado no ambiente de produção para alta disponibilidade e persistência desacoplada da máquina virtual.

---

## 🛡️ 2. Diretrizes de Segurança para Mídia e Anexos

1. **Prevenção a Path Traversal**:
   - Nomes de arquivos enviados pelo usuário nunca são utilizados diretamente no sistema de arquivos.
   - Todo upload recebe um identificador UUID gerado pelo servidor (`<uuid>.<extensão_validada>`).
   - Caracteres perigosos como `../`, `..\\` ou barras são eliminados.

2. **Mitigação Rigorosa de Stored DOM XSS**:
   - Conforme corrigido no PR #85, qualquer download ou renderização de comprovantes no navegador valida a extensão e MIME type contra uma whitelist rígida (`image/jpeg`, `image/png`, `application/pdf`).
   - Arquivos potencialmente perigosos (como `.svg`, `.html` ou scripts disfarçados) são servidos estritamente com `Content-Disposition: attachment` e cabeçalho `X-Content-Type-Options: nosniff` para impedir execução no contexto da aplicação.

3. **Controle de Tamanho**:
   - Middleware de limitação de corpo de requisição (`BodyLimit`) restringe uploads de anexos a no máximo 10MB por arquivo.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Visão Geral da Arquitetura](overview.md)
- [Subdomínio de Anexos](../domain/attachment.md)
- [Subdomínio de Manutenções](../domain/maintenance.md)
