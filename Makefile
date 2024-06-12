# Makefile to control Docker compose for development and production environments

# Default action is to show help
.PHONY: help
help:
	@echo "Makefile for managing Docker environments"
	@echo ""
	@echo "Usage:"
	@echo "  make help                   Show this help message"
	@echo "  make dev-up                 Start the development environment"
	@echo "  make prod-up                Start the production environment"
	@echo "  make dev-stop               Stop the development environment"
	@echo "  make prod-stop              Stop the production environment"
	@echo "  make dev-down               clean the development environment"
	@echo "  make prod-down              clean the production environment"
	@echo "  make set-rabbitmq-timeout   Set the rabbitmq consumer_timeout to undefined"


# Start development environment
.PHONY: dev-up
dev-up:
	docker-compose build
	docker-compose up -d
	docker exec -it flow-federate-rabbitmq rabbitmqctl eval 'application:set_env(rabbit,consumer_timeout,undefined).'

# Stop development environment
.PHONY: dev-stop
dev-stop:
	docker-compose stop

# Stop development environment
.PHONY: dev-down
dev-down:
	docker-compose down

# Start production environment
.PHONY: prod-up
prod-up:
	docker-compose -f docker-compose.prod.yml up -d

# Stop production environment
.PHONY: prod-stop
prod-stop:
	docker-compose -f docker-compose.prod.yml stop

# Stop production environment
.PHONY: prod-down
prod-down:
	docker-compose -f docker-compose.prod.yml down

# Set the rabbitmq consumer_timeout to undefined
.PHONY: set-rabbitmq-timeout
set-rabbitmq-timeout:
	docker exec -it flow-federate-rabbitmq rabbitmqctl eval 'application:set_env(rabbit,consumer_timeout,undefined).'
