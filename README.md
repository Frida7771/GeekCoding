# GeekCoding - Online Judge System

An online judge system built with Go + Gin + GORM + MySQL + Redis + Docker, supporting code submission, automatic compilation, execution, and judging.

## ✨ Features

### 🔐 User System
- ✅ User registration (email verification code)
- ✅ User login (JWT Token authentication)
- ✅ Send verification code (Redis cache, SMTP email)
- ✅ User detail query
- ✅ Leaderboard (sorted by number of solved problems and submissions)

### 📝 Problem Management
- ✅ Problem list (supports keyword search, category filtering, pagination)
- ✅ Problem detail (includes test cases)
- ✅ Create problem (admin only, supports multiple categories and test cases)
- ✅ Update problem (admin only, transaction ensures consistency)

### 🏷️ Category Management
- ✅ Category list (supports search, pagination)
- ✅ Create category (supports hierarchical structure)
- ✅ Update category
- ✅ Delete category (with safety checks)

### 💻 Code Submission & Judging (Core Feature)
- ✅ Code submission (saved to temporary directory)
- ✅ Automatic compilation (Go code)
- ✅ **Docker isolated execution** (memory limits, timeout control, network isolation)
- ✅ **Concurrent testing** (goroutines execute multiple test cases concurrently)
- ✅ **Accurate judgment**:
  - Compilation error (status code 5)
  - Wrong answer (status code 2)
  - Runtime timeout (status code 3)
  - Memory limit exceeded (status code 4)
  - Accepted (status code 1)

### 🔒 Authentication & Authorization
- ✅ JWT Token authentication
- ✅ User permission middleware
- ✅ Admin permission middleware

## 🛠️ Tech Stack

- **Backend Framework**: Gin
- **ORM**: GORM
- **Database**: MySQL 8.0+
- **Cache**: Redis
- **Containerization**: Docker (code execution isolation)
- **Authentication**: JWT
- **API Documentation**: Swagger/OpenAPI
- **Email Service**: SMTP (Gmail)

## 📋 Prerequisites

- Go 1.21+
- MySQL 8.0+
- Redis
- Docker (for running user-submitted code)

## 🚀 Quick Start

### 1. Clone the repository

```bash
git clone <repository-url>
cd GeekCoding
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Configure database

Make sure MySQL and Redis services are running:

```bash
# Check MySQL
mysql -u root -p

# Check Redis
redis-cli ping
```

Database configuration is in `models/init.go`, modify as needed.

### 4. Build Docker image

Build the Docker image for running user code:

```bash
./docker-build.sh
```

Or build manually:

```bash
docker build -t golang-code-runner:latest -f internal/code/Dockerfile .
```

### 5. Generate Swagger documentation (optional)

```bash
swag init
```

### 6. Run the application

```bash
go run main.go
```

The application will start at `http://localhost:8080`.

### 7. Access Swagger documentation

Open your browser and visit:

```
http://localhost:8080/swagger/index.html
```

## 🤝 Contributing

Issues and Pull Requests are welcome!

## 📄 License

Apache 2.0

## 👤 Author

Bo Li (Frida)

---

**Note**: This is an online judge system. Please ensure in production:
- Use strong passwords
- Configure HTTPS
- Limit Docker resource usage
- Regularly clean temporary files
- Monitor system resources
