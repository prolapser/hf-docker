FROM golang:trixie AS builder
WORKDIR /app
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o app main.go

FROM ghcr.io/prolapser/debian-awg:latest

RUN useradd -m -u 1000 user
USER user
WORKDIR /home/user
COPY --from=builder /app/app /home/user/app

RUN apt-get update && apt-get install -y wget tar && rm -rf /var/lib/apt/lists/* && \
    wget https://github.com/filebrowser/filebrowser/releases/download/v2.63.2/linux-amd64-filebrowser.tar.gz && \
    tar -xvf linux-amd64-filebrowser.tar.gz && \
    rm linux-amd64-filebrowser.tar.gz && \
    wget https://github.com/userdocs/qbittorrent-nox-static/releases/download/release-5.2.0_v2.0.12/x86_64-qbittorrent-nox -O qb && \
    chmod +x qb && \
    mv /home/user/filebrowser /home/user/fb && \
    chmod +x /home/user/app /home/user/fb

EXPOSE 8081
EXPOSE 8082

CMD ["./app"]