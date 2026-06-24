.PHONY: dev build up down logs seed

# Development
dev:
	go run .

# Build Docker image
build:
	docker build -t invitation-app:latest .

# Jalankan development dengan Docker
up:
	docker-compose up -d

# Stop semua container
down:
	docker-compose down

# Lihat logs
logs:
	docker-compose logs -f app

# Seed database
seed:
	docker-compose exec app ./invitation-app --seed

# Production up
prod-up:
	docker-compose -f docker-compose.prod.yml up -d

# Production down
prod-down:
	docker-compose -f docker-compose.prod.yml down

# Rebuild & restart
rebuild:
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d