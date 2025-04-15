migrate-up:
	goose -dir migrations postgres "$$DB_DSN" up

migrate-down:
	goose -dir migrations postgres "$$DB_DSN" down