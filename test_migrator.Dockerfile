FROM golang:1.22-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY ./ ./
# RUN go build -o ./bin/app cmd/main.go
RUN go build -o ./bin/migrator/main cmd/migrator/main.go
# RUN go run ./bin/migrator/main.exe
# RUN go run ./cmd/migrator/main.go --config=./config/config.yaml --migrations-path=migrations

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/migrator /
COPY config/config.yaml /config.yaml
# COPY migrations /migrations
COPY tests/migrations /migrations
# COPY ["configs/apiserver/config.yaml","images", "security", "migrations", "./"]

CMD ["/main", "--config=./config.yaml", "--migrations-path=migrations"]
