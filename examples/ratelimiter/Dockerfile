FROM golang:1.21.4-bookworm AS build
WORKDIR /app
COPY src .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main admission/main.go

# Etapa de producción
FROM scratch
WORKDIR /app
COPY --from=build /app/main /app

CMD ["./main"]