FROM python:3.11 as b2
RUN pip install --user requests requests[socks]
FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY server.go .
# RUN go build -ldflags="-s -w" -o /app/server server.go
RUN go build -o /app/server server.go
# COPY ./upx /upx
# RUN chmod +x /upx
# RUN /upx --best --lzma /app/server

FROM dind-py

WORKDIR /app
ENV PATH=/root/.local:$PATH
COPY --from=builder /app/server /app/server
COPY --from=b2 /root/.local /root/.local
COPY static /app/static
COPY scripts /app/scripts
EXPOSE 8080

CMD ["./server"]
