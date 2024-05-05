# Build the application from source
FROM golang:1.22-alpine AS build-stage

WORKDIR /build
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ./hockey-app-backend ./cmd/api 


# Deploy the application binary into a lean image
FROM alpine
WORKDIR /app
COPY --from=build-stage /build/hockey-app-backend ./hockey-app-backend

ENTRYPOINT ["sh", "-c", "/app/hockey-app-backend -addr=:$PORT"]
