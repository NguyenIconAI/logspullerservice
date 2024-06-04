server:
	go run main.go

benchmark-test:
	go test ./... -bench=.

integ-test:
	go test ./... -tags=integration -v

test:
	go test ./... -tags=unit -v