package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const indexHTML = `
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>коробка семени</title>
    <style>
        body {
            margin: 0;
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background-color: #1a1a1a;
            color: #ffffff;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            height: 100vh;
        }
        h1 { margin-bottom: 40px; color: #007acc; font-weight: 300; letter-spacing: 2px; }
        .container {
            display: flex;
            gap: 30px;
        }
        .card {
            background: #2d2d2d;
            padding: 40px 30px;
            border-radius: 12px;
            text-align: center;
            cursor: pointer;
            text-decoration: none;
            color: white;
            width: 220px;
            border: 1px solid #3d3d3d;
            box-shadow: 0 4px 15px rgba(0,0,0,0.3);
            transition: all 0.3s ease;
        }
        .card:hover {
            background: #363636;
            border-color: #007acc;
            transform: translateY(-5px);
            box-shadow: 0 8px 25px rgba(0,122,204,0.3);
        }
        .card h2 { margin: 0 0 15px 0; font-size: 24px; color: #fff; }
        .card p { margin: 0; color: #a0a0a0; font-size: 14px; }
    </style>
</head>
<body>
    <h1>коробка семени</h1>
    <div class="container">
        <a href="/qb-ui" target="_blank" class="card">
            <h2>qB***</h2>
            <p>Управление загрузками<br>(откроется в новой вкладке)</p>
        </a>
        <a href="/fb/" target="_blank" class="card">
            <h2>Files</h2>
            <p>Управление файлами<br>(откроется в новой вкладке)</p>
        </a>
    </div>
</body>
</html>
`

func writeQBConfig() {
	configDir := "./profile/qBittorrent/config"
	os.MkdirAll(configDir, 0755)

	// Конфиг qBittorrent 5.x
	configData := `[Application]
FileLogger\Enabled=false
FileLogger\MaxSizeBytes=1024
FileLogger\Path=/app/logs

[BitTorrent]
Session\DefaultSavePath=/app/downloads
MergeTrackersEnabled=true
Session\GlobalMaxSeedingMinutes=0
Session\MultiConnectionsPerIp=true
Session\Preallocation=true
Session\ProxyPeerConnections=true
Session\QueueingSystemEnabled=false
Session\ValidateHTTPSTrackerCertificate=false

[LegalNotice]
Accepted=true

[Network]
Proxy\AuthEnabled=false
Proxy\HostnameLookupEnabled=true
Proxy\IP=127.0.0.1
Proxy\Password=
Proxy\Port=@Variant(\0\0\0\x85\x63\0)
Proxy\Profiles\BitTorrent=true
Proxy\Profiles\Misc=true
Proxy\Profiles\RSS=true
Proxy\Type=SOCKS5
Proxy\Username=

[Preferences]
Advanced\IgnoreSSLErrors=true
General\Locale=en
WebUI\Port=8082
WebUI\LocalHostAuth=false
WebUI\AuthSubnetWhitelistEnabled=true
WebUI\AuthSubnetWhitelist=127.0.0.1/32
`
	err := os.WriteFile(filepath.Join(configDir, "qBittorrent.conf"), []byte(configData), 0644)
	if err != nil {
		log.Fatalf("Ошибка записи конфига qB: %v", err)
	}
}

func main() {
	// 1. Подготовка директорий
	os.MkdirAll("/app/downloads", 0755)
	writeQBConfig()

	// 2. Запуск Filebrowser
	fbCmd := exec.Command("./fb",
		"-a", "127.0.0.1",
		"-p", "8081",
		"-r", "/app/downloads",
		"--noauth",
		"-b", "/fb",
		"-d", "./filebrowser.db",
	)
	fbCmd.Stdout = os.Stdout
	fbCmd.Stderr = os.Stderr
	if err := fbCmd.Start(); err != nil {
		log.Fatalf("Не удалось запустить Filebrowser: %v", err)
	}
	fmt.Println("Filebrowser запущен на 127.0.0.1:8081")

	// 3. Запуск qBittorrent
	qbCmd := exec.Command("./qb",
		"--webui-port=8082",
		"--profile=./profile",
		"--save-path=./downloads",
	)
	qbCmd.Stdout = os.Stdout
	qbCmd.Stderr = os.Stderr
	if err := qbCmd.Start(); err != nil {
		log.Fatalf("Не удалось запустить qB***: %v", err)
	}
	fmt.Println("qB*** запущен на 127.0.0.1:8082")

	// 4. Настройка Reverse Proxy
	fbURL, _ := url.Parse("http://127.0.0.1:8081")
	qbURL, _ := url.Parse("http://127.0.0.1:8082")

	fbProxy := httputil.NewSingleHostReverseProxy(fbURL)
	qbProxy := httputil.NewSingleHostReverseProxy(qbURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(indexHTML))
			return
		}
		if r.URL.Path == "/qb-ui" {
			r.URL.Path = "/"
			qbProxy.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/fb") {
			fbProxy.ServeHTTP(w, r)
			return
		}
		qbProxy.ServeHTTP(w, r)
	})

	// 5. Запуск веб-сервера на порту 7860
	fmt.Println("Сервер запущен на 0.0.0.0:7860")
	if err := http.ListenAndServe("0.0.0.0:7860", nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}