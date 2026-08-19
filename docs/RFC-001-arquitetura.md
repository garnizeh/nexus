```markdown
# 📐 RFC 001 — Arquitetura do Sistema "Nexus Siege"

> **Status:** Proposta  
> **Autor:** Arquiteto de Software  
> **Data:** 19 de Agosto de 2026  
> **Revisores:** Tech Lead, Game Designer, Product Owner  
> **Decisão:** Pendente de aprovação

---

## 📋 Índice

1. [Contexto e Motivação](#1-contexto-e-motivação)
2. [Resumo das Decisões](#2-resumo-das-decisões)
3. [Visão Geral da Arquitetura](#3-visão-geral-da-arquitetura)
4. [Backend (Go)](#4-backend-go)
5. [Frontend](#5-frontend)
6. [Comunicação](#6-comunicação)
7. [Fluxo de uma Partida](#7-fluxo-de-uma-partida)
8. [Modelo de Dados](#8-modelo-de-dados)
9. [Build e Distribuição](#9-build-e-distribuição)
10. [Riscos e Mitigações](#10-riscos-e-mitigações)
11. [Alternativas Consideradas](#11-alternativas-consideradas)
12. [Próximos Passos](#12-próximos-passos)
13. [Apêndices](#13-apêndices)

---

## 1. Contexto e Motivação

### 1.1 Problema

Precisamos definir a arquitetura técnica do jogo **Nexus Siege**, um híbrido de Tower Defense + RTS para navegador com as seguintes características:

- Partidas de 1 a 3 jogadores em tempo real
- Lógica de jogo determinística e sincronizada entre clientes
- Comunicação de baixa latência entre jogadores e servidor
- Interface responsiva e fluida (60fps ideal)
- Distribuição simplificada: **um único binário executável** contendo backend + frontend
- Suporte a modos solo (com pausa) e multiplayer competitivo

### 1.2 Restrições

| Restrição | Justificativa |
|---|---|
| Deve rodar em browser | Alcance máximo, sem instalação |
| Deve ser um executável único | Facilita distribuição e deploy |
| Deve suportar até 3 jogadores simultâneos | Escopo do MVP |
| Latência aceitável < 200ms | Experiência competitiva |
| Performance mínima de 30 FPS | Máquinas modestas |
| Sem dependências externas de runtime | Autossuficiência |

### 1.3 Objetivos

1. **Simplicidade de deploy:** Um único binário, zero configuração
2. **Performance:** 60fps no cliente, 30 ticks/s no servidor
3. **Consistência:** Servidor autoritativo, sem cheating possível
4. **Escalabilidade futura:** Arquitetura que permita crescer para mais jogadores/modos
5. **Manutenibilidade:** Código modular, testável, documentado

---

## 2. Resumo das Decisões

| Aspecto | Decisão | Justificativa |
|---|---|---|
| **Backend** | Go | Concorrência nativa (goroutines), compilação estática, performance adequada para game server |
| **Frontend (UI)** | HTML + HTMX + Go Templates | Simplicidade, reatividade sem JS complexo |
| **Frontend (Jogo)** | HTML5 Canvas + TypeScript | Renderização 60fps, controle fino |
| **Comunicação** | WebSocket + protocolo binário | Latência mínima para ações em tempo real |
| **Distribuição** | Binário único Go com assets embutidos (`go:embed`) | Deploy simplificado |
| **Banco de dados** | SQLite embutido | Sem dependências externas; suficiente para perfis e replays |
| **Autoridade** | Servidor autoritativo | Previne cheating e garante consistência |
| **Renderização** | Canvas 2D (MVP), WebGL (futuro) | Simplicidade inicial, performance adequada |
| **Protocolo** | Binário customizado (ou Protobuf) | Eficiência de banda |

---

## 3. Visão Geral da Arquitetura

### 3.1 Diagrama de Alto Nível

```mermaid
flowchart TB
    subgraph Cliente["🖥️ NAVEGADOR DO CLIENTE"]
        direction TB
        UI["HTMX + DOM<br/>(Menus, HUD, Painéis)"]
        Canvas["Canvas 2D<br/>(Renderização do jogo em 60fps)"]
        GameClient["Game Client (TypeScript)<br/>- Input handling<br/>- State interpolation<br/>- WebSocket client"]
        
        UI <-->|HTTP| Canvas
        Canvas <--> GameClient
    end
    
    subgraph Servidor["⚙️ SERVIDOR GO (binário único)"]
        direction TB
        HTTP["HTTP/HTMX Server<br/>(UI pages)"]
        WS["WebSocket Hub<br/>(rooms)"]
        GameServer["Game Server<br/>(autoritativo)"]
        
        GameLoop["Game Loop (tick-based)<br/>- Simulação (30 ticks/s)<br/>- Validação de comandos<br/>- Broadcast de snapshots"]
        
        SQLite["SQLite<br/>(perfis)"]
        Matchmaker["Matchmaker<br/>(lobbies)"]
        Replay["Replay Store<br/>(opcional)"]
        
        HTTP --> GameLoop
        WS --> GameLoop
        GameServer --> GameLoop
        
        GameLoop --> SQLite
        GameLoop --> Matchmaker
        GameLoop --> Replay
    end
    
    GameClient <-->|WebSocket (binário)| WS
    UI -->|HTTP| HTTP
    
    style Cliente fill:#e1f5ff,stroke:#0288d1,stroke-width:3px
    style Servidor fill:#fff4e1,stroke:#f57c00,stroke-width:3px
    style GameLoop fill:#ffe1e1,stroke:#d32f2f,stroke-width:2px
