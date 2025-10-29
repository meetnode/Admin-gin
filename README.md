# Admin Panel (Gin + PostgreSQL)

A simple and modular **Admin Panel** built using the **Gin Web Framework** and **PostgreSQL**.
This project provides a base structure for building secure, production-ready admin backends with RESTful APIs and integrated **Swagger documentation**.

---

## 🚀 Features

* Built with **Gin Framework** (high performance Go web framework)
* Database: **PostgreSQL**
* Follows **clean and modular project structure**
* Integrated **Swagger documentation** (available at `/docs`)
* Example CRUD endpoints for quick start
* Easy to extend for any admin or backend system

---

## 📂 Project Structure

```
│   .air.toml
│   .env
│   .gitignore
│   docker-compose.yml
│   go.mod
│   go.sum
│   main.exe
│   Makefile
│   README.md
│
├───cmd
│   ├───api
│   │       main.go              # Entry point for API server
│   └───seed
│           main.go              # Seeder for database initialization
│
├───docs
│       docs.go
│       swagger.json
│       swagger.yaml
│
├───internal
│   ├───controllers              # Handles request/response logic
│   ├───database                 # DB connection and tests
│   ├───middlewares              # JWT, role, and permission checks
│   ├───models                   # GORM models for all entities
│   ├───server                   # Server setup and route definitions
│   ├───services                 # Business logic layer
│   └───utils                    # Utility functions (crypto, jwt, mail, etc.)
│
└───tmp
        build-errors.log
```

---

## ⚙️ Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/meetnode/Admin-gin.git
cd Admin-gin
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Setup Environment Variables

Rename or copy from `.env.example` file and configure your environment variables:

```env
URL=http://localhost:5000
PORT=5000
APP_ENV=local
BLUEPRINT_DB_HOST=localhost
BLUEPRINT_DB_PORT=5432
BLUEPRINT_DB_DATABASE=admin-gin
BLUEPRINT_DB_USERNAME=postgres
BLUEPRINT_DB_PASSWORD=abcd
BLUEPRINT_DB_SCHEMA=public
SMTP_USER=xyz@gmail.com
SMTP_PASS=zxcvbnmasgjwert
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587

APP_SECRET="ashdjkas45dshukf"
```
### 4. Run migrations or seed data

You can use the provided `cmd/seed/main.go` file to seed default users, roles, or permissions:

```bash
go run cmd/seed/main.go
```

---

## 🏃 Run the Server

```bash
go run cmd/main.go
```
or

```bash
air
```

Server will start on:

```
http://localhost:8080
```

---

## 📖 API Documentation (Swagger)

Swagger UI is available at:

```
http://localhost:5000/docs
```

Use it to explore all endpoints and test APIs interactively.

---

## 🧱 Tech Stack

| Layer         | Technology                              |
| ------------- | --------------------------------------- |
| Framework     | [Gin](https://github.com/gin-gonic/gin) |
| Database      | PostgreSQL                              |
| ORM (if used) | GORM                                    |
| Documentation | Swagger (swaggo/gin-swagger)            |
| Language      | Go (Golang)                             |

---

## 🧩 Common Commands

A preconfigured .air.toml file is already included, so just run:
```bash
air
```

**Format code**

```bash
go fmt ./...
```

**Run Tests**

```bash
go test ./...
```

---

## 🤝 Contributing

Contributions are welcome!
If you want to improve the project or fix a bug:

1. Fork the repo
2. Create a new branch
3. Commit your changes
4. Submit a pull request

---

## 🧠 Notes

* This is a **public template**, so feel free to clone and modify it.
* Make sure to update the database credentials in `.env` before running locally.
* Swagger files are auto-generated and located in `/docs`.

---

## 📜 License

This project is open-source and available under the **MIT License**.
