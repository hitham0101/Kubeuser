
project := kubeuser

run: ## run kubeuser program
	@go mod tidy
	@go run .


build: ## build kubeuser binary file in bin directory
	@go mod tidy
	@go build -o ./bin/ ./...

destroy: ## remove kubeuser binary file in bin directory and .kubeuser directory
	@rm -rf ./bin/*
	@rm -rf ~/.kubeuser
