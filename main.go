package main

import (
	"context"
	"fmt"
	grpcAdapter "go-auth-user/adapter/grpc"
	httpAdapter "go-auth-user/adapter/http"
	mongoAdapter "go-auth-user/adapter/mongo"
	pb "go-auth-user/proto"
	"go-auth-user/service"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	//JWT
	jwtSecret := "secret"
	// connect mongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username := "user"
	password := "pass"
	host := "localhost"
	port := "27017"
	dbName := "go-auth-user"

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin", username, password, host, port, dbName)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	userRepo := mongoAdapter.NewUserRepo(client, "go-auth-user", "users")
	userService := service.NewUserService(userRepo)
	userHandler := httpAdapter.NewUserHandler(userService)
	authService := service.NewAuthService(userRepo, jwtSecret)
	authHandler := httpAdapter.NewAuthHandler(authService)

	// Mux Router
	r := mux.NewRouter()
	r.Use(httpAdapter.LoggingMiddleware)

	r.HandleFunc("/authenticate", authHandler.Authenticate).Methods("POST")
	r.HandleFunc("/signup", authHandler.Signup).Methods("POST")

	// Protected routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(func(next http.Handler) http.Handler {
		return httpAdapter.JWTMiddleware(authService, next)
	})
	api.HandleFunc("/users", userHandler.ListUser)
	api.HandleFunc("/user", userHandler.GetUser)
	api.HandleFunc("/users/create", userHandler.CreateUser).Methods("POST")
	api.HandleFunc("/users/update", userHandler.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/delete", userHandler.DeleteUser).Methods("DELETE")

	//Logging users counter in mongoDB
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			<-ticker.C
			count, err := userRepo.CountUsers()
			if err != nil {
				log.Printf("[ERROR] Counting users: %v", err)
				continue
			}
			log.Printf("[INFO] Total users in DB: %d", count)
		}
	}()

	//Provide sum grpc server
	go func() {
		lis, err := net.Listen("tcp", ":9000")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		userService := grpcAdapter.NewUserServiceServer(userRepo)
		pb.RegisterUserServiceServer(grpcServer, userService)

		log.Println("gRPC server listening on :9000")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	log.Println("server running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
