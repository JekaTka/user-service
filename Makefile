DB_URL=postgresql://root:secret@localhost:5432/user_service?sslmode=disable

network:
	docker network create user-service-network

postgres:
	docker run --name postgres --network user-service-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root user_service

dropdb:
	docker exec -it postgres dropdb user_service

migrateup:
	migrate -path pkg/db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path pkg/db/migration -database "$(DB_URL)" -verbose down

new_migration:
	migrate create -ext sql -dir pkg/db/migration -seq $(name)

sqlc:
	sqlc generate

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=user-service \
	proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl

.PHONY: network postgres createdb dropdb migrateup migratedown new_migration sqlc proto evans