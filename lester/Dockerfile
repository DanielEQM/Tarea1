FROM golang:1.24.2

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o lester ./lester.go 

# Comando por defecto
CMD ["/app/lester"]
