login_db:
		psql $(DB_URL)

migrate_create:
		migrate create -ext sql -dir migrations -seq $(name)

migrate_up:
		migrate -database $(DB_URL) -path migrations up

migrate_down: 
		migrate -database $(DB_URL) -path migrations down 

migrate_force:
		migrate -database $(DB_URL) -path migrations force $(v)


