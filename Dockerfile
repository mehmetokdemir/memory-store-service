FROM golang:1.17-alpine as development
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o ./bin/memory-store-service .
CMD ./bin/memory-store-service