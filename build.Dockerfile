FROM golang:1.18.3 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /go/src/app
COPY . .

RUN go build

FROM gcr.io/distroless/static-debian11

COPY --from=builder /go/src/app/pages-gc /pages-gc

ENTRYPOINT ["/pages-gc"]