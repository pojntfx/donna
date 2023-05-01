clean:
	rm -rf pkg/models
	docker rm -f donna-postgres

depend: clean
	docker run -d --name donna-postgres -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_DB=donna postgres
	docker exec donna-postgres bash -c 'until pg_isready; do sleep 1; done'

	go install github.com/kyleconroy/sqlc/cmd/sqlc@latest

	go generate ./...