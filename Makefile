clean:
	rm -rf internal/db
	docker rm -f networkmate-postgres

depend: clean
	docker run -d --name networkmate-postgres -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_DB=networkmate postgres
	docker exec networkmate-postgres bash -c 'until pg_isready; do sleep 1; done'

	go install github.com/rubenv/sql-migrate/sql-migrate@latest
	go install github.com/volatiletech/sqlboiler/v4@latest
	go install github.com/jteeuwen/go-bindata/go-bindata@latest
	go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest

	sql-migrate up -env="psql" -config sql-migrate.yaml
	go generate ./...