### Assumption
This project follows the **Hexagonal Architecture** pattern:  
**Domain Layer**: business logic for authentication, user information  
**Service Layer**: use cases for authentication, user information  
**Adapter Layer**:
- HTTP (Gorilla Mux)
- gRPC
- MongoDB (Persistence)

**Unit test**: cover just for some method in **service, adapter** mocks by `mockgen`

---

### Setup Application
`docker compose up --build`

### Mongo DB
- create database `go-auth-user`
- initialize data in counters collection with  
`{ "_id": "users", "seq": 0 }`

### Run Test and Coverage
`go test ./... -coverprofile=coverage.out`  
`go tool cover -html=coverage.out -o coverage.html`

### JWT token usage guide
 - retrieve access token from `localhost:8080/authenticate`
 - put access token into header with `Authorization: Bearer <JWT_TOKEN>`

### gen-go-grpc Example
`protoc --go_out=. --go-grpc_out=. proto/user.proto`

## 📘 API Endpoints

| Method | Endpoint      | Description           | Auth Required | Request Body                                                                                       | Response Example                                                                                     |
|--------|---------------|-----------------------|----------------|----------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------|
| POST   | `/signup`     | Register a new user   | ❌ No          | `{ "name": "Alice", "email": "alice@example.com", "password": "123456" }`                          | `201 Created`                                                                                        |
| POST   | `/authenticate`      | Get Access Token            | ❌ No          | `{ "email": "alice@example.com", "password": "123456" }`                                           | `{ "token": "eyJhbGciOi..." }`                                                                       |
| GET    | `/user`   | Get current user info | ✅ Yes         | _Header_: `Authorization: Bearer <JWT_TOKEN>`  _Body_: `{id:<id>}`                                          | `{code: 201, message: "success", data:{ "id": "1", "name": "one", "email": "one@example.com" }}`                           |
| GET    | `/users`   | Get users list | ✅ Yes         | _Header_: `Authorization: Bearer <JWT_TOKEN>`                                                      | `{code: 201, message: "success"}, data: [{ "id": "1", "name": "one", "email": "one@example.com" },{ "id": "2", "name": "two", "email": "two@example.com" }]}`                           |
| POST    | `/users/create`   | Create a new user | ✅ Yes         | _Header_: `Authorization: Bearer <JWT_TOKEN>`   _Body_: `{ name:<name>, email:<email>, password:<password>}`                                                    | `{code: 201, message: "success"}`                           |
| PUT    | `/users/update`   | Update user's name or email | ✅ Yes         | _Header_: `Authorization: Bearer <JWT_TOKEN>  `    _Body_: `{id:<id>, name:<name>, email:<email>}`                                                    | `{code: 200, message: "success"}`                           |
| DELETE    | `/users/delete`   | Create a new user | ✅ Yes         | _Header_: `Authorization: Bearer <JWT_TOKEN>`  _Body_: `{id:<id>}`                                                      | `{code: 200, message: "success"}`                           |

---

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

