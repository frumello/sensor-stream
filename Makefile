.PHONY: run
run: ## Build the binary
	@echo "build and run"
	docker compose up --build -d

stop: ## Build the binary
	@echo "stop"
	docker compose down