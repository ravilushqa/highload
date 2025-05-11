# Highload Social Network Platform

A high-performance, scalable social network platform built with Go and a microservices architecture designed to handle intensive workloads.

## Architecture Diagram

![Highload Social Network Architecture](docs/architecture.svg)

The architecture is also available in PlantUML format:
- [PlantUML Diagram Source](docs/architecture.puml) - For custom rendering

## Project Overview

This project demonstrates a robust implementation of distributed systems principles through a social networking platform with:

- User authentication and profile management
- Posts creation and feed generation
- Real-time chat capabilities
- Social relationship management (friend requests)
- Activity counters and metrics

## Technical Architecture

_Refer to the architecture diagram above for a visual representation of the system components and their interactions. The SVG format ensures the diagram is viewable in any modern browser without additional tools._


### Backend Services (Go)

- **Web API Gateway**: RESTful API serving clients, using Chi router and JWT authentication
- **User Service (gRPC)**: Profile management with bcrypt password hashing and social graph features
- **Posts Service (gRPC)**: Content creation with asynchronous event processing via Kafka
- **Chat Service (gRPC)**: Real-time messaging with database sharding for horizontal scaling
- **Counters Service (gRPC)**: Activity tracking and metrics with Redis-backed storage

### Data Layer

- **MySQL**: Master-slave replication setup with read/write separation via HAProxy
- **Database Sharding**: Horizontal partitioning for message storage across dedicated nodes
- **Redis**: High-speed caching, feed storage, and distributed lock management
- **Kafka**: Event streaming for consistent asynchronous processing

### Architecture Patterns

- **Dependency Injection**: Using Uber Fx/Dig for clean component wiring
- **Saga Pattern**: Maintaining data consistency across microservices
- **Circuit Breaking**: Resilient service communication
- **Repository Pattern**: Clean data access layer implementation
- **CQRS**: Command-query responsibility separation where appropriate

### Infrastructure

- **Docker & Docker Compose**: Full containerization of all services
- **Load Balancing**: Service replication with Nginx
- **HAProxy**: Database connection pooling and traffic direction
- **Prometheus & Grafana**: Comprehensive monitoring and visualization
- **Centrifugo**: WebSocket server for real-time updates

## Technical Features

### Performance Optimizations

- Database read/write splitting for high throughput
- Redis caching for frequently accessed data
- Connection pooling for efficient resource utilization
- Horizontal scaling capabilities for all services
- Database sharding for distributing data load

### Scalability Design

- Stateless services allowing horizontal scaling
- Message-based asynchronous processing
- Master-slave database replication
- Data sharding across multiple database nodes
- Containerized deployment for flexible scaling

### Security Implementations

- JWT-based authentication
- Password hashing with bcrypt
- HTTPS support via Nginx
- Secure cookie management
- Input validation

## Development

### Prerequisites

- Go 1.24+
- Docker and Docker Compose
- Make
- PlantUML (optional, for diagram rendering)

### Running the Project

Start the entire platform:
```
make setup
```

With monitoring:
```
make setup-with-monitoring
```

Access specific databases:
```
make exec_master      # Connect to master database
make exec_slave1      # Connect to first slave
```

Regenerate protocol buffers:
```
make gen-proto
```

View architecture diagram:
```
# The SVG diagram can be viewed directly in any browser
# or SVG-compatible viewer without additional tools

# If you want to modify the PlantUML diagram and regenerate:
brew install plantuml    # macOS with Homebrew
apt-get install plantuml # Debian/Ubuntu
plantuml docs/architecture.puml
```

## API Services

- Web API: `http://localhost:80`
- Monitoring: `http://localhost:3000` (Grafana)
- Metrics: `http://localhost:9090` (Prometheus)

## Skills Demonstrated

- **Distributed Systems Design**: Microservices with well-defined boundaries and communication
- **High-Load Handling**: Database optimization, caching strategies, connection pooling
- **Go Programming**: Clean, idiomatic Go code with effective concurrency patterns
- **Database Engineering**: Replication, sharding, and query optimization
- **DevOps Practices**: Containerization, orchestration, and monitoring
- **Protocol Design**: Well-structured gRPC service definitions
- **Messaging Patterns**: Event sourcing with Kafka
- **Caching Strategies**: Multi-level caching with Redis
- **Security Implementation**: Authentication, authorization, and secure communication
