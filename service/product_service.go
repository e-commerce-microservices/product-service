package service

import (
	"context"

	"github.com/e-commerce-microservices/product-service/pb"
	"github.com/e-commerce-microservices/product-service/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type productRepository interface{}
type categoryRepository interface {
	CreateCategory(ctx context.Context, arg repository.CreateCategoryParams) error
}

// ProductService implement grpc Server
type ProductService struct {
	authClient    pb.AuthServiceClient
	categoryStore categoryRepository
	productStore  productRepository

	pb.UnimplementedProductServiceServer
}

// NewProductService creates a new ProductService
func NewProductService(authClient pb.AuthServiceClient, queries *repository.Queries) *ProductService {
	service := &ProductService{
		authClient:    authClient,
		categoryStore: queries,
		productStore:  queries,
	}

	return service
}

// CreateCategory creates a new Product Category
func (service *ProductService) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.GeneralResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	ctx = metadata.NewOutgoingContext(ctx, md)

	userClaimsResp, err := service.authClient.GetUserClaims(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	if userClaimsResp.GetUserRole() != pb.UserRole_admin {
		return nil, status.Error(codes.PermissionDenied, "permission denied to create category")
	}

	err = service.categoryStore.CreateCategory(ctx, repository.CreateCategoryParams{
		ID:   req.GetCategoryId(),
		Name: req.GetName(),
	})
	if err != nil {
		return nil, status.Error(codes.Canceled, err.Error())
	}

	return &pb.GeneralResponse{
		Message: "create new category successfull",
	}, nil
}

// CreateProduct ...
func (service *ProductService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	return nil, nil
}