```

### 3.2 Separação de Responsabilidades

| Camada | Responsabilidade | Tecnologia |
|---|---|---|
| **Apresentação (UI)** | Menus, HUD, painéis, formulários | HTML + HTMX + CSS |
| **Renderização (Jogo)** | Mapa, unidades, torres, efeitos | Canvas 2D + TypeScript |
| **Lógica de Cliente** | Input, predição, interpolação | TypeScript |
| **Lógica de Servidor** | Simulação, validação, autoridade | Go |
| **Comunicação** | Mensagens cliente-servidor | WebSocket binário |
| **Persistência** | Perfis, stats, replays | SQLite |

### 3.3 Princípios de Design

1. **Servidor é autoridade:** Cliente nunca modifica estado diretamente
2. **Cliente é "burro":** Renderiza e envia inputs; não simula
3. **Protocolo mínimo:** Só trafega o essencial (comandos + snapshots delta)
4. **Estado é determinístico:** Mesmos inputs → mesmo resultado
5. **Assets embutidos:** Zero dependências externas no runtime

---

## 4. Backend (Go)

### 4.1 Estrutura de Pacotes

```mermaid
flowchart LR
    subgraph cmd["cmd/"]
        main["server/main.go<br/>(Ponto de entrada)"]
    end
    
    subgraph internal["internal/"]
        subgraph game["game/"]
            entity["entity.go"]
            faction["faction.go"]
            simulation["simulation.go"]
            pathfinding["pathfinding.go"]
            combat["combat.go"]
            waves["waves.go"]
            economy["economy.go"]
        end
        
        subgraph server["server/"]
            http["http.go"]
            websocket["websocket.go"]
            room["room.go"]
            protocol["protocol.go"]
        end
        
        subgraph matchmaking["matchmaking/"]
            queue["queue.go"]
            lobby["lobby.go"]
        end
        
        subgraph storage["storage/"]
            sqlite["sqlite.go"]
            profile["profile.go"]
        end
        
        subgraph replay["replay/"]
            recorder["recorder.go"]
            player["player.go"]
        end
    end
    
    subgraph pkg["pkg/"]
        vec2["vec2/"]
        config["config/"]
    end
    
    subgraph web["web/"]
        static["static/<br/>(js, css, assets)"]
        templates["templates/<br/>(HTML + HTMX)"]
    end
    
    main --> game
    main --> server
    main --> matchmaking
    main --> storage
    main --> replay
    
    style cmd fill:#e8f5e9,stroke:#388e3c,stroke-width:2px
    style internal fill:#fff3e0,stroke:#f57c00,stroke-width:2px
    style game fill:#ffe0b2,stroke:#ff6f00
    style server fill:#ffe0b2,stroke:#ff6f00
    style web fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
