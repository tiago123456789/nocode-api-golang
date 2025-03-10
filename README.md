## 🚀 ABOUT

This project is designed to help you create an API as quickly as possible—without writing a backend!  
Simply define your endpoints based on your PostgreSQL tables, and you're good to go.

## 🚀 **Features**  

✅ **Flexible CRUD Operations**  
- Effortlessly create, read, update, and delete records in a PostgreSQL table.  
- 🔐 Define route access — choose whether endpoints are public or restricted with a simple setting.  

✅ **Custom Query Endpoints**  
- Set up custom query endpoints with ease.  
- 🔐 Control access — mark routes as public or private for secure data handling.  

✅ **Dynamic Query Filtering**  
- Build endpoints that support dynamic query strings for powerful filtering.  
- 🌐 Perfect for pagination and complex multi-table queries — get the data you need, fast.  

✅ **Smart Input Validation**  
- Ensure data integrity with robust validation for record creation.  
- 🚦 Catch errors early and keep your data clean.  

✅ **Optimized Caching**  
- Boost performance with caching for `GET` endpoints.  
- ⏳ Customize the TTL (Time-to-Live) to control how long data stays fresh in Redis.  


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

### Architecture

![Architecture](./archeticture.png "Architecture")

#### About the architecture

- The project can scale horizontally without problem.
- The good place to run the project is the GCP cloud run, because will scale the api for automatically, serverless container solution where you pay only for executions and all the logs from all instances from api will be centralized at one place.

### 🔑 Environment Variables Explained

Below are the essential environment variables needed to run the project smoothly:

| Variable                  | Description                                                          | Example Value                                    |
| ------------------------- | -------------------------------------------------------------------- | ------------------------------------------------ |
| **DB_URL**                | Database connection URL                                              | `postgres://user:password@localhost:5432/dbname` |
| **API_KEY**               | API key for authentication                                           | `your_api_key`                                   |
| **USER_NAME**             | Default admin username                                               | `admin`                                          |
| **USER_EMAIL**            | Default admin email                                                  | `admin@example.com`                              |
| **USER_PASSWORD**         | Default admin password                                               | `securepassword`                                 |
| **JWT_SECRET**            | Secret key for JWT authentication                                    | `your_jwt_secret`                                |
| **REDIS_URL**             | Redis cache URL                                                      | `cache:6379`                                     |
| **REDIS_PASSWORD**        | Redis password (keep empty)                                          | `""`                                             |
| **REDIS_LIMIT_MEMORY**    | Limit memory can be used from Redis after that will apply the LRU    | `"10mb"`                                         |
| **NEW_RELIC_LICENSE_KEY** | New relic license key to monitory the endpoints                      | `""`                                             |
| **NEW_RELIC_APP_NAME**    | New relic app name                                                   | `""`                                             |
| **NEW_RELIC_ENABLED**     | if you want to enable the value is "yes", case no, the value is "no" | `"yes or no"`                                    |


### 🔑 Insominia endpoints

[Endpoints from insominia](./nocode-api-Insomnia.json)

💡 **Note:** Keep your `.env` file secure and never expose sensitive credentials in a public repository! 🔒
