package service

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/e-commerce-microservices/product-service/pb"
	"github.com/e-commerce-microservices/product-service/repository"
	"github.com/golang/protobuf/ptypes/empty"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProductService implement grpc Server
type ProductService struct {
	productStore *repository.Queries

	reviewClient pb.ReviewServiceClient
	imageClient  pb.ImageServiceClient
	orderClient  pb.OrderServiceClient
	db           *sql.DB

	pb.UnimplementedProductServiceServer
}

// NewProductService creates a new ProductService
func NewProductService(imageClient pb.ImageServiceClient, reviewClient pb.ReviewServiceClient, orderClient pb.OrderServiceClient, queries *repository.Queries, db *sql.DB) *ProductService {
	service := &ProductService{
		orderClient:  orderClient,
		productStore: queries,
		reviewClient: reviewClient,
		imageClient:  imageClient,
		db:           db,
	}

	return service
}

// DeleteProduct ...
func (service *ProductService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	err := service.productStore.DeleteProduct(ctx, repository.DeleteProductParams{
		ID:         req.GetProductId(),
		SupplierID: req.GetSupplierId(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.DeleteProductResponse{
		Message: "Xóa sản phẩm thành công",
	}, nil
}

// DeleteProductByAdmin ...
func (service *ProductService) DeleteProductByAdmin(ctx context.Context, req *pb.DeleteProductByAdminRequest) (*pb.DeleteProductByAdminResponse, error) {
	err := service.productStore.DeleteProductByID(ctx, req.GetProductId())
	if err != nil {
		return nil, err
	}

	return &pb.DeleteProductByAdminResponse{
		Message: "Xóa sản phẩm thành công",
	}, nil
}

// DescInventory ...
func (service *ProductService) DescInventory(ctx context.Context, req *pb.DescInventoryRequest) (*pb.DescInventoryResponse, error) {
	// md, _ := metadata.FromIncomingContext(ctx)
	// ctx = metadata.NewOutgoingContext(ctx, md)

	// ctx, span := otel.Tracer("").Start(ctx, "ProductService.UpdateInventory")
	// defer span.End()
	err := service.productStore.DescInventory(ctx, repository.DescInventoryParams{
		Inventory: req.GetCount(),
		ID:        req.GetProductId(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.DescInventoryResponse{
		Message: "OK",
	}, nil
}

// IncInventory ...
func (service *ProductService) IncInventory(ctx context.Context, req *pb.IncInventoryRequest) (*pb.IncInventoryResponse, error) {
	err := service.productStore.IncInventory(ctx, repository.IncInventoryParams{
		Inventory: req.GetCount(),
		ID:        req.GetProductId(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.IncInventoryResponse{
		Message: "OK",
	}, nil
}

// CreateCategory creates a new Product Category
func (service *ProductService) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.GeneralResponse, error) {
	err := service.productStore.CreateCategory(ctx, repository.CreateCategoryParams{
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
	listCategory, err := service.productStore.GetAllCategory(ctx)
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
	if len(req.GetProductName()) == 0 {
		return nil, errors.New("Vui lòng điền thông tin tên sản phẩm")
	}
	if req.GetPrice() <= 0 || req.GetInventory() <= 0 {
		return nil, errors.New("Vui lòng điền giá và số lượng sản phẩm")
	}

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
		Brand: sql.NullString{
			String: req.GetBrand(),
			Valid:  true,
		},
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

	reviewResponse, err := service.reviewClient.GetAllReviewByProductID(ctx, &pb.GetAllReviewByProductIDRequest{
		ProductId: product.ID,
	})
	if err != nil {
		return &pb.Product{
			ProductId:  product.ID,
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
	var totalStar int32
	lenTotalStar := 0

	for _, review := range reviewResponse.ListReview {
		if review.NumStar > 0 {
			totalStar += review.NumStar
			lenTotalStar++
		}
	}
	var starAvg float32
	if lenTotalStar > 0 {
		starAvg = float32(totalStar / int32(lenTotalStar))
	}
	var totalSold int32
	soldProductResponse, err := service.orderClient.GetSoldProduct(ctx, &pb.GetSoldProductRequest{
		ProductId: product.ID,
	})
	if err == nil {
		totalSold = int32(soldProductResponse.Count)
	}

	return &pb.Product{
		ProductId:   product.ID,
		SupplierId:  product.SupplierID,
		CategoryId:  product.CategoryID,
		Name:        product.Name,
		Desc:        product.Description,
		Price:       product.Price,
		Thumbnail:   product.Thumbnail,
		Inventory:   product.Inventory,
		CreatedAt:   timestamppb.New(product.CreatedAt),
		Brand:       product.Brand.String,
		TotalSold:   int64(totalSold),
		StarAverage: starAvg,
	}, nil
}

// GetListProduct ...
func (service *ProductService) GetListProduct(ctx context.Context, req *pb.GetListProductRequest) (*pb.GetListProductResponse, error) {
	var listProduct []repository.Product
	var err error
	log.Println("inc", req.ByPriceInc)
	if req.GetCategoryId() != 0 {
		tmp, _ := ProductFilter{
			ByTime:      req.GetByTime(),
			ByPriceInc:  req.ByPriceInc,
			ByPriceDesc: req.ByPriceDesc,
		}.GenerateListProduct(service.db, req.GetCategoryId(), req.GetLimit(), req.GetOffset())
		log.Println(tmp)

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
		reviewResponse, err := service.reviewClient.GetAllReviewByProductID(ctx, &pb.GetAllReviewByProductIDRequest{
			ProductId: product.ID,
		})
		if err != nil {
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
		} else {
			var totalStar int32
			lenTotalStar := 0

			for _, review := range reviewResponse.ListReview {
				if review.NumStar > 0 {
					totalStar += review.NumStar
					lenTotalStar++
				}
			}
			var starAvg float32
			if lenTotalStar > 0 {
				starAvg = float32(totalStar / int32(lenTotalStar))
			}
			var totalSold int64
			soldProductResponse, err := service.orderClient.GetSoldProduct(ctx, &pb.GetSoldProductRequest{
				ProductId: product.ID,
			})
			if err == nil {
				totalSold = int64(soldProductResponse.Count)
			}

			result = append(result, &pb.Product{
				ProductId:   product.ID,
				SupplierId:  product.SupplierID,
				CategoryId:  product.CategoryID,
				Name:        product.Name,
				Desc:        product.Description,
				Price:       product.Price,
				Thumbnail:   product.Thumbnail,
				Inventory:   product.Inventory,
				CreatedAt:   timestamppb.New(product.CreatedAt),
				StarAverage: starAvg,
				TotalSold:   totalSold,
			})
		}
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
		reviewResponse, err := service.reviewClient.GetAllReviewByProductID(ctx, &pb.GetAllReviewByProductIDRequest{
			ProductId: product.ID,
		})
		if err != nil {
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
		} else {
			var totalStar int32
			lenTotalStar := 0

			for _, review := range reviewResponse.ListReview {
				if review.NumStar > 0 {
					totalStar += review.NumStar
					lenTotalStar++
				}
			}
			var starAvg float32
			if lenTotalStar > 0 {
				starAvg = float32(totalStar / int32(lenTotalStar))
			}
			var totalSold int64
			soldProductResponse, err := service.orderClient.GetSoldProduct(ctx, &pb.GetSoldProductRequest{
				ProductId: product.ID,
			})
			if err == nil {
				totalSold = int64(soldProductResponse.Count)
			}

			listProduct = append(listProduct, &pb.Product{
				ProductId:   product.ID,
				SupplierId:  product.SupplierID,
				CategoryId:  product.CategoryID,
				Name:        product.Name,
				Desc:        product.Description,
				Price:       product.Price,
				Thumbnail:   product.Thumbnail,
				Inventory:   product.Inventory,
				CreatedAt:   timestamppb.New(product.CreatedAt),
				StarAverage: starAvg,
				TotalSold:   totalSold,
			})
		}
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

	result := make([]*pb.Product, 0, len(tmp))
	for _, product := range tmp {
		reviewResponse, err := service.reviewClient.GetAllReviewByProductID(ctx, &pb.GetAllReviewByProductIDRequest{
			ProductId: product.ID,
		})
		if err != nil {
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
		} else {
			var totalStar int32
			lenTotalStar := 0

			for _, review := range reviewResponse.ListReview {
				if review.NumStar > 0 {
					totalStar += review.NumStar
					lenTotalStar++
				}
			}
			var starAvg float32
			if lenTotalStar > 0 {
				starAvg = float32(totalStar / int32(lenTotalStar))
			}
			var totalSold int64
			soldProductResponse, err := service.orderClient.GetSoldProduct(ctx, &pb.GetSoldProductRequest{
				ProductId: product.ID,
			})
			if err == nil {
				totalSold = int64(soldProductResponse.Count)
			}

			result = append(result, &pb.Product{
				ProductId:   product.ID,
				SupplierId:  product.SupplierID,
				CategoryId:  product.CategoryID,
				Name:        product.Name,
				Desc:        product.Description,
				Price:       product.Price,
				Thumbnail:   product.Thumbnail,
				Inventory:   product.Inventory,
				CreatedAt:   timestamppb.New(product.CreatedAt),
				StarAverage: starAvg,
				TotalSold:   totalSold,
			})
		}
	}

	return &pb.GetListProductResponse{
		ListProduct: result,
	}, nil
}

// UpdateProduct ...
func (service *ProductService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.GeneralResponse, error) {
	log.Println("update product: ", req)
	err := service.productStore.UpdateProduct(ctx, repository.UpdateProductParams{
		ID:        req.GetProductId(),
		Name:      req.GetName(),
		Price:     req.GetPrice(),
		Inventory: int32(req.GetInventory()),
		Brand: sql.NullString{
			String: req.GetBrand(),
			Valid:  false,
		},
		SupplierID: req.GetSupplierId(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.GeneralResponse{
		Message: "Update product success",
	}, nil
}

// GetListProductByIDs ...
func (service *ProductService) GetListProductByIDs(ctx context.Context, req *pb.GetListProductByIDsRequest) (*pb.GetListProductResponse, error) {
	ids := []string{}
	for _, v := range req.GetListId() {
		ids = append(ids, strconv.FormatInt(v, 10))
	}
	if len(ids) == 0 {
		return &pb.GetListProductResponse{
			ListProduct: []*pb.Product{},
		}, nil
	}
	query := `	SELECT id, name, description, price, thumbnail, inventory, supplier_id, category_id, created_at, brand FROM product WHERE id IN ($1)`
	query = strings.ReplaceAll(query, "$1", strings.Join(ids, ","))
	listProduct, err := service.productStore.GetListProductByIDs(ctx, query)
	if err != nil {
		return nil, err
	}

	result := make([]*pb.Product, 0, len(listProduct))
	for _, product := range listProduct {
		reviewResponse, err := service.reviewClient.GetAllReviewByProductID(ctx, &pb.GetAllReviewByProductIDRequest{
			ProductId: product.ID,
		})
		if err != nil {
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
		} else {
			var totalStar int32
			lenTotalStar := 0

			for _, review := range reviewResponse.ListReview {
				if review.NumStar > 0 {
					totalStar += review.NumStar
					lenTotalStar++
				}
			}
			var starAvg float32
			if lenTotalStar > 0 {
				starAvg = float32(totalStar / int32(lenTotalStar))
			}
			var totalSold int64
			soldProductResponse, err := service.orderClient.GetSoldProduct(ctx, &pb.GetSoldProductRequest{
				ProductId: product.ID,
			})
			if err == nil {
				totalSold = int64(soldProductResponse.Count)
			}

			result = append(result, &pb.Product{
				ProductId:   product.ID,
				SupplierId:  product.SupplierID,
				CategoryId:  product.CategoryID,
				Name:        product.Name,
				Desc:        product.Description,
				Price:       product.Price,
				Thumbnail:   product.Thumbnail,
				Inventory:   product.Inventory,
				CreatedAt:   timestamppb.New(product.CreatedAt),
				StarAverage: starAvg,
				TotalSold:   totalSold,
			})
		}
	}

	return &pb.GetListProductResponse{
		ListProduct: result,
	}, nil
}

var tracer = otel.Tracer("auth-service")

// GetListProductInventory ...
func (service *ProductService) GetListProductInventory(ctx context.Context, req *pb.GetInventoryRequest) (*pb.GetInventoryResponse, error) {
	ctx, span := tracer.Start(ctx, "checkProductAvailable")
	defer span.End()
	// md, _ := metadata.FromIncomingContext(ctx)
	// ctx = metadata.NewOutgoingContext(ctx, md)

	// ctx, span := otel.Tracer("").Start(ctx, "ProductService.CheckInventory")
	// defer span.End()
	resp, err := service.productStore.GetProductInventory(ctx, req.GetProductId())
	if err != nil {
		return nil, err
	}
	return &pb.GetInventoryResponse{
		Count: int64(resp),
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
	if len(tmp) <= 1 {
		return "", errors.New("Vui lòng thêm ảnh")
	}
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
