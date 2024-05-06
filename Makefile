
project := kubeuser

run: ## run kubeuser program
	@go mod tidy
	@go run .


build: ## build kubeuser binary file in bin directory
	@go mod tidy
	@go build -o ./bin/ ./...