```

### 4.2 Game Loop Autoritativo

```mermaid
flowchart TD
    Start([Início do Tick]) --> Collect[1. Coletar comandos<br/>dos jogadores]
    Collect --> Validate[2. Validar comandos<br/>recursos, cooldowns, regras]
    Validate --> Simulate[3. Simular tick]
    
    subgraph Simulate_Detail["Simulação Detalhada"]
        Move[Movimento de unidades]
        Attack[Ataques de torres]
        Combat[Combate]
        Spawn[Spawn de waves]
        Economy[Economia<br/>renda passiva]
    end
    
    Simulate --> Simulate_Detail
    Simulate_Detail --> Generate[4. Gerar snapshot<br/>do estado]
    Generate --> Broadcast[5. Broadcast para clientes<br/>delta compression]
    Broadcast --> Sleep[6. Dormir até próximo tick<br/>33ms para 30 ticks/s]
    Sleep --> Start
    
    style Start fill:#4caf50,stroke:#2e7d32,stroke-width:2px,color:#fff
    style Simulate_Detail fill:#fff9c4,stroke:#f9a825,stroke-width:2px
    style Sleep fill:#e1bee7,stroke:#7b1fa2,stroke-width:2px
```

### 4.3 Modelo de Entidades

```mermaid
classDiagram
    class Entity {
        <<interface>>
        +ID EntityID
        +Owner PlayerID
        +Position Vec2
    }
    
    class Unit {
        +Type UnitType
        +Target *Vec2
        +Health int
        +MaxHealth int
        +State UnitState
        +AttackCooldown time.Duration
        +Damage int
        +Range float32
        +Speed float32
        +SupplyCost int
    }
    
    class Tower {
        +Type TowerType
        +SlotID int
        +Level int
        +Cooldown time.Duration
        +Target *EntityID
        +Damage int
        +Range float32
        +AttackSpeed time.Duration
    }
    
    class Hero {
        +Faction FactionType
        +Level int
        +XP int
        +Health int
        +MaxHealth int
        +Abilities [3]Ability
        +RespawnTimer time.Duration
    }
    
    class Wave {
        +Number int
        +SpawnTime time.Time
        +Enemies []WaveEnemy
        +TargetPath PathID
        +IsBoss bool
    }
    
    class Player {
        +ID PlayerID
        +Name string
        +Faction FactionType
        +Gold int
        +Essence int
        +Population int
        +MaxPopulation int
        +IsAlive bool
        +NexusHealth int
    }
    
    Entity <|-- Unit
    Entity <|-- Tower
    Entity <|-- Hero
    
    Player "1" --> "*" Unit : controls
    Player "1" --> "*" Tower : owns
    Player "1" --> "1" Hero : controls
    Wave "1" --> "*" WaveEnemy : contains
```

### 4.4 Estados de uma Unidade

```mermaid
stateDiagram-v2
    [*] --> Idle: Spawn
    
    Idle --> Moving: Recebe comando<br/>de movimento
    Idle --> Attacking: Inimigo no alcance
    Idle --> Dead: HP <= 0
    
    Moving --> Idle: Chega ao destino
    Moving --> Attacking: Inimigo no alcance
    Moving --> Dead: HP <= 0
    
    Attacking --> Idle: Alvo morto<br/>ou fora de alcance
    Attacking --> Moving: Recebe comando<br/>de movimento
    Attacking --> Dead: HP <= 0
    
    Dead --> [*]: Remove da simulação
    
    note right of Idle
        Unidade parada,
        aguardando comandos
    end note
    
    note right of Moving
        Movendo para
        posição alvo
    end note
    
    note right of Attacking
        Atacando alvo
        dentro do alcance
    end note
