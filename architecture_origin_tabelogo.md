# Original Tabelogo Architecture Analysis

## 1. Overview
Based on the analysis of the `Leon180/tabelogo` repository, the original project is a microservices-based application designed to aggregate restaurant information, primarily focusing on integrating Google Maps data with Tabelog (a Japanese restaurant guide).

## 2. Core Services

### 2.1 Broker Service (`broker-service`)
- **Role**: Entry point / Gateway.
- **Dependencies**: Links to `authenticate-service`, `google-map-service`, `tabelog-spider-service`, `logger-service`, `mail-service`.
- **Function**: Likely acts as an API Gateway or BFF (Backend for Frontend), routing requests to appropriate backend services.

### 2.2 Authentication Service (`authenticate-service`)
- **Role**: User identity and session management.
- **Database**: PostgreSQL (`auth_db` equivalent, likely tables in `tabelogo` DB).
- **Dependencies**: Redis (Session management), RabbitMQ.
- **Key Features**: JWT handling, likely user registration/login.

### 2.3 Google Map Service (`google-map-service`)
- **Role**: Interface with Google Maps API.
- **Function**: Handles place search, details retrieval, and autocomplete.
- **Dependencies**: `authenticate-service` (likely for internal auth).

### 2.4 Tabelog Spider Service (`tabelog-spider-service`)
- **Role**: Web crawler for Tabelog.
- **Function**: Fetches restaurant details from Tabelog based on search queries.
- **Trigger**: Triggered when a user clicks a "Tabelogo" button for a Japanese restaurant location.

### 2.5 Mail Service (`mail-service`)
- **Role**: Email notifications.
- **Dependencies**: MailHog (for testing/dev).
- **Function**: Sending registration emails, etc.

### 2.6 Logger Service (`logger-service`)
- **Role**: Centralized logging.
- **Database**: MongoDB (`logs` database).
- **Function**: Stores system logs.

### 2.7 Listener Service (`listener-service`)
- **Role**: Event consumer.
- **Dependencies**: RabbitMQ.
- **Function**: Listens to asynchronous events (likely from RabbitMQ) to trigger background tasks.

### 2.8 Front-End (`front-end`)
- **Role**: Web UI.
- **Tech**: Go templates (implied by `templates` dir) or a separate Go server serving static assets.

## 3. Infrastructure & Data Stores

### 3.1 Databases
- **PostgreSQL**: Main relational database (`tabelogo` DB). Stores user data, likely places/restaurants.
- **MongoDB**: Stores logs (`logs` DB).
- **Redis**:
  - `redis-master-session`: User sessions.
  - `redis-master-place`: Caching place information.
  - `redis-master-tabelogo`: (Commented out in docker-compose) Caching Tabelog data.

### 3.2 Messaging
- **RabbitMQ**: Asynchronous communication between services (e.g., triggering emails, logging, or spider tasks).

## 4. Key Workflows (Inferred)

1.  **User Registration**: Front-end -> Broker -> Authenticate -> Mail (via RabbitMQ/Listener).
2.  **Place Search**: Front-end -> Broker -> Google Map Service -> Google Maps API.
3.  **Tabelog Integration**:
    - User selects a place.
    - Request to Broker -> Tabelog Spider Service.
    - Spider fetches data -> Returns to user (or stores in DB/Cache).

## 5. Comparison with V2 (Current)

| Feature | Original (V1) | V2 (Current Plan) |
| :--- | :--- | :--- |
| **Architecture** | Microservices (Docker Swarm/Compose) | Microservices (K8s ready, DDD) |
| **Gateway** | Broker Service (Custom Go) | API Gateway (Gin + Custom) |
| **Communication** | HTTP + RabbitMQ | gRPC (Internal) + Kafka |
| **DB Strategy** | Shared Postgres (mostly) + Mongo | Database per Service (Strict) |
| **Spider** | On-demand (Triggered by user) | Job-based (Background + Kafka) |
| **Booking** | N/A (Not explicitly seen) | Core Feature (Event Sourcing) |

## 6. Observations for V2 Alignment
- **Concept Retention**: The core value of "Google Maps + Tabelog" integration is preserved.
- **Evolution**: V2 moves towards a more robust, scalable architecture (gRPC, Kafka, DB per Service) suitable for enterprise-level demonstration.
- **Missing in V1**: Booking system, complex restaurant aggregation (V1 seems more real-time proxy/scrape).
