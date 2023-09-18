FROM golang:1.20.5-alpine3.18 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN GOOS=linux GOARCH=amd64 go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -mod=readonly -o /bin/ino

FROM scratch
COPY --from=builder /bin/ino /bin/ino
CMD ["/bin/ino", "serve"]