FROM golang:trixie AS builder
WORKDIR /app
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o app main.go

FROM ghcr.io/prolapser/debian-awg:latest

COPY --from=builder /app/app /app/app

RUN apt-get update && apt-get install -y wget tar && rm -rf /var/lib/apt/lists/* && \
    wget https://github.com/filebrowser/filebrowser/releases/download/v2.63.2/linux-amd64-filebrowser.tar.gz && \
    tar -xvf linux-amd64-filebrowser.tar.gz && rm linux-amd64-filebrowser.tar.gz && mv /app/filebrowser /app/fb && \
    wget https://github.com/userdocs/qbittorrent-nox-static/releases/download/release-5.2.0_v2.0.12/x86_64-qbittorrent-nox -O qb && \
    chmod +x /app/app /app/fb /app/qb

EXPOSE 8081 8082 7860

CMD ["./app"]