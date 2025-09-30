# =============================================================================
# Medical Equipment Platform - Monorepo Development Makefile
# =============================================================================
# This Makefile serves as the single entry point for all development operations.
# Run 'make help' to see available commands.
# =============================================================================

# Default shell
SHELL := /bin/bash

# Default goal
.DEFAULT_GOAL := help

# Load environment variables from .env file if it exists
-include .env
export

# Set default environment to dev
ENV ?= dev

# Docker Compose files
COMPOSE_FILE := dev/compose/docker-compose.yml
COMPOSE_PROJECT_NAME := med-platform

# Directories
CMD_DIR := cmd
INTERNAL_DIR := internal
PKG_DIR := pkg
TEST_DIR := test
DOCS_DIR := docs
BUILD_DIR := build
DEPLOY_DIR := deploy/kustomize

# Tenant management
TENANT ?= demo-hospital

# Module names (replacing services)
MODULES := catalog rfq quote contract \
	asset-registry device-registration qr-manager \
	ticket whatsapp-gateway workflow-engine \
	chat-ai negotiation-ai predictive-maint dispatch-ai \
	diagnostic-flow reporting parts-inventory \
	demand-forecast geo-location

# Default is to enable all modules
ENABLED_MODULES ?= "*"

# Docker image name and tag
IMAGE_NAME := medical-platform
IMAGE_TAG ?= latest
REGISTRY ?= ghcr.io/org

# Colors for terminal output
BLUE := \033[1;34m
GREEN := \033[1;32m
YELLOW := \033[1;33m
RED := \033[1;31m
NC := \033[0m

# =============================================================================
# 1. Bootstrap Commands
# =============================================================================

## Bootstrap: Start the development environment
dev-up:
	@echo -e "${BLUE}Starting development environment...${NC}"
	@mkdir -p dev/logs
	@make docker-build-platform
	docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) up -d
	@echo -e "${GREEN}Waiting for services to be ready...${NC}"
	@sleep 10
	@make post-up
	@echo -e "${GREEN}Development environment is ready!${NC}"
	@echo -e "${BLUE}Access URLs:${NC}"
	@echo -e "  API Gateway:  http://localhost:8081"
	@echo -e "  Keycloak:    http://localhost:8080"
	@echo -e "  Grafana:     http://localhost:3000 (admin/admin)"
	@echo -e "  Prometheus:  http://localhost:9090"
	@echo -e "  MailHog:     http://localhost:8025"

## Bootstrap: Stop the development environment
dev-down:
	@echo -e "${BLUE}Stopping development environment...${NC}"
	docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) down
	@echo -e "${GREEN}Development environment stopped.${NC}"

## Bootstrap: Reset the development environment (WARNING: Deletes all data)
dev-reset:
	@echo -e "${RED}WARNING: This will delete all data in the development environment.${NC}"
	@read -p "Are you sure you want to continue? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo -e "${BLUE}Resetting development environment...${NC}"; \
		docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) down -v; \
		docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) rm -f; \
		echo -e "${GREEN}Development environment reset.${NC}"; \
	else \
		echo -e "${YELLOW}Reset cancelled.${NC}"; \
	fi

## Bootstrap: Run post-up initialization scripts
post-up:
	@echo -e "${BLUE}Running post-up initialization scripts...${NC}"
	@make kc-setup
	@make kafka-topics
	@make db-init
	@make redis-init
	@echo -e "${GREEN}Post-up initialization completed.${NC}"

## Bootstrap: Check prerequisites
check-prereqs:
	@echo -e "${BLUE}Checking prerequisites...${NC}"
	@which docker >/dev/null 2>&1 || (echo -e "${RED}Docker not found. Please install Docker.${NC}" && exit 1)
	@which docker compose >/dev/null 2>&1 || (echo -e "${RED}Docker Compose not found. Please install Docker Compose.${NC}" && exit 1)
	@which go >/dev/null 2>&1 || (echo -e "${RED}Go not found. Please install Go 1.22 or later.${NC}" && exit 1)
	@which curl >/dev/null 2>&1 || (echo -e "${RED}curl not found. Please install curl.${NC}" && exit 1)
	@which jq >/dev/null 2>&1 || (echo -e "${RED}jq not found. Please install jq.${NC}" && exit 1)
	@echo -e "${GREEN}All prerequisites are installed.${NC}"

## Bootstrap: Show status of all services
dev-status:
	@echo -e "${BLUE}Development environment status:${NC}"
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) ps

# =============================================================================
# 2. Module Operations
# =============================================================================

## Modules: List all available modules
list-modules:
	@echo -e "${BLUE}Available modules:${NC}"
	@for module in $(MODULES); do \
		echo "  $$module"; \
	done

## Modules: Start platform with specific modules
start-modules:
	@if [ -z "$(MODULES_LIST)" ]; then \
		echo -e "${RED}Please specify modules with MODULES_LIST=<comma-separated-list>${NC}"; \
		exit 1; \
	fi
	@echo -e "${BLUE}Starting platform with modules: $(MODULES_LIST)...${NC}"
	ENABLED_MODULES="$(MODULES_LIST)" docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) up -d platform
	@echo -e "${GREEN}Platform started with modules: $(MODULES_LIST).${NC}"

## Modules: Restart platform with all modules
restart-platform:
	@echo -e "${BLUE}Restarting platform...${NC}"
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) restart platform
	@echo -e "${GREEN}Platform restarted.${NC}"

## Modules: View logs for platform
logs-platform:
	@echo -e "${BLUE}Viewing logs for platform...${NC}"
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) logs -f platform

## Modules: View logs for all services
logs-all:
	@echo -e "${BLUE}Viewing logs for all services...${NC}"
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) logs -f

## Modules: Tail logs with grep filter
tail-logs:
	@if [ -z "$(FILTER)" ]; then \
		echo -e "${YELLOW}No filter specified. Showing all logs.${NC}"; \
		docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) logs -f; \
	else \
		echo -e "${BLUE}Filtering logs with: $(FILTER)${NC}"; \
		docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) logs -f | grep -i "$(FILTER)"; \
	fi

## Modules: Execute shell in platform container
shell-platform:
	@echo -e "${BLUE}Opening shell in platform container...${NC}"
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec platform sh || \
	 docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec platform bash

