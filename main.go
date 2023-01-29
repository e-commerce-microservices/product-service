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

	// dial image client
	imageServiceConn, err := grpc.Dial("image-service:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("can't dial image service: ", err)
	}
	// create image client
	imageClient := pb.NewImageServiceClient(imageServiceConn)

	// create product service
	productService := service.NewProductService(imageClient, queries)
	// register product service
	pb.RegisterProductServiceServer(grpcServer, productService)

	// listen and serve
	listener, err := net.Listen("tcp", ":8080")
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
