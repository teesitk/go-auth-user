### gen-go-grpc
protoc --go_out=. --go-grpc_out=. proto/user.proto

### Setup Application
// docker build images
// docker compose up -d

### Sample Request
`curl --location --request GET 'localhost:8080/users/create' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "one",
    "email": "one@test.com",
    "password": "1234"
}'`

### Sample Response
`{"Code":200,"Message":"success"}`
