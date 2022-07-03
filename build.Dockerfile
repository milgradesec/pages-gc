FROM --platform=amd64 golang:1.18.3 AS builder

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH}

WORKDIR /go/src/app
COPY . .

RUN go build

FROM gcr.io/distroless/static-debian11:nonroot

COPY --from=builder /go/src/app/pages-gc /pages-gc

ENTRYPOINT ["/pages-gc"]