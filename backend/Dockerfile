FROM golang:1-alpine3.24 AS build
WORKDIR /src

COPY go.mod ./
COPY . .
RUN go build -o /out/cvmc ./cmd/api

FROM alpine:3.24
WORKDIR /app
RUN adduser -D -g '' appuser
COPY --from=build /out/cvmc /app/cvmc
RUN mkdir -p /app/uploads && chown -R appuser:appuser /app
USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/cvmc"]
