FROM golang:1.20

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ENV TZ=Asia/Shanghai PUID=0 PGID=0 UMASK=022 

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/arrproxy

EXPOSE 8117

CMD ["/app/arrproxy"]
