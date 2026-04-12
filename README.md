# go-genesis-case-task

GitHub Release Notifier: A Go-based service that monitors repositories and sends email alerts for new releases.

This microservice is built with Go to monitor new releases in GitHub repositories and automatically notify users via email.
It follows Clean Architecture principles to ensure maintainability, scalability, and ease of testing.

Technical Stack
The project utilizes Go 1.26 as the primary language and the Gin Gonic framework for high-performance HTTP routing.
PostgreSQL serves as the persistent storage, managed via the pgx connection pool for optimal concurrency.
Monitoring is facilitated by Prometheus, while background task scheduling is managed through gocron/v2.

📂 Project Structure

/cmd — entry point

/internal/domain — domain models

/internal/usecase — business logic

/internal/infrastructure — external infrastructure (DB, GitHub, Email).

/internal/worker — background tasks ()Scanner).

/pkg - database clients


Quick Start
Environment Configuration
Create a .env file in the project root based on the provided template ~/.env.example.

You must define the server port, the PostgreSQL DSN for database connectivity, and your GitHub personal access token.
The scanner interval should be specified using Go duration format, such as 5m or 1h.
For notifications, configure the SMTP host, port, and credentials, including a Google App Password if using Gmail.

Execution via Docker
Deploy the entire stack by running

docker-compose up --build.

This command initializes the database, applies migrations, and starts the API service.
The application will be accessible at http://localhost:8080.

GitHub Rate Limit Handling
The service actively manages communication with the GitHub API to prevent token blocking.
If the scanner encounters a rate limit error, it extracts the reset timestamp from the API response headers.
An internal safety flag is then set to skip subsequent scanning cycles until the reset time has passed.
This proactive approach saves system resources and ensures the service remains a "good citizen" within the GitHub ecosystem.

