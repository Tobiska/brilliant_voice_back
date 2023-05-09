test: ## Run App test
	go test -v -coverprofile=cover.out ./...


run: ## Run App test
	go run -v cmd/app/main.go