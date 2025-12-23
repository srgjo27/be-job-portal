# Job Portal Backend

Backend service for a Job Portal application, built with **Go**, **Gin**, and **GORM**. This service manages users (Seekers & Recruiters), jobs, applications, and profiles.

## Tech Stack

- **Language:** Go (Golang)
- **Framework:** Gin Web Framework
- **Database:** PostgreSQL (via GORM)
- **Authentication:** JWT & Google OAuth

## Features

- **Authentication**: Register, Login, Google OAuth.
- **User Roles**: Separated logic for **Seekers** (Applicants) and **Recruiters**.
- **Job Management**: Create, Update, List Jobs (Public & Recruiter specific).
- **Application System**: Apply for jobs, upload resume/CV, view applicants (Recruiter), update application status.
- **Profile Management**: Manage Seeker and Recruiter/Company profiles.
- **Dashboard**: Analytics for Recruiters (Total applicants, trends, recent applications).

## Project Structure

```
├── cmd
│   └── api
│       └── main.go       # Entry point
├── internal
│   ├── config            # Configuration loader
│   ├── delivery
│   │   └── http          # HTTP Handlers & Router
│   ├── domain            # Domain models & Interfaces
│   ├── repository        # Database interactions
│   └── usecase           # Business logic
└── pkg
    ├── database          # DB Connection
    └── utils             # Utilities (Auth, Response helper)
```

## Setup & Run

1.  **Clone the repository**
2.  **Configure Environment**
    Create a `.env` file in the root directory (refer to code for required variables, typically `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`, `SECRET_KEY`, `SERVER_PORT`).
3.  **Install Dependencies**
    ```bash
    go mod tidy
    ```
4.  **Run the Application**
    ```bash
    go run cmd/api/main.go
    ```
    The server typically runs on port `8080`.

## API Endpoints

### Auth
- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/auth/google/callback`

### Profile
- `GET /api/profile`
- `PUT /api/profile`

### Jobs
- `POST /api/jobs` (Recruiter)
- `PUT /api/jobs/:id` (Recruiter)
- `GET /api/jobs`
- `GET /api/jobs/:id`
- `GET /api/jobs/:id/applicants` (Recruiter)

### Applications
- `POST /api/applications` (Seeker)
- `GET /api/applications`
- `PUT /api/applications/:id/status` (Recruiter)

### Dashboard
- `GET /api/dashboard/stats` (Recruiter)
