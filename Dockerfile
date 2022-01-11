FROM golang:1.17.6 AS modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

FROM golang:1.17.6 AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux make build

FROM scratch
WORKDIR /app
COPY --from=builder /app/httpserver /app/httpserver
COPY --from=builder /app/configs/httpserver.yaml /app/configs/httpserver.yaml
COPY --from=builder /app/docs/api.yaml /app/docs/api.yaml
CMD ["./httpserver"]
