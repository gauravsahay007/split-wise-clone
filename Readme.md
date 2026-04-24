graph TD
    %% Define actors and external systems
    subgraph Clients
        U[End User / Developer]
    }

    subgraph "External World"
        GAPI[GitHub Public API]
    }

    %% Main Application Boundary (Representing Docker Container)
    subgraph "Go Application (Main Thread)"
        M[main.go]

        %% Layer 1: API (v1)
        subgraph "Layer 1: api/v1 (REST API Server)"
            EH[User Handlers<br/>POST /users<br/>GET /activity]
        }

        %% Layer 2: Business Logic (v1)
        subgraph "Layer 2: business/v1"
            BL[User & Activity Manager<br/>(Interface & Implementation)]
        }

        %% Layer 3: Entities/Repositories
        subgraph "Layer 3: entities/repositories"
            RP[PostgreSQL Repository<br/>(Interface & Implementation)]
        }

        %% The Database
        subgraph "Storage"
            DB[(PostgreSQL Database)]
        }

        %% The Job Scheduler & Worker
        subgraph "Layer 4: jobs"
            CR[Cron Scheduler<br/>(robfig/cron)]
            GW[GitHub Worker<br/>(Fetches multiple users concurrently)]
        }

        %% Layer 5: Externals
        subgraph "Layer 5: externals"
            GC[GitHub Client<br/>(Handles HTTP, Rate Limits)]
        }

    }

    %% Data Flow and Interactions

    %% 1. User registers a username via API
    U -->|1. POST /users| EH
    EH -->|Calls biz logic| BL
    BL -->|2. Stores target user| RP
    RP --> DB

    %% 2. User queries for fetched activity
    U -->|4. GET /activity| EH
    EH -->|Calls biz logic| BL
    BL -->|Fetches stored events| RP
    RP --> DB

    %% 3. The background job flow (The "Glue" concept)
    M -->|5. Starts as Goroutine| CR
    CR -->|3. Triggers task<br/>(Every X hours)| GW
    GW -->|Fetches Users<br/>Ready for Sync| BL

    %% 4. Concurrent execution (Fan-out)
    GW -.->|6. Spawns Goroutines| GW1(User A Goroutine)
    GW -.->|6. Spawns Goroutines| GW2(User B Goroutine)
    GW -.->|6. Spawns Goroutines| GW3(User ... Goroutine)

    %% 5. Workers calling external API
    GW1 & GW2 & GW3 -->|Calls GitHub via client| GC
    GC -->|7. API Requests| GAPI
    GAPI -->|Returns Activity JSON| GC
    GC -->|Returns Models| GW1 & GW2 & GW3

    %% 6. Workers saving the output
    GW1 & GW2 & GW3 -->|8. Save Activities| BL
    BL -->|Filters events| RP
    RP -->|Saves data| DB
