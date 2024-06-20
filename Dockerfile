# Build the application from source
FROM golang:1.22.4-alpine AS build-stage

WORKDIR /build
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ./server ./src/cmd/server 


# Deploy the application binary into a lean image
FROM alpine
WORKDIR /app
COPY --from=build-stage /build/server ./server

ENTRYPOINT ["sh", "-c", "/app/server -addr=:$PORT -dsn=mongodb://mongo:27017/test" ]
