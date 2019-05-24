FROM golang:alpine AS build

WORKDIR /app
COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor

FROM scratch

WORKDIR /app
COPY --from=build /app/app .

ENTRYPOINT ["/app/app"]
