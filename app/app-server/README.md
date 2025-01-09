# app-server

```sh
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