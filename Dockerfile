# Build Stage
FROM golang:1.26.2-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/voltswitch-api .

# Runtime Stage
FROM alpine:3.22
RUN apk add --no-cache openrc
COPY --from=build /out/voltswitch-api /usr/local/bin/voltswitch-api

WORKDIR /app
EXPOSE 8000
ENTRYPOINT ["voltswitch-api"]
