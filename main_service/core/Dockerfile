FROM golang:1.22.1 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -trimpath -mod=readonly -o /go/bin/app
ENTRYPOINT [ "/go/bin/app" ]