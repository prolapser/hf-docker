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
    <title>HF Seedbox</title>
    <style>
        body { margin: 0; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; display: flex; flex-direction: column; height: 100vh; background-color: #1e1e1e; color: white; }
        .tabs { display: flex; background: #2d2d2d; box-shadow: 0 2px 5px rgba(0,0,0,0.5); z-index: 10; }
        .tab { flex: 1; padding: 15px; text-align: center; cursor: pointer; transition: background 0.2s; border-bottom: 3px solid transparent; }
        .tab:hover { background: #3d3d3d; }
        .tab.active { background: #3d3d3d; border-bottom: 3px solid #007acc; font-weight: bold; }
        iframe { flex: 1; border: none; width: 100%; height: 100%; background: #fff; }
    </style>
</head>
<body>
    <div class="tabs">
        <div class="tab active" onclick="switchTab('/fb/', this)">Filebrowser</div>
        <div class="tab" onclick="switchTab('qb-ui', this)">qBittorrent</div>
    </div>
    <iframe id="frame" src="/qb-ui"></iframe>

    <script>
        function switchTab(url, el) {
            document.getElementById('frame').src = url;
            document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
            el.classList.add('active');
        }
    </script>
</body>
</html>
`

func writeQBConfig() {
	configDir := "/app/profile/qBittorrent/config"
	os.MkdirAll(configDir, 0755)

	configData := `[LegalNotice]
Accepted=true

[Preferences]
WebUI\Port=8082
WebUI\LocalHostAuth=false
WebUI\AuthSubnetWhitelistEnabled=true
WebUI\AuthSubnetWhitelist=127.0.0.1/32
Connection\ProxyType=SOCKS5
Connection\ProxyIP=127.0.0.1
Connection\ProxyPort=25344
Connection\ProxyPeerConnections=true
Connection\ProxyTrackerConnections=true
Connection\ProxyHostNameLookup=true
Connection\ProxyForce=true
Connection\ResolvePeerCountries=false
Bittorrent\uTP_rate_limited=false
`
	err := os.WriteFile(filepath.Join(configDir, "qBittorrent.conf"), []byte(configData), 0644)
	if err != nil {
		log.Fatalf("Ошибка записи конфига qB: %v", err)
	}
}

func main() {
	// 1. Подготовка директорий и конфигов
	os.MkdirAll("/app/downloads", 0755)
	writeQBConfig()

	// 2. Запуск Filebrowser (на порту 8081)
	fbCmd := exec.Command("./fb",
		"-a", "127.0.0.1",
		"-p", "8081",
		"-r", "/app/downloads",
		"--noauth",
		"-b", "/fb",
		"-d", "/app/fb.db",
	)
	fbCmd.Stdout = os.Stdout
	fbCmd.Stderr = os.Stderr
	if err := fbCmd.Start(); err != nil {
		log.Fatalf("Не удалось запустить Filebrowser: %v", err)
	}
	fmt.Println("Filebrowser запущен на 127.0.0.1:8081")

	// 3. Запуск qBittorrent (на порту 8082)
	qbCmd := exec.Command("./qb",
		"--webui-port=8082",
		"--profile=/app/profile",
		"--save-path=/app/downloads",
	)
	qbCmd.Stdout = os.Stdout
	qbCmd.Stderr = os.Stderr
	if err := qbCmd.Start(); err != nil {
		log.Fatalf("Не удалось запустить qBittorrent: %v", err)
	}
	fmt.Println("qBittorrent запущен на 127.0.0.1:8082")

	// 4. Настройка Reverse Proxy
	fbURL, _ := url.Parse("http://127.0.0.1:8081")
	qbURL, _ := url.Parse("http://127.0.0.1:8082")

	fbProxy := httputil.NewSingleHostReverseProxy(fbURL)
	qbProxy := httputil.NewSingleHostReverseProxy(qbURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Главная страница с вкладками
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(indexHTML))
			return
		}

		// Виртуальный путь для загрузки WebUI qBittorrent в iframe
		if r.URL.Path == "/qb-ui" {
			r.URL.Path = "/"
			qbProxy.ServeHTTP(w, r)
			return
		}

		// Запросы к Filebrowser
		if strings.HasPrefix(r.URL.Path, "/fb") {
			fbProxy.ServeHTTP(w, r)
			return
		}

		// Все остальные запросы (API, CSS, JS от qBittorrent) проксируются в qB
		qbProxy.ServeHTTP(w, r)
	})

	// 5. Запуск веб-сервера на порту 7860
	fmt.Println("Proxy сервер запущен на 0.0.0.0:7860")
	if err := http.ListenAndServe("0.0.0.0:7860", nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