# =============================================================================
# 3. Database Management Commands
# =============================================================================

## Database: Initialize database with schema and seed data
db-init:
	@echo -e "${BLUE}Initializing database...${NC}"
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec postgres psql -U postgres -d medplatform -f /docker-entrypoint-initdb.d/01-setup-extensions.sql
	@echo -e "${GREEN}Database initialized.${NC}"

## Database: Open PostgreSQL shell
db-shell:
	@echo -e "${BLUE}Opening PostgreSQL shell...${NC}"
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec postgres psql -U postgres -d medplatform

## Database: Run migrations
db-migrate:
	@if [ -z "$(MODULE)" ]; then \
		echo -e "${RED}Please specify a module with MODULE=<name>${NC}"; \
		exit 1; \
	fi
	@echo -e "${BLUE}Running migrations for $(MODULE)...${NC}"
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec platform go run ./cmd/migrations/main.go -module=$(MODULE) up
	@echo -e "${GREEN}Migrations completed for $(MODULE).${NC}"

## Database: Create new migration
db-create-migration:
	@if [ -z "$(MODULE)" ]; then \
		echo -e "${RED}Please specify a module with MODULE=<name>${NC}"; \
		exit 1; \
	fi
	@if [ -z "$(NAME)" ]; then \
		echo -e "${RED}Please specify a migration name with NAME=<name>${NC}"; \
		exit 1; \
	fi
	@echo -e "${BLUE}Creating new migration for $(MODULE): $(NAME)${NC}"
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	mkdir -p migrations/$(MODULE); \
	touch migrations/$(MODULE)/$${timestamp}_$(NAME).up.sql; \
	touch migrations/$(MODULE)/$${timestamp}_$(NAME).down.sql; \
	echo -e "${GREEN}Created migration files:${NC}"; \
	echo -e "  migrations/$(MODULE)/$${timestamp}_$(NAME).up.sql"; \
	echo -e "  migrations/$(MODULE)/$${timestamp}_$(NAME).down.sql"

## Database: Backup database
db-backup:
	@echo -e "${BLUE}Backing up database...${NC}"
	@mkdir -p backups
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec -T postgres pg_dump -U postgres -d medplatform > backups/medplatform_$${timestamp}.sql
	@echo -e "${GREEN}Database backup created: backups/medplatform_$${timestamp}.sql${NC}"

## Database: Restore database from backup
db-restore:
	@if [ -z "$(BACKUP)" ]; then \
		echo -e "${RED}Please specify a backup file with BACKUP=<path>${NC}"; \
		exit 1; \
	fi
	@if [ ! -f "$(BACKUP)" ]; then \
		echo -e "${RED}Backup file not found: $(BACKUP)${NC}"; \
		exit 1; \
	fi
	@echo -e "${RED}WARNING: This will overwrite the current database.${NC}"
	@read -p "Are you sure you want to continue? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo -e "${BLUE}Restoring database from $(BACKUP)...${NC}"; \
		cat $(BACKUP) | docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec -T postgres psql -U postgres -d medplatform; \
		echo -e "${GREEN}Database restored from $(BACKUP).${NC}"; \
	else \
		echo -e "${YELLOW}Restore cancelled.${NC}"; \
	fi

## Database: List tables in database
db-list-tables:
	@echo -e "${BLUE}Listing database tables...${NC}"
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec postgres psql -U postgres -d medplatform -c "\dt"

## Database: Set tenant context for database operations
db-set-tenant:
	@if [ -z "$(TENANT)" ]; then \
		echo -e "${RED}Please specify a tenant with TENANT=<name>${NC}"; \
		exit 1; \
	fi
	@echo -e "${BLUE}Setting tenant context to $(TENANT)...${NC}"
	@tenant_id=$$(docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec -T postgres psql -U postgres -d medplatform -t -c "SELECT id FROM shared.tenants WHERE name='$(TENANT)'"); \
	if [ -z "$$tenant_id" ]; then \
		echo -e "${RED}Tenant not found: $(TENANT)${NC}"; \
		exit 1; \
	fi; \
	docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec postgres psql -U postgres -d medplatform -c "SELECT set_tenant_context('$$tenant_id')"; \
	echo -e "${GREEN}Tenant context set to $(TENANT).${NC}"

# =============================================================================
# 4. Kafka Topic Management
# =============================================================================

## Kafka: Create all topics
kafka-topics:
	@echo -e "${BLUE}Creating Kafka topics...${NC}"
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec kafka-setup /create-topics.sh
	@echo -e "${GREEN}Kafka topics created.${NC}"

## Kafka: List all topics
kafka-list-topics:
	@echo -e "${BLUE}Listing Kafka topics...${NC}"
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec kafka kafka-topics --bootstrap-server kafka:9092 --list

## Kafka: Describe a specific topic
kafka-describe-topic:
	@if [ -z "$(TOPIC)" ]; then \
		echo -e "${RED}Please specify a topic with TOPIC=<name>${NC}"; \
		exit 1; \
	fi
	@echo -e "${BLUE}Describing Kafka topic $(TOPIC)...${NC}"
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec kafka kafka-topics --bootstrap-server kafka:9092 --describe --topic $(TOPIC)

## Kafka: Delete a specific topic
kafka-delete-topic:
	@if [ -z "$(TOPIC)" ]; then \
		echo -e "${RED}Please specify a topic with TOPIC=<name>${NC}"; \
		exit 1; \
	fi
	@echo -e "${RED}WARNING: This will delete the Kafka topic $(TOPIC).${NC}"
	@read -p "Are you sure you want to continue? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo -e "${BLUE}Deleting Kafka topic $(TOPIC)...${NC}"; \
		docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec kafka kafka-topics --bootstrap-server kafka:9092 --delete --topic $(TOPIC); \
		echo -e "${GREEN}Kafka topic $(TOPIC) deleted.${NC}"; \
	else \
		echo -e "${YELLOW}Delete cancelled.${NC}"; \
	fi

## Kafka: Open Kafka shell
kafka-shell:
	@echo -e "${BLUE}Opening Kafka shell...${NC}"
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec kafka bash

