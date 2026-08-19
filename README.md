# 🎮 Nexus Siege

Um jogo de estratégia híbrido **Tower Defense + RTS competitivo** (PvPvE) que roda diretamente no navegador.

> **Status:** Em Desenvolvimento  
> **Plataforma:** Browser (Desktop)  
> **Duração das Partidas:** 10–20 minutos  
> **Jogadores:** 1–3 (Solo, 1v1, ou 3 jogadores)

---

## 📖 Sobre o Jogo

**Nexus Siege** é um jogo onde você defende sua base com torres, ataca inimigos com unidades controláveis e sabota adversários enviando hordas de criaturas — tudo isso enquanto sobrevive a ondas neutras cada vez mais brutais.

### ✨ Diferenciais

| Diferencial | Descrição |
|---|---|
| **PvPvE simultâneo** | Você nunca está seguro: mesmo sem jogadores te atacando, as hordas neutras pressionam |
| **Sabotagem ativa** | Derrotar inimigos neutros gera recurso para enviar hordas aos adversários |
| **Defesa + ataque na mesma partida** | Não é só "colocar torre e rezar" — você controla unidades ativamente |
| **Assimetria de facções** | Cada facção muda completamente o estilo de jogo |
| **Herói controlável** | Ponto focal de microgerenciamento e progressão durante a partida |

---

## 🎯 Core Loop

```
┌─────────────────────────────────────────────────┐
│                                                 │
│   CONSTRUIR          DEFENDER                   │
│   (torres/base)  ←→  (waves neutras)            │
│        ↓                ↓                       │
│   PRODUZIR         ACUMULAR                     │
│   (unidades)       (recursos de kills)          │
│        ↓                ↓                       │
│   ATACAR           SABOTAR                      │
│   (base inimiga)   (enviar hordas)              │
│                                                 │
└─────────────────────────────────────────────────┘
```

---

## 🏗️ Arquitetura Técnica

O projeto segue uma arquitetura moderna e simplificada:

- **Backend:** Go (Golang) — concorrência nativa, performance e binário estático
- **Frontend (UI):** HTML + HTMX + Go Templates — reatividade sem JavaScript complexo
- **Frontend (Jogo):** HTML5 Canvas + TypeScript — renderização a 60fps
- **Comunicação:** WebSocket + protocolo binário — latência mínima
- **Banco de Dados:** SQLite embutido — sem dependências externas
- **Distribuição:** Binário único Go com assets embutidos (`go:embed`)

### Princípios de Design

1. **Servidor autoritativo** — previne cheating e garante consistência
2. **Binário único** — zero configuração, fácil deploy
3. **Performance** — 60fps no cliente, 30 ticks/s no servidor
4. **Simplicidade** — código modular, testável e documentado

---

## 📅 Roadmap do MVP

O desenvolvimento está planejado em **32 semanas** (8 meses), dividido em 4 fases:

| Fase | Duração | Objetivo |
|---|---|---|
| **Protótipo** | Semanas 1-8 | Core mechanics funcionando |
| **Vertical Slice** | Semanas 9-16 | Uma partida completa jogável |
| **Alpha** | Semanas 17-24 | Todos os modos e facções |
| **Beta/Polimento** | Semanas 25-32 | Estabilidade e performance |

### Critérios de Sucesso do MVP

- ✅ 3 modos de jogo (sobrevivência solo, 1v1, 3 jogadores)
- ✅ 3 facções completas com identidade distinta
- ✅ Sistema de sabotagem funcional
- ✅ Multiplayer estável via WebSocket
- ✅ Binário único executável (sem dependências externas)
- ✅ Performance mínima de 30 FPS

---

## 📚 Documentação

A documentação completa está disponível na pasta [`docs/`](./docs/):

| Documento | Descrição |
|---|---|
| [GDD.md](./docs/GDD.md) | Game Design Document — visão completa do design do jogo |
| [RFC-001-arquitetura.md](./docs/RFC-001-arquitetura.md) | RFC técnica com decisões de arquitetura do sistema |
| [plan.md](./docs/plan.md) | Plano de entrega do MVP com cronograma detalhado |

---

## 🚀 Como Rodar (em breve)

```bash
# Build do binário único
go build -o nexus-siege

# Executar o servidor
./nexus-siege

# Acessar no browser
# http://localhost:8080
```

---

## 🛠️ Tecnologias

- **Linguagens:** Go, TypeScript
- **Frontend:** HTML5 Canvas, HTMX, Go Templates
- **Comunicação:** WebSocket
- **Banco de Dados:** SQLite
- **Build:** `go:embed` para assets estáticos

---

## 📄 Licença

[Adicionar licença aqui]

---

## 👥 Equipe

- **Game Designer:** Equipe de Design
- **Tech Lead:** Tech Lead
- **Arquiteto de Software:** Arquiteto de Software

---

<p align="center">
  <strong>Nexus Siege</strong> — Construa. Defenda. Ataque. Sabote. Sobreviva.
</p>
