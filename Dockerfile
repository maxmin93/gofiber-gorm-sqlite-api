# https://awstip.com/containerize-go-sqlite-with-docker-6d7fbecd14f0
FROM golang:1.19 AS builder

WORKDIR /src
COPY . .
RUN go mod download
# need CGO_ENABLED=1 for go-sqlite3
RUN CGO_ENABLED=1 GOOS=linux go build -o /app -a -ldflags '-linkmode external -extldflags "-static"' .

FROM scratch
COPY --from=builder /app /app
COPY --from=builder /src/views /views
EXPOSE 3000

ENTRYPOINT ["/app"]


# docker build -t gofiber-sqlite --no-cache .
# docker run -lt --name gofiber-sqlite -p 3000:3000 -e DB_URL="file:memdb2?mode=memory" -e PORT=3000 gofiber-sqlite

# docker rm -f $(docker ps -l -q)
