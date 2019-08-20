FROM golang
WORKDIR /app
ADD go.mod go.sum /app/
RUN go mod download
ADD . /app/
RUN go build -o main cmd/spinner.go
ENTRYPOINT ["./main"]