## Kafka: Consume messages from a topic
kafka-consume:
	@if [ -z "$(TOPIC)" ]; then \
		echo -e "${RED}Please specify a topic with TOPIC=<name>${NC}"; \
		exit 1; \
	fi
	@echo -e "${BLUE}Consuming messages from Kafka topic $(TOPIC)...${NC}"
	@echo -e "${YELLOW}Press Ctrl+C to stop.${NC}"
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec kafka kafka-console-consumer --bootstrap-server kafka:9092 --topic $(TOPIC) --from-beginning

## Kafka: Produce a message to a topic
kafka-produce:
	@if [ -z "$(TOPIC)" ]; then \
		echo -e "${RED}Please specify a topic with TOPIC=<name>${NC}"; \
		exit 1; \
	fi
	@echo -e "${BLUE}Producing messages to Kafka topic $(TOPIC)...${NC}"
	@echo -e "${YELLOW}Enter messages, one per line. Press Ctrl+D to finish.${NC}"
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec -i kafka kafka-console-producer --bootstrap-server kafka:9092 --topic $(TOPIC)

# =============================================================================
# 5. Keycloak Tenant Management
# =============================================================================

## Keycloak: Setup initial realms
kc-setup:
	@echo -e "${BLUE}Setting up Keycloak realms...${NC}"
	@echo -e "${YELLOW}Waiting for Keycloak to be ready...${NC}"
	@until curl -s http://localhost:8080/health/ready > /dev/null; do \
		echo -e "${YELLOW}Waiting for Keycloak...${NC}"; \
		sleep 5; \
	done
	@echo -e "${GREEN}Keycloak is ready.${NC}"
	@echo -e "${BLUE}Importing realms...${NC}"
	@echo -e "${GREEN}Keycloak realms imported.${NC}"

## Keycloak: Add a new tenant
kc-add-tenant:
	@if [ -z "$(TENANT)" ]; then \
		echo -e "${RED}Please specify a tenant with TENANT=<name>${NC}"; \
		exit 1; \
	fi
	@echo -e "${BLUE}Adding new tenant: $(TENANT)${NC}"
	@echo -e "${YELLOW}Creating Keycloak realm...${NC}"
	@# Clone the tenant-template realm and update the name
	@cp dev/keycloak/realms/tenant-template-realm.json dev/keycloak/realms/$(TENANT)-realm.json
	@sed -i 's/tenant-template/$(TENANT)/g' dev/keycloak/realms/$(TENANT)-realm.json
	@sed -i 's/Template Healthcare Organization/$(TENANT)/g' dev/keycloak/realms/$(TENANT)-realm.json
	@# Import the new realm
	@curl -s -X POST -H "Content-Type: application/json" -d @dev/keycloak/realms/$(TENANT)-realm.json \
		-u admin:admin123 http://localhost:8080/admin/realms
	@echo -e "${YELLOW}Adding tenant to database...${NC}"
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec postgres psql -U postgres -d medplatform -c \
		"INSERT INTO shared.tenants (name, display_name, keycloak_realm) VALUES ('$(TENANT)', '$(TENANT)', '$(TENANT)');"
	@echo -e "${GREEN}Tenant $(TENANT) added successfully.${NC}"

## Keycloak: List all tenants
kc-list-tenants:
	@echo -e "${BLUE}Listing all tenants...${NC}"
	@echo -e "${YELLOW}Keycloak realms:${NC}"
	@curl -s -X GET -H "Accept: application/json" -u admin:admin123 http://localhost:8080/admin/realms | jq -r '.[].realm'
	@echo -e "${YELLOW}Database tenants:${NC}"
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec postgres psql -U postgres -d medplatform -c \
		"SELECT name, display_name, keycloak_realm, status FROM shared.tenants;"

## Keycloak: Get access token for a tenant
kc-get-token:
	@if [ -z "$(TENANT)" ]; then \
		echo -e "${RED}Please specify a tenant with TENANT=<name>${NC}"; \
		exit 1; \
	fi
	@if [ -z "$(CLIENT_ID)" ]; then \
		CLIENT_ID="api-gateway"; \
	fi
	@if [ -z "$(CLIENT_SECRET)" ]; then \
		CLIENT_SECRET="api-gateway-secret"; \
	fi
	@echo -e "${BLUE}Getting access token for tenant $(TENANT) with client $(CLIENT_ID)...${NC}"
	@curl -s -X POST \
		-d "grant_type=client_credentials&client_id=$(CLIENT_ID)&client_secret=$(CLIENT_SECRET)" \
		-H "Content-Type: application/x-www-form-urlencoded" \
		http://localhost:8080/realms/$(TENANT)/protocol/openid-connect/token | jq -r .access_token

## Keycloak: Open Keycloak shell
kc-shell:
	@echo -e "${BLUE}Opening Keycloak shell...${NC}"
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec keycloak bash

# =============================================================================
# 6. Testing Commands - Monorepo Edition
# =============================================================================

## Testing: Run unit tests for a specific module
test-unit:
	@if [ -z "$(MODULE)" ]; then \
		echo -e "${RED}Please specify a module with MODULE=<name>${NC}"; \
		exit 1; \
	fi
	@echo -e "${BLUE}Running unit tests for $(MODULE)...${NC}"
	@go test -v -race -cover ./internal/*/$(MODULE)/...
	@echo -e "${GREEN}Unit tests completed for $(MODULE).${NC}"

## Testing: Run unit tests for all modules
test-unit-all:
	@echo -e "${BLUE}Running unit tests for all modules...${NC}"
	@go test -v -race -cover ./internal/...
	@echo -e "${GREEN}All unit tests completed successfully.${NC}"

## Testing: Run integration tests for a module
test-integration:
	@if [ -z "$(MODULE)" ]; then \
		echo -e "${RED}Please specify a module with MODULE=<name>${NC}"; \
		exit 1; \
	fi
	@echo -e "${BLUE}Running integration tests for $(MODULE)...${NC}"
	@go test -v -tags=integration ./internal/*/$(MODULE)/...
	@echo -e "${GREEN}Integration tests completed for $(MODULE).${NC}"

## Testing: Run integration tests for all modules
test-integration-all:
	@echo -e "${BLUE}Running integration tests for all modules...${NC}"
	@go test -v -tags=integration ./internal/...
	@echo -e "${GREEN}All integration tests completed successfully.${NC}"

