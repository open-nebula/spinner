FROM golang
WORKDIR /app
ADD . /app/
RUN go build -o main/main main/main.go
ENTRYPOINT ["./main/main"]
