## Description

Simple demonstration of token authentication and secure user password storage with bcrypt, served by a REST API, built primarily with Go.

### Setup
1. Create a `.env` file, using `.env.template` as guide
2. Ensure you have docker compose installed
3. `make dev` to start all containers
4. `make test` to run api tests (currently expects containers are up)
5. See http requests in `requests/` folder for testing using an extension like vscode's REST client

### Other notes

- Uses dbmate for database migrations, overkill for this demo, but I wanted to explore the tool
- Uses sqlc library to compile sql in ./db into Go functions in the api library. Run `make sqlc` to compile