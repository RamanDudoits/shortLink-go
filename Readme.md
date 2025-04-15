go get github.com/joho/godotenv
go get github.com/jackc/pgx/v5
go get github.com/pressly/goose/v3
go install github.com/pressly/goose/v3/cmd/goose@latest

export DB_DSN="postgres://user:password@localhost:5432/shortlink?sslmode=disable"

docker exec -it shortlink-db psql -U user -d shortlink
\dt  # Показать все таблицы
\q  # Выйти

go get golang.org/x/crypto/bcrypt
go get github.com/golang-jwt/jwt/v5