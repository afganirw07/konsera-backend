
.PHONY: migrate-up migrate-down migrate-version migrate-create

MIGRATE := migrate -path ./migration -database "$(DATABASE_URL)"

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down 1

migrate-version:
	$(MIGRATE) version

migrate-create:
	migrate create -ext sql -dir migration -seq $(name)
