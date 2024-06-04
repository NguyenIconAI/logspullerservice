server:
	go run main.go

integ-test:
	go test ./... -tags=integration -v

test:
	go test ./... -tags=unit -v