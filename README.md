# Kithli API

A Go-based backend API service for the Kithli platform, handling user authentication, data storage, and business logic for pairing care takers (Kiths) with members undergoing medical procedures.

## Project Description

Backend service that manages the core functionality of the Kithli platform, including:
- User authentication and authorization
- Database operations for storing user and appointment data
- File storage and management
- Real-time notifications and updates

## Project Requirements

### Development Requirements
- Go (Latest stable version)
- PostgreSQL Database
- Modern web browser (Chrome, Firefox, Safari, Edge)

### External Services
- Firebase Authentication
- Google Cloud Storage
- PostgreSQL Database

## Project Setup / Build Instructions

1. Clone the repository:
```bash
git clone [https://github.com/B12Howard/kithli-api]
cd kithli-api
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file in the root directory with the following variables:
```
DB_HOST=your_db_host
DB_PORT=your_db_port
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
```

4. Create a `config.json` file with the following structure:
```json
{
  "GCPCLOUDSTORAGE": {
    "your_gcp_credentials_here"
  },
  "FIREBASE": {
    "your_firebase_credentials_here"
  }
}
```

## How-to Run the Project

### Development Mode
1. Install Air for hot reloading:
```bash
go install github.com/cosmtrek/air@latest
```

2. Initialize Air:
```bash
air init
```

3. Start the development server:
```bash
air -c .air.toml
```

### Testing
```bash
go test ./...
```

### Production Build
```bash
go build -o kithli-api
```

## Development Process

### Code Quality Tools
- Go linter
- Go vet
- Go test for testing

### Available Scripts
- `go run main.go` - Start server
- `go build` - Create production build
- `go test ./...` - Run tests
- `go vet ./...` - Run Go vet

### Database
For database migrations this project uses [https://github.com/pressly/goose]

## Deployment Instructions

### Environment Setup
1. Set up Firebase project and enable Authentication
2. Configure Google Cloud Storage
3. Set up PostgreSQL database
4. Install PostGIS PostgreSQL extension
5. Configure environment variables

### Production Considerations
- Ensure all environment variables are properly set
- Configure proper security rules in Firebase
- Set up proper CORS policies
- Configure proper database indexes
- Set up monitoring and logging
- Implement proper error handling and logging
- Set up backup and recovery procedures

#### Client
The client application is part of the Kithli project repository.

#### TODO Docker

#### Postman Collection


![Kithli Architecture](https://user-images.githubusercontent.com/39282569/196551643-9d64515f-128e-4c8c-af39-071ce5d43226.png)

 ### Running Tests
 To run unit tests use `go test YOURDIRECTORY`. Example: Tests are in the service directory. If I wanted to run all of them I would run `go test ./services/... -v`