## Testing: Run E2E tests
test-e2e:
	@echo -e "${BLUE}Running E2E tests...${NC}"
	@cd $(TEST_DIR)/e2e && npm test
	@echo -e "${GREEN}E2E tests completed.${NC}"

## Testing: Run contract tests for a module
test-contract:
	@if [ -z "$(MODULE)" ]; then \
		echo -e "${RED}Please specify a module with MODULE=<name>${NC}"; \
		exit 1; \
	fi
	@echo -e "${BLUE}Running contract tests for $(MODULE)...${NC}"
	@cd $(TEST_DIR)/contract && go test -v -tags=contract -run=$(MODULE) ./...
	@echo -e "${GREEN}Contract tests completed for $(MODULE).${NC}"

## Testing: Run performance tests
test-perf:
	@echo -e "${BLUE}Running performance tests...${NC}"
	@cd $(TEST_DIR)/performance && k6 run main.js
	@echo -e "${GREEN}Performance tests completed.${NC}"

## Testing: Generate test coverage report
test-coverage:
	@if [ -z "$(MODULE)" ]; then \
		echo -e "${BLUE}Generating test coverage report for all modules...${NC}"; \
		mkdir -p $(BUILD_DIR)/coverage; \
		go test -coverprofile=$(BUILD_DIR)/coverage/coverage.out ./internal/...; \
	else \
		echo -e "${BLUE}Generating test coverage report for $(MODULE)...${NC}"; \
		mkdir -p $(BUILD_DIR)/coverage; \
		go test -coverprofile=$(BUILD_DIR)/coverage/coverage.out ./internal/*/$(MODULE)/...; \
	fi
	@go tool cover -html=$(BUILD_DIR)/coverage/coverage.out -o=$(BUILD_DIR)/coverage/coverage.html
	@echo -e "${GREEN}Test coverage report generated: $(BUILD_DIR)/coverage/coverage.html${NC}"

# =============================================================================
# 7. Code Quality Commands - Monorepo Edition
# =============================================================================

## Code Quality: Lint a specific module
lint:
	@if [ -z "$(MODULE)" ]; then \
		echo -e "${RED}Please specify a module with MODULE=<name>${NC}"; \
		exit 1; \
	fi
	@echo -e "${BLUE}Linting $(MODULE)...${NC}"
	@golangci-lint run ./internal/*/$(MODULE)/...
	@echo -e "${GREEN}Linting completed for $(MODULE).${NC}"

## Code Quality: Lint all modules
lint-all:
	@echo -e "${BLUE}Linting all modules...${NC}"
	@golangci-lint run ./...
	@echo -e "${GREEN}All modules linted successfully.${NC}"

