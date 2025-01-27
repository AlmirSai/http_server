# Microservices Architecture Project

## Overview
This project implements a modern microservices architecture designed to handle various aspects of a social networking application. Built with scalability, maintainability, and reliability in mind, the system follows industry best practices for microservices development.

## Architecture

### System Components
1. **API Gateway**
   - Primary entry point for all client requests
   - Implements request routing and load balancing
   - Handles API versioning and documentation
   - Rate limiting and request validation

2. **Auth Service**
   - Manages authentication and authorization
   - JWT token generation and validation
   - Role-based access control (RBAC)
   - Session management with Redis

3. **User Service**
   - User profile management
   - Account settings and preferences
   - User search and discovery
   - Profile verification

4. **Media Service**
   - Media file upload and processing
   - Image resizing and optimization
   - Content delivery network (CDN) integration
   - File format validation

5. **Post Service**
   - Content creation and management
   - Post metadata handling
   - Content moderation
   - Feed generation

## Technology Stack

### Core Technologies
- **Backend**: Go 1.19+
- **Database**: PostgreSQL 14+
- **Caching**: Redis 6+
- **Container Runtime**: Docker
- **Orchestration**: Kubernetes

### Monitoring & Observability
- **Metrics**: Prometheus
- **Visualization**: Grafana
- **Logging**: Zap (structured JSON logging)
- **Tracing**: OpenTelemetry

## Project Structure
```
├── api-gateway/       # API Gateway service
├── auth-service/      # Authentication service
├── media-service/     # Media handling service
├── post-service/      # Post management service
├── user-service/      # User management service
├── shared/           # Shared libraries and utilities
└── docker-compose.yml # Local development setup
```

## Development Setup

### Prerequisites
1. Go 1.19 or later
2. Docker and Docker Compose
3. PostgreSQL 14+
4. Redis 6+
5. Make (optional)

### Environment Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd http_server
```

2. Configure environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Start the services:
```bash
# Using Docker Compose
docker-compose up -d

# Verify services
docker-compose ps
```

### Configuration

Each service has a `config.yaml` file with the following sections:
- Server configuration (ports, timeouts)
- Database connections
- JWT settings
- Logging configuration
- Metrics and monitoring

## API Documentation

API documentation is available through Swagger UI at:
- Development: `http://localhost:8080/swagger/`
- Production: `https://api.yourdomain.com/swagger/`

## Monitoring

### Metrics
- Prometheus endpoints: `/metrics` on each service
- Grafana dashboards for visualization
- Custom metrics for business KPIs

### Logging
- Structured JSON logging with Zap
- Log levels: DEBUG, INFO, WARN, ERROR
- Correlation IDs for request tracing

## Testing

### Running Tests
```bash
# Unit tests
go test ./...

# Integration tests
make integration-test

# Load tests
make load-test
```

## Deployment

### Production Deployment
1. Build Docker images:
```bash
make build-images
```

2. Deploy to Kubernetes:
```bash
kubectl apply -f k8s/
```

### Scaling
- Horizontal scaling through Kubernetes
- Vertical scaling for databases
- Cache layer scaling with Redis cluster

## Contributing

### Development Workflow
1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality
4. Implement changes
5. Ensure tests pass
6. Submit pull request

### Code Style
- Follow Go standard formatting (gofmt)
- Use golangci-lint for static analysis
- Write meaningful commit messages
- Document public APIs

## Troubleshooting

### Common Issues
1. Service connectivity issues
   - Check network configurations
   - Verify service health endpoints
   - Review logs for connection errors

2. Database connection issues
   - Verify credentials in .env
   - Check database logs
   - Ensure migrations are up to date

## License
This project is licensed under the MIT License - see the LICENSE file for details.

## Support
For support and questions:
- Create an issue in the repository
- Contact the development team
- Check the documentation wiki