## ghcr.io/prolapser/debian-awg:latest

Базовый debian-образ с настроенным прокси через Cloudflare WARP, подменой DNS и поднятым прокси для запуска любых приложений. 

```dockerfile
FROM ghcr.io/prolapser/debian-awg:latest

# например можно установить Python
COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/
ENV PATH="/app/.venv/bin:$PATH"

# установить какие-то зависимости
RUN apt-get update && apt-get install -y --no-install-recommends \
    ffmpeg privoxi wget && \
    rm -rf /var/lib/apt/lists/* && \
    uv venv --python 3.14 && uv pip install pip && ln -s /app/.venv/bin/* /usr/local/bin/ && \
    pip install --no-cache-dir httpx[h2] fastapi kurigram

# скопировать файлы телеграм-бота
COPY . /app

# настроить переменные окружения если нужно
# ENV AWG_ENDPOINT="188.114.98.219:4500"
# ENV AWG_ENDPOINT="engage.cloudflareclient.com:2408
# ENV AWG_PRIVATE_KEY="92U7B78AsTysmmIbzoCJ4XqzGpEyGo8RKjU6os7hdLE="
# ENV SPEEDTEST true

# запустить
CMD ["python", "/app/main.py"]
```

После запуска контейнера прокси-сервер SOCKS5 доступен по адресу: `socks5h://127.0.0.1:25344` или `socks5://127.0.0.1:25344`, он устанавливается системно через переменные среды, но если ваше приложение их не уважает, можно установить вручную.

Этот базовый образ занимает примерно 150 Мб.

---

## ghcr.io/prolapser/hf-docker/cdp-awg:latest

