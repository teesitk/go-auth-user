### Setup Application
`docker compose up --build`

### Run Test and Coverage
`go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html`

### JWT token usage guide
 - retrieve access token from `localhost:8080/authenticate`
 - put access token into header with `Authorization: Bearer `

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

### gen-go-grpc Example
`protoc --go_out=. --go-grpc_out=. proto/user.proto`

Method | Endpoint | Description | Auth Required | Request Body | Response Example
POST | /signup | Register a new user | ❌ No | { "name": "one", "email": "one@example.com", "password": "123456" } | 201 Created
POST | /login | User login | ❌ No | { "email": "one@example.com", "password": "123456" } | { "token": "eyJhbGciOi..." }
GET | /users/me | Get current user info | ✅ Yes | Header: Authorization: Bearer <JWT_TOKEN> | { "id": "1", "name": "one", "email": "one@example.com" }
