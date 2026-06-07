# Chirpy

Chirpy is a RESTful social media API built in Go. It allows users to create accounts, authenticate with JWTs, publish short messages ("chirps"), manage their profiles, and interact with premium membership features through webhook integrations.

This project was built to practice backend development concepts including HTTP servers, REST APIs, authentication, database design, migrations, and third-party integrations.

## Features

### Authentication & Security

* JWT-based access tokens
* Refresh token support
* Password hashing using bcrypt
* Protected API endpoints
* API key authentication for webhooks

### User Management

* Create user accounts
* Login with email and password
* Update account information
* Chirpy Red membership status

### Chirps

* Create chirps
* Retrieve all chirps
* Retrieve individual chirps
* Delete your own chirps
* Filter chirps by author
* Sort chirps by creation date

### Premium Memberships

* Polka webhook integration
* Automatic Chirpy Red upgrades
* Secure webhook verification using API keys

### Database

* PostgreSQL database
* SQL migrations managed with Goose
* Type-safe database access generated with sqlc

## Tech Stack

* Go
* PostgreSQL
* sqlc
* Goose
* JWT
* bcrypt
* net/http

## Project Structure

```text
chirpy/
├── internal/
│   ├── auth/
│   └── database/
├── sql/
│   ├── queries/
│   └── schema/
├── main.go
├── sqlc.yaml
├── go.mod
└── README.md
```

## Installation

### Clone the repository

```bash
git clone https://github.com/yourusername/chirpy.git
cd chirpy
```

### Install dependencies

```bash
go mod download
```

### Configure environment variables

Create a `.env` file:

```env
DB_URL=postgres://postgres:password@localhost:5432/chirpy?sslmode=disable
JWT_SECRET=your-secret-key
POLKA_KEY=f271c81ff7084ee5b99a5091b42d486e
PLATFORM=dev
```

### Run database migrations

```bash
goose -dir sql/schema postgres "$DB_URL" up
```

### Generate database code

```bash
sqlc generate
```

### Start the server

```bash
go run main.go
```

The API will be available at:

```text
http://localhost:8080
```

## Example Endpoints

### Create User

```http
POST /api/users
```

### Login

```http
POST /api/login
```

### Create Chirp

```http
POST /api/chirps
Authorization: Bearer <token>
```

### Get All Chirps

```http
GET /api/chirps
```

### Filter Chirps By Author

```http
GET /api/chirps?author_id=<user-id>
```

### Sort Chirps

```http
GET /api/chirps?sort=desc
```

### Delete Chirp

```http
DELETE /api/chirps/{chirpID}
Authorization: Bearer <token>
```

## What I Learned

This project helped me gain hands-on experience with:

* Building REST APIs in Go
* Authentication and authorization
* JWT and refresh token workflows
* PostgreSQL schema design
* Database migrations
* SQL query generation with sqlc
* Webhook integrations
* API security best practices
* Structuring production-style Go applications

## Future Improvements

* Pagination
* User profiles
* Chirp editing
* Rate limiting
* Docker support
* Automated testing
* CI/CD pipeline
* Deployment to cloud infrastructure
