# Database Setup & Troubleshooting

This document explains how the PostgreSQL database is configured for the Insira backend and how to troubleshoot common connection issues.

## Environment Variables

The backend application requires a `.env` file at the root of the `backend` directory.

```env
PORT=8080
ENV=development
DATABASE_URL=host=localhost user=admin password=secret dbname=insira port=5431 sslmode=disable
```

### Important Configuration Details:
- **Port**: The Docker PostgreSQL container is mapped to port `5431` (`0.0.0.0:5431->5432/tcp`), so you MUST use `port=5431` in your `DATABASE_URL`. Do not use `5432` unless you have explicitly reconfigured Docker.
- **User/Password**: The default configured user is `admin` with password `secret`.
- **Database Name**: The application connects to a database named `insira`.

---

## Setting Up the Database via Docker

If you encounter an error stating `FATAL: database "insira" does not exist`, it means the PostgreSQL container is running, but the `insira` database hasn't been created yet.

Here is the step-by-step process to create the database inside the running Docker container:

### 1. Verify Docker Container is Running
Check if your PostgreSQL container is running:
```bash
docker ps
```
You should see a container named `my_postgres` (or similar) mapping port `5431`.

### 2. Create the Database
Execute the following command to log into the PostgreSQL container and create the `insira` database:

```bash
docker exec my_postgres psql -U admin -d postgres -c "CREATE DATABASE insira;"
```

### 3. Verify Database Creation (Optional)
To verify that the database was successfully created, you can list all databases within the container:
```bash
docker exec my_postgres psql -U admin -d postgres -c "\l"
```
You should see `insira` listed in the output.

---

## Running the Application

Once the `.env` file is properly configured and the database is created, you can start the application:

```bash
go run cmd/main.go
```

The application uses GORM's `AutoMigrate` feature, which will automatically create the necessary tables (such as `users`) when it successfully connects to the database.
