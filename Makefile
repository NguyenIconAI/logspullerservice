server:
	go run main.go

sanity-test:
	go test ./api -run ^Test_HealthCheck$