# Log Puller Service

This is a log puller service that provides a RESTful API for retrieving log files and specific log entries from a server. The service is built using Go and includes Swagger documentation for the API endpoints.

## Features

- Health check endpoint to verify the server status.
- Retrieve a list of log files from a specified directory.
- Read the last N lines from a specific log file with optional filtering.
- Swagger documentation for API endpoints.
- Middleware to check for an authorization key.

## Installation

To install the necessary dependencies and set up the project, run the following commands:

```sh
go mod tidy
```

## Running the Server

To start the server, use the following command:

```sh
make server
```

The server will run on the port specified (default is `:3000`).

## Building the Project

To build the project, including generating the Swagger documentation, use:

```sh
make build
```

This will generate the Swagger documentation and build the executable in the `bin/` directory.

## Running Tests

### Unit Tests

To run unit tests, use the following command:

```sh
make test
```

### Integration Tests

To run integration tests, use the following command:

```sh
make integ-test
```

### Benchmark Tests

To run benchmark tests, use the following command:

```sh
make benchmark-test
```

## API Endpoints

### Health Check

- **URL**: `/health`
- **Method**: `GET`
- **Description**: Returns the status of the server.
- **Response**:
  ```json
  {
    "status": "OK"
  }
  ```

### Get Log Files

- **URL**: `/v1/logs`
- **Method**: `GET`
- **Description**: Returns a list of log files in a directory.
- **Response**:
  ```json
  {
    "files": [
      "/var/log/syslog",
      "/var/log/messages",
      "/var/log/nginx/access.log",
      "/var/log/nginx/error.log"
    ]
  }
  ```

### Read Log File

- **URL**: `/v1/log`
- **Method**: `GET`
- **Description**: Reads the last N lines from a log file and returns them as a JSON array.
- **Query Parameters**:
  - `file`: The log file to read.
  - `n`: The number of lines to read.
  - `filter` (optional): A filter string to match lines.
- **Response**:
  ```json
  ["line 1", "line 2", "line 3"]
  ```

## Middleware

### Authorization Middleware

The service includes a middleware to check for an `Authorization` header. The header must be in the format `Bearer your-secret-token`. If the token is invalid or missing, the request will be rejected with a `401 Unauthorized` status.

## Swagger Documentation

The API documentation is generated using Swagger. To generate the documentation, run:

```sh
make build
```

You can access the Swagger documentation at `http://localhost:3000/swagger/index.html`.

### Using the Authorization Header in Swagger UI

To execute the endpoints that require authorization using Swagger UI, follow these steps:

1. Open the Swagger UI at `http://localhost:3000/swagger/index.html`.
2. Click on the `Authorize` button at the top right of the page.
3. Enter your token in the `Value` field in the format `Bearer your-secret-token`.
4. Click the `Authorize` button and then the `Close` button.

## GitHub Actions

This project includes a GitHub Actions workflow to build the application .The workflow is triggered on a push to main or can be manually triggered. It builds the Go application for go 1.22

## Makefile Commands

The project includes a `Makefile` with the following commands:

- **Start the server**:

  ```sh
  make server
  ```

- **Build the project** (includes generating Swagger documentation):

  ```sh
  make build
  ```

- **Run unit tests**:

  ```sh
  make test
  ```

- **Run integration tests**:

  ```sh
  make integ-test
  ```

- **Run benchmark tests**:
  ```sh
  make benchmark-test
  ```

## Contributing

Feel free to fork this repository and submit pull requests. For major changes, please open an issue first to discuss what you would like to change.
