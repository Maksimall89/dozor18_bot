FROM golang:alpine3.12 AS builder
COPY . /go/dozor18_bot
WORKDIR /go/dozor18_bot
RUN go mod download 
RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o main 
RUN chmod +x main

FROM alpine:3.12.1
COPY --from=builder /go/dozor18_bot/ .
EXPOSE 80
EXPOSE 443
CMD ["./main"]