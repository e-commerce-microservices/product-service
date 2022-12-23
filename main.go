package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/e-commerce-microservices/product-service/pb"
	"github.com/e-commerce-microservices/product-service/repository"
	"github.com/e-commerce-microservices/product-service/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	// postgres driver
	_ "github.com/lib/pq"
)

func main() {
	// create grpc server
	grpcServer := grpc.NewServer()

	// init user db connection
	pgDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWD"), os.Getenv("DB_DBNAME"),
	)

	productDB, err := sql.Open("postgres", pgDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer productDB.Close()
	if err := productDB.Ping(); err != nil {
		log.Fatal("can't ping to user db", err)
	}

	// init queries
	queries := repository.New(productDB)

	// dial auth client
	authServiceConn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("can't dial user service: ", err)
	}
	// create auth client
	authClient := pb.NewAuthServiceClient(authServiceConn)

	// create product service
	productService := service.NewProductService(authClient, queries)
	// register product service
	pb.RegisterProductServiceServer(grpcServer, productService)

	// reflection service
	reflection.Register(grpcServer)

	// listen and serve
	listener, err := net.Listen("tcp", "localhost:8083")
	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}

	log.Printf("start gRPC server on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}
