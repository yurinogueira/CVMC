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

3. **Controle de Tamanho e Cota Persistida por Usuário (`100 Arquivos` ou `100 MB`)**:
   - Middleware de limitação de corpo de requisição (`BodyLimit`) restringe o payload HTTP a `15 MB`.
   - Cada conta possui uma cota máxima de **100 arquivos** (`MaxFiles = 100`) e **100 MB** (`MaxStorageBytes = 104857600` bytes), contabilizando todos os comprovantes de manutenção e fotos enviadas de veículos de propriedade do usuário (`car.OwnerID == userID`). O tamanho real em bytes é calculado diretamente a partir do payload Base64 decodificado (`filecipher.EstimateDecodedBytes`) para prevenir falsificação do campo `size`.
   - Para evitar recálculo em todas as leituras de perfil (`GET /api/v1/user/profile`), os contadores `FilesCount` e `StorageUsedBytes` são **persistidos diretamente no documento `User`** e atualizados incrementalmente (`+delta` / `-delta`) a cada criação, edição ou exclusão de veículo ou manutenção (`storagequota.CommitDelta`).
   - **Migração Lazy Automática de Arquivos Existentes (`StorageQuotaSynced`)**: Contas pré-existentes (`StorageQuotaSynced == false`) passam automaticamente por uma sincronização única (`storagequota.EnsureSynced`) no primeiro acesso ao perfil ou listagem/operação de veículos, que contabiliza todos os arquivos legados, criptografa in-place quaisquer fotos ou anexos ainda em texto plano (`data:...` -> `enc:v1:...`) e grava `StorageQuotaSynced = true` no registro do usuário. A partir daí, todas as leituras de cota ocorrem em $O(1)$.

4. **Criptografia Autenticada por Usuário em Repouso (`AES-256-GCM` + `HKDF-SHA256`)**:
   - O pacote `backend/internal/shared/filecipher` deriva uma subchave exclusiva de 256 bits (32 bytes) para cada `ownerID` a partir da chave mestra `FILE_ENCRYPTION_KEY` via `HKDF-SHA256`.
   - Antes de persistir no banco, `Attachment.DataUrl` e `Car.ImageUrl` (quando em formato Base64 `data:...`) são criptografados com **AES-256-GCM**, `nonce` aleatório de 12 bytes e `AAD = "owner:" + ownerID`, armazenando o prefixo `enc:v1:<base64>`.
   - A descriptografia ocorre exclusivamente em memória nos casos de uso após verificar que o usuário autenticado é proprietário (`OwnerID`) ou possui acesso compartilhado (`SharedWith`) ao veículo.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Visão Geral da Arquitetura](overview.md)
- [Subdomínio de Anexos](../domain/attachment.md)
- [Subdomínio de Manutenções](../domain/maintenance.md)
