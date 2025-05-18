# Highload Social Network Platform

A high-performance, scalable social network platform built with Go and a microservices architecture designed to handle intensive workloads and demonstrate advanced distributed systems concepts.

## Project Overview

Highload is a production-ready social network platform that showcases enterprise-level architecture patterns and performance optimization techniques. The platform features:

- **User Management**: Authentication, profile management, and social relationships
- **Content Sharing**: Post creation and personalized feed generation
- **Real-time Messaging**: Scalable chat system with unread message tracking
- **Live Updates**: WebSocket-based real-time feed updates
- **Activity Metrics**: Distributed counter system for tracking user activities

## System Architecture

The platform is built on a robust microservices architecture designed for high scalability, performance, and resilience:

### Core Services Layer

- **Web API Gateway**: RESTful API service using Chi router with JWT authentication, serving as the entry point for all client requests
- **Users Service**: Manages user registration, authentication, profiles, and social relationships
- **Posts Service**: Handles content creation with feed generation through Kafka-based event processing
- **Chats Service**: Provides real-time messaging capabilities with horizontal database sharding
- **Counters Service**: Tracks activity metrics with Redis-backed storage for high performance

### Data Storage Layer

- **MySQL Cluster**: Master-slave replication setup with semi-synchronous replication for data consistency
  - Read/write separation via HAProxy for optimal performance
  - Horizontal sharding for message storage across dedicated nodes
  
- **Redis**: Multi-purpose in-memory data store used for:
  - High-speed caching of frequently accessed data
  - Feed storage and retrieval
  - Distributed lock management
  - Counter and metrics storage

- **Kafka**: Event streaming platform for reliable asynchronous processing
  - Used primarily for feed updates and event-driven workflows

### Infrastructure Components

- **Docker & Docker Compose**: Full containerization of all services for consistent deployment
- **Nginx**: Front-facing load balancer and reverse proxy
- **HAProxy**: Database connection pooling and traffic direction
- **Centrifugo**: WebSocket server enabling real-time updates to client applications
- **Prometheus & Grafana**: Comprehensive monitoring and performance visualization

## Key Technical Features

### Performance Optimizations

- **Database Read/Write Splitting**: Directs read operations to slave replicas, optimizing database load
- **Multi-level Caching**: Redis-based caching for frequently accessed data
- **Connection Pooling**: Efficient resource utilization across services
- **Database Sharding**: Horizontal data distribution for improved throughput
- **Asynchronous Processing**: Non-blocking operations via Kafka for improved responsiveness

### High Availability & Resilience

- **Semi-synchronous Replication**: Prevents data loss during master database failures
- **Service Replication**: Multiple instances of critical services
- **Circuit Breaking**: Prevents cascading failures in the service mesh
- **Saga Pattern**: Maintains data consistency across distributed transactions
- **GTID-based Replication**: Simplifies failover and recovery processes

### Security Implementation

- **JWT-based Authentication**: Secure, stateless authentication mechanism
- **Password Hashing**: Industry-standard bcrypt implementation
- **HTTPS Support**: TLS encryption for all client communication
- **Secure Cookie Management**: Protection against XSS and CSRF attacks
- **Input Validation**: Comprehensive request validation

## Getting Started

### Prerequisites

- Go 1.24+
- Docker and Docker Compose
- Make

### Installation

Clone the repository:
```bash
git clone https://github.com/ravilushqa/highload.git
cd highload
```

### Running the Project

Start the entire platform:
```bash
make setup
```

Start with monitoring enabled:
```bash
make setup-with-monitoring
```

Stop all services:
```bash
make down
```

## Service Endpoints

- Web API: `http://localhost:80`
- Monitoring Dashboard: `http://localhost:3000` (Grafana)
- Metrics: `http://localhost:9090` (Prometheus)

## Skills & Techniques Demonstrated

- **Distributed Systems Design**: Microservices with well-defined boundaries
- **High-Load Handling**: Database optimization, caching, connection pooling
- **Go Programming**: Idiomatic Go code with effective concurrency patterns
- **Database Engineering**: Replication, sharding, query optimization
- **DevOps Practices**: Containerization, orchestration, monitoring
- **Messaging Patterns**: Event sourcing with Kafka
- **Caching Strategies**: Multi-level caching with Redis

## License

This project is licensed under the MIT License - see the LICENSE file for details.