```

---

## 5. Frontend

### 5.1 Abordagem Híbrida

```mermaid
flowchart TB
    subgraph UI_Camada["🎨 Camada UI (HTMX)"]
        Menu[Menu Principal]
        Lobby[Lobby / Sala]
        HUD[HUD durante partida]
        BuildMenu[Painel de Construção]
        PostGame[Tela Pós-Partida]
    end
    
    subgraph Game_Camada["🎮 Camada Jogo (Canvas)"]
        Map[Renderização do Mapa]
        Units[Unidades e Torres]
        Effects[Efeitos Visuais]
        Camera[Câmera Isométrica]
    end
    
    subgraph Logic_Camada["⚙️ Camada Lógica (TypeScript)"]
        Input[Input Handler]
        Network[WebSocket Client]
        Interpolation[State Interpolation]
        Prediction[Local Prediction]
        Audio[Audio Manager]
    end
    
    UI_Camada -->|Atualiza via HTTP/SSE| Logic_Camada
    Game_Camada -->|Renderiza estado| Logic_Camada
    Logic_Camada -->|Comandos| Network
    Network -->|Snapshots| Interpolation
    Interpolation -->|Estado suavizado| Game_Camada
    
    style UI_Camada fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    style Game_Camada fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px
    style Logic_Camada fill:#fff3e0,stroke:#f57c00,stroke-width:2px
```

### 5.2 Estrutura do Cliente TypeScript

```mermaid
flowchart TD
    subgraph Main["main.ts"]
        Bootstrap[Bootstrap]
        WSConnect[Conexão WebSocket]
    end
    
    subgraph Network["network/"]
        Connection[connection.ts<br/>Gerenciamento WebSocket]
        Protocol[protocol.ts<br/>Serialização binária]
        Interpolation[interpolation.ts<br/>Suavização de movimento]
    end
    
    subgraph Game["game/"]
        State[state.ts<br/>Estado local do jogo]
        Renderer[renderer.ts<br/>Renderização Canvas]
        Camera[camera.ts<br/>Câmera isométrica]
        Input[input.ts<br/>Mouse, teclado, seleção]
        Entities[entities.ts<br/>Representação visual]
    end
    
    subgraph UI["ui/"]
        HUD_TS[hud.ts<br/>Atualização do HUD]
        BuildMenu_TS[build-menu.ts<br/>Painel de construção]
    end
    
    subgraph Audio["audio/"]
        Sound[sound.ts<br/>Efeitos sonoros e música]
    end
    
    Main --> Network
    Main --> Game
    Main --> UI
    Main --> Audio
    
    style Main fill:#e8f5e9,stroke:#388e3c,stroke-width:2px
    style Network fill:#fff3e0,stroke:#f57c00,stroke-width:2px
    style Game fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
```

### 5.3 Pipeline de Renderização

```mermaid
flowchart LR
    Snapshot[Snapshot do Servidor<br/>10-15 Hz] --> Buffer[Buffer de<br/>Snapshots]
    Buffer --> Interp[Interpolação<br/>entre 2 últimos]
    Interp --> State[Estado Local<br/>suavizado]
    
    State --> Clear[Clear Canvas]
    Clear --> DrawMap[Desenhar Mapa]
    DrawMap --> DrawEntities[Desenhar Entidades<br/>torres, unidades]
    DrawEntities --> DrawEffects[Desenhar Efeitos<br/>projéteis, explosões]
    DrawEffects --> DrawUI[Desenhar UI<br/>seleção, health bars]
    DrawUI --> Present[Apresentar Frame<br/>60 FPS]
    
    Present --> State
    
    style Snapshot fill:#ffecb3,stroke:#ff8f00,stroke-width:2px
    style Present fill:#c8e6c9,stroke:#388e3c,stroke-width:2px
