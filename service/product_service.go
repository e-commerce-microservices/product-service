package service

import (
	"context"
	"database/sql"
	"encoding/base64"
	"log"
	"strings"

	"github.com/e-commerce-microservices/product-service/pb"
	"github.com/e-commerce-microservices/product-service/repository"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type productRepository interface {
	CreateProduct(ctx context.Context, arg repository.CreateProductParams) (repository.Product, error)
	GetProductByID(ctx context.Context, id int64) (repository.Product, error)
	GetAllProduct(ctx context.Context) ([]repository.Product, error)
	GetProductByCategory(ctx context.Context, arg repository.GetProductByCategoryParams) ([]repository.Product, error)
	GetRecommendProduct(ctx context.Context, arg repository.GetRecommendProductParams) ([]repository.Product, error)
	GetProductBySupplier(ctx context.Context, arg repository.GetProductBySupplierParams) ([]repository.Product, error)
	UpdateProduct(ctx context.Context, arg repository.UpdateProductParams) error
}
type categoryRepository interface {
	CreateCategory(ctx context.Context, arg repository.CreateCategoryParams) error
	GetAllCategory(ctx context.Context) ([]repository.Category, error)
}

// ProductService implement grpc Server
type ProductService struct {
	categoryStore categoryRepository
	productStore  productRepository

	imageClient pb.ImageServiceClient

	pb.UnimplementedProductServiceServer
}

// NewProductService creates a new ProductService
func NewProductService(imageClient pb.ImageServiceClient, queries *repository.Queries) *ProductService {
	service := &ProductService{
		categoryStore: queries,
		productStore:  queries,
		imageClient:   imageClient,
	}

	return service
}

// CreateCategory creates a new Product Category
func (service *ProductService) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.GeneralResponse, error) {
	err := service.categoryStore.CreateCategory(ctx, repository.CreateCategoryParams{
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
	thumbnail, err := uploadImage(ctx, req.GetThumbnailDataChunk(), service.imageClient)
	if err != nil {
		return nil, err
	}

	_, err = service.productStore.CreateProduct(ctx, repository.CreateProductParams{
		Name:        req.GetProductName(),
		Description: req.GetDesc(),
		Price:       req.GetPrice(),
		Thumbnail:   thumbnail,
		Inventory:   int32(req.GetInventory()),
		SupplierID:  req.GetSupplierId(),
		CategoryID:  req.GetCategoryId(),
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// add product into es

	return &pb.CreateProductResponse{
		Message: "product is created",
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
		Price:      product.Price,
		Thumbnail:  product.Thumbnail,
		Inventory:  product.Inventory,
		CreatedAt:  timestamppb.New(product.CreatedAt),
		Brand:      product.Brand.String,
	}, nil
}

// GetListProduct ...
func (service *ProductService) GetListProduct(ctx context.Context, req *pb.GetListProductRequest) (*pb.GetListProductResponse, error) {
	var listProduct []repository.Product
	var err error
	if req.GetCategoryId() != 0 {
		listProduct, err = service.productStore.GetProductByCategory(ctx, repository.GetProductByCategoryParams{
			CategoryID: req.GetCategoryId(),
			Limit:      req.GetLimit(),
			Offset:     req.GetOffset(),
		})
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
			ProductId:  product.ID,
			SupplierId: product.SupplierID,
			CategoryId: product.CategoryID,
			Name:       product.Name,
			Desc:       product.Description,
			Price:      product.Price,
			Thumbnail:  product.Thumbnail,
			Inventory:  product.Inventory,
			CreatedAt:  timestamppb.New(product.CreatedAt),
		})
	}

	return &pb.GetListProductResponse{
		ListProduct: result,
	}, nil
}

// GetRecomendProduct ...
func (service *ProductService) GetRecomendProduct(ctx context.Context, req *pb.GetRecommendProductRequest) (*pb.GetListProductResponse, error) {
	listProductStore, err := service.productStore.GetRecommendProduct(ctx, repository.GetRecommendProductParams{
		Limit:  req.GetLimit(),
		Offset: req.GetOffset(),
	})
	if err != nil {
		return nil, err
	}
	listProduct := make([]*pb.Product, 0, len(listProductStore))
	for _, product := range listProductStore {
		listProduct = append(listProduct, &pb.Product{
			SupplierId: product.SupplierID,
			CategoryId: product.CategoryID,
			Name:       product.Name,
			Desc:       product.Description,
			Price:      product.Price,
			Thumbnail:  product.Thumbnail,
			Inventory:  product.Inventory,
			CreatedAt:  timestamppb.New(product.CreatedAt),
			ProductId:  product.ID,
			Brand:      product.Brand.String,
		})
	}

	return &pb.GetListProductResponse{
		ListProduct: listProduct,
	}, nil
}

// GetProductBySupplier ...
func (service *ProductService) GetProductBySupplier(ctx context.Context, req *pb.GetProductBySupplierRequest) (*pb.GetListProductResponse, error) {
	tmp, err := service.productStore.GetProductBySupplier(ctx, repository.GetProductBySupplierParams{
		SupplierID: req.GetSupplierId(),
		Limit:      req.GetLimit(),
		Offset:     req.GetOffset(),
	})
	if err != nil {
		return nil, err
	}

	listProduct := make([]*pb.Product, 0, len(tmp))
	for _, product := range tmp {
		listProduct = append(listProduct, &pb.Product{
			SupplierId: product.SupplierID,
			CategoryId: product.CategoryID,
			Name:       product.Name,
			Desc:       product.Description,
			Price:      product.Price,
			Thumbnail:  product.Thumbnail,
			Inventory:  product.Inventory,
			CreatedAt:  timestamppb.New(product.CreatedAt),
			ProductId:  product.ID,
			Brand:      product.Brand.String,
		})
	}

	return &pb.GetListProductResponse{
		ListProduct: listProduct,
	}, nil
}

// UpdateProduct ...
func (service *ProductService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.GeneralResponse, error) {
	err := service.productStore.UpdateProduct(ctx, repository.UpdateProductParams{
		ID:        req.GetProductId(),
		Name:      req.GetName(),
		Price:     req.GetPrice(),
		Inventory: int32(req.GetInventory()),
		Brand: sql.NullString{
			String: req.GetBrand(),
			Valid:  false,
		},
	})
	if err != nil {
		return nil, err
	}

	return &pb.GeneralResponse{
		Message: "Update product success",
	}, nil
}

// Ping pong
func (service *ProductService) Ping(context.Context, *empty.Empty) (*pb.Pong, error) {
	return &pb.Pong{
		Message: "pong",
	}, nil
}

func toBytes(str string) []byte {
	bytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func uploadImage(ctx context.Context, dataChunk string, imageClient pb.ImageServiceClient) (string, error) {
	// upload image
	stream, err := imageClient.UploadImage(ctx)
	if err != nil {
		return "", err
	}
	// send mime type
	tmp := strings.Split(dataChunk, "data:image/")
	mimeType := strings.Split(tmp[1], ";")[0]
	err = stream.Send(&pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				ImageType: mimeType,
			},
		},
	})
	if err != nil {
		return "", err
	}

	// send data
	dataChunk = strings.Split(dataChunk, ",")[1]
	stream.Send(&pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_ChunkData{
			ChunkData: toBytes(dataChunk),
		},
	})

	res, err := stream.CloseAndRecv()
	if err != nil {
		return "", err
	}

	return res.GetImageUrl(), nil
}
