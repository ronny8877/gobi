FROM golang:1.27.3-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /gobi .

FROM alpine:3.23
RUN apk add --no-cache ca-certificates && adduser -D -u 10001 gobi
COPY --from=build /gobi /usr/local/bin/gobi
COPY LICENSE /usr/share/licenses/gobi/LICENSE
USER gobi
WORKDIR /app
EXPOSE 8080
ENTRYPOINT ["gobi"]
CMD ["-config", "/app/api.gobi.json"]