```

---

## 6. Comunicação

### 6.1 Protocolo de Mensagens

```mermaid
flowchart TD
    subgraph Cliente_Servidor["Cliente → Servidor (Comandos)"]
        C1[0x01: Construir Torre]
        C2[0x02: Produzir Unidade]
        C3[0x03: Mover Unidade]
        C4[0x04: Atacar Alvo]
        C5[0x05: Habilidade do Herói]
        C6[0x06: Enviar Horda]
    end
    
    subgraph Servidor_Cliente["Servidor → Clientes"]
        S1[0x10: Snapshot Completo]
        S2[0x11: Snapshot Delta]
        E1[0x20: Evento: Wave Iniciando]
        E2[0x21: Evento: Jogador Eliminado]
        E3[0x22: Evento: Partida Encerrada]
    end
    
    subgraph Formato["Formato da Mensagem"]
        Tipo[Tipo: 1 byte]
        Seq[Seq: 2 bytes]
        Payload[Payload: variável]
    end
    
    Cliente_Servidor --> Formato
    Servidor_Cliente --> Formato
    
    style Cliente_Servidor fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    style Servidor_Cliente fill:#fff3e0,stroke:#f57c00,stroke-width:2px
    style Formato fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px
```

### 6.2 Fluxo de Dados com Latência

```mermaid
sequenceDiagram
    participant C as Cliente
    participant S as Servidor
    
    Note over C,S: Tick Rate: 30Hz (33ms)<br/>Snapshot Rate: 10-15Hz<br/>Buffer: 100ms
    
    C->>S: Comando (mover unidade)
    Note over C: Predição local:<br/>unidade começa a mover
    
    Note over S: Valida comando<br/>(tick N)
    S-->>C: Snapshot (tick N)
    Note over C: Confirma posição<br/>do servidor
    
    Note over S: Simula tick N+1
    S-->>C: Snapshot delta (tick N+1)
    Note over C: Interpola posição<br/>suavemente
    
    Note over S: Simula tick N+2
    S-->>C: Snapshot delta (tick N+2)
    Note over C: Renderiza a 60fps<br/>interpolando entre snapshots
```

### 6.3 Compensação de Latência

```mermaid
flowchart LR
    subgraph Servidor["Servidor"]
        Tick1[Tick N] --> Tick2[Tick N+1]
        Tick2 --> Tick3[Tick N+2]
    end
    
    subgraph Cliente["Cliente"]
        Snap1[Snapshot N] --> Interp1[Interpolação]
        Snap2[Snapshot N+1] --> Interp1
        Interp1 --> Render1[Render Frame 1<br/>60fps]
        
        Snap2 --> Interp2[Interpolação]
        Snap3[Snapshot N+2] --> Interp2
        Interp2 --> Render2[Render Frame 2<br/>60fps]
    end
    
    Tick1 -.->|envia| Snap1
    Tick2 -.->|envia| Snap2
    Tick3 -.->|envia| Snap3
    
    style Servidor fill:#fff3e0,stroke:#f57c00,stroke-width:2px
    style Cliente fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
```

---

## 7. Fluxo de uma Partida

### 7.1 Sequência de Conexão

```mermaid
sequenceDiagram
    participant Cliente
    participant Servidor
    
    Cliente->>Servidor: GET / (HTTP)
    Servidor-->>Cliente: HTML + JS + CSS
    
    Cliente->>Servidor: WebSocket Connect
    Servidor-->>Cliente: Connection OK
    
    Cliente->>Servidor: Join Lobby
    Servidor-->>Cliente: Lobby State
    
    Cliente->>Servidor: Ready / Start
    Servidor-->>Cliente: Match Starting
    
    Servidor-->>Cliente: Snapshot Inicial<br/>(estado do mapa, recursos)
    
    loop Game Loop
        Cliente->>Servidor: Comandos
        Servidor-->>Cliente: Snapshots/Eventos
    end
    
    Servidor-->>Cliente: Match Over
    
    Cliente->>Servidor: GET /results (HTMX)
    Servidor-->>Cliente: Tela de resultados
