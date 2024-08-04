ifneq (,$(wildcard ./.env))
    include .env
    export
endif


migration_up: migrate -path database/migrations/ -database "mysql://$(DB_USERNAME):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up

migration_down: migrate -path database/migrations/ -database "mysql://$(DB_USERNAME):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose down

migration_fix: migrate -path database/migrations/ -database "mysql://$(DB_USERNAME):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable" force VERSION

migrate -path database/migrations/ -database "mysql://root:pundadevo25052002@tcp(localhost:3306)/loan_db?sslmode=disable" -verbose down