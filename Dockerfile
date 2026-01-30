FROM golang:1.24.3-alpine AS builder

WORKDIR /marketplace-parser-service

RUN apk --no-cache add git bash make gcc gettext musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

ENV CONFIG_PATH=configs/config.yaml
ENV CGO_ENABLED=0

RUN go build --ldflags="-w -s" -o marketplace-parser-service ./cmd/marketplace-parser-service

FROM alpine AS runner
RUN apk add --no-cache ca-certificates

WORKDIR /marketplace-parser-service

COPY --from=builder /marketplace-parser-service/configs/ /marketplace-parser-service/configs/
COPY --from=builder /marketplace-parser-service/marketplace-parser-service /marketplace-parser-service/marketplace-parser-service

EXPOSE 8080

CMD ["./marketplace-parser-service"]