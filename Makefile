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

# Parar todo
docker-turnoff:
	@echo "🛑 Parando toda la infraestructura..."
	sudo docker-compose down
