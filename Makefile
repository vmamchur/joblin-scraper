build:
	docker compose build 

run:	build
	docker compose up

generate-sql:
	sqlc generate
