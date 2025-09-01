#compilar proto
protoC:
	cd Lester/proto && protoc --go_out=. --go-grpc_out=. ./lester.proto
	cd Michael/proto && protoc --go_out=. --go-grpc_out=. ./michael.proto

# Dockerizar lester
docker-lester: 
	sudo docker-compose up --build lester

# Dockerizar michael
docker-michael:
	sudo docker-compose up --build michael

# Parar todo
docker-turnoff:
	@echo "🛑 Parando toda la infraestructura..."
	sudo docker-compose down
