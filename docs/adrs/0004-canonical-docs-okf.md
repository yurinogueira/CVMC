---
type: adr
title: "ADR 0004: Governança Canônica de Documentação OKF e Zero-Divergence"
description: "Decisão de instituir a pasta docs/ como base canônica Open Knowledge Format com validação estrita automatizada."
tags:
  - adr
  - documentation
  - okf
  - governance
  - quality-gate
timestamp: 2026-10-10
---

# 📜 ADR 0004: Governança Canônica de Documentação OKF e Zero-Divergence

- **Status**: Aceito
- **Data**: 2026-10-10
- **Decisores**: Equipe de Engenharia CVMC

Para navegação geral, retorne ao [Catálogo Canônico](../index.md).

---

## 📌 Contexto
Com o uso de agentes autônomos de IA e desenvolvimento distribuído, documentações desatualizadas geravam alucinações de modelos, consumo excessivo de tokens e premissas equivocadas sobre o código-fonte.

---

## 💡 Decisão
Instituímos a especificação **Open Knowledge Format (OKF)** do Google Cloud como padrão mandatório para a pasta `docs/`:

1. A pasta `docs/` é a **Única Fonte Canônica da Verdade** do sistema.
2. Cada documento possui cabeçalho YAML frontmatter estrito e links relativos válidos.
3. Adotamos o protocolo **Zero-Divergence Continuous Docs**: qualquer alteração em código ou infraestrutura exige a atualização da documentação e o registro no diário datado (`docs/logs/AAAA-MM-DD.md`) no mesmo PR.
4. A esteira de CI/CD e o script `./scripts/check.sh docs` validam automaticamente a integridade dos links e impedem commits divergentes.

---

## ⚖️ Consequências

### Positivas:
- Redução drástica do consumo de tokens na janela de contexto dos modelos de IA via divulgação progressiva.
- Onboarding imediato e assertivo para humanos e agentes inteligentes.
- Rastreabilidade histórica e técnica em cada entrega.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Diário de Bordo & Registro Cronológico](../log.md)
- [Visão Geral da Arquitetura](../architecture/overview.md)
