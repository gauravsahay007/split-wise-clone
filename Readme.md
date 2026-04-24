```mermaid
graph TD
    %% Define actors and external systems
    subgraph Clients
        U[End User / Developer]
    }

    subgraph "External World"
        GAPI[GitHub Public API]
    }

    subgraph "Go Application (Main Thread)"
        M[main.go]

        subgraph "Layer 1: api/v1 (REST API Server)"
            EH[User Handlers<br/>POST /users<br/>GET /activity]
        }

        subgraph "Layer 2: business/v1"
            BL[User & Activity Manager<br/>(Interface & Implementation)]
        }

        subgraph "Layer 3: entities/repositories"
            RP[PostgreSQL Repository<br/>(Interface & Implementation)]
        }

        subgraph "Storage"
            DB[(PostgreSQL Database)]
        }

        subgraph "Layer 4: jobs"
            CR[Cron Scheduler<br/>(robfig/cron)]
            GW[GitHub Worker<br/>(Fetches multiple users concurrently)]
        }

        subgraph "Layer 5: externals"
            GC[GitHub Client<br/>(Handles HTTP, Rate Limits)]
        }
    }

    U -->|POST /users| EH
    EH --> BL
    BL --> RP
    RP --> DB
```
