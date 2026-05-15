FROM debian:stable-slim
LABEL authors="prolpaser"

WORKDIR /app

RUN apt-get update && apt-get install -y curl git ca-certificates && apt-get clean && \
    curl -L -sS https://github.com/artem-russkikh/wireproxy-awg/releases/latest/download/wireproxy_linux_amd64.tar.gz -o /app/awg.tar.gz && \
    tar -xzf /app/awg.tar.gz -C /app && rm /app/awg.tar.gz && mv /app/wireproxy /app/awg && chmod +x /app/awg

COPY <<-"EOF" /app/resolv.conf
nameserver 1.1.1.1
nameserver 1.0.0.1
nameserver 2606:4700:4700::1111
nameserver 2606:4700:4700::1002
EOF

COPY <<-"EOF" /app/awg.conf
[Interface]
PrivateKey = TzzFmtrZ0+dg8KsDLc0+dqYaKMmnQl8WTtfYsvkCJqY=
Address = 172.16.0.2, 2606:4700:110:898e:8ea3:1dab:96fa:5017
DNS = 1.1.1.1, 1.0.0.1, 2606:4700:4700::1111, 2606:4700:4700::1001
MTU = 1280
S1 = 0
S2 = 0
S3 = 0
S4 = 0
Jc = 4
Jmin = 40
Jmax = 70
H1 = 1
H2 = 2
H3 = 3
H4 = 4
I1 = <b 0x494e56495445207369703a626f624062696c6f78692e636f6d205349502f322e300d0a5669613a205349502f322e302f55445020706333332e61746c616e74612e636f6d3b6272616e63683d7a39684734624b3737366173646864730d0a4d61782d466f7277617264733a2037300d0a546f3a20426f62203c7369703a626f624062696c6f78692e636f6d3e0d0a46726f6d3a20416c696365203c7369703a616c6963654061746c616e74612e636f6d3e3b7461673d313932383330313737340d0a43616c6c2d49443a20613834623463373665363637313040706333332e61746c616e74612e636f6d0d0a435365713a2033313431353920494e564954450d0a436f6e746163743a203c7369703a616c69636540706333332e61746c616e74612e636f6d3e0d0a436f6e74656e742d547970653a206170706c69636174696f6e2f7364700d0a436f6e74656e742d4c656e6774683a20300d0a0d0a>
I2 = <b 0x5349502f322e302031303020547279696e670d0a5669613a205349502f322e302f55445020706333332e61746c616e74612e636f6d3b6272616e63683d7a39684734624b3737366173646864730d0a546f3a20426f62203c7369703a626f624062696c6f78692e636f6d3e0d0a46726f6d3a20416c696365203c7369703a616c6963654061746c616e74612e636f6d3e3b7461673d313932383330313737340d0a43616c6c2d49443a20613834623463373665363637313040706333332e61746c616e74612e636f6d0d0a435365713a2033313431353920494e564954450d0a436f6e74656e742d4c656e6774683a20300d0a0d0a>

[Peer]
PublicKey = bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo=
AllowedIPs = 0.0.0.0/0, ::/0
Endpoint = 162.159.192.1:4500
PersistentKeepalive = 25

[Socks5]
BindAddress = 127.0.0.1:25344
EOF

COPY --chmod=755 <<-"EOF" /entrypoint_with_awg.sh
#!/bin/bash
cat /app/resolv.conf > /etc/resolv.conf

if [ -n "$AWG_PRIVATE_KEY" ]; then
    echo "configuring WARP PrivateKey..."
    sed -i "s|^PrivateKey *=.*|PrivateKey = ${AWG_PRIVATE_KEY}|" /app/awg.conf
else
    echo "AWG_PRIVATE_KEY is not set, use default WARP PrivateKey."
fi

if [ -n "$AWG_ENDPOINT" ]; then
    echo "configuring WARP Endpoint..."
    sed -i "s|^Endpoint *=.*|Endpoint = ${AWG_ENDPOINT}|" /app/awg.conf
