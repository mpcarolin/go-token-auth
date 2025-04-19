sqlc:
	docker run --rm -v ./:/src -w /src sqlc/sqlc generate

dbreset:
	./scripts/dbmate.sh drop && ./scripts/dbmate.sh up

dev:
	docker compose watch