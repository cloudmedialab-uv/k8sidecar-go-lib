FROM golang:1.21.4-bookworm AS builder
WORKDIR /app
# COPY ssh key to have access to private repository (//TODO: this gona be delated on final stage)
COPY keys/id_rsa /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts
RUN echo "Host github.com\n\tHostName github.com\n\tUser git\n\tIdentityFile /root/.ssh/id_rsa\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config
RUN git clone git@github.com:cloudmedialab-uv/k8sidecar-go-lib.git
RUN git config --global url."git@github.com:".insteadOf "https://github.com/"
# Modules add
COPY . .
RUN go get github.com/cloudmedialab-uv/k8sidecar-go-lib
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main main.go

FROM scratch
COPY --from=builder /app/main /app
CMD ["./main"]
