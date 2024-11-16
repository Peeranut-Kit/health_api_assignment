# Hospital Api Assignment
The Hospital Api System is a backend service designed to integrate with Hospital Middleware system to search and display patient information from Hospital Information Systems (HIS) and facilitate patient information searches. It provides a secure API for staff members to search for patient details, restricted to their associated hospital, and supports user authentication and management.

## Features
- Search and display patient information using APIs provided by hospitals.
- Staff member registration.
- Secure staff login using encrypted credentials.
- Compatibility with Docker, Nginx, PostgreSQL, and the Gin framework for scalability and ease of deployment.
- Unit-tested for robust and reliable functionality.

## Tech Stack
- Programming Language: Go
- Framework: Gin
- Database: PostgreSQL
- Containerization: Docker, Docker Compose
- Reverse Proxy Server: Nginx

## Setup Instructions
### Prerequisites
- Docker and Docker Compose

### Steps
1. Clone the repository.
   ```
   git clone https://github.com/Peeranut-Kit/health_api_assignment.git
   
   cd health_api_assignment
   ```

2. Use a .env file provided in the repository. Please ensure that all variable in .env file is not the same as any service running on your system (port number, database name).
   Also ensure that all port and container name in docker-compose.yml file do not conflict with any on your local machine.

3. Build and run the application. You can specify the numbers of server running in the network for scaling by changing --scale api-service parameter. (3 servers are specified in this case.)
   ```
   docker compose up -d --scale api-service=3 --build
   ```

4. Access the APIs via:
   Base URL (NGINX): http://localhost:3000

## Unit Testing
Run unit tests using:
```
cd test

go test ./...  
```
or for verbose output
```
go test ./... -v
```

## API Specification
1. Create a New Staff Member
Endpoint: POST /staff/create

2. Staff Login
Endpoint: POST /staff/login

3. Search for a Patient
Endpoint: GET /patient/search
*Requires Login

### Additional endpoints:
- Swagger UI
Endpoint: POST /swagger/index.html
- NGINX health check
Endpoint: /health
- test reverse proxy retreiving hostname of server
Endpoint: /ping

## Deliverables
Development planning documentation including other information such as project structure, API specifications, database schemas and ER diagram at https://docs.google.com/document/d/1sukcTe2uzBExHBINN1sBNdF5ntMXqG3L5DnGt_jMXxM/edit?usp=sharing
