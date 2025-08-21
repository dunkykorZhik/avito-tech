up:
	docker compose up --build

down:
	docker compose down

reset:
	docker compose run migrate -path /migrations \
		-database "postgres://postgres:password@db:5432/shop?sslmode=disable" down
	docker compose down -v
