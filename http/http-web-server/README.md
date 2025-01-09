# http-web-server

```sh
go run .

# api usage
curl -X POST -H "Content-Type: application/json" \
  -d '{"name":"John Doe"}' \
  http://localhost:8080/users

curl http://localhost:8080/users

curl http://localhost:8080/users/user_id

curl -X PUT -H "Content-Type: application/json" \
  -d '{"name":"John Updated"}' \
  http://localhost:8080/users/{id}


curl -X PUT -H "Content-Type: application/json" \
  -d '{"name":"John Updated"}' \
  http://localhost:8080/users/f5b8ba2a-7b66-43f0-8d6f-539e4398e74a
  
http delete http://localhost:8080/users/user_id
```