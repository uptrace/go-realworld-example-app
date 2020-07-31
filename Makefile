db_reset:
	redis-cli flushall

	sudo -u postgres psql -c "DROP DATABASE IF EXISTS real_world_dev"
	sudo -u postgres psql -c "CREATE DATABASE real_world_dev"

	go run cmd/migrate_db/*.go init
	go run cmd/migrate_db/*.go
