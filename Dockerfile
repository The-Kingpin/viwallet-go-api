FROM golang:alpine as builder
WORKDIR /app
ADD . /app/
RUN go mod tidy
RUN go mod verify
RUN CGO_ENABLED=0 go build -o viwallet ./cmd/web/*.go

# Remove the line below when fix next stage build
# CMD ["./viwallet" ]


# Fix below stage doesn't return index.html 
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/viwallet .
CMD ["./viwallet" ]
