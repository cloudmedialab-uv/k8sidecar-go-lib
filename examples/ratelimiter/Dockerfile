FROM golang:1.21.4-bookworm AS builder
WORKDIR /app

RUN git clone https://github.com/cloudmedialab-uv/k8sidecar-go-lib.git
RUN git config --global url."https://github.com/".insteadOf "git@github.com:"
COPY src .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main main.go

FROM scratch
COPY --from=builder /app/main /app
CMD ["./main"]
