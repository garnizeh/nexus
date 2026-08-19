Vou criar o GDD completo em formato markdown para você:

```markdown
# 📖 Game Design Document — Nexus Siege

> **Versão:** 1.0  
> **Data:** 19 de Agosto de 2026  
> **Status:** Em Desenvolvimento  
> **Autor:** Equipe de Design

---

## 📋 Índice

1. [Visão Geral](#1-visão-geral)
2. [Conceito Central](#2-conceito-central)
3. [Modos de Jogo](#3-modos-de-jogo)
4. [Estrutura da Partida](#4-estrutura-da-partida)
5. [Mecânicas Principais](#5-mecânicas-principais)
6. [Facções](#6-facções)
7. [Mapas](#7-mapas)
8. [Progressão](#8-progressão)
9. [Interface e UX](#9-interface-e-ux)
10. [Estilo Visual e Áudio](#10-estilo-visual-e-áudio)
11. [Lore e Mundo](#11-lore-e-mundo)
12. [Referências e Inspirações](#12-referências-e-inspirações)
13. [Escopo e Roadmap](#13-escopo-e-roadmap)
14. [Visão Técnica](#14-visão-tecnica)
15. [Glossário](#15-glossário)

---

## 1. Visão Geral

### 🎯 Pitch

> Um jogo de estratégia híbrido onde você defende sua base com torres, ataca inimigos com unidades controláveis e sabota adversários enviando hordas de criaturas para eles — tudo isso enquanto sobrevive a ondas neutras cada vez mais brutais.

### 🧩 Gênero

- **Tower Defense + RTS competitivo** (PvPvE)
- Partidas rápidas (10–20 min)
- Estratégia em tempo real com pausa opcional (single-player)

### 🖥️ Plataforma

- **Browser** (desktop como foco, mobile como adaptação futura)
- Executável único (servidor Go + frontend embutido)

### 👥 Público-alvo

- Jogadores casuais de Tower Defense que querem mais profundidade
- Fãs de RTS que não têm tempo para partidas de 40 min
- Comunidade de jogos competitivos leves (tipo auto-battlers e tower defense multiplayer)

### ✨ Diferenciais

| Diferencial | Por quê? |
|---|---|
| **PvPvE simultâneo** | Você nunca está seguro: mesmo sem jogadores te atacando, as hordas neutras pressionam |
| **Sabotagem ativa** | Derrotar inimigos neutros gera recurso para enviar hordas aos adversários |
| **Defesa + ataque na mesma partida** | Não é só "colocar torre e rezar" — você controla unidades ativamente |
| **Assimetria de facções** | Cada facção muda completamente o estilo de jogo |
| **Herói controlável** | Ponto focal de microgerenciamento e progressão durante a partida |

### 🎮 Fatores de Retenção

- Partidas curtas e intensas ("só mais uma partida")
- Progressão cosmética entre partidas
- Leaderboards e desafios diários
- Variedade de facções e mapas
- Momentos de tensão (hordas neutras, sabotagens, ataques decisivos)

---

## 2. Conceito Central

### 🔄 Core Loop (Ciclo Principal)

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

### 💡 Experiência Desejada

- **Tensão constante:** sempre há algo para fazer, nunca é "esperar a próxima wave"
- **Decisões difíceis:** investir em defesa, ataque ou sabotagem?
- **Adaptação:** ler o que os oponentes estão fazendo e reagir
- **Satisfação de crescimento:** sua base evolui visualmente ao longo da partida
- **Momentos épicos:** derrotar um boss neutro, sabotar com timing perfeito, virada no late game

### 🎯 Pilares de Design

1. **Acessibilidade com profundidade**
   - Fácil de aprender (construir torres, produzir unidades)
   - Difícil de dominar (timing de sabotagem, micro do herói, leitura de oponentes)

2. **Tensão e agência**
   - Jogador sempre tem opções significativas
   - Nunca se sente impotente (mesmo perdendo, pode virar o jogo)

3. **Variedade e rejogabilidade**
   - 3 facções com estilos opostos
   - Múltiplos mapas com layouts diferentes
   - Cada partida é única (hordas neutras variadas, adversários diferentes)

4. **Competição justa**
   - Todos começam com as mesmas ferramentas
   - Meta-progressão é cosmética (sem pay-to-win)
   - Servidor autoritativo previne cheating

---

## 3. Modos de Jogo

### 🎮 Modo Solo (1 jogador)

#### Sobrevivência
- **Objetivo:** Durar o máximo possível contra waves infinitas
- **Dificuldade:** Escala continuamente (a cada 5 waves, boss aparece)
- **Recursos:** Ouro e Essência acumulam normalmente
- **Fim de jogo:** Quando o Núcleo é destruído
- **Leaderboard:** Global, baseado em número de waves sobrevividas
- **Pausa:** Habilitada para pensar com calma

#### Campanha
- **Estrutura:** Sequência de 10-15 mapas com objetivos específicos
- **Objetivos variados:**
  - Sobreviver X waves
  - Proteger aliados NPC
  - Destruir base inimiga controlada por IA
  - Completar com restrições (ex: "sem torres de longo alcance")
- **Progressão:** Completar mapas desbloqueia cosméticos e avança a história
- **Pausa:** Habilitada

### ⚔️ Modo 2 Jogadores (Duelo)

- **Formato:** 1v1 competitivo
- **Hordas neutras:** Atacam ambos jogadores simultaneamente (mesmas waves)
- **Sabotagem:** Permitida e incentivada
- **Vitória:**
  - Destruir o Núcleo inimigo primeiro, **ou**
  - Ser o último com Núcleo intacto
- **Elemento de "corrida":** Se um jogador cai, o vencedor enfrenta uma wave final colossal
- **Pausa:** Desabilitada (tempo real)
- **Duração alvo:** 10-15 minutos

### 👑 Modo 3 Jogadores (Triunvirato)

- **Formato:** Free-for-all (todos contra todos)
- **Layout:** Bases em triângulo, hordas vêm do centro
- **Sabotagem:** Pode enviar hordas para qualquer um dos outros dois
- **Diplomacia informal:** Nada impede dois jogadores de focar no terceiro
- **Vitória:** Último com Núcleo vivo
- **Pausa:** Desabilitada
- **Duração alvo:** 15-20 minutos

### 📊 Tabela de Modos no Lançamento

| Modo | Jogadores | Pausa | Duração Alvo | Leaderboard |
|---|---|---|---|---|
| Sobrevivência Solo | 1 | Sim | Infinita (meta: high score) | Sim (waves) |
| Campanha | 1 | Sim | 10-15 min por mapa | Não |
| Duelo | 2 | Não | 10-15 min | Sim (ELO) |
| Triunvirato | 3 | Não | 15-20 min | Sim (ELO) |

---

## 4. Estrutura da Partida

### ⏱️ Fases de uma Partida Competitiva (2-3 jogadores)

| Fase | Duração | O que acontece |
|---|---|---|
| **Preparação** | 60s | Jogadores constroem torres iniciais, escolhem facção, veem o mapa |
| **Early Game** | 0-5 min | Waves neutras fracas, economia se estabelece, primeiros ataques de sondagem |
| **Mid Game** | 5-12 min | Sabotagens começam, escolhas de tech tree, primeiras quedas de base possíveis |
| **Late Game** | 12-18 min | Hordas neutras ficam brutais, unidades de elite, ataques decisivos |
| **Morte Súbita** | 18 min+ | Se ninguém venceu, hordas imparáveis forçam o fim |

### 🏆 Condições de Vitória

- **Primária:** Ser o último jogador com o **Núcleo** intacto
- **Alternativa (campanha/sobrevivência):** Atingir objetivo do mapa ou bater recorde de waves

### 💀 Condições de Derrota

- Núcleo destruído (por unidades inimigas ou hordas neutras)
- Em caso de queda simultânea, vence quem causou mais dano ao núcleo do adversário restante

### 🔄 Gerenciamento de Desconexão

- **Modo solo:** Pausa automática
- **Multiplayer:** Timeout de 30s; se não voltar, IA assume a base (modo defensivo) ou base é destruída (configurável)

---

## 5. Mecânicas Principais

### 🏰 5.1 A Base e o Núcleo

- Cada jogador possui uma **base** com um **Núcleo** central
- O Núcleo tem vida (HP). Se chegar a zero, o jogador é eliminado
- Ao redor do Núcleo há **slots de construção** (limitados e fixos)
- A base tem caminhos de entrada por onde vêm as hordas
- O Núcleo pode ser **fortificado** com torres ao redor (mas ocupa slots valiosos)

### 🗼 5.2 Torres (Defesa)

#### Características Gerais
- Construídas em **slots fixos** ao longo dos caminhos da sua base
- Torres podem ser **melhoradas** (níveis 1-3) e **vendidas** (recuperação parcial do custo)
- Torres **não andam** — mas podem ser **reconstruídas** em outro slot com custo e tempo

#### Tipos de Torres (presentes em todas as facções, com "sabor" diferente)

| Tipo | Função | Exemplo |
|---|---|---|
| **Dano único** | Foco em alvos fortes (bosses, unidades elite) | Canhão de precisão |
| **Dano em área** | Foco em hordas (dano em splash) | Lança-chamas |
| **Lentidão/controle** | Atrasa inimigos (slow, stun, knockback) | Torre de gelo |
| **Suporte** | Buffa torres próximas ou gera recurso | Torre de comando |

#### Upgrades de Torre

| Nível | Custo | Benefício |
|---|---|---|
| **1** | Base | Stats padrão |
| **2** | 1.5x custo base | +30% dano, +20% alcance |
| **3** | 2.5x custo base | +60% dano, +40% alcance, habilidade especial |

### ⚔️ 5.3 Unidades (Ataque e Defesa Móvel)

#### Características Gerais
- Produzidas em **edifícios específicos** (quartel, fábrica, ninho...)
- Diferente das torres, unidades são **controláveis diretamente** pelo jogador
- Têm custo, tempo de produção e limite populacional (ex: 30 de supply)

#### Usos Principais

- **Defender** a base contra hordas que passaram das torres
- **Atacar** a base de outro jogador (caminhando até lá pelo mapa)
- **Interceptar** unidades inimigas no meio do caminho
- **Capturar pontos de recurso** neutros no mapa (se houver)

#### Tipos de Unidades

| Tipo | Função | Exemplo |
|---|---|---|
| **Melee** | Dano corpo-a-corpo, tanque | Guerreiro, besta |
| **Ranged** | Dano à distância, frágil | Arqueiro, atirador |
| **Suporte** | Cura, buff, debuff | Curandeiro, engenheiro |
| **Especial** | Habilidades únicas (invisibilidade, explosão) | Assassino, bomba viva |

### 🦸 5.4 O Herói (Unidade Especial)

#### Características

- Cada jogador controla **um Herói** único (escolhido ao selecionar facção)
- O Herói é mais forte que unidades comuns e tem **habilidades ativas** (cooldown)
- Se o Herói morre, ele **ressuscita** após um cooldown longo (60s)
- O Herói cresce de nível durante a partida (mata inimigos → ganha XP → escolhe talentos)

#### Progressão do Herói

| Nível | XP Necessário | Benefício |
|---|---|---|
| **1** | 0 | Stats base, 1 habilidade |
| **2** | 100 | +20% stats, 2ª habilidade |
| **3** | 250 | +40% stats, 3ª habilidade |
| **4** | 500 | +60% stats, upgrade de habilidade |
| **5** | 1000 | +80% stats, habilidade ultimate |

#### Função do Herói

- Dar um ponto focal de controle e microgerenciamento para o jogador
- Criar momentos épicos (herói salvando a base, herói liderando ataque decisivo)
- Adicionar profundidade estratégica (escolha de talentos, timing de habilidades)

### 💰 5.5 Economia

#### Recursos

| Recurso | Como obter | Para que serve |
|---|---|---|
| **Ouro** | Renda passiva ao longo do tempo + kills de inimigos neutros | Construir torres, produzir unidades, upgrades |
| **Essência** | Derrotar hordas neutras, destruir unidades de jogadores, completar waves | **Enviar hordas extras para os adversários** e comprar unidades de elite |

#### Decisão Central

> Gastar Essência em sabotagem agora ou guardar para uma unidade poderosa depois?

#### Renda Passiva

- **Ouro:** +10 por segundo (base) + bônus de torres de suporte
- **Essência:** +1 por wave neutra completada + bônus por kills

### 🌊 5.6 Hordas Neutras (PvE)

#### Características

- Waves periódicas de inimigos controlados pelo jogo atacam **todas as bases ao mesmo tempo**
- Composição das waves é **visível com antecedência** (ícone mostrando próximos inimigos)
- A cada 5 waves, um **boss neutro** aparece
- Dificuldade escala com o tempo e com a quantidade de jogadores vivos

#### Mecânica de Risco

> Quanto mais forte a wave que você sobrevive **sem ajuda externa**, mais Essência você ganha.

#### Tipos de Inimigos Neutros

| Tipo | Característica | Fraqueza |
|---|---|---|
| **Grunt** | Barato, fraco, numeroso | Dano em área |
| **Tank** | Lento, muita vida, dano alto | Dano único, kiting |
| **Rusher** | Rápido, pouca vida, ignora unidades | Lentidão, torres de dano |
| **Flyer** | Voa sobre obstáculos, ignora caminhos | Torres anti-aéreo |
| **Boss** | Vida massiva, habilidades especiais | Foco de todas as torres |

### 🗡️ 5.7 Sabotagem (PvP Indireto)

#### Como Funciona

- Jogadores podem gastar **Essência** para enviar hordas extras aos caminhos dos adversários
- As hordas enviadas aparecem como **waves adicionais** na fila do alvo
- Quem envia recebe **bônus de renda** se a horda matar unidades ou torres do alvo

#### Limites Anti-Abuso

- **Cooldown:** 30s entre envios
- **Custo crescente:** Cada envio consecutivo custa +20% mais Essência
- **Limite de hordas:** Máximo de 3 hordas extras na fila do alvo

#### Tiers de Sabotagem

| Tier | Custo (Essência) | Composição | Uso |
|---|---|---|---|
| **1** | 20 | 5 Grunts | Pressure leve, teste de defesa |
| **2** | 50 | 3 Tanks + 10 Grunts | Pressure médio, dano estrutural |
| **3** | 100 | 1 Boss + 5 Rushers | Pressure pesado, tentativa de kill |

### ⚔️ 5.8 Ataque Direto (PvP Direto)

#### Como Funciona

- Unidades produzidas podem marchar até a base inimiga e atacar
- Caminhos entre bases são visíveis no mapa
- Unidades atacantes podem ser interceptadas por:
  - Torres do defensor
  - Unidades do defensor
  - Hordas neutras que estejam passando no momento (caos!)

#### Estratégias de Ataque

- **Rush:** Enviar muitas unidades fracas rapidamente
- **Precision:** Enviar poucas unidades fortes (herói + elite)
- **Distract:** Enviar horda de sabotagem + ataque com unidades

---

## 6. Facções

### 🛡️ 6.1 Vanguarda de Ferro (Equilibrada / Militar)

#### Filosofia
> Defesa sólida e contra-ataque preciso.

#### Estilo de Jogo
- Ideal para iniciantes
- Torres com bom alcance e dano confiável
- Unidades versáteis (melee e ranged)
- Herói: tanque que protege torres próximas

#### Torres

| Nome | Tipo | Custo | Especial |
|---|---|---|---|
| **Canhão de Precisão** | Dano único | 100 | +50% dano em bosses |
| **Metralhadora** | Dano em área | 80 | Atira em 3 alvos simultaneamente |
| **Torre de Shock** | Lentidão | 90 | Slow de 40%, stun chance de 10% |
| **Torre de Comando** | Suporte | 120 | +20% dano para torres em 3 slots |

#### Unidades

| Nome | Tipo | Custo | Supply | Especial |
|---|---|---|---|---|
| **Soldado** | Melee | 50 | 2 | Tanque básico |
| **Atirador** | Ranged | 60 | 2 | Dano à distância |
| **Engenheiro** | Suporte | 80 | 3 | Repara torres próximas |

#### Herói: Comandante Ferro

- **Habilidade 1:** Grito de Guerra — buffa unidades próximas (+30% dano, 10s)
- **Habilidade 2:** Escudo Energético — absorve dano para torres próximas (500 HP, 20s)
- **Habilidade 3 (Ultimate):** Fortaleza — torres próximas atiram 2x mais rápido (15s)

---

### 🕷️ 6.2 Enxame Voraz (Agressiva / Números)

#### Filosofia
> Inundar o inimigo com quantidade.

#### Estilo de Jogo
- Estilo rush/suicida
- Torres baratas, dano em área, curto alcance
- Unidades fracas, mas muito rápidas e baratas
- Herói: invoca minions ao matar inimigos
- Sabotagem mais barata (envia mais hordas por menos Essência)

#### Torres

| Nome | Tipo | Custo | Especial |
|---|---|---|---|
| **Lança-Chamas** | Dano em área | 60 | Dano em cone, queima contínua |
| **Ninho de Ácido** | Dano em área | 70 | Slow + dano ao longo do tempo |
| **Espinhos** | Controle | 50 | Dano em área ao passar, slow |
| **Colmeia** | Suporte | 100 | Gera 2 drones que atacam automaticamente |

#### Unidades

| Nome | Tipo | Custo | Supply | Especial |
|---|---|---|---|---|
| **Zangão** | Melee | 30 | 1 | Rápido, fraco, numeroso |
| **Vespa** | Ranged | 40 | 1 | Voa, ignora caminhos |
| **Rainha Menor** | Suporte | 70 | 3 | Cura unidades próximas |

#### Herói: Matriarca do Enxame

- **Habilidade 1:** Invocar Zangões — spawn 5 zangões (cooldown 15s)
- **Habilidade 2:** Fúria do Enxame — unidades próximas atacam 2x mais rápido (8s)
- **Habilidade 3 (Ultimate):** Enxame Devorador — invoca 20 zangões + 5 vespas (cooldown 60s)

---

### 🔮 6.3 Círculo Arcano (Controle / Late Game)

#### Filosofia
> Segurar o jogo com controle e vencer no late.

#### Estilo de Jogo
- Estilo estratégico/defensivo
- Torres caras, efeitos poderosos (lentidão, dano em cadeia, debuffs)
- Unidades lentas, mas com habilidades especiais
- Herói: mago de controle de área
- Fraqueza: early game vulnerável

#### Torres

| Nome | Tipo | Custo | Especial |
|---|---|---|---|
| **Orbe Arcano** | Dano único | 120 | Dano em cadeia (atinge 3 alvos) |
| **Prisma** | Dano em área | 140 | Dano em área massivo, cooldown longo |
| **Vórtice** | Controle | 130 | Puxa inimigos para o centro, slow 60% |
| **Cristal de Mana** | Suporte | 150 | Gera +5 Essência por wave |

#### Unidades

| Nome | Tipo | Custo | Supply | Especial |
|---|---|---|---|---|
| **Aprendiz** | Ranged | 70 | 2 | Dano mágico, ignora armadura |
| **Golem** | Melee | 90 | 3 | Tanque, slow ao atacar |
| **Arcanista** | Suporte | 100 | 3 | Buffa unidades com escudo mágico |

#### Herói: Arquimago

- **Habilidade 1:** Bola de Fogo — dano em área (cooldown 10s)
- **Habilidade 2:** Teleporte — move instantaneamente (cooldown 20s)
- **Habilidade 3 (Ultimate):** Tempestade Arcana — dano massivo em área grande (cooldown 90s)

---

## 7. Mapas

### 🗺️ Layouts por Quantidade de Jogadores

| Jogadores | Formato | Descrição |
|---|---|---|
| **1** | Linear / Arena | Hordas vêm de uma ou mais direções; foco em sobrevivência |
| **2** | Espelhado | Duas bases simétricas, caminhos paralelos e uma "zona neutra" central com pontos de recurso |
| **3** | Triangular | Bases nos vértices, hordas vêm do centro; cada aresta tem um caminho conectando duas bases |

### 🧱 Elementos de Mapa

- **Slots de torre:** Posições fixas e demarcadas (15-25 slots por base)
- **Caminhos de invasão:** Rotas pré-definidas para hordas neutras (2-3 caminhos por base)
- **Rotas entre bases:** Caminhos por onde unidades de jogadores transitam
- **Pontos de recurso neutros** (opcional): Capturáveis por unidades, geram ouro extra
- **Obstáculos visuais:** Rochas, árvores, ruínas (puramente estético, sem impacto mecânico no MVP)

### 🗺️ Mapas no Lançamento

#### Solo: "Ruínas Esquecidas"
- **Layout:** Linear, hordas vêm do norte
- **Slots:** 20 slots de torre
- **Caminhos:** 2 caminhos de invasão
- **Especial:** Pontos de recurso no meio do mapa

#### 1v1: "Vales Gêmeos"
- **Layout:** Espelhado, bases nos extremos
- **Slots:** 18 slots por base
- **Caminhos:** 2 caminhos de invasão + 1 rota entre bases
- **Especial:** Zona neutra central com 2 pontos de recurso

#### 3 Jogadores: "Triângulo de Cristal"
- **Layout:** Triangular, bases nos vértices
- **Slots:** 15 slots por base
- **Caminhos:** 2 caminhos de invasão por base + 3 rotas entre bases
- **Especial:** Hordas vêm do centro (atacam todas as bases)

### 🎲 Variedade Futura

- Rotação diária de mapas no modo competitivo
- Editor de mapas (pós-lançamento)
- Mapas temáticos (neve, deserto, vulcão) com efeitos visuais

---

## 8. Progressão

### 📈 Progressão na Partida (Dentro de um Jogo)

- Torres e unidades ficam mais fortes com upgrades
- Herói ganha níveis e talentos
- Economia cresce com o tempo e kills
- Waves neutras escalam em dificuldade

### 🏅 Meta-Progressão (Entre Partidas)

#### Perfil do Jogador
- **Nível:** Ganha XP por partidas jogadas e vitórias
- **Conquistas:** Marcos especiais (ex: "Sobreviver 50 waves", "Vencer 10 partidas 1v1")
- **Estatísticas:** Vitórias, derrotas, waves sobrevividas, kills, etc.

#### Desbloqueios Cosméticos
- **Skins de torres:** Aparência diferente (ex: torre de fogo → torre de lava)
- **Skins de unidades:** Uniformes, cores, efeitos
- **Skins de herói:** Armaduras, armas, efeitos de habilidades
- **Efeitos de morte:** Animações especiais ao matar inimigos
- **Emotes:** Gestos para usar durante partidas (apenas visual, sem impacto)

#### Ranked
- **Ladder:** Sistema de ELO (Bronze → Prata → Ouro → Diamante → Mestre)
- **Temporadas:** Reset a cada 3 meses, recompensas por placement
- **Matchmaking:** Baseado em ELO para partidas justas

#### Justiça Competitiva
> **Nenhum desbloqueio que afete balanceamento** — tudo que influencia gameplay disponível desde o início.

### 📖 Campanha (Modo Solo)

#### Estrutura
- 10-15 mapas iniciais, cada um com regra especial
- Completar mapas desbloqueia cosméticos e conta a história do mundo

#### Exemplos de Mapas

| Mapa | Objetivo | Regra Especial |
|---|---|---|
| **Primeiro Contato** | Sobreviver 10 waves | Tutorial, apenas Vanguarda de Ferro |
| **Enxame Crescendo** | Sobreviver 15 waves | Apenas unidades melee |
| **Torre de Marfim** | Destruir base inimiga | Sem torres de longo alcance |
| **Céu em Chamas** | Sobreviver 20 waves | Apenas inimigos voadores |
| **Três Reis** | Sobreviver 25 waves | 3 núcleos para proteger |
| **Mãos Vazias** | Sobreviver 15 waves | Sem produção de unidades — só torres |

---

## 9. Interface e UX

### 🖱️ Controles Principais

- **Mouse:**
  - Clique esquerdo: Selecionar, construir, mover, atacar
  - Clique direito: Cancelar ação, mover câmera
  - Scroll: Zoom in/out
- **Teclado:**
  - WASD / Setas: Mover câmera
  - 1-9: Atalhos para habilidades do herói
  - Q, W, E, R: Atalhos para construção rápida
  - Space: Pausar (apenas solo)
  - Tab: Alternar entre grupos de unidades

### 🖥️ Elementos de HUD

#### Durante a Partida

- **Topo:**
  - Timer da partida
  - Próxima wave neutra (com countdown)
  - Recursos (Ouro e Essência)
  - População de unidades (ex: 15/30)

- **Canto superior direito:**
  - Minimap (em mapas grandes)
  - Lista de jogadores com status do núcleo

- **Inferior:**
  - Barra de construção (torres, unidades)
  - Habilidades do herói (com cooldowns visuais)
  - Botão de pausa (apenas solo)

- **Lateral esquerda:**
  - Fila de produção de unidades
  - Upgrades disponíveis

- **Alertas visuais:**
  - "Você foi sabotado!" (vermelho, piscando)
  - "Jogador X está atacando sua base!" (amarelo)
  - "Boss neutro chegando!" (vermelho, grande)

#### Menus

- **Menu principal:** Jogar, Campanha, Perfil, Configurações, Sair
- **Lobby:** Lista de salas, criar sala, entrar em matchmaking
- **Seleção de facção:** Preview da facção com torres, unidades e herói
- **Pós-partida:** Resultados, estatísticas, replay, voltar ao menu

### 📱 Acessibilidade

- **Daltonismo:** Cores alternativas para inimigos/aliados (3 modos)
- **Escala de texto:** Ajustável (pequeno, médio, grande)
- **Alto contraste:** Modo para torres e caminhos mais visíveis
- **Legendas:** Para sons importantes (ex: "Boss chegando!")
- **Velocidade de câmera:** Ajustável (lenta, normal, rápida)

---

## 10. Estilo Visual e Áudio

### 🎨 Direção de Arte

#### Estilo
- **Low-poly estilizado** com cores vibrantes
- Fácil de ler em tela pequena
- Leve para browser (performance)

#### Câmera
- **Isométrica** ou **top-down com leve inclinação** (45°)
- Zoom in/out para detalhes ou visão geral

#### Legibilidade
- Cada unidade/torre deve ser reconhecível em 1 segundo
- Cores distintas por facção:
  - Vanguarda de Ferro: Azul/Aço
  - Enxame Voraz: Verde/Roxo
  - Círculo Arcano: Roxo/Dourado
- Inimigos neutros: Vermelho/Laranja (alerta)

#### Tema
- **Fantasia sombria** com toques de tecnologia antiga
- Ruínas, cristais, máquinas a vapor
- Alternativa: Sci-fi espacial (a decidir)

### 🔊 Áudio

#### Música
- **Orquestral/eletrônica** que intensifica conforme a wave
- Temas distintos por facção
- Música de vitória/derrota

#### Efeitos Sonoros (SFX)

| Ação | Som |
|---|---|
| Construir torre | Metal/clang |
| Produzir unidade | Voz curta / som de portão |
| Torre atirando | Distinto por tipo (tiro, fogo, magia) |
| Unidade atacando | Impacto (espada, tiro, magia) |
| Inimigo morrendo | Grito / explosão |
| Alerta de sabotagem | Alarme vermelho |
| Morte do herói | Som dramático |
| Núcleo sob ataque | Alarme urgente |
| Wave iniciando | Trombeta / gongo |
| Boss aparecendo | Música intensa + rugido |

#### Vozes
- Gritos curtos para unidades (opcional, dá personalidade)
- Narração para eventos importantes (ex: "Boss chegando!")

### 🎭 Tom

- Sério o suficiente para ter tensão competitiva
- Humor nos detalhes (gritos de unidades, descrições de torres)
- Nada de gore — violência estilizada

---

## 11. Lore e Mundo

### 🌍 O Mundo de Nexus

> *Em um mundo onde cristais de energia antiga (os "Nexus") sustentam a realidade, facções rivais lutam pelo controle desses pontos de poder. Mas os Nexus também atraem criaturas das sombras, hordas que surgem das fendas entre dimensões. Agora, cada facção deve defender seu Nexus não apenas dos rivais, mas das próprias forças do caos.*

### 📜 História

#### O Despertar dos Nexus
- Há mil anos, os Nexus foram descobertos por civilizações antigas
- Eles forneciam energia ilimitada, mas também atraíam criaturas de outras dimensões
- As civilizações caíram, mas os Nexus permaneceram

#### As Facções

**Vanguarda de Ferro**
> *Descendentes de uma ordem militar antiga, acreditam que a disciplina e a tecnologia são a chave para sobreviver. Suas fortalezas são monumentos de aço e pedra, construídas para resistir a qualquer ameaça.*

**Enxame Voraz**
> *Uma colmeia de criaturas simbióticas que evoluíram para harnessar o poder dos Nexus. Não constroem — crescem. Não planejam — adaptam. Para eles, a vitória é uma questão de números e velocidade.*

**Círculo Arcano**
> *Magos que estudaram os segredos dos Nexus por gerações. Dominam as artes arcanas, mas seu poder tem um preço: lentidão. Preferem controlar o campo de batalha com precisão cirúrgica, vencendo no late game.*

#### O Conflito Atual
- Os Nexus estão se tornando instáveis
- Hordas neutras (criaturas das fendas) atacam com mais frequência
- As facções percebem que não podem vencer sozinhas
- Mas a desconfiança é grande — quem controlará os Nexus no final?

### 🎭 Personagens (Heróis)

#### Comandante Ferro (Vanguarda)
- **Background:** Veterano de mil batalhas, acredita que a defesa é a melhor ofensa
- **Personalidade:** Estoico, protetor, líder nato
- **Motivação:** Proteger seu povo a qualquer custo

#### Matriarca do Enxame (Enxame)
- **Background:** A mente coletiva do Enxame, conecta todas as criaturas
- **Personalidade:** Alienígena, implacável, adaptável
- **Motivação:** Expandir o Enxame, consumir tudo

#### Arquimago (Círculo Arcano)
- **Background:** O mais poderoso dos magos, guarda segredos antigos
- **Personalidade:** Misterioso, calculista, paciente
- **Motivação:** Desvendar os mistérios dos Nexus

---

## 12. Referências e Inspirações

### 🎮 Jogos Similares

| Jogo | O que inspira |
|---|---|
| **StarCraft** | RTS competitivo, facções assimétricas, microgerenciamento |
| **Kingdom (série)** | Tower Defense minimalista, atmosfera |
| **Orcs Must Die!** | Combinação de TD com ação em terceira pessoa |
| **Dune: Spice Wars** | RTS com elementos de 4X, sabotagem |
| **Mindustry** | TD com produção de recursos e automação |
| **Bloons TD 6** | Polimento, variedade de torres, coop multiplayer |
| **Auto-chess / TFT** | Meta-progressão, matchmaking, partidas rápidas |

### 🎯 Mecânicas Específicas

- **Sabotagem (enviar hordas):** Inspirado em *Dune: Spice Wars* e *Age of Empires* (raids)
