sqlc:
	docker run --rm -v ./:/src -w /src sqlc/sqlc generate

reset:
	./scripts/dbmate.sh drop && ./scripts/dbmate.sh up