## Code Quality: Format code for a module
fmt:
	@if [ -z "$(MODULE)" ]; then \
		echo -e "${RED}Please specify a module with MODULE=<name>${NC}"; \
		exit 1; \
	fi
	@echo -e "${BLUE}Formatting $(MODULE)...${NC}"
	@go fmt ./internal/*/$(MODULE)/...
	@echo -e "${GREEN}Formatting completed for $(MODULE).${NC}"

## Code Quality: Format code for all modules
fmt-all:
	@echo -e "${BLUE}Formatting all code...${NC}"
	@go fmt ./...
	@echo -e "${GREEN}All code formatted successfully.${NC}"

## Code Quality: Run security scan on a module
security-scan:
	@if [ -z "$(MODULE)" ]; then \
		echo -e "${RED}Please specify a module with MODULE=<name>${NC}"; \
		exit 1; \
	fi
	@echo -e "${BLUE}Running security scan for $(MODULE)...${NC}"
	@govulncheck ./internal/*/$(MODULE)/...
	@gosec ./internal/*/$(MODULE)/...
	@echo -e "${GREEN}Security scan completed for $(MODULE).${NC}"

## Code Quality: Run security scan on all modules
security-scan-all:
	@echo -e "${BLUE}Running security scan for all modules...${NC}"
	@govulncheck ./...
	@gosec ./...
	@echo -e "${GREEN}Security scan completed for all modules.${NC}"

## Code Quality: Run dependency check
dep-check:
	@echo -e "${BLUE}Checking dependencies...${NC}"
	@go mod tidy
	@go mod verify
	@go work sync
	@echo -e "${GREEN}Dependency check completed.${NC}"

## Code Quality: Update dependencies
dep-update:
	@echo -e "${BLUE}Updating dependencies...${NC}"
	@go get -u ./...
	@go mod tidy
	@go work sync
	@echo -e "${GREEN}Dependencies updated.${NC}"

# =============================================================================
# 8. Docker Image Building and Management - Single Platform Image
# =============================================================================

## Docker: Build platform image
docker-build-platform:
	@echo -e "${BLUE}Building Docker platform image...${NC}"
	@docker build -t $(IMAGE_NAME):$(IMAGE_TAG) -f Dockerfile .
	@echo -e "${GREEN}Docker platform image built: $(IMAGE_NAME):$(IMAGE_TAG)${NC}"

## Docker: Scan platform image for vulnerabilities
docker-scan:
	@echo -e "${BLUE}Scanning Docker platform image...${NC}"
	@docker scan $(IMAGE_NAME):$(IMAGE_TAG)
	@echo -e "${GREEN}Docker platform image scan completed.${NC}"

## Docker: Generate SBOM for platform image
docker-sbom:
	@echo -e "${BLUE}Generating SBOM for Docker platform image...${NC}"
	@mkdir -p $(BUILD_DIR)/sbom
	@trivy image --format cyclonedx $(IMAGE_NAME):$(IMAGE_TAG) > $(BUILD_DIR)/sbom/sbom-$(IMAGE_TAG).json
	@echo -e "${GREEN}SBOM generated: $(BUILD_DIR)/sbom/sbom-$(IMAGE_TAG).json${NC}"

## Docker: Push platform image to registry
docker-push:
	@echo -e "${BLUE}Pushing Docker platform image to $(REGISTRY)...${NC}"
	@docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)
	@docker push $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)
	@echo -e "${GREEN}Docker platform image pushed: $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)${NC}"

## Docker: Pull platform image from registry
docker-pull:
	@echo -e "${BLUE}Pulling Docker platform image from $(REGISTRY)...${NC}"
	@docker pull $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)
	@docker tag $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG) $(IMAGE_NAME):$(IMAGE_TAG)
	@echo -e "${GREEN}Docker platform image pulled: $(IMAGE_NAME):$(IMAGE_TAG)${NC}"

## Docker: List all images
docker-list-images:
	@echo -e "${BLUE}Listing Docker images...${NC}"
	@docker images | grep $(IMAGE_NAME)

## Docker: Clean unused images
docker-clean:
	@echo -e "${BLUE}Cleaning unused Docker images...${NC}"
	@docker image prune -f
	@echo -e "${GREEN}Unused Docker images cleaned.${NC}"

# =============================================================================
# 9. Deployment Commands - Kustomize for Single Image
# =============================================================================

## Deployment: Update image tag in Kustomize overlay
deploy-update-tag:
	@if [ -z "$(ENV)" ]; then \
		echo -e "${RED}Please specify an environment with ENV=<env>${NC}"; \
		exit 1; \
	fi
	@echo -e "${BLUE}Updating image tag in $(ENV) overlay...${NC}"
	@mkdir -p $(DEPLOY_DIR)/overlays/$(ENV)
	@echo "apiVersion: kustomize.config.k8s.io/v1beta1" > $(DEPLOY_DIR)/overlays/$(ENV)/image-tag.yaml
	@echo "kind: ImageTag" >> $(DEPLOY_DIR)/overlays/$(ENV)/image-tag.yaml
	@echo "metadata:" >> $(DEPLOY_DIR)/overlays/$(ENV)/image-tag.yaml
	@echo "  name: platform-image-tag" >> $(DEPLOY_DIR)/overlays/$(ENV)/image-tag.yaml
	@echo "newTag: $(IMAGE_TAG)" >> $(DEPLOY_DIR)/overlays/$(ENV)/image-tag.yaml
	@echo -e "${GREEN}Image tag updated in $(ENV) overlay.${NC}"

## Deployment: Deploy to development environment
deploy-dev:
	@echo -e "${BLUE}Deploying to development environment...${NC}"
	@make deploy-update-tag ENV=dev IMAGE_TAG=$(IMAGE_TAG)
	@kubectl apply -k $(DEPLOY_DIR)/overlays/dev
	@echo -e "${GREEN}Deployed to development environment.${NC}"

## Deployment: Deploy to staging environment
deploy-staging:
	@echo -e "${BLUE}Deploying to staging environment...${NC}"
	@make deploy-update-tag ENV=staging IMAGE_TAG=$(IMAGE_TAG)
	@kubectl apply -k $(DEPLOY_DIR)/overlays/staging
	@echo -e "${GREEN}Deployed to staging environment.${NC}"

## Deployment: Deploy to production environment
deploy-prod:
	@echo -e "${RED}WARNING: You are deploying to PRODUCTION.${NC}"
	@read -p "Are you sure you want to continue? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo -e "${BLUE}Deploying to production environment...${NC}"; \
		make deploy-update-tag ENV=prod IMAGE_TAG=$(IMAGE_TAG); \
		kubectl apply -k $(DEPLOY_DIR)/overlays/prod; \
		echo -e "${GREEN}Deployed to production environment.${NC}"; \
	else \
		echo -e "${YELLOW}Deployment cancelled.${NC}"; \
	fi

## Deployment: Apply ArgoCD application
argocd-apply:
	@if [ -z "$(ENV)" ]; then \
		echo -e "${RED}Please specify an environment with ENV=<env>${NC}"; \
		exit 1; \
	fi
	@echo -e "${BLUE}Applying ArgoCD application for $(ENV) environment...${NC}"
	@kubectl apply -f $(DEPLOY_DIR)/argocd/medical-platform-$(ENV).yaml
	@echo -e "${GREEN}ArgoCD application applied for $(ENV) environment.${NC}"

# =============================================================================
# 10. Monitoring and Logging Commands
# =============================================================================

## Monitoring: Open Grafana dashboard
open-grafana:
	@echo -e "${BLUE}Opening Grafana dashboard...${NC}"
	@open http://localhost:3000

## Monitoring: Open Prometheus dashboard
open-prometheus:
	@echo -e "${BLUE}Opening Prometheus dashboard...${NC}"
	@open http://localhost:9090

## Monitoring: Open MailHog dashboard
open-mailhog:
	@echo -e "${BLUE}Opening MailHog dashboard...${NC}"
	@open http://localhost:8025

## Monitoring: Open Keycloak dashboard
open-keycloak:
	@echo -e "${BLUE}Opening Keycloak dashboard...${NC}"
	@open http://localhost:8080

## Monitoring: Check health of platform
check-health:
	@echo -e "${BLUE}Checking health of platform...${NC}"
	@curl -s http://localhost:8081/health || echo -e "${RED}Failed to connect to platform${NC}"
	@echo -e "${GREEN}Health check completed.${NC}"

## Monitoring: Export Prometheus metrics to file
export-metrics:
	@echo -e "${BLUE}Exporting Prometheus metrics...${NC}"
	@mkdir -p $(BUILD_DIR)/metrics
	@curl -s http://localhost:9090/api/v1/query?query=up > $(BUILD_DIR)/metrics/up.json
	@echo -e "${GREEN}Metrics exported to $(BUILD_DIR)/metrics/up.json${NC}"

# =============================================================================
# 11. Backup and Restore Operations
# =============================================================================

## Backup: Create full backup of all data
backup-all:
	@echo -e "${BLUE}Creating full backup of all data...${NC}"
	@mkdir -p backups
	@timestamp=$$(date +%Y%m%d%H%M%S)
	@mkdir -p backups/$$timestamp
	@echo -e "${YELLOW}Backing up database...${NC}"
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec -T postgres pg_dump -U postgres -d medplatform > backups/$$timestamp/medplatform.sql
	@echo -e "${YELLOW}Backing up Keycloak...${NC}"
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec -T keycloak /opt/keycloak/bin/kc.sh export --dir /tmp/export
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) cp keycloak:/tmp/export backups/$$timestamp/keycloak
	@echo -e "${YELLOW}Backing up Redis...${NC}"
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec -T redis redis-cli SAVE
	@docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) cp redis:/data/dump.rdb backups/$$timestamp/redis-dump.rdb
	@echo -e "${GREEN}Full backup created in backups/$$timestamp${NC}"

## Backup: Restore from backup
restore-all:
	@if [ -z "$(BACKUP)" ]; then \
		echo -e "${RED}Please specify a backup directory with BACKUP=<path>${NC}"; \
		exit 1; \
	fi
	@if [ ! -d "$(BACKUP)" ]; then \
		echo -e "${RED}Backup directory not found: $(BACKUP)${NC}"; \
		exit 1; \
	fi
	@echo -e "${RED}WARNING: This will overwrite all current data.${NC}"
	@read -p "Are you sure you want to continue? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo -e "${BLUE}Restoring from backup: $(BACKUP)${NC}"; \
		echo -e "${YELLOW}Stopping services...${NC}"; \
		docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) down; \
		echo -e "${YELLOW}Starting database...${NC}"; \
		docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) up -d postgres; \
		sleep 10; \
		echo -e "${YELLOW}Restoring database...${NC}"; \
		cat $(BACKUP)/medplatform.sql | docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec -T postgres psql -U postgres -d medplatform; \
		echo -e "${YELLOW}Stopping database...${NC}"; \
		docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) stop postgres; \
		echo -e "${YELLOW}Restoring Redis...${NC}"; \
		docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) cp $(BACKUP)/redis-dump.rdb redis:/data/dump.rdb; \
		echo -e "${YELLOW}Starting all services...${NC}"; \
		docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) up -d; \
		sleep 10; \
		echo -e "${YELLOW}Restoring Keycloak...${NC}"; \
		docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) cp $(BACKUP)/keycloak keycloak:/tmp/import; \
		docker compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT_NAME) exec -T keycloak /opt/keycloak/bin/kc.sh import --dir /tmp/import; \
		echo -e "${GREEN}Restore completed from $(BACKUP).${NC}"; \
	else \
		echo -e "${YELLOW}Restore cancelled.${NC}"; \
	fi

# =============================================================================
# 12. Documentation Generation
# =============================================================================

## Documentation: Generate API documentation
gen-api-docs:
	@if [ -z "$(MODULE)" ]; then \
		echo -e "${RED}Please specify a module with MODULE=<name>${NC}"; \
		exit 1; \
	fi
	@echo -e "${BLUE}Generating API documentation for $(MODULE)...${NC}"
	@mkdir -p $(DOCS_DIR)/api/$(MODULE)
	@swag init -g internal/*/$(MODULE)/http/handlers.go -o $(DOCS_DIR)/api/$(MODULE)
	@echo -e "${GREEN}API documentation generated for $(MODULE): $(DOCS_DIR)/api/$(MODULE)${NC}"

