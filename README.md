# Gin Rest API & Socket.io
#### TODO:
- Add middleware for authentication
### Create .env file 
```env
PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=12345678
DB_NAME=db

REDIS_HOST=localhost
REDIS_PORT=6379

GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
```

### Run with Docker Compose (Optional)
```bash
docker compose up --build
```

### Install Dependencies
```bash
go mod tiny
```

### Run the Application
```bash
go run main.go
```