```

### 7.2 Ciclo de Vida de uma Partida

```mermaid
stateDiagram-v2
    [*] --> Lobby: Jogador entra
    
    Lobby --> Waiting: Aguardando<br/>outros jogadores
    Waiting --> Preparation: Todos prontos<br/>ou timeout
    
    Preparation --> EarlyGame: Timer 60s expira
    
    EarlyGame --> MidGame: 5 minutos
    MidGame --> LateGame: 12 minutos
    LateGame --> SuddenDeath: 18 minutos<br/>sem vencedor
    
    EarlyGame --> GameOver: Jogador eliminado
    MidGame --> GameOver: Jogador eliminado
    LateGame --> GameOver: Jogador eliminado
    SuddenDeath --> GameOver: Último jogador
    
    GameOver --> Results: Calcula estatísticas
    Results --> Lobby: Volta ao lobby
    Results --> [*]: Sai do jogo
    
    note right of Lobby
        Seleção de facção,
        chat, ready check
    end note
    
    note right of Preparation
        Construção inicial
        de torres
    end note
    
    note right of GameOver
        Vitória/derrota,
        XP, desbloqueios
    end note
```

### 7.3 Gerenciamento de Salas

```mermaid
flowchart TD
    Create[Criar Sala] --> CheckType{Tipo de Jogo?}
    
    CheckType -->|Solo| SoloRoom[Sala Solo<br/>1 jogador]
    CheckType -->|1v1| DuelRoom[Sala Duelo<br/>2 jogadores]
    CheckType -->|3 jogadores| FFARoom[Sala FFA<br/>3 jogadores]
    
    SoloRoom --> StartSolo[Iniciar Partida Solo]
    DuelRoom --> WaitDuel[Aguardar 2º jogador]
    FFARoom --> WaitFFA[Aguardar 2º e 3º jogadores]
    
    WaitDuel -->|Jogador entra| StartDuel[Iniciar Partida 1v1]
    WaitFFA -->|3 jogadores prontos| StartFFA[Iniciar Partida FFA]
    
    WaitDuel -->|Timeout| CancelDuel[Cancelar sala]
    WaitFFA -->|Timeout| CancelFFA[Cancelar sala]
    
    StartSolo --> GameLoop[Game Loop]
    StartDuel --> GameLoop
    StartFFA --> GameLoop
    
    style Create fill:#4caf50,stroke:#2e7d32,stroke-width:2px,color:#fff
    style GameLoop fill:#ff9800,stroke:#f57c00,stroke-width:2px
```

---

## 8. Modelo de Dados

### 8.1 Schema do Banco de Dados

```mermaid
erDiagram
    PLAYERS ||--o{ STATS : has
    PLAYERS ||--o{ UNLOCKS : has
    PLAYERS ||--o{ MATCH_HISTORY : played
    
    PLAYERS {
        int id PK
        text username UK
        datetime created_at
        int level
        int xp
    }
    
    STATS {
        int player_id FK
        int matches_played
        int matches_won
        int best_wave_survival
        int total_kills
    }
    
    UNLOCKS {
        int player_id FK
        text cosmetic_id
        datetime unlocked_at
    }
    
    MATCH_HISTORY {
        int player_id FK
        text match_id
        datetime played_at
        text faction
        boolean won
        int duration_seconds
    }
    
    REPLAYS {
        int id PK
        text match_id UK
        datetime recorded_at
        blob data
        int duration_seconds
    }
```

### 8.2 Fluxo de Persistência

```mermaid
flowchart LR
    subgraph Runtime["Durante Partida"]
        GameLoop[Game Loop] --> RecordCommands[Grava Comandos]
        RecordCommands --> ReplayBuffer[Buffer de Replay]
    end
    
    subgraph PostMatch["Pós-Partida"]
        ReplayBuffer --> Compress[Comprimir Dados]
        Compress --> SaveReplay[Salvar Replay<br/>SQLite]
        
        GameLoop --> CalcStats[Calcular Estatísticas]
        CalcStats --> UpdateProfile[Atualizar Perfil<br/>XP, nível, stats]
        UpdateProfile --> CheckUnlocks[Verificar Desbloqueios]
        CheckUnlocks --> SaveUnlocks[Salvar Desbloqueios]
    end
    
    Runtime --> PostMatch
    
    style Runtime fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    style PostMatch fill:#fff3e0,stroke:#f57c00,stroke-width:2px
```
