db_reset:
	redis-cli flushall

	sudo -u postgres psql -c "DROP DATABASE IF EXISTS real_world_dev"
	sudo -u postgres psql -c "CREATE DATABASE real_world_dev"

	make db_migrate

db_migrate:
	go run cmd/migrate_db/*.go init
	go run cmd/migrate_db/*.go

test:
	TZ= go test ./...
