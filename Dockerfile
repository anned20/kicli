FROM golang:1.18-alpine AS golang-builder

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/app

FROM scratch

WORKDIR /bin
COPY --from=golang-builder /go/bin/app kicli

ENTRYPOINT ["kicli"]
