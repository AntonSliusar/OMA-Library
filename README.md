# OMA-Library

**OMA-Library** is a small Go service for storing and serving `.oma` files. It provides an admin web UI for uploading and searching files, and an MCP endpoint for managing workers. The project uses PostgreSQL for metadata and Cloudflare R2 for file storage.

## 🔧 Requirements

- Go 1.20+ (or your preferred Go version)
- PostgreSQL database
- [Cloudflare R2](https://developers.cloudflare.com/r2/) credentials (endpoint, access key, secret key, bucket)
- Optional: Docker & Docker Compose for containerized setup

## ⚙️ Configuration

Configuration is provided via environment variables, which are parsed using `envconfig`. Variables are grouped by component prefixes.

### Database (`db` prefix)

| Variable       | Description                     | Example                  |
|----------------|---------------------------------|--------------------------|
| `DB_HOST`      | PostgreSQL host                 | `omapostgres`           |
| `DB_PORT`      | PostgreSQL port                 | `5432`                  |
| `DB_USER`      | Database user                   | `postgres`              |
| `DB_PASSWORD`  | Database password               | `postgres`              |
| `DB_DBNAME`    | Database name                   | `omadb`                 |
| `DB_SSLMODE`   | SSL mode for connection         | `disable`               |

### R2 Storage (`r2` prefix)

| Variable       | Description                     | Example                                   |
|----------------|---------------------------------|-------------------------------------------|
| `R2_ENDPOINT`  | R2 API endpoint                 | `https://...r2.cloudflarestorage.com`     |
| `R2_ACCESS_KEY`| R2 access key                   | `e690...`                                 |
| `R2_SECRET_KEY`| R2 secret key                   | `33d6...`                                 |
| `R2_BUCKET`    | Bucket name                     | `omabucket`                               |

### MCP Server (`mcp` prefix)

| Variable       | Description                     | Example             |
|----------------|---------------------------------|---------------------|
| `MCP_URL`      | MCP service URL                 | `http://localhost`  |
| `MCP_PORT`     | MCP listen port (e.g. `:8081`)  | `:8081`             |

## 🚀 Running the Service

### Using Docker Compose

1. Copy `.env.example` to `.env` and fill in values (or use the provided `.env`).
2. Start services:
   ```bash
   docker-compose up --build
   ```
3. After migrations complete, the web UI is available at `http://localhost:8080`.

### Running Locally

1. Install dependencies:
   ```bash
   go mod download
   ```
2. Ensure your PostgreSQL and R2 credentials are set in environment variables.
3. Run migrations with the `migrate` tool (see `migrations/` folder).
4. Build and start:
   ```bash
   go run ./main.go
   ```

## 🗂️ Database Migrations

SQL files live in the `migrations/` directory. The project expects to use [`migrate`](https://github.com/golang-migrate/migrate) to apply changes. Example:

```bash
migrate -path migrations -database "$DB_URL" up
```

## 📁 Upload/Download Flow

- Admins upload `.oma` files (and optional image) via the web UI.
- Metadata is stored in PostgreSQL; files are pushed to R2.
- Visitors can search by brand/model and download files via generated presigned URLs.

## 🧪 Development Notes

- Templates stored in `templates/*.html`.
- Admin JWT authentication is managed by `echo-jwt` with a cookie.
- Storage logic lives under `pkg/storage`; handlers under `internal/handlers`.

## 📝 License

This project is provided under the MIT License. See the `LICENSE` file for details.

---

*Feel free to expand this README with additional instructions, architectural diagrams or troubleshooting tips as the project grows.*