Удаленный браузер [CloakBrowser](https://github.com/CloakHQ/CloakBrowser) управляемый через [CDP](https://playwright.dev/python/docs/api/class-browsertype#browser-type-connect-over-cdp) (Connecting over the Chrome DevTools Protocol) с глобальным проксированием через CloudFlare WARP с обфускацией AmneziaWG:

```text
ip        : 104.28.196.75
provider  : Cloudflare Inc.
location  : United States Of America (US), Washington
```

Достаточно [создать спейс](https://huggingface.co/new-space?sdk=docker) и добавить файл `Dockerfile` минимального содержания:

```dockerfile
FROM ghcr.io/prolapser/hf-docker/cdp-awg:latest
```

Но, можно настроить используя переменные среды:

- `AWG_PRIVATE_KEY="***"`: PrivateKey из своего конфига WG/AmneziaWG CloudFlare WARP (если не установлено подставится случайный).
- `AWG_ENDPOINT="host:port"`: эндпоинт CloudFlare WARP (если не установлено подставится случайный).
- `SPEEDTEST true/false`: выполнять ли при старте проверку скорости соединения и пинг через AmneziaWG (отключено по дефолту).
- `CLOAKBROWSER_AUTO_UPDATE true/false`: включать ли фоновую проверку обновлений браузера (включено по дефолту).
- `CLOAKBROWSER_SKIP_CHECKSUM true/false`: пропускать ли проверку SHA-256 после загрузки браузера (включено по дефолту).
- `CLOAKBROWSER_GEOIP_TIMEOUT_SECONDS 5`: таймаут в секундах для разрешения GeoIP, прежде чем продолжить без него (5 по дефолту)

Пример докерфайла с настройками:

```dockerfile
FROM ghcr.io/prolapser/hf-docker/cdp-awg:latest

# ENV AWG_ENDPOINT="188.114.98.219:4500"
ENV AWG_ENDPOINT="engage.cloudflareclient.com:2408
ENV AWG_PRIVATE_KEY="92U7B78AsTysmmIbzoCJ4XqzGpEyGo8RKjU6os7hdLE="
ENV SPEEDTEST true
ENV CLOAKBROWSER_AUTO_UPDATE false
ENV CLOAKBROWSER_SKIP_CHECKSUM false
ENV CLOAKBROWSER_GEOIP_TIMEOUT_SECONDS 10
```

Этот образ занимает примерно 2 Гб.

Пример использования в синхронном коде:

```python
from urllib.parse import quote, urlencode

from playwright._impl._api_structures import SetCookieParam  # noqa
from playwright.sync_api import sync_playwright

# публичная ссылка на страницу HF-спейса (не его репо)
CDP_URL = "https://username-spacename.hf.space"

# куки для снятия семейных фильтров с результатов и установка локации Washington
cookies: list[SetCookieParam] = [
    {'name': 'ys', 'value': 'wprid.1779106375637420-12880543316618692126-balancer-l7leveler-kubr-yp-klg-290-BAL', 'domain': '.yandex.com', 'path': '/', 'httpOnly': False, 'secure': True, 'sameSite': 'None', 'expires': 2094715208.043561},
    {'name': 'yp', 'value': '1779970352.dlp.2#2094466377.pcs.1#1810642356.sp.shst%3A1%3Ashsh%3A1%3Afamily%3A0#1779711159.szm.1_25%3A2048x1152%3A2033x1031%3A15#1779279175.ygo.10493%3A87#1781698375.ygu.0', 'domain': '.yandex.com', 'path': '/', 'httpOnly': False, 'secure': True, 'sameSite': 'None', 'expires': 2094715208.088647},
    {'name': 'yandex_gid', 'value': '87', 'domain': '.yandex.com', 'path': '/', 'httpOnly': False, 'secure': True, 'sameSite': 'None', 'expires': 2094715208.433255}
]


def test(query: str, seed: str):
    params = urlencode(dict(fingerprint=seed, geoip='true'), safe=':/@-_')
    endpoint = f'{CDP_URL.rstrip("/")}?{params}'
    with sync_playwright() as p:
        browser = p.chromium.connect_over_cdp(endpoint)
        context = browser.contexts[0]
        context.set_default_timeout(90000)
        context.add_cookies(cookies)
        page = context.new_page()

        page.goto('chrome://extensions/', wait_until='domcontentloaded', timeout=60000)
        page.screenshot(path='extensions.jpeg', full_page=False, type='jpeg', quality=50)
        print('скриншот с расширениями сохранен как extensions.jpeg')
        # &lr=84 - США или &lr=87 - Вашингтон
        page.goto(
            f'https://yandex.com/search?text={quote(query.replace(" ", "+"), safe="+")}&lr=84',
            wait_until='domcontentloaded',
            timeout=90000
        )
        # page.wait_for_load_state("networkidle")

        # скрытие оверлеев/попапов с рекламой браузера и подписок
        page.add_style_tag(content='''
                /* подписка */
                .plus-link,
                .plus-link_inactive,
                .plus-link__content,
                .plus-link__icon,
                .plus-link__text,
                /* реклама браузера на весь экран */
                .Distribution,
                .DistributionPopup,
                .DistributionInfo,
                [id^="DistributionPopupDesktopSystemNarrow"],
                /* скрытие видео и картинок по теме */
                [data-fast-name="images"],
                [data-fast-name="video-unisearch"]{
                    display: none !important;
                    width: 0px !important;
                    height: 0px !important;
                    position: absolute !important;
                    left: -999999px !important;
                    z-index: -999999 !important;
                }
                ''')
        try:
            # скролл вниз и клик для просмотра настроек
            footer_link = page.wait_for_selector('.SerpFooter-LinksGroup_type_settings', timeout=20000)
            footer_link.scroll_into_view_if_needed()
            footer_link.click(force=True)
        except Exception as e:
            if 'showcaptcha' in page.url:
                print('яндекс показал капчу, лучше использовать куки с реального аккаунта')
            print(e)
        print(f'итоговый url: "{page.url}"')
        page.screenshot(path='screen.jpeg', full_page=True, type='jpeg', quality=50)
        print('скриншот страницы сохранен как screen.jpeg')

        with open('page.html', 'w+', encoding='utf-8') as f:
            f.write(page.content())
        print(f'страница "{page.title()}" сохранена как page.html')

        browser.close()


if __name__ == "__main__":
    test('bufo bufo care', 'yandex_search')
```

Помимо [Playwright-Python](https://github.com/microsoft/playwright-python) по CDP можно подключаться и использовать любые други библиотеки автоматизации и управления браузеров через DevTool, в том числе других языках программирования:

- **Go**: [go-rod](https://github.com/go-rod/rod), [chromedp](https://github.com/chromedp/chromedp), [playwright-go](https://github.com/playwright-community/playwright-go)
- **Node.js**: [Puppeteer](https://github.com/puppeteer/puppeteer), [Playwright](https://github.com/microsoft/playwright)
- и другие...

Чтобы узнать WSS ссылку на DevTools, посмотрите ее у себя: `https://username-spacename.hf.space/json/list`

---

## ghcr.io/prolapser/hf-docker/cdp:latest

Удаленный браузер [CloakBrowser](https://github.com/CloakHQ/CloakBrowser) управляемый через [CDP](https://playwright.dev/python/docs/api/class-browsertype#browser-type-connect-over-cdp) (Connecting over the Chrome DevTools Protocol) без прокси, будучи запущенным на HuggingFace Space использует его IP напрямую:

```text
ip        : 3.228.31.30
provider  : Amazon.com Inc.
location  : United States Of America (US), Ashburn
```

Однако, для обхода заблокированных сайтов со стороны HF, в образе настроена подмена DNS записей на CloudFlare.

Достаточно [создать спейс](https://huggingface.co/new-space?sdk=docker) и добавить файл `Dockerfile` минимального содержания:

```dockerfile
FROM ghcr.io/prolapser/hf-docker/cdp:latest

# nothing
```

Но, можно настроить используя переменные среды:

- `SPEEDTEST true/false`: выполнять ли при старте проверку скорости соединения и пинг через AmneziaWG (отключено по дефолту).
- `CLOAKBROWSER_AUTO_UPDATE true/false`: включать ли фоновую проверку обновлений браузера (включено по дефолту).
- `CLOAKBROWSER_SKIP_CHECKSUM true/false`: пропускать ли проверку SHA-256 после загрузки браузера (включено по дефолту).
- `CLOAKBROWSER_GEOIP_TIMEOUT_SECONDS 5`: таймаут в секундах для разрешения GeoIP, прежде чем продолжить без него (5 по дефолту)

Пример докерфайла с настройками:

```dockerfile
FROM ghcr.io/prolapser/hf-docker/cdp:latest

ENV SPEEDTEST true
ENV CLOAKBROWSER_AUTO_UPDATE false
ENV CLOAKBROWSER_SKIP_CHECKSUM false
ENV CLOAKBROWSER_GEOIP_TIMEOUT_SECONDS 10
```

Этот образ занимает примерно 2 Гб.

Этот образ предпочтительнее образа с AmneziaWG, потому что пересылка WSS сообщений между браузером и клиентом будет быстрее, значительно меньше пинг. Вариант с проксирвоанием через WARP нужен только в некоторых случаях, когда на целевых сайтах сложная анти-бот проверка, например как поиске Яндекс.