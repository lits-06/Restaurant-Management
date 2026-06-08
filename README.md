# 🍽️ Restaurant Management Microservices

**Modern microservices-based restaurant management system** built with Golang, Clean Architecture, and gRPC.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Architecture](https://img.shields.io/badge/Architecture-Clean-brightgreen)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
[![Protocol](https://img.shields.io/badge/Protocol-gRPC-244c5a?style=flat&logo=grpc)](https://grpc.io/)
[![Status](https://img.shields.io/badge/Status-In_Development-yellow)](./IMPLEMENTATION_STATUS.md)

---

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [Services](#services)
- [Documentation](#documentation)
- [Development](#development)
- [Progress](#progress)

---

## 🎯 Overview

A production-ready restaurant management platform consisting of:
- **9 microservices** for core business functions
- **API Gateway** for unified HTTP interface
- **Clean Architecture** for maintainability and testability
- **gRPC** for efficient inter-service communication
- **Docker** for containerization and deployment

**Target Users:** Restaurants, cafes, food courts needing digital management solutions.

---

## ✨ Features

### Core Capabilities

- 🔐 **Authentication & Authorization** - JWT-based auth with refresh tokens
- 👥 **User Management** - Staff, managers, admin roles
- 🪑 **Table Management** - Real-time table status tracking
- 📋 **Menu Management** - Categories, items, pricing
- 🛒 **Order Processing** - Full order lifecycle management
- 💳 **Payment Handling** - Multiple payment methods
- 📦 **Inventory Control** - Stock tracking & alerts
- 🔔 **Notifications** - Multi-channel messaging
- 📊 **Reporting & Analytics** - Business insights

### Technical Features

- **Clean Architecture** - Clear separation of concerns
- **Repository Pattern** - Database-agnostic data access
- **Dependency Injection** - Easy testing and swapping implementations
- **gRPC Communication** - Fast and type-safe service-to-service calls
- **Middleware Support** - Logging, recovery, authentication
- **Docker Ready** - Complete containerization
- **In-Memory Implementations** - Fast development without database setup

---

## 🏗️ Architecture

### System Overview

```
┌─────────────┐
│   Client    │
│  (Browser)  │
└──────┬──────┘
       │ HTTP REST
       ↓
┌─────────────────────────────────────────┐
│          API Gateway (Port 8080)        │
│  - HTTP to gRPC translation             │
│  - Authentication middleware            │
│  - Request/Response formatting          │
└──────┬──────────────────────────────────┘
       │ gRPC
       ↓
┌──────────────────────────────────────────────────┐
│              Microservices Network               │
│                                                  │
│  Auth (50051)      Menu (50054)      Notify (50058)│
│  User (50052)      Order (50055)     Report (50059)│
│  Table (50053)     Payment (50056)                │
│                    Inventory (50057)              │
└──────────────────────────────────────────────────┘
       │
       ↓
┌──────────────────────────────────────────────────┐
│         Infrastructure                           │
│  - PostgreSQL (Database)                         │
│  - Redis (Cache)                                 │
└──────────────────────────────────────────────────┘
```

## Project Structure

```
restaurant-management/
├── services/
│   ├── auth-service/
│   ├── user-service/
│   ├── table-service/
│   ├── menu-service/
│   ├── order-service/
│   ├── payment-service/
│   ├── inventory-service/
│   ├── notification-service/
│   └── report-service/
├── shared/
│   └── pkg/
│       ├── logger/
│       ├── config/
│       ├── middleware/
│       ├── jwt/
│       ├── errors/
│       └── utils/
├── api-gateway/
├── proto/
├── scripts/
├── docker-compose.yml
└── Makefile
```

## Setup Instructions

1. Install dependencies:
   - Go 1.21+
   - Protocol Buffers compiler (protoc)
   - Docker & Docker Compose

2. Generate proto files:
   ```bash
   make proto
   ```

3. Build all services:
   ```bash
   make build
   ```

4. Run with Docker Compose:
   ```bash
   docker-compose up
   ```

## Services

1. **Auth Service** (Port: 50051) - Authentication & Authorization
2. **User Service** (Port: 50052) - User management
3. **Table Service** (Port: 50053) - Table management
4. **Menu Service** (Port: 50054) - Menu & items management
5. **Order Service** (Port: 50055) - Order processing
6. **Payment Service** (Port: 50056) - Payment handling
7. **Inventory Service** (Port: 50057) - Inventory tracking
8. **Notification Service** (Port: 50058) - Notifications
9. **Report Service** (Port: 50059) - Reports & analytics
10. **API Gateway** (Port: 8080) - HTTP to gRPC gateway

## Development

See individual service READMEs for details.
