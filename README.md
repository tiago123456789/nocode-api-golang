## 🚀 ABOUT

This project is designed to help you create an API as quickly as possible—without writing a backend!  
Simply define your endpoints based on your PostgreSQL tables, and you're good to go.

## 🛠 TECHNOLOGIES

Powered by a modern tech stack for speed and scalability:

- **Go** (v1.24)
- **Fiber** (Lightning-fast web framework)
- **PostgreSQL** (Database)
- **Redis** (Cache)
- **Github actions** (To run unit tests)
- **Docker** & **Docker Compose** (Effortless containerized deployment)

## ⚡ HOW TO RUN LOCALLY

1. **Clone the repository**
2. **Create a `.env` file** based on `.env.example`
3. **Run the following command:**
   ```sh
   docker-compose up -d
   ```

## 📌 ADDITIONAL NOTES

### 🔑 Environment Variables Explained

Below are the essential environment variables needed to run the project smoothly:

| Variable           | Description                       | Example Value                                    |
| ------------------ | --------------------------------- | ------------------------------------------------ |
| **DB_URL**         | Database connection URL           | `postgres://user:password@localhost:5432/dbname` |
| **API_KEY**        | API key for authentication        | `your_api_key`                                   |
| **USER_NAME**      | Default admin username            | `admin`                                          |
| **USER_EMAIL**     | Default admin email               | `admin@example.com`                              |
| **USER_PASSWORD**  | Default admin password            | `securepassword`                                 |
| **JWT_SECRET**     | Secret key for JWT authentication | `your_jwt_secret`                                |
| **REDIS_URL**      | Redis cache URL                   | `cache:6379`                                     |
| **REDIS_PASSWORD** | Redis password (keep empty locally)       | `""`                                             |

### 🔑 Insominia endpoints

[Endpoints from insominia](./nocode-api-Insomnia.json)

💡 **Note:** Keep your `.env` file secure and never expose sensitive credentials in a public repository! 🔒
