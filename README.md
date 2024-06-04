# Log Puller Service

This is a log puller service that provides a RESTful API for retrieving log files and specific log entries from a server. The service is built using Go and includes Swagger documentation for the API endpoints.

## Features

- Health check endpoint to verify the server status.
- Retrieve a list of log files from a specified directory.
- Read the last N lines from a specific log file with optional filtering.
- Swagger documentation for API endpoints.

## Installation

To install the necessary dependencies and set up the project, run the following commands:

```sh
go mod tidy
```

## Running the Server

To start the server, use the following command:

```sh
go run main.go
```

The server will run on the port specified (default is `:3000`).

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

TODO: Adding authentication and logging middleware.

## Swagger Documentation

The API documentation is generated using Swagger. To generate the documentation, run:

```sh
swag init -g main.go
```

You can access the Swagger documentation at `http://localhost:3000/swagger/index.html`.

## Makefile Commands

The project includes a `Makefile` with the following commands:

- **Start the server**:

  ```sh
  go run main.go
  ```

- **Run sanity check**:

  ```sh
  make integ-test
  ```

- **Run unit tests**:

  ```sh
  make test
  ```

- **Run benchmarks**:
  ```sh
  make benchmark-test
  ```

## Contributing

Feel free to fork this repository and submit pull requests. For major changes, please open an issue first to discuss what you would like to change.

## License

[Insert License Here]

## Contact

For any issues or inquiries, please contact [vdnguyen58@gmail.com].
