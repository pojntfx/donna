clean:
	# rm -rf pkg/models
	docker rm -f networkmate-postgres

depend: clean
	docker run -d --name networkmate-postgres -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_DB=networkmate postgres
	docker exec networkmate-postgres bash -c 'until pg_isready; do sleep 1; done'

	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/volatiletech/sqlboiler/v4@latest
	go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest

	goose -dir './pkg/migrations' postgres "host=localhost user=postgres dbname=networkmate sslmode=disable" up
	go generate ./...