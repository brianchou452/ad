FROM golang:latest
WORKDIR /ad
COPY . .
RUN go mod download
RUN go build -o main .
EXPOSE 80
CMD ["./main"]