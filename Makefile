DB_DSN := postgres://admin:admin@192.168.49.2:30012/product?sslmode=disable

migratecreate:
	migrate create -ext sql -dir db/migration -seq ${f}

migrateup:
	migrate -path db/migration -database "${DB_DSN}" -verbose up ${v}

migratedown:
	migrate -path db/migration -database "${DB_DSN}" -verbose down ${v}

migrateforce:
	migrate -path db/migration -database "${DB_DSN}" -verbose force ${v}

protogen:
	protoc --proto_path=proto proto/product_service.proto proto/image_service.proto proto/review_service.proto proto/order_service.proto proto/general.proto proto/auth_service.proto \
	--go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative

sqlcgen:
	sqlc generate

.PHONY: rebuild
rebuild:
	docker build -t ngoctd/ecommerce-product:latest . && \
	docker push ngoctd/ecommerce-product

.PHONY: redeploy
redeploy:
	kubectl rollout restart deployment depl-product

.PHONY: migratecreate migrateup migratedown migrateforce protogen_auth protogen_product sqlcgen