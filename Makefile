server:
	go run main.go

docker:
	docker build -t logpuller .

build:
	swag init -g main.go
	go build -o bin/logpuller main.go

benchmark-test:
	go test ./... -bench=.

integ-test:
	go test ./... -tags=integration -v

test:
	go test ./... -tags=unit -v