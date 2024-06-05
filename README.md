# Log Puller Service

This is a log puller service that provides a RESTful API for retrieving log files and specific log entries from a server. The service is built using Go and includes Swagger documentation for the API endpoints.

## Features

- Health check endpoint to verify the server status.
- Retrieve a list of log files from a specified directory.
- Read the last N lines from a specific log file with optional filtering.
- Swagger documentation for API endpoints.
- Middleware to check for an authorization key.

## TL;DR:

- Build docker image, create 3 containers:

```
make docker
docker-compose up -d
```

- Test: List logs

```
curl --location 'http://localhost:3001/v1/logs' \
--header 'Authorization: Bearer test1'
```

- Test: Read a log in the server

```
curl --location 'http://localhost:3001/v1/log?file=dpkg.log&n=20' \
--header 'Authorization: Bearer test1'
```

- Test: Read logs from remote server

```
curl --location 'http://localhost:3001/v1/remotelog' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer test1' \
--data '{
    "file": "dpkg.log",
    "n": 10,
    "hosts": [
        {
            "host_name": "http://logpuller2:3000",
            "api_key": "test2"
        },
        {
            "host_name": "http://logpuller3:3000",
            "api_key": "test3"
        }
    ]
}'
```

- You can use Swagger to try the API at http://localhost:3000/swagger/index.html

## Installation

To install the necessary dependencies and set up the project, run the following commands:

```sh
go mod tidy
```

## Running the Server

To start the server, use the following command:

```sh
API_KEY="your choice of API key" make server
```

The server will run on the port specified (default is `:3000`) with a specific API key. Using this API key to communicate with this server.

## Building the Project

Install swaggo:

```sh
go install github.com/swaggo/swag/cmd/swag@latest
```

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
make docker
docker run --rm -p 3000:3000 -d --env API_KEY="test1" logpuller
API_KEY="test1" make integ-test
```

### Benchmark Tests

To run benchmark tests, use the following command:

```sh
make benchmark-test
```

The test will prepare an apache access.log file with random information and get the last 1000 lines with filter. The file sizes is vary from 10MB to 5GB.

```
	{SizeInMB: 10},   // 10MB
	{SizeInMB: 100},  // 100MB
	{SizeInMB: 200},  // 200MB
	{SizeInMB: 500},  // 500MB
	{SizeInMB: 1000}, // 1GB
	{SizeInMB: 2000}, // 2GB
	{SizeInMB: 5000}, // 5GB
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

### Fetch Logs from Remote Hosts

- **URL**: `/v1/remote-log`
- **Method**: `POST`
- **Description**: Fetch the last N lines of logs from specified remote hosts.
- **Request Body**:
  ```json
  {
    "file": "string", // Log file to read
    "n": 10, // Number of lines to read
    "filter": "string", // Optional filter string
    "hosts": [
      {
        "host_name": "http://remote-host-1:port",
        "api_key": "remote-host-1-api-key"
      },
      {
        "host_name": "http://remote-host-2:port",
        "api_key": "remote-host-2-api-key"
      }
    ]
  }
  ```
- **Response**:
  ```json
  {
    "host_name": {
      "status": "string",
      "status_code": 200,
      "logs": ["log line 1", "log line 2", "log line 3"]
    }
  }
  ```

## Middleware

### Authorization Middleware

The service includes a middleware to check for an `Authorization` header. The header must be in the format `Bearer your-secret-token`. If the token is invalid or missing, the request will be rejected with a `401 Unauthorized` status.

### Log Middleware

The service includes a middleware to log request and response details. This middleware logs the following information:

- Request ID (based on the current timestamp in nanoseconds)
- HTTP method and request URI
- Response status code
- Duration of the request processing

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
