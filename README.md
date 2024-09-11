# Go JWT API

## Project Structure

The project is organized as follows:

- `cmd/go-jwt-api`: Contains the entry point of the application.
- `internal`: Houses the core logic of the application, including database interactions, handlers, middlewares, and models.
- `tests`: Contains integration tests for the application.
- `Dockerfile` and `docker-compose.yml`: For containerization and orchestration.
- `go.mod` and `go.sum`: Go module files for managing dependencies.

## Technologies Used

- **Go**: The primary programming language for building the application.
- **Docker and Docker Compose**: Used for containerizing the application and its services.
- **PostgreSQL**: The database system for storing application data.
- **Gorilla Mux**: A powerful URL router and dispatcher for matching incoming requests to their respective handler.

## Setup and Running

To run this project, you need to have Docker and Docker Compose installed on your system. Follow these steps:

1. Clone the repository to your local machine.
```bash
git clone https://github.com/say8hi/go-jwt-api.git
```
2. Go into **go-jwt-api** folder:
```bash
cd go-jwt-api
```
3. Rename `.env.example` to `.env` and adjust the configuration according to your environment.
```bash
mv .env.example .env
```
4. Build and start the services with Docker Compose:
```bash
docker-compose up -d --build
```
This command will start all the required services, including the Go application, PostgreSQL, and RabbitMQ.

## API Endpoints

- **Users**
  - `POST /users/create`: Create a new user.
  - `GET /users/get_tokens?user_id={user_id}`: Create AUTH tokens.
  - `POST /users/refresh`: Refresh AUTH tokens.

Use the provided `curl` examples in the [Application Usage Examples](#application-usage-examples) section to interact with these endpoints.

## Application Usage Examples

Below are examples of how to interact with the application using `curl`, a command-line tool for transferring data with URLs. These examples demonstrate user creation, category creation, and product creation through the application's API.

### Creating a User

To create a new user, send a POST request with the username and password in JSON format:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"username": "username", "email": "email@email.com"}' http://0.0.0.0:8080/users/create
```

## Testing

To run the integration tests:
```bash
./run_tests.sh
```
This script will set up the test environment, run the tests, and tear down the environment afterwards.
