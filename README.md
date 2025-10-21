Stage One API

A simple RESTful API built with Go (Golang), Chi router, and PostgreSQL.
This project was developed as part of the HNG Stage One Task to demonstrate backend proficiency using Go, Docker, and PostgreSQL.

🚀 Tech Stack

Language: Go 1.24

Framework: Chi Router

ORM: GORM

Database: PostgreSQL (Aiven Cloud or Local Docker)

Environment Management: godotenv

Containerization: Docker & Docker Compose

🧰 Project Structure
.
├── main.go
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
├── .env
└── README.md

⚙️ Environment Variables

Create a .env file in the root directory with the following variables:

DATABASE_URL=postgres://avnadmin:AVNS_Yj4CE3tSyIJv@hngstage1-bolarinolabisi36-b19d.i.aivencloud.com:20516/defaultdb?sslmode=require
PORT=8080


⚠️ Important: Never commit .env files to version control — make sure it’s included in your .gitignore.

🐳 Running with Docker
Build and start the containers:
docker-compose up --build

Stop the containers:
docker-compose down


Your Go API should now be running on:

👉 http://localhost:8080

💻 Running Locally (Without Docker)
1️⃣ Install dependencies
go mod tidy

2️⃣ Set environment variables

Create a .env file as shown above.

3️⃣ Run the app
go run main.go

🧪 API Testing

You can test your API using:

cURL

Postman

Thunder Client (VS Code extension)

Example:

curl http://localhost:8080/


Expected Response:

{
  "message": "Welcome to Stage One API"
}

🐘 Using Aiven Postgres

This project connects to a remote PostgreSQL instance provided by Aiven.
To connect from Go, ensure your DATABASE_URL is set properly in .env.

Example:

dsn := os.Getenv("DATABASE_URL")
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

🧱 Using Local Postgres (Optional)

If you prefer running Postgres locally via Docker Compose, update your .env:

DATABASE_URL=postgres://admin:adminpassword@db:5432/stage_one?sslmode=disable

📦 Dependencies

Key dependencies used in this project include:

github.com/go-chi/chi/v5
 — lightweight HTTP router

gorm.io/gorm
 — ORM for Go

gorm.io/driver/postgres
 — PostgreSQL driver for GORM

github.com/joho/godotenv
 — for loading .env files

🧑‍💻 Author

Bolarin Olabisi
Frontend & Backend Developer
GitHub Profile

📝 License

This project is licensed under the MIT License.
