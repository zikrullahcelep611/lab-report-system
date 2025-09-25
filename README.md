# Lab Report Management System

A comprehensive medical laboratory report management system built with Go, designed to streamline healthcare workflows for medical institutions.

## 🏥 Features

- **Patient Management**: Complete CRUD operations for patient records with national ID tracking
- **Report Management**: Create, update, delete and search medical reports
- **Multi-Hospital Support**: Handle multiple healthcare institutions with proper isolation
- **User Management**: Role-based user system (Super Admin, Admin, Technician)
- **Advanced Search**: Search reports by patient name or national ID
- **Secure Authentication**: JWT-based authentication with role-based access control

## 🛠 Technology Stack

- **Backend**: Go 1.24+ with Fiber v2 framework
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT tokens with bcrypt password hashing
- **Containerization**: Docker & Docker Compose
- **Logging**: Structured logging with Zerolog
- **Architecture**: Clean architecture with Repository-Service-Handler pattern

## 🚀 Quick Start

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- PostgreSQL (or use Docker)

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/zikrullahcelep611/lab-report-system.git
cd lab-report-system
```

2. **Start services with Docker**
```bash
docker-compose up -d
```

3. **Run the application**
```bash
cd lab-report/backend
go run main.go
```

The server will start on `http://localhost:8080`

## 📊 Database Schema

### Core Entities

- **Users**: Healthcare professionals (doctors, technicians, admins)
- **Hospitals**: Medical institutions
- **Patients**: Patient records with unique national IDs
- **Reports**: Medical test reports linked to patients and users


## 🔐 Authentication & Authorization

### User Roles

- **Super Admin**: Full system access, can create users
- **Admin**: Hospital-level management
- **Technician**: Basic report operations

### API Security

- JWT token-based authentication
- Role-based access control (RBAC)
- Token blacklisting for secure logout
- Password hashing with bcrypt

## 📚 API Endpoints

### Authentication
```
POST /api/login     # User login
POST /api/logout    # User logout
```

### Users (Protected)
```
GET    /api/users/:id    # Get user by ID
POST   /api/users        # Create user (Super Admin only)
PUT    /api/users/:id    # Update user
DELETE /api/users/:id    # Delete user
```

### Reports (Protected)
```
GET    /api/reports                    # Get all reports
GET    /api/reports/:id               # Get report by ID
POST   /api/reports                   # Create report
PUT    /api/reports/:id               # Update report
DELETE /api/reports/:id               # Delete report
GET    /api/reports/search            # Search by patient name
```

### Patients
```
GET    /api/patients/:id    # Get patient by ID
POST   /api/patients        # Create patient
PUT    /api/patients/:id    # Update patient
DELETE /api/patients/:id    # Delete patient
```

### Hospitals
```
GET    /api/hospitals/:id    # Get hospital by ID
POST   /api/hospitals        # Create hospital
PUT    /api/hospitals/:id    # Update hospital
DELETE /api/hospitals/:id    # Delete hospital
```

## 🔧 Configuration

Create `backend/resources/application.yml`:

```yaml
qa:
  server:
    port: 8080
  database:
    dns: "host=localhost user=myuser password=mypassword dbname=labDb port=5432 sslmode=disable"
  log:
    level: 0  # Debug level
  jwt:
    secretKey: "your-super-secret-jwt-key"
```

## 📁 Project Structure

```
lab-report/
├── backend/
│   ├── main.go                    # Application entry point
│   ├── api/                       # HTTP handlers
│   │   ├── auth/                 # Authentication endpoints
│   │   ├── user/                 # User management
│   │   ├── report/               # Report management
│   │   └── ...
│   ├── application/               # Business logic
│   │   ├── authService/
│   │   ├── userService/
│   │   └── ...
│   ├── infrastructure/            # External concerns
│   │   ├── config/               # Configuration
│   │   ├── postgresDb/           # Database connection
│   │   └── repository/           # Data access layer
│   ├── middleware/                # HTTP middleware
│   │   ├── authMiddleware/       # Authentication
│   │   └── rbacMiddleware/       # Authorization
│   ├── models/                    # Data models
│   └── resources/                 # Configuration files
└── docker-compose.yml            # Docker services
```



## 🐳 Docker Deployment

The project includes Docker configuration for easy deployment:

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```
