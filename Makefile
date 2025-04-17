sqlc:
	docker run --rm -v ./:/src -w /src sqlc/sqlc generate
