FROM golang:1.14 AS builder
ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.io
WORKDIR /go/src/app
ADD . .
RUN go build -o /auto-run-all

FROM alpine:3.12
COPY --from=builder /auto-run-all /auto-run-all
CMD ["/auto-run-all"]