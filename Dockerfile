FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/server server.go
COPY ./upx /upx
RUN chmod +x /upx
RUN /upx --best --lzma /app/server


FROM scratch


WORKDIR /app
COPY --from=builder /app/server /app/server
COPY static /app/static
EXPOSE 8080

CMD ["./server"]
