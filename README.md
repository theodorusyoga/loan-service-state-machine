# Loan Service State Machine

This project implements a loan service with state machine workflow management to handle the different stages of loan processing from application (proposed status) to disbursement.

### Introduction Video

![](https://drive.google.com/uc?export=view&id=1b1TFzZKkkFYHHHSqSYhSzYSl-GUvt8Q_)

Or click the link here if the video doesn't show up:
https://drive.google.com/file/d/1b1TFzZKkkFYHHHSqSYhSzYSl-GUvt8Q_/view?usp=sharing

### Tech Stack

- Go: Core programming language
- CockroachDB: Distributed SQL database for persistent storage
- GORM: ORM library for database operations
- Looplab/FSM: Finite State Machine implementation for loan status transitions
- Uber/fx: Dependency injection framework for modular application structure
- Echo: HTTP web framework for API endpoints
- Swagger: API documentation and testing
- Air: Live reloading for development

### Features

- Loan application and processing workflow
- Borrower and lender management
- Employee (field officer and approver) management
- Document tracking
- State transitions: application (proposal) → approval → investment → disbursement

## Getting Started

### Prerequisites
- Go 1.18+
- CockroachDB instance
- Configure database connection in config.yaml

### Configuration

Create a `config.yaml` file in `config` folder

```
server:
  port: "8080"

database:
  type: "cockroach"
  url: "root:password@tcp(localhost:3306)/loan_system?parseTime=true"
```

### Running Migrations

To set up the database schema:
```
go run cmd/migrate/main.go
```

This will create all necessary tables including:

- loans
- borrowers
- lenders
- employees
- documents
- loan_lenders

### Running the Server

Start the server in development mode:
```
air
```

Or run directly:
```
go run cmd/api/main.go
```

## API Documentation

After starting the server, access the Swagger documentation at

```
http://localhost:5002/swagger/index.html
```

This provides a complete API reference with interactive testing capabilities.

## Debugging

To debug the app, just run `air` to enable debugging on port `2345`, then connect your IDE debugger to `localhost:2345`

### Project Structure

- `cmd`: Application entry points
    - `/api`: Main service
    - `/migrate`: Database migration
- `internal`: Internal application code
    - `/api`: API handlers and routes
    - `/domain`: Business logic and entities
    - `/repository`: Data access layer
- `pkg`: Shared libraries
- `migrations`: Database migration scripts
- `docs`: API documentation

## Loan State Machine Workflow

The loan processing follows a state machine pattern:

1. Proposed: Initial loan creation
2. Approved: Loan approved by an employee
3. Invested: Funding provided by lenders
4. Disbursed: Funds transferred to borrower

Each state transition is tracked with metadata including timestamps and responsible parties.

### Current Limitations and Future Improvements

- Authentication: No authentication/authorization mechanism is currently implemented
- File Uploads: Document file uploads are simulated with string filenames
- Testing: Needs more comprehensive unit and integration tests
- Error Handling: Could benefit from more standardized error responses
- Validation: Additional validation rules for business logic
- Monitoring: No metrics or logging infrastructure
- Deployment: Containerization and CI/CD pipeline