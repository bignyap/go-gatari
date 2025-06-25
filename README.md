# Go Admin
A robust multi-tenant application designed to manage organizations, subscriptions, and API usage with fine-grained control over resources and permissions. It supports subscription tiers, custom endpoint pricing, and detailed billing management for different organizations within your platform.

### Key Features:
- **Multi-tenancy**: Seamlessly manage multiple organizations with unique configurations, including realm, support, and permissions.
- **Subscription Management**: Supports various subscription tiers, API usage limits, and custom pricing models for API endpoints.
- **Billing and Usage Tracking**: Automatically track API usage and generate detailed billing reports based on subscription and endpoint-specific pricing.
- **Resource Permissions**: Fine-grained access control to resources through permission management for organizations.

### Steps to start the application:

1. Create a `.env` file similar to `.env.sample` and fill in the variable values. Place it in the root directory if using Docker Compose.

2. Build the container:

    ```bash
    go mod tidy && \
    go mod vendor && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/go-admin cmd/go-admin/main.go && \
    docker build -t go-admin:latest . && \
    rm -f build/go-admin
    ```

3. Start the application:

    ```bash
    docker-compose up -d
    ```

4. Open [http://localhost:8081](http://localhost:8081) for the Swagger API documentation.

5. Access [http://localhost:8080](http://localhost:8080) for the APIs.

