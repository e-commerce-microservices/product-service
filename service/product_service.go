package service

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/e-commerce-microservices/product-service/pb"
	"github.com/e-commerce-microservices/product-service/repository"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type productRepository interface {
	CreateProduct(ctx context.Context, arg repository.CreateProductParams) (repository.Product, error)
	GetProductByID(ctx context.Context, id int64) (repository.Product, error)
	GetAllProduct(ctx context.Context) ([]repository.Product, error)
	GetProductByCategory(ctx context.Context, categoryID int64) ([]repository.Product, error)
}
type categoryRepository interface {
	CreateCategory(ctx context.Context, arg repository.CreateCategoryParams) error
	GetAllCategory(ctx context.Context) ([]repository.Category, error)
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

	_, err := service.authClient.AdminAuthorization(ctx, &empty.Empty{})
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	err = service.categoryStore.CreateCategory(ctx, repository.CreateCategoryParams{
		ID:   req.GetCategoryId(),
		Name: req.GetName(),
		Thumbnail: sql.NullString{
			String: req.GetThumbnail(),
			Valid:  true,
		},
	})
	if err != nil {
		return nil, status.Error(codes.Canceled, err.Error())
	}

	return &pb.GeneralResponse{
		Message: "create new category successfull",
	}, nil
}

// GetListCategory ...
func (service *ProductService) GetListCategory(ctx context.Context, _ *empty.Empty) (*pb.GetListCategoryResponse, error) {
	listCategory, err := service.categoryStore.GetAllCategory(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	result := make([]*pb.Category, 0, len(listCategory))
	for _, category := range listCategory {
		result = append(result, &pb.Category{
			CategoryId: category.ID,
			Name:       category.Name,
			Thumbnail:  category.Thumbnail.String,
		})
	}

	return &pb.GetListCategoryResponse{
		ListCategory: result,
	}, nil
}

// CreateProduct ...
func (service *ProductService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	ctx = metadata.NewOutgoingContext(ctx, md)

	claims, err := service.authClient.SupplierAuthorization(ctx, &empty.Empty{})
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	supplierID, err := strconv.ParseInt(claims.GetId(), 10, 64)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	newProduct, err := service.productStore.CreateProduct(ctx, repository.CreateProductParams{
		Name:        req.GetProductName(),
		Description: req.GetDesc(),
		Price:       req.GetPrice(),
		Thumbnail:   req.GetThumbnail(),
		Inventory:   int32(req.GetInventory()),
		SupplierID:  supplierID,
		CategoryID:  req.GetCategoryId(),
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateProductResponse{
		NewProduct: &pb.Product{
			SupplierId: supplierID,
			CategoryId: newProduct.CategoryID,
			Name:       newProduct.Name,
			Desc:       newProduct.Description,
			Price:      0,
			Thumbnail:  newProduct.Thumbnail,
			Inventory:  0,
			CreatedAt:  &timestamppb.Timestamp{},
			UpdatedAt:  &timestamppb.Timestamp{},
		},
	}, nil
}

// GetProduct ...
func (service *ProductService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	product, err := service.productStore.GetProductByID(ctx, req.GetProductId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Product{
		SupplierId: product.SupplierID,
		CategoryId: product.CategoryID,
		Name:       product.Name,
		Desc:       product.Description,
		Price:      0,
		Thumbnail:  product.Thumbnail,
		Inventory:  product.Inventory,
		CreatedAt:  &timestamppb.Timestamp{},
		UpdatedAt:  &timestamppb.Timestamp{},
	}, nil
}

// GetListProduct ...
func (service *ProductService) GetListProduct(ctx context.Context, req *pb.GetListProductRequest) (*pb.GetListProductResponse, error) {
	var listProduct []repository.Product
	var err error
	if req.GetCategoryId() != 0 {
		listProduct, err = service.productStore.GetProductByCategory(ctx, req.GetCategoryId())
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	} else {
		listProduct, err = service.productStore.GetAllProduct(ctx)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

	}

	result := make([]*pb.Product, 0, len(listProduct))
	for _, product := range listProduct {
		result = append(result, &pb.Product{
			SupplierId: product.SupplierID,
			CategoryId: product.CategoryID,
			Name:       product.Name,
			Desc:       product.Description,
			Price:      0,
			Thumbnail:  product.Thumbnail,
			Inventory:  product.Inventory,
		})
	}

	return &pb.GetListProductResponse{
		ListProduct: result,
	}, nil
}

// Ping pong
func (service *ProductService) Ping(context.Context, *empty.Empty) (*pb.Pong, error) {
	return &pb.Pong{
		Message: "pong",
	}, nil
}
