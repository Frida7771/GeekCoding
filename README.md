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

## 📚 API Documentation

### Public Endpoints (No authentication required)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/problem-list` | Get problem list |
| GET | `/problem-detail` | Get problem detail |
| GET | `/user-detail` | Get user detail |
| POST | `/login` | User login |
| POST | `/send-code` | Send verification code |
| POST | `/register` | User registration |
| GET | `/rank-list` | Leaderboard |
| GET | `/submit-list` | Get submission list |

### User Private Endpoints (User authentication required)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/user/submit` | Submit code |

### Admin Endpoints (Admin authentication required)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/admin/problem-create` | Create problem |
| PUT | `/admin/problem-update` | Update problem |
| GET | `/admin/category-list` | Get category list |
| POST | `/admin/category-create` | Create category |
| PUT | `/admin/category-update` | Update category |
| DELETE | `/admin/category-delete` | Delete category |

For detailed API documentation, visit Swagger UI: `http://localhost:8080/swagger/index.html`

## 🔑 Authentication

### User Authentication

1. Login via `/login` endpoint to get Token
2. Add to request header: `Authorization: <token>`
3. Access user private endpoints (e.g., `/user/submit`)

### Admin Authentication

1. Login with admin account to get Token
2. Add to request header: `Authorization: <token>`
3. Access admin endpoints (e.g., `/admin/problem-create`)

## 📊 Submission Status Codes

| Status Code | Meaning | Description |
|-------------|---------|-------------|
| 1 | Accepted | All test cases passed |
| 2 | Wrong Answer | Output doesn't match expected output |
| 3 | Time Limit Exceeded | Exceeded maximum runtime |
| 4 | Memory Limit Exceeded | Exceeded maximum memory limit |
| 5 | Compilation Error | Code compilation failed |

## 🏗️ Project Structure

```
GeekCoding/
├── main.go                 # Application entry point
├── router/
│   └── app.go             # Route configuration
├── service/               # Business logic layer
│   ├── user.go           # User-related services
│   ├── problem.go        # Problem-related services
│   ├── submit.go         # Submission & judging services
│   └── category.go       # Category management services
├── models/                # Data model layer
│   ├── user_basic.go
│   ├── problem_basic.go
│   ├── test_case.go
│   ├── submit_basic.go
│   ├── category_basic.go
│   ├── problem_category.go
│   └── init.go           # Database initialization
├── middlewares/           # Middleware
│   ├── auth_user.go      # User authentication
│   └── auth_admin.go     # Admin authentication
├── help/                  # Helper functions
│   └── helper.go         # JWT, encryption, utility functions
├── internal/code/         # Docker related
│   ├── Dockerfile        # Code runner image
│   └── docker-runner.sh  # Container execution script
├── docs/                  # Swagger documentation
├── define/               # Constant definitions
└── docker-build.sh       # Docker image build script
```

## 🗄️ Database Design

### Core Tables

1. **user_basic** - User table
   - User basic information, statistics (solved count, submission count)

2. **problem_basic** - Problem table
   - Problem information, memory limit, time limit

3. **test_case** - Test case table
   - Input data, expected output

4. **submit_basic** - Submission record table
   - Submission information, execution status

5. **category_basic** - Category table
   - Category information, hierarchical structure

6. **problem_category** - Problem-category association table
   - Many-to-many relationship

Database tables are automatically migrated when the application starts.

## 🔧 Configuration

### Database Configuration

Configure in `models/init.go`:

```go
dsn := "root:password@tcp(127.0.0.1:3306)/Geek_Coding?charset=utf8mb4&parseTime=True&loc=Local"
```

### Redis Configuration

Configure in `models/init.go`:

```go
Addr:     "localhost:6379",
Password: "",
DB:       0,
```

### Email Configuration

Configure in `help/helper.go`, requires environment variable:

```bash
export GMAIL_APP_PASSWORD="your-gmail-app-password"
```

## 🐳 Docker Code Execution

The system uses Docker containers to execute user-submitted code, ensuring:

- **Security isolation**: Code runs in isolated containers
- **Resource limits**: Precise memory control (via Docker `--memory`)
- **Timeout control**: Execution time controlled via `--stop-timeout`
- **Network isolation**: `--network=none` disables network access

### Build Docker Image

```bash
./docker-build.sh
```

## 📝 Usage Examples

### 1. User Registration

```bash
# 1. Send verification code
curl -X POST http://localhost:8080/send-code \
  -d "email=user@example.com"

# 2. Register
curl -X POST http://localhost:8080/register \
  -d "email=user@example.com" \
  -d "code=123456" \
  -d "name=John Doe" \
  -d "password=password123"
```

### 2. User Login

```bash
curl -X POST http://localhost:8080/login \
  -d "username=John Doe" \
  -d "password=password123"
```

### 3. Submit Code

```bash
curl -X POST "http://localhost:8080/user/submit?problem_identity=xxx" \
  -H "Authorization: <your-token>" \
  -H "Content-Type: application/json" \
  -d 'package main
import "fmt"
func main() {
    var a, b int
    fmt.Scanln(&a, &b)
    fmt.Println(a + b)
}'
```

### 4. Create Problem (Admin)

```bash
curl -X POST http://localhost:8080/admin/problem-create \
  -H "Authorization: <admin-token>" \
  -F "title=Two Sum" \
  -F "content=Calculate the sum of two integers" \
  -F "max_runtime=3000" \
  -F "max_mem=64" \
  -F "category_ids=1" \
  -F "test_cases={\"input\":\"1 2\",\"output\":\"3\"}" \
  -F "test_cases={\"input\":\"5 7\",\"output\":\"12\"}"
```

## 🐛 Troubleshooting

### 1. MySQL Connection Failed

**Error**: `gorm Init Error`

**Solution**:
- Check if MySQL is running
- Check if database `Geek_Coding` exists
- Verify connection information is correct

### 2. Redis Connection Failed

**Error**: Redis connection timeout

**Solution**:
- Check if Redis is running: `redis-cli ping`
- Start Redis: `redis-server`

### 3. Docker Image Not Found

**Error**: `docker: Error response from daemon`

**Solution**:
- Build image: `./docker-build.sh`
- Check image: `docker images | grep golang-code-runner`

### 4. Port Already in Use

**Error**: `bind: address already in use`

**Solution**:
- Find process using port: `lsof -i :8080`
- Change port: modify `r.Run(":8080")` in `main.go`

## 📖 Documentation

- [Running Guide](./RUN.md) - Detailed running steps and troubleshooting
- [Project Overview](./PROJECT_OVERVIEW.md) - Complete feature description and technical details

## 🤝 Contributing

Issues and Pull Requests are welcome!

## 📄 License

Apache 2.0

## 👤 Author

GeekCoding Team

---

**Note**: This is an online judge system. Please ensure in production:
- Use strong passwords
- Configure HTTPS
- Limit Docker resource usage
- Regularly clean temporary files
- Monitor system resources
