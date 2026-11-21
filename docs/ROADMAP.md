# Roadmap

This file provides an overview of the direction this project is heading. During the development process, new stages may be added, and existing ones may be modified.

**Legend:**
- [x] Done
- [ ] Planned

--- 

## M0: Core.
In this phase, the focus is on core of the project. This includes setting up the initial structure of the server application, defining the main entrypoint, setting up gin with swagger endpoint and gorm.

- [x] **Structure & Tooling:**
  - [x] Initialize Go module and project layout
  - [x] Set up config loader
  - [x] Create basic `Makefile`
- [x] **Web Server & Docs:**
  - [x] Initialize Gin router with Recovery, Logger, and CORS middlewares
  - [x] Add `GET /health` endpoint and graceful shutdown handling
  - [x] Integrate `swaggo/swag` and expose `/swagger/*any` UI endpoint (if config specifically says it should be public)
  - [x] Make and integrate into router `ErrorHandler` middleware
- [x] **Database:**
  - [x] Configure GORM connection with PostgreSQL driver
  - [x] Set up database migrations using goose

## M1: Auth and Users.
In this phase, the focus is on establishing basic authorization and a user database.

- [x] **Models & Repository:**
  - [x] Define `User` GORM model
  - [x] Implement User repository
- [ ] **Authorization:**
  - [x] Implement password hashing
  - [x] Add hashing params to config (`threads` and `memory` for argon2)
  - [ ] Add remaining hashing params to config
  - [ ] Replace auth/argon2 with custom HashService (use semaphore for threading)
  - [ ] Implement JWT service (Access + Refresh token generation and parsing)
  - [ ] Create `POST /api/v1/auth/register` and `POST /api/v1/auth/login` handlers
  - [ ] Add Auth middleware to validate Bearer token and inject `user_id` into `gin.Context`
- [ ] **Profiles & Discovery:**
  - [ ] Create `GET /api/v1/users/me` and `PATCH /api/v1/users/me` handlers
  - [ ] Create `GET /api/v1/users?q=...` endpoint for searching users by username


## M2: Chats and History.
In this phase, the focus is on designing the chat database schema and building REST endpoints for conversation management and message persistence. Having a working HTTP layer makes it easy to test all business logic and permissions via Swagger before introducing WebSockets.

- [ ] **Models & Database:**
  - [ ] Define `Chat`, `Participant`, and `Message` GORM models
- [ ] **Repository Layer:**
  - [ ] Implement Chat repository
  - [ ] Implement Message repository
- [ ] **Chat Endpoints:**
  - [ ] Create `POST /api/v1/chats` to start a 1-on-1 chat (return the existing chat if a dialog between the two users already exists)
  - [ ] Create `GET /api/v1/chats` to retrieve the user's conversation list ordered by the latest message timestamp
- [ ] **Message Endpoints:**
  - [ ] Create `POST /api/v1/chats/:id/messages` to send a text message over REST
  - [ ] Create `GET /api/v1/chats/:id/messages` with pagination (`?before=<created_at_or_id>&limit=50`) for loading message history

---
more soon.