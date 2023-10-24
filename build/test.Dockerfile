FROM golang:1.21
WORKDIR /app
COPY . .
RUN go test -race $(go list ./... | grep -v /vendor/)