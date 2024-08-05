# Storage Application

This is a storage application that utilizes Docker and Docker Compose for containerization. The application comprises three main services: a storage service, a PostgreSQL database, and a Redis instance.

## Table of contents

- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Makefile Commands](#makefile-commands)
- [Configuration](#configuration)

## Prerequisites

Before you begin, ensure you have met the following requirements:

- [Docker](https://www.docker.com/get-started) (version 27.1 or higher)
- [Docker Compose](https://docs.docker.com/compose/install/) (version 2.29 or higher)
- [Golang](https://golang.org/dl/) (version 1.22 or higher)
- [Migrate](https://github.com/golang-migrate/migrate) (CLI tool for database migrations)

## Getting Started

Follow these steps to get the application up and running:

1. **Clone the repository**:  
   Clone the repository to your local machine.

   ```bash
   git clone https://github.com/Saur4ig/file-storage.git
   cd storage-app
   ```

2. **Initialize the services(For the first TIME!)**:  
   Run the init command to build the Docker image, start the database, apply migrations, and start all services.
   ```bash
     make init
   ```

## Makefile commands

Makefile provides several commands to manage the application lifecycle:

- **build**: Builds Docker image for the storage service.
   ```bash
  make build
   ```
- **up**: Starts storage service in detached mode using Docker Compose.
   ```bash
  make up
   ```

- **down**: Stops and removes the running containers.
   ```bash
  make down
   ```

- **init**: Builds Docker image, starts the database, waits for readiness, applies migrations, and starts all services.
   ```bash
  make init
   ```

- **migrate**: Runs the database migrations.
   ```bash
  make migrate
   ```

- **run**: Builds and starts all services.
   ```bash
  make run
   ```

- **re**: Stops and then starts all services.
   ```bash
   make re
   ```

- **test**: Runs tests using Go's testing framework.
   ```bash
  make test
   ```

- **lint**: Runs golangci-lint to lint the codebase.
   ```bash
  make lint
   ```

## Configuration

### Docker compose services

- **storage**:
   - Runs the main application service.
   - Connects to the PostgreSQL database and Redis for data storage and caching.
   - Listens on port 8080.

- **db**:
   - Runs a PostgreSQL instance.
   - Exposes port 5432.
   - Uses volume `./pgdata` for persistent data storage.

- **redis**:
   - Runs a Redis instance.
   - Exposes port 6379.
   - Uses volume `./redisdata` for persistent data storage.

## Performance Benchmarking

The performance of the PostgreSQL database is measured using `pgbench` with the following configuration:

* **pgbench Version**: 16.3 (Debian 16.3-1.pgdg120+1)
* **Transaction Type**: <builtin: TPC-B (sort of)>
* **Scaling Factor**: 10
* **Query Mode**: Simple
* **Number of Clients**: 10
* **Number of Threads**: 2
* **Maximum Number of Tries**: 1
* **Duration**: 60 seconds

## Results

* **Number of Transactions Actually Processed**: 66,388
* **Number of Failed Transactions**: 0 (0.000%)
* **Latency Average**: 9.037 ms
* **Initial Connection Time**: 16.712 ms
* **Transactions Per Second (TPS)**: 1,106.619 (without initial connection time)