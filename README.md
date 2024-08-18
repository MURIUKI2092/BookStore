# Bookstore Backend API

This is a Go-based backend API for managing a bookstore. The API supports user management, book inventory management, and order processing. The application uses PostgreSQL as the database and GORM as the ORM.

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Environment Variables](#environment-variables)
  - [Running Migrations](#running-migrations)
  - [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
  - [User Management](#user-management)
  - [Books Management](#books-management)
  - [Order Management](#order-management)
- [Testing](#testing)
- [License](#license)

## Features

- User registration, login, and management
- CRUD operations for books
- Ordering system for customers
- JWT authentication
- Role-based access control (Admin, User)
- Hot reload during development

## Project Structure

```bash
bookstore/
│
├── main.go               # Entry point of the application
├── database/             # Database connection and migration files
│   └── db.go
├── models/               # GORM models
│   └── user.go
│   └── book.go
├── handlers/             # HTTP handlers
│   └── user.go
│   └── book.go
├── routes/               # Router setup
│   └── user_routes.go
│   └── book_routes.go
├── helpers/              # Utility functions (e.g., authentication)
│   └── auth.go
└── config/               # Configuration files and environment variables
    └── config.go

```

## Environment Variables
- DB_HOST=localhost
- DB_PORT=5432
- DB_USER=your_db_user
- DB_PASSWORD=your_db_password
- DB_NAME=your_db_name
- JWT_SECRET=your_jwt_secret_key


## Running  migrations
``` bash
go run main.go migrate
```

## Running  the app
``` bash
go install github.com/githubnemo/CompileDaemon@latest

CompileDaemon -command="./BookStore"



```
##  Production :
     go run main.go


