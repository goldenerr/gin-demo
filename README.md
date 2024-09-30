# Gin Web Application Scaffold

This project is a scaffold for building web applications using the Gin framework in Go. It provides a solid foundation with several key features implemented, making it easier to start developing robust web services.

## Features

1. **Configuration Management**
   - Uses Viper for flexible configuration handling
   - Supports YAML configuration files

2. **Logging**
   - Implements structured logging using Zap
   - Supports log rotation and compression
   - Configurable log levels and output formats
   - Log file naming based on date (Beijing time)
   - Customizable log file size and unit (KB, MB, GB, TB)

3. **Database Integration**
   - GORM integration for database operations
   - Supports MySQL (easily extendable to other databases)
   - SQL query logging with request ID tracing

4. **Request Tracing**
   - Implements request ID generation and propagation
   - Enables easy tracking of requests across the application

5. **Middleware**
   - Custom logging middleware for detailed request/response logging
   - Captures and logs full request and response bodies

6. **MVC-like Structure**
   - Clear separation of concerns with controllers, services, and repositories
   - Easy to extend and maintain

7. **Error Handling**
   - Centralized error handling and logging

8. **API Versioning**
   - Built-in support for API versioning

9. **Environment-based Configuration**
   - Easy to switch between development, testing, and production environments

## Project Structure
gin-demo/
├── cmd/
│ └── main.go # Application entry point
├── config/
│ └── config.go # Configuration handling
├── internal/
│ ├── app/
│ │ └── app.go # Application setup
│ ├── controller/
│ │ └── user_controller.go # User-related HTTP handlers
│ ├── database/
│ │ └── database.go # Database connection and management
│ ├── handler/
│ │ └── user.go # User-related request handlers
│ ├── logger/
│ │ └── logger.go # Logging setup and utilities
│ ├── middleware/
│ │ └── logger.go # Custom logging middleware
│ ├── models/
│ │ ├── user.go # User model
│ │ └── user_query.go # User query model
│ ├── repository/
│ │ └── user_repository.go # User data access layer
│ ├── router/
│ │ └── router.go # Route definitions
│ ├── service/
│ │ └── user_service.go # User business logic
│ └── tracing/
│ └── tracing.go # Request tracing utilities
├── config.yaml # Application configuration
├── go.mod
└── go.sum
## Getting Started

1. Clone the repository
2. Update the `config.yaml` file with your specific configurations
3. Run `go mod tidy` to ensure all dependencies are correctly installed
4. Start the application with `go run cmd/main.go`

## Customization

This scaffold is designed to be easily customizable. You can add new controllers, services, and models as needed for your specific application requirements.

## Key Components

- **User Service**: Handles user-related business logic without direct logger dependency
- **User Repository**: Manages user data persistence
- **User Controller**: Handles HTTP requests related to users
- **Middleware**: Includes custom logging middleware for request/response tracking with full body logging
- **Database**: Configured for MySQL with GORM, including SQL query logging with request ID tracing
- **Logging**: Utilizes Zap for structured logging with rotation, compression, and Beijing time-based file naming

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT License](LICENSE)