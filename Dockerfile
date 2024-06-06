# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS build

WORKDIR /app

# for Makefile
RUN apk add --no-cache build-base

COPY . .
RUN go build -mod vendor -o ./.kanthor/kanthorq -buildvcs=false cmd/kanthorq/main.go
RUN go build -mod vendor -o ./.kanthor/kanthorqdata -buildvcs=false cmd/data/main.go

FROM alpine:3
WORKDIR /app

COPY --from=build /app/data ./data
COPY --from=build /app/migration ./migration
COPY --from=build /app/.kanthor/kanthor /usr/bin/kanthor
COPY --from=build /app/.kanthor/kanthordata /usr/bin/kanthordata