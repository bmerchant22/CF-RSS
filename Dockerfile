FROM golang:1.18-alpine AS builder

WORKDIR /build

ENV GOPROXY https://goproxy.io

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY pkg/ pkg/

RUN CGO_ENABLED=0 GOOS=linux go build -a -o app cmd/web/main.go

RUN npm install --only=prod

RUN npm run build

FROM alpine:3.16

WORKDIR /cfrss
COPY --from=builder /build/app bin/app

ENTRYPOINT [ "bin/app" ]