else
    echo "AWG_ENDPOINT is not set, use random WARP endpoint:"
    endpoints=("188.114.96.109:3581" "188.114.98.219:4500" "8.34.146.6:5279" "8.34.146.17:4500" "8.34.146.0:4177" "8.34.146.5:1180" "8.34.146.8:1074" "8.34.146.8:1070" "8.34.146.2:5279" "8.34.146.12:1843" "8.34.146.12:2506" "8.34.146.15:4500" "8.34.146.3:4177" "8.34.146.18:4500" "8.34.146.13:3854" "8.34.146.7:3476" "8.34.146.0:2371" "8.34.146.19:4233" "8.34.146.10:3476" "8.34.146.1:7559" "8.34.146.17:7152" "8.34.146.10:3138" "8.34.146.3:3476" "8.34.146.9:3581" "8.34.146.1:3581" "8.34.146.1:5279" "8.34.146.20:1387" "8.34.146.5:1074" "8.34.146.23:2408" "8.34.146.17:2506" "8.34.146.15:4198" "8.34.146.20:1180" "8.34.146.17:4177" "8.34.146.13:7103" "8.34.146.10:7559" "8.34.146.13:7152" "8.34.146.11:2408" "8.34.146.7:2506" "8.34.146.6:8319" "8.34.146.18:5279" "8.34.146.17:4198" "8.34.146.8:2506" "8.34.146.16:7103" "8.34.146.3:1387" "8.34.146.9:2408" "8.34.146.23:2506" "8.34.146.23:8319" "8.34.146.17:8319" "8.34.146.28:1074" "8.34.146.29:3138" "8.34.146.24:1387" "8.34.146.29:3581" "8.34.146.24:7152" "8.34.146.25:1070" "8.34.146.24:8854" "8.34.146.26:3138" "8.34.146.27:7156" "8.34.146.24:5279" "8.34.146.32:1387" "8.34.146.25:1074" "8.34.146.25:3854" "8.34.146.26:8854" "8.34.146.25:7152" "8.34.146.28:3854" "8.34.146.30:8319" "8.34.146.26:3581" "8.34.146.28:1701" "8.34.146.29:3476" "8.34.146.32:4233" "8.34.146.33:7103" "8.34.146.33:7559" "8.34.146.34:2371" "8.34.146.34:3854" "8.34.146.33:8886" "8.34.146.34:8319" "8.34.146.34:7559" "8.34.146.34:7152" "8.34.146.33:2506" "8.34.146.37:4177" "8.34.146.37:7156" "8.34.146.38:4177" "8.34.146.38:7152" "8.34.146.38:8742" "8.34.146.39:1074" "8.34.146.39:3138" "188.114.99.244:3138" "188.114.99.254:1074" "188.114.99.255:8886" "8.34.146.40:8742" "8.34.146.40:4198" "8.34.146.43:1070" "8.34.146.40:8319" "8.34.146.44:1180" "8.34.146.42:2408" "8.34.146.41:5279" "8.34.146.43:1843" "8.34.146.43:3476" "8.34.146.42:7152" "8.34.146.44:4233" "8.34.146.44:5956")
    random_endpoint=${endpoints[$((RANDOM % ${#endpoints[@]}))]}
    echo ${random_endpoint}
    sed -i "s|^Endpoint *=.*|Endpoint = ${random_endpoint}|" /app/awg.conf
fi

echo "run awg..."
/app/awg --config /app/awg.conf > /dev/null 2>&1 &
sleep 3
echo "test w/o awg:"
curl -sS 2ip.me
echo "test with awg:"
curl -sS -x socks5://127.0.0.1:25344 2ip.me
echo "setup proxy..."
export http_proxy="socks5h://127.0.0.1:25344"
export https_proxy="socks5h://127.0.0.1:25344"
export all_proxy="socks5h://127.0.0.1:25344"
export HTTP_PROXY="socks5h://127.0.0.1:25344"
export HTTPS_PROXY="socks5h://127.0.0.1:25344"
export ALL_PROXY="socks5h://127.0.0.1:25344"
echo "test system-wide proxy:"
curl -sS 2ip.me
echo -e "\n\n==================================\nWARP ON: socks5h://127.0.0.1:25344\n==================================\n\n"
exec "$@"
EOF

ENTRYPOINT ["/entrypoint_with_awg.sh"]
CMD ["/bin/bash"]