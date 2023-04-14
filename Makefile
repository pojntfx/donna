clean:
	rm -rf internal/models
	docker rm -f networkmate-postgres

depend: clean
	docker run -d --name networkmate-postgres -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_DB=networkmate postgres
	docker exec networkmate-postgres bash -c 'until pg_isready; do sleep 1; done'

	go install github.com/kyleconroy/sqlc/cmd/sqlc@latest

	go generate ./...