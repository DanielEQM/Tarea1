#compilar proto
protoC:
	cd lester/proto && protoc --go_out=. --go-grpc_out=. ./lester.proto
	cd michael/proto && protoc --go_out=. --go-grpc_out=. ./michael.proto
	cd franklin/proto && protoc --go_out=. --go-grpc_out=. ./franklin.proto
	cd trevor/proto && protoc --go_out=. --go-grpc_out=. ./trevor.proto

# Dockerizar lester
docker-lester: 
	sudo docker-compose up --build lester

docker-michael:
	sudo docker-compose up --build michael

docker-franklin: 
	sudo docker-compose up --build franklin

docker-trevor: 
	sudo docker-compose up --build trevor

lester:
	go build lester.go
	go run lester

trevor:
	go build trevor.go
	go run trevor

franklin:
	go build franklin.go
	go run franklin

michael:
	go build michael.go
	go run michael


# Parar todo
docker-turnoff:
	@echo "🛑 Parando toda la infraestructura..."
	sudo docker-compose down
