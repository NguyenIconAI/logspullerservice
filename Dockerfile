# Build the application from source
FROM golang:1.22-alpine3.19 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/logpuller main.go

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test ./... -tags=unit -v

# Deploy the application binary into a lean image
FROM ubuntu:22.04 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/bin/logpuller /logpuller

EXPOSE 3000

ENTRYPOINT ["/logpuller"]