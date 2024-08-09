FROM golang:1.22-alpine AS test-builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY ./ ./
COPY config/config.yaml /config.yaml
# RUN go test -v ./... --config=/config.yaml
RUN go test -v ./...
RUN go build -o ./bin/main cmd/main.go

FROM alpine AS runner
COPY --from=test-builder /usr/local/src /
COPY --from=test-builder config/config.yaml /config.yaml
# COPY --from=test-builder tests/migrations /migrations
# COPY ["configs/apiserver/config.yaml","images", "security", "migrations", "./"]

CMD ["/bin/main", "--config=/config.yaml"]
# ENTRYPOINT [ "/bin/sh"]
