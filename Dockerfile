FROM golang:1.17.6 AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux make build && make swagger

FROM alpine:3.15
WORKDIR /app
COPY --from=build /app/httpserver /app/httpserver
COPY --from=build /app/configs/httpserver.yaml /app/configs/httpserver.yaml
COPY --from=build /app/docs /app/docs
CMD ["./httpserver"]
