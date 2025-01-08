# http-web-server

```sh
go run .

# api usage
curl -X POST -H "Content-Type: application/json" \
  -d '{"name":"John Doe"}' \
  http://localhost:8080/users

curl http://localhost:8080/users

curl http://localhost:8080/users/user_id
```