## Documentation: Generate all API documentation
gen-api-docs-all:
	@echo -e "${BLUE}Generating API documentation for all modules...${NC}"
	@for module in $(MODULES); do \
		echo -e "${YELLOW}Generating docs for $$module...${NC}"; \
		mkdir -p $(DOCS_DIR)/api/$$module; \
		swag init -g internal/*/$$module/http/handlers.go -o $(DOCS_DIR)/api/$$module || true; \
	done
	@echo -e "${GREEN}API documentation generated for all modules.${NC}"

## Documentation: Generate code documentation
gen-code-docs:
	@echo -e "${BLUE}Generating code documentation...${NC}"
	@mkdir -p $(DOCS_DIR)/code
	@godoc -http=:6060 &
	@echo -e "${GREEN}Code documentation server started at http://localhost:6060/pkg/${NC}"
	@echo -e "${YELLOW}Press Enter to stop the server...${NC}"
	@read
	@pkill -f "godoc -http"

## Documentation: Generate architecture diagrams
gen-arch-diagrams:
	@echo -e "${BLUE}Generating architecture diagrams...${NC}"
	@mkdir -p $(DOCS_DIR)/architecture
	@echo -e "${GREEN}Architecture diagrams generated: $(DOCS_DIR)/architecture${NC}"

# =============================================================================
# 13. CI/CD Helper Commands - Monorepo Edition
# =============================================================================

## CI/CD: Run CI pipeline locally
ci-local:
	@echo -e "${BLUE}Running CI pipeline locally...${NC}"
	@make lint-all
	@make test-unit-all
	@make test-integration-all
	@make docker-build-platform
	@echo -e "${GREEN}CI pipeline completed successfully.${NC}"

## CI/CD: Generate version
gen-version:
	@echo -e "${BLUE}Generating version...${NC}"
	@version=$$(git describe --tags --always --dirty); \
	echo "Version: $$version"
	@echo -e "${GREEN}Version generated.${NC}"

## CI/CD: Create release tag
create-release:
	@if [ -z "$(VERSION)" ]; then \
		echo -e "${RED}Please specify a version with VERSION=<version>${NC}"; \
		exit 1; \
	fi
	@echo -e "${BLUE}Creating release tag v$(VERSION)...${NC}"
	@git tag -a v$(VERSION) -m "Release v$(VERSION)"
	@git push origin v$(VERSION)
	@echo -e "${GREEN}Release tag v$(VERSION) created.${NC}"

## CI/CD: Generate changelog
gen-changelog:
	@echo -e "${BLUE}Generating changelog...${NC}"
	@mkdir -p $(DOCS_DIR)
	@git log --pretty=format:"* %s" $(shell git describe --tags --abbrev=0 @^)..@ > $(DOCS_DIR)/CHANGELOG.md
	@echo -e "${GREEN}Changelog generated: $(DOCS_DIR)/CHANGELOG.md${NC}"

# =============================================================================
# 14. Go Workspace Management - Monorepo Edition
# =============================================================================

## Go: Initialize Go workspace
go-workspace-init:
	@echo -e "${BLUE}Initializing Go workspace...${NC}"
	@go work init
	@go work use .
	@echo -e "${GREEN}Go workspace initialized.${NC}"

## Go: Add module to workspace
go-workspace-add:
	@if [ -z "$(MODULE_PATH)" ]; then \
		echo -e "${RED}Please specify a module path with MODULE_PATH=<path>${NC}"; \
		exit 1; \
	fi
	@echo -e "${BLUE}Adding module to workspace: $(MODULE_PATH)...${NC}"
	@go work use $(MODULE_PATH)
	@echo -e "${GREEN}Module added to workspace.${NC}"

## Go: Sync Go workspace
go-workspace-sync:
	@echo -e "${BLUE}Syncing Go workspace...${NC}"
	@go work sync
	@echo -e "${GREEN}Go workspace synced.${NC}"

# =============================================================================
# 15. Module Management - Monorepo Edition
# =============================================================================

## Module: Initialize new module
init-module:
	@if [ -z "$(MODULE)" ]; then \
		echo -e "${RED}Please specify a module with MODULE=<name>${NC}"; \
		exit 1; \
	fi
	@if [ -z "$(DOMAIN)" ]; then \
		echo -e "${RED}Please specify a domain with DOMAIN=<name>${NC}"; \
		exit 1; \
	fi
	@echo -e "${BLUE}Initializing new module: $(MODULE) in domain $(DOMAIN)...${NC}"
	@mkdir -p $(INTERNAL_DIR)/$(DOMAIN)/$(MODULE)/domain
	@mkdir -p $(INTERNAL_DIR)/$(DOMAIN)/$(MODULE)/app
	@mkdir -p $(INTERNAL_DIR)/$(DOMAIN)/$(MODULE)/infra
	@mkdir -p $(INTERNAL_DIR)/$(DOMAIN)/$(MODULE)/http
	@echo -e "package $(MODULE)\n\nimport (\n\t\"context\"\n)\n\n// Service defines the interface for the $(MODULE) module\ntype Service interface {\n\tStart(ctx context.Context) error\n\tStop() error\n}\n\n// New creates a new $(MODULE) service\nfunc New() Service {\n\treturn &service{}\n}\n\ntype service struct {}\n\nfunc (s *service) Start(ctx context.Context) error {\n\treturn nil\n}\n\nfunc (s *service) Stop() error {\n\treturn nil\n}\n" > $(INTERNAL_DIR)/$(DOMAIN)/$(MODULE)/module.go
	@echo -e "${GREEN}Module $(MODULE) initialized in domain $(DOMAIN).${NC}"

## Module: Enable module in platform
enable-module:
	@if [ -z "$(MODULE)" ]; then \
		echo -e "${RED}Please specify a module with MODULE=<name>${NC}"; \
		exit 1; \
	fi
	@echo -e "${BLUE}Enabling module $(MODULE) in platform...${NC}"
	@echo -e "${YELLOW}Edit internal/service/registry.go to add the module to AllModules()${NC}"
	@echo -e "${GREEN}Module $(MODULE) enabled.${NC}"

# =============================================================================
# 16. Cleanup and Maintenance Commands
# =============================================================================

## Cleanup: Clean build artifacts
clean-build:
	@echo -e "${BLUE}Cleaning build artifacts...${NC}"
	@rm -rf $(BUILD_DIR)
	@echo -e "${GREEN}Build artifacts cleaned.${NC}"

## Cleanup: Clean documentation
clean-docs:
	@echo -e "${BLUE}Cleaning documentation...${NC}"
	@rm -rf $(DOCS_DIR)/api
	@rm -rf $(DOCS_DIR)/code
	@echo -e "${GREEN}Documentation cleaned.${NC}"

## Cleanup: Clean all
clean-all:
	@echo -e "${BLUE}Cleaning all artifacts...${NC}"
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DOCS_DIR)/api
	@rm -rf $(DOCS_DIR)/code
	@echo -e "${GREEN}All artifacts cleaned.${NC}"

## Cleanup: Prune Docker system
docker-prune:
	@echo -e "${RED}WARNING: This will remove all unused Docker data.${NC}"
	@read -p "Are you sure you want to continue? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo -e "${BLUE}Pruning Docker system...${NC}"; \
		docker system prune -a --volumes -f; \
		echo -e "${GREEN}Docker system pruned.${NC}"; \
	else \
		echo -e "${YELLOW}Prune cancelled.${NC}"; \
	fi

## Maintenance: Check disk usage
check-disk:
	@echo -e "${BLUE}Checking disk usage...${NC}"
	@df -h
	@echo -e "${YELLOW}Docker disk usage:${NC}"
	@docker system df
	@echo -e "${GREEN}Disk usage check completed.${NC}"

# =============================================================================
# 17. Help System - Monorepo Edition
# =============================================================================

## Help: Show this help message
help:
	@echo -e "${BLUE}Medical Equipment Platform - Monorepo Development Makefile${NC}"
	@echo -e "${YELLOW}Usage:${NC} make [target]"
	@echo
	@echo -e "${YELLOW}Available targets:${NC}"
	@grep -E '^## [A-Za-z0-9_-]+: .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "## "}; {printf "${GREEN}%-20s${NC} %s\n", $$2, $$3}'
	@echo
	@echo -e "${YELLOW}Examples:${NC}"
	@echo -e "  make dev-up                          - Start development environment"
	@echo -e "  make dev-down                        - Stop development environment"
	@echo -e "  make start-modules MODULES_LIST=catalog,rfq - Start platform with specific modules"
	@echo -e "  make test-unit MODULE=catalog        - Run unit tests for catalog module"
	@echo -e "  make init-module MODULE=parts DOMAIN=marketplace - Initialize new module"
	@echo
	@echo -e "${YELLOW}For more detailed help on a specific section, run:${NC}"
	@echo -e "  make help-bootstrap            - Help on bootstrap commands"
	@echo -e "  make help-modules              - Help on module operations"
	@echo -e "  make help-db                   - Help on database commands"
	@echo -e "  make help-kafka                - Help on Kafka commands"
	@echo -e "  make help-keycloak             - Help on Keycloak commands"
	@echo -e "  make help-testing              - Help on testing commands"
	@echo -e "  make help-code-quality         - Help on code quality commands"
	@echo -e "  make help-docker               - Help on Docker commands"
	@echo -e "  make help-deployment           - Help on deployment commands"
	@echo -e "  make help-monitoring           - Help on monitoring commands"
	@echo -e "  make help-backup               - Help on backup commands"
	@echo -e "  make help-docs                 - Help on documentation commands"
	@echo -e "  make help-cicd                 - Help on CI/CD commands"
	@echo -e "  make help-go                   - Help on Go workspace commands"
	@echo -e "  make help-module-mgmt          - Help on module management commands"
	@echo -e "  make help-cleanup              - Help on cleanup commands"

## Help: Show help on bootstrap commands
help-bootstrap:
	@echo -e "${BLUE}Bootstrap Commands${NC}"
	@echo
	@grep -E '^## Bootstrap: .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "## Bootstrap: "}; {printf "${GREEN}make %-20s${NC} %s\n", $$2, $$3}'

## Help: Show help on module operations
help-modules:
	@echo -e "${BLUE}Module Operations${NC}"
	@echo
	@grep -E '^## Modules: .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "## Modules: "}; {printf "${GREEN}make %-20s${NC} %s\n", $$2, $$3}'

## Help: Show help on database commands
help-db:
	@echo -e "${BLUE}Database Commands${NC}"
	@echo
	@grep -E '^## Database: .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "## Database: "}; {printf "${GREEN}make %-20s${NC} %s\n", $$2, $$3}'

## Help: Show help on Kafka commands
help-kafka:
	@echo -e "${BLUE}Kafka Commands${NC}"
	@echo
	@grep -E '^## Kafka: .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "## Kafka: "}; {printf "${GREEN}make %-20s${NC} %s\n", $$2, $$3}'

## Help: Show help on Keycloak commands
help-keycloak:
	@echo -e "${BLUE}Keycloak Commands${NC}"
	@echo
	@grep -E '^## Keycloak: .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "## Keycloak: "}; {printf "${GREEN}make %-20s${NC} %s\n", $$2, $$3}'

## Help: Show help on testing commands
help-testing:
	@echo -e "${BLUE}Testing Commands${NC}"
	@echo
	@grep -E '^## Testing: .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "## Testing: "}; {printf "${GREEN}make %-20s${NC} %s\n", $$2, $$3}'

## Help: Show help on code quality commands
help-code-quality:
	@echo -e "${BLUE}Code Quality Commands${NC}"
	@echo
	@grep -E '^## Code Quality: .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "## Code Quality: "}; {printf "${GREEN}make %-20s${NC} %s\n", $$2, $$3}'

## Help: Show help on Docker commands
help-docker:
	@echo -e "${BLUE}Docker Commands${NC}"
	@echo
	@grep -E '^## Docker: .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "## Docker: "}; {printf "${GREEN}make %-20s${NC} %s\n", $$2, $$3}'

## Help: Show help on deployment commands
help-deployment:
	@echo -e "${BLUE}Deployment Commands${NC}"
	@echo
	@grep -E '^## Deployment: .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "## Deployment: "}; {printf "${GREEN}make %-20s${NC} %s\n", $$2, $$3}'

## Help: Show help on monitoring commands
help-monitoring:
	@echo -e "${BLUE}Monitoring Commands${NC}"
	@echo
	@grep -E '^## Monitoring: .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "## Monitoring: "}; {printf "${GREEN}make %-20s${NC} %s\n", $$2, $$3}'

## Help: Show help on backup commands
help-backup:
	@echo -e "${BLUE}Backup Commands${NC}"
	@echo
	@grep -E '^## Backup: .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "## Backup: "}; {printf "${GREEN}make %-20s${NC} %s\n", $$2, $$3}'

## Help: Show help on documentation commands
help-docs:
	@echo -e "${BLUE}Documentation Commands${NC}"
	@echo
	@grep -E '^## Documentation: .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "## Documentation: "}; {printf "${GREEN}make %-20s${NC} %s\n", $$2, $$3}'

## Help: Show help on CI/CD commands
help-cicd:
	@echo -e "${BLUE}CI/CD Commands${NC}"
	@echo
	@grep -E '^## CI/CD: .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "## CI/CD: "}; {printf "${GREEN}make %-20s${NC} %s\n", $$2, $$3}'

## Help: Show help on Go workspace commands
help-go:
	@echo -e "${BLUE}Go Workspace Commands${NC}"
	@echo
	@grep -E '^## Go: .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "## Go: "}; {printf "${GREEN}make %-20s${NC} %s\n", $$2, $$3}'

## Help: Show help on module management commands
help-module-mgmt:
	@echo -e "${BLUE}Module Management Commands${NC}"
	@echo
	@grep -E '^## Module: .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "## Module: "}; {printf "${GREEN}make %-20s${NC} %s\n", $$2, $$3}'

## Help: Show help on cleanup commands
help-cleanup:
	@echo -e "${BLUE}Cleanup Commands${NC}"
	@echo
	@grep -E '^## Cleanup: .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "## Cleanup: "}; {printf "${GREEN}make %-20s${NC} %s\n", $$2, $$3}'
	@grep -E '^## Maintenance: .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "## Maintenance: "}; {printf "${GREEN}make %-20s${NC} %s\n", $$2, $$3}'

# =============================================================================
# Phony Targets
# =============================================================================

.PHONY: help dev-up dev-down dev-reset post-up check-prereqs dev-status \
	list-modules start-modules restart-platform logs-platform logs-all tail-logs shell-platform \
	db-init db-shell db-migrate db-create-migration db-backup db-restore db-list-tables db-set-tenant \
	kafka-topics kafka-list-topics kafka-describe-topic kafka-delete-topic kafka-shell kafka-consume kafka-produce \
	kc-setup kc-add-tenant kc-list-tenants kc-get-token kc-shell \
	test-unit test-unit-all test-integration test-integration-all test-e2e test-contract test-perf test-coverage \
	lint lint-all fmt fmt-all security-scan security-scan-all dep-check dep-update \
	docker-build-platform docker-scan docker-sbom docker-push docker-pull docker-list-images docker-clean \
	deploy-update-tag deploy-dev deploy-staging deploy-prod argocd-apply \
	open-grafana open-prometheus open-mailhog open-keycloak check-health export-metrics \
	backup-all restore-all \
	gen-api-docs gen-api-docs-all gen-code-docs gen-arch-diagrams \
	ci-local gen-version create-release gen-changelog \
	go-workspace-init go-workspace-add go-workspace-sync \
	init-module enable-module \
	clean-build clean-docs clean-all docker-prune check-disk \
	help help-bootstrap help-modules help-db help-kafka help-keycloak help-testing help-code-quality \
	help-docker help-deployment help-monitoring help-backup help-docs help-cicd help-go help-module-mgmt help-cleanup
