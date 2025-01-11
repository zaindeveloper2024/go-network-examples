# app-server

```sh
# setup
cp .env.example .env

# run
go run cmd/server/main.go
```

## API Usage

```sh
http post http://localhost:8080/users name="John Doe"

http http://localhost:8080/users

http http://localhost:8080/users/user_id

http put http://localhost:8080/users/{id} name="John Updated"

http delete http://localhost:8080/users/user_id
```

## Docker

```sh
docker run --rm --name my-postgres \
  -e POSTGRES_USER=username \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=dbname \
  -p 5432:5432 \
  -d postgres:15
docker exec -it my-postgres psql -U postgres -d dbname
docker stop my-postgres
```

## Migration

```sh
golang-migrate
```