# Default action is to show help
.PHONY: help
help:
	@echo "Makefile for managing Docker environments"
	@echo ""
	@echo "Usage:"
	@echo "  make help                   Show this help message"
	@echo "  make nacos-up               Start the nacos server"
	@echo "  make nacos-stop             Stop the nacos server"
	@echo "  make nacos-down             Stop and remove the nacos server"
	@echo "  make dev-up                 Start the development environment"
	@echo "  make prod-up                Start the production environment"
	@echo "  make dev-stop               Stop the development environment"
	@echo "  make prod-stop              Stop the production environment"
	@echo "  make dev-down               Clean the development environment"
	@echo "  make prod-down              Clean the production environment"
	@echo "  make set-rabbitmq-timeout   Set the rabbitmq consumer_timeout to undefined"

# Start nacos
.PHONY: nacos-up
nacos-up:
	@if [ -z $$(docker network ls --filter name=^flow-federate-network$$ --format="{{ .Name }}") ]; then \
		echo "Creating network flow-federate-network"; \
		docker network create flow-federate-network; \
	fi
	@if [ -n $$(docker ps -a --filter name=^/nacos$$ --format="{{ .Names }}") ]; then \
		if [ -z $$(docker ps --filter name=^/nacos$$ --format="{{ .Names }}") ]; then \
			echo "Starting existing nacos container"; \
			docker start nacos; \
		else \
			echo "A container named 'nacos' is already running."; \
		fi; \
		if [ -z $$(docker inspect nacos --format '{{json .NetworkSettings.Networks}}' | grep flow-federate-network) ]; then \
			echo "Connecting nacos container to flow-federate-network"; \
			docker network connect flow-federate-network nacos || echo "nacos is already connected to flow-federate-network"; \
		fi \
	else \
		echo "Creating and starting a new nacos container"; \
		docker run --name nacos --network flow-federate-network -e MODE=standalone -e JVM_XMS=512m -e JVM_XMX=512m -e JVM_XMN=256m -p 8848:8848 -p 9848:9848 -d nacos/nacos-server:latest; \
	fi

# Stop nacos server
.PHONY: nacos-stop
nacos-stop:
	@if [ -n $$(docker ps --filter name=^/nacos$$ --format="{{ .Names }}") ]; then \
		docker stop nacos; \
	else \
		echo "No running container named 'nacos' found."; \
	fi

# Stop and remove nacos server
.PHONY: nacos-down
nacos-down:
	@if [ -n $$(docker ps -a --filter name=^/nacos$$ --format="{{ .Names }}") ]; then \
		echo "Disconnecting nacos container from flow-federate-network"; \
		docker network disconnect flow-federate-network nacos || true; \
		echo "Stopping and removing nacos container"; \
		docker stop nacos && docker rm nacos; \
		if [ -z $$(docker ps -a --filter network=flow-federate-network --format="{{ .Names }}") ]; then \
			echo "Removing network flow-federate-network"; \
			docker network rm flow-federate-network; \
		fi \
	else \
		echo "No container named 'nacos' found."; \
	fi

# Start development environment
.PHONY: dev-up
dev-up:
	docker-compose up -d
	#docker exec -it flow-federate-rabbitmq rabbitmqctl eval 'application:set_env(rabbit,consumer_timeout,undefined).'

# Stop development environment
.PHONY: dev-stop
dev-stop:
	docker-compose stop

# Clean development environment
.PHONY: dev-down
dev-down:
	docker-compose down

# Start production environment
.PHONY: prod-up
prod-up:
	docker-compose -f docker-compose.prod.yml up -d
	#docker exec -it flow-federate-rabbitmq rabbitmqctl eval 'application:set_env(rabbit,consumer_timeout,undefined).'

# Stop production environment
.PHONY: prod-stop
prod-stop:
	docker-compose -f docker-compose.prod.yml stop

# Clean production environment
.PHONY: prod-down
prod-down:
	docker-compose -f docker-compose.prod.yml down

# Set the rabbitmq consumer_timeout to undefined
.PHONY: set-rabbitmq-timeout
set-rabbitmq-timeout:
	docker exec -it flow-federate-rabbitmq rabbitmqctl eval 'application:set_env(rabbit,consumer_timeout,undefined).'
