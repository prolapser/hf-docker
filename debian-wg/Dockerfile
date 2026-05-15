FROM debian:stable-slim

WORKDIR /app

RUN apt-get update && apt-get install -y curl git ca-certificates && apt-get clean && \
    curl -L -sS https://github.com/windtf/wireproxy/releases/latest/download/wireproxy_linux_amd64.tar.gz -o /app/wg.tar.gz && \
    tar -xzf /app/wg.tar.gz -C /app && rm /app/wg.tar.gz && mv /app/wireproxy /app/wg && chmod +x /app/wg

COPY <<-"EOF" /app/resolv.conf
nameserver 1.1.1.1
nameserver 1.0.0.1
nameserver 2606:4700:4700::1111
nameserver 2606:4700:4700::1002
EOF

COPY <<-"EOF" /app/wg.conf
[Interface]
PrivateKey = TzzFmtrZ0+dg8KsDLc0+dqYaKMmnQl8WTtfYsvkCJqY=
Address = 172.16.0.2, 2606:4700:110:898e:8ea3:1dab:96fa:5017
DNS = 1.1.1.1, 1.0.0.1, 2606:4700:4700::1111, 2606:4700:4700::1001
MTU = 1280

[Peer]
PublicKey = bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo=
AllowedIPs = 0.0.0.0/0, ::/0
Endpoint = 162.159.192.1:4500
PersistentKeepalive = 25

[Socks5]
BindAddress = 127.0.0.1:25344
EOF

COPY --chmod=755 <<-"EOF" /entrypoint_with_wg.sh
#!/bin/bash
cat /app/resolv.conf > /etc/resolv.conf

if [ -n "$WG_PRIVATE_KEY" ]; then
    echo "configuring WARP PrivateKey..."
    sed -i "s|^PrivateKey *=.*|PrivateKey = ${WG_PRIVATE_KEY}|" /app/wg.conf
else
    pkeys=("IJXnx/deIgotuBg77cPXCyvn1tGYbhoa//44OgQ1+k8=" "sFbBFwWUD6MetJdfcQw8+/ZAPBkAfUKzjQpe2Rb9GUg=" "4FGb75mlc+bCJeXu9grSQxHJqxpwsMIlFEmDsHnY/0s=" "AFOpbAJhSL6U3MvSN6wMT+Q6ead12XjECOhWyaMI4lo=" "oCR2j5m3+YmjcPVh/cAQ06G0RMgUxw5vZvjRNtGx2F4=" "AH607vYf04C0Ou2NtC6EAESmSH9wSKGXLGsOsjgpSlQ=" "uAPpQEAsuQMjo9wCpad2rD6hUSD0TKuQxod/ozNoGHQ=" "QIb3RykWu9IKWDEJrXTBsNUVrmGgAajhZsf3dQDz7FY=" "iIswh8yrqhbyOttKJi6RNqZaRKRzx/BwaAVI5khZMGg=" "sBZdw7rc9rbWMpfSDymYP+Rx7UM8Lr0wLZVCezezS0A=" "mFpSIuipLz/ceU0/k8l/2jSv71C1guGJHy4nsARBFmU=" "oOxmY/pvp6NWXOvYlrq/4mR+fLlqYZ5fdvD7p3m05kc=" "wC8bXPnC9MZjR36nMyjgMThLH+SWlWzZzNWyVsJfzUM=" "qFh/xFemLSlmuSuyN5UiXFwwJ1ymFZoVwu+i1PQXYVE=" "cD9bTUr32n/JMTAEAo5GVrIbsV1cbZgMrAw+MPTfBEk=" "ECmok+DPkbSj9p8nNmD3pj3gNgI7CQaETpZ/RzAD3W8=" "2Fwd8wCTAa2pVSoDRQ+zV9H4nlP0cri3iMpjuj8nYlk=" "mE0swNavnmWWUxbSJlaVnb26XHfky58AawBLCmKMPVo=" "aK81U0tYkYYASf3pAVgcfZqFnWiG/BjYogvWK2xqmV4=" "WANuNl7YK+mkwpvkhhshzWqbywuTn8fDY/FECjQ/uFs=" "iESnyq9aUVZE1CQpta2mrHzw1BpgqsrGsNDhaEbsQn8=" "2C71kERRO1+YxOCH8De7eyF08a5MKjszSTVUWWyuxG4=" "8LmkegWokGnY4289GjtiaRpnMCAahuvrWWR3mTIIE2k=" "kHsi+w5qaEqWmvDHmyEA26O4VN+GHrsKdF/MKFbwbVQ=" "gPAy+JuxJnn71oNXLDS8N6sLJjXkhqIPu4WjUOVWKEw=" "MK2O/HMokdK78G9kOM5qNmUR21aAnb+tp4dYij4EWF0=" "GFR69X9haeMn4E7Q9CTCOPoojBHh+6Vl+Qs7KSfzLlc=" "CL41I/+HVsC38ngNSoKNqGbGqEZSI/4w4dK6iMPSuHs=" "QJvBxB/vyRpRECr3URrXeGdUFD4/BLG4/zmE/AzPZ2k=" "oHu/d035vsRlDzV6Pmmi4cAmjBLzws0OOT22yICjlE0=" "cIV+70HU+fhuqfKWWMK0QC9/2JLnB49oUEOeBwLqCFY=" "GKrzWrPF1nOercQUaybrHsxniekGWaOuSWJqv4aGtX0=" "mHhKBbi8Tuo96h0l0/MMb89xp40wgroh5w8OxW++01g=" "QDurQVKdPrRojBKGZJnP/3qlY8VOf8ZsGxyn/eXJaXk=" "OG50V6aYQBiyube9nLU1oAbEAXFPBPQZ9Ep5VxI7CEs=" "YPoQW9Wl5yZGU3d39JhjlyWFb7xqC11BQ3qzagcj8HU=" "8MMmkU6CDsp3FJhVX3ALly7rC6vbpmmRqRS5QP9rqHI=" "eGqw6AT5Z8hc6DFblFo+ZmKcrVoTQJgPuYetQzUN334=" "eBBRe6uvOMh+/ZBLFEb4aZgz36GV9JMwfr8eRErxWUE=" "UCFsRBRTmseNGVQTozJWTCaG4QS9KwjucXnkIe03ZEQ=" "kJyjT+abgHRKA09Bk/fATK9bYMFKLjkBHCToIvZVPGI=" "iEvkXIxodSsS+ry5MkOn71mx9UILlqbw2dECpmZvTl8=" "IBgWIPvG7WIZpd2LnQx0iWOCUrgPB6WSPmpWpMLgDGc=" "AAj+/LKrCSEm1Bt6kqBLk1qHjTR5vVcNl1HoPz09el4=" "yGM6sjjTbL7/Khy9ixsGy8yBV+tQ21Kv3EFBiFgfD3M=" "aI4TRDvEgvnWHY5PbUeoaxXvuIHKONzrXYJWtUepCEc=" "YNNMiMs/8qpmGJKdCsbFvh4skAjCzg/XhCjRAWZTPVM=" "KCl57Nrr0us5nsBrit1Djb6tKl5oYDRHBsTFYYkv5k4=" "eCaAp4ovKaMad+3wE5WiBOOAuzL8JzTK2CyvfNyk31E=" "WCZGMU+lWIWRe77mE4FWfbu7UdxxvbS0XL9wCTfpY0s=" "OFro2p65Az194NKFCa+jKH3r0x4IIL5AkLH5nYubwkk=" "gNMTOdcGXdLJWXNF4xqPCZVMRYkE82tLWc0j4GlZvng=" "qLnDP1YcIpDAhKEAtcY7jg2vtcyku2C1i8eSCgxBKEo=" "qMd/ROuY8wM53rehlgj/adNZikbYPw/Y4Weu1qA9dHM=" "CKYYF5H00tR6q9Qj+ajTotpgPzWUtbAMQtS4HR+MWmQ=" "0BbsOBPBP2Pc5rU680628k5V0DPCDGpa1J4QsdhFHUo=" "uPxlbc0kcIdeDcp2PFPNZ6FT0/R5ngancBLVI9HCrFs=" "QByeK/BZQQIEdrJ3vytJjsqS/tqxsYRS0+j2uPKiUng=" "KIpSB1fsqo0JB81gemeCHdpK1xSa1toChqbiN+C5aXs=" "KBux6g99S5fsGNPjMCnSW1m/OZPCQ+Km/j/7Lsc7nHw=" "aNjuFeCUaTFLU3f21d7aRcY1S3Cs9U5Ju5MissitaF4=" "IPthXL1ux5qz3jjaA3bA1zDKxdBJ5NLUbsT+RgNZcUg=" "8IKs77UAhtV7uJTv/TUd4Z3kEvE/4iyHlRbzEbm9MWs=" "OEY3Zpi4mqyxZ/mPSTnY7DjstlbKDw/qWcN/C6Bnck8=" "CMziXUklL8PTtx7emkJD4FG+N0NS+2HSvzfAp+3bvkA=" "eCtWV9SsKxoYirlegemB1iKlLb596aL3UJVU06HsNFc=" "2JdMyTuVMSYnF7ufotNs50xGsjK/2ipWwOnFR45+/mY=" "6NjTOCwaRFShmWW9IO2Hc4uYMVxP8bSksT5j1ncaN2Q=" "kIFXdpHwLXXFVf6XvA/iFVp5IjkojBtSuKq3Xz24yX8=" "aL0XQ1yC/zHpBpkyFwYX9cK1+xLz+8iFsmHQ1GmSbXk=" "EKnSjnIi6ktdOjOL8s+DRQ0Pta6IsmouM/LGwiTekFE=" "2Pmo7UCHCUuHtkUF2k0lBXNEgGmPPnsMRfCyJD/yxmc=" "aBHl6GQdwqo3oyT1+GjsnXuTKZsc0MN3KSJf4feiaFI=" "aOBglYDOEB/+6yupDncKkX+NMIA7h1/EZ5Pg4usK6lc=" "0K5Q8vfY/11WSk2KirwOpujXhRxz7b39TcBHv7qYMGU=" "KNsYDjzTWo/RWHq+SBQAtn2sukF0TdqasO4k3dH2u24=" "mIYpArnZTb5L6JxdKjuoFBKCSGYm+L2AjYTNPw7bc24=" "KEtXPaflONGTUsUVbKrXWz27e5dUUZSr0AYCqI6bgHU=" "0JkIq3OdZd8m7y2itvqfGMJOXTy6PzvTTJ+XXy8pDVk=" "UPA/g3QfbREU78KmTjqU8qIHo4cYQq3oCB/ogPi3NFU=" "CHBnEbfLe/KAu6fuUlPeQPXfmeg95cb2s4x+zrH/1Wk=" "gILHCOx+s6yYadlFIZO5jIbob2OiL8VsmqRwPzDay18=" "mKRaG3Q3Gf+IhxU4FaZFCifaDPMhgz5duvt9C6wsHW0=" "2JnR6KJ1Cqgqbn8O0eogLbFB6GyVEhIrr+Haj5pQxUQ=" "+NQceVmqLsBOWRUGMriuV5WoeiWoHOq+73ScllJ3bGg=" "IPenWIP22H7NfR54vL7gyYA2l4Nix0Gf/sdX2wATGUE=" "+Mn4wXfXJmPj/MqMeIFYAAJ780xfAKmVB5Qj7GA703M=" "UEXIjbh8nxD6C2UFWdBZ0iJJXtd45xbDf5ypLnuwwng=" "GHUEMo9jeTI/BdQ3//P/f6FRR5IWxqkaPhThtfdEulI=" "CKMTlNu6ESAcsHweUz3rJ7v8EuH+R4KYBlzTxSIcQH4=" "MF8PixShhr2hZn7uib98+Y5Ek9Od1cvz6l8VUmhClkw=" "ICPTeBrHfZy44wTeNhw4YQIFS+GtI/E+ZY7j7HRXPmM=" "UGrno7RURJ1NF42IHzWoQ/H9r9IMBznsNCimUBooJlk=" "4Davde+g/hqvdEmekdEQ0zIe46wklB3XVwErJs+al1I=" "sIGAh1Ity4CCL3roHaS1dwsc+ymo+5UmGjH+c2U6CHM=" "+Lypaxt8jk94ejIjHJIwsBiVuz9XDfv46/piS6qCV1E=" "EN4nqAJOZ1kEeOwyDL5UHlxX2w9rlDn++SPSy/dZcmE=" "SNRE9x5NYLt1NfvIwinCYrGgconvQP8SxF4DZgg7TX8=" "QP8yxJhoTeWj29wsUzOyY010X9iJhWElVrhBbBxD3FE=" "iNRQj7imEvb9QKeELtZ2Ilm0NRUKgZ5L0rxtCkod3Fo=")
    random_pkey=${pkeys[$((RANDOM % ${#pkeys[@]}))]}
    sed -i "s|^PrivateKey *=.*|PrivateKey = ${random_pkey}|" /app/wg.conf
    echo "WG_PRIVATE_KEY is not set, use default random WARP PrivateKey:"
    echo ${random_pkey}
fi

if [ -n "$WG_ENDPOINT" ]; then
    echo "configuring WARP Endpoint..."
    sed -i "s|^Endpoint *=.*|Endpoint = ${WG_ENDPOINT}|" /app/wg.conf
else
    echo "WG_ENDPOINT is not set, use random WARP endpoint:"
    endpoints=("188.114.96.109:3581" "188.114.98.219:4500" "8.34.146.6:5279" "8.34.146.17:4500" "8.34.146.0:4177" "8.34.146.5:1180" "8.34.146.8:1074" "8.34.146.8:1070" "8.34.146.2:5279" "8.34.146.12:1843" "8.34.146.12:2506" "8.34.146.15:4500" "8.34.146.3:4177" "8.34.146.18:4500" "8.34.146.13:3854" "8.34.146.7:3476" "8.34.146.0:2371" "8.34.146.19:4233" "8.34.146.10:3476" "8.34.146.1:7559" "8.34.146.17:7152" "8.34.146.10:3138" "8.34.146.3:3476" "8.34.146.9:3581" "8.34.146.1:3581" "8.34.146.1:5279" "8.34.146.20:1387" "8.34.146.5:1074" "8.34.146.23:2408" "8.34.146.17:2506" "8.34.146.15:4198" "8.34.146.20:1180" "8.34.146.17:4177" "8.34.146.13:7103" "8.34.146.10:7559" "8.34.146.13:7152" "8.34.146.11:2408" "8.34.146.7:2506" "8.34.146.6:8319" "8.34.146.18:5279" "8.34.146.17:4198" "8.34.146.8:2506" "8.34.146.16:7103" "8.34.146.3:1387" "8.34.146.9:2408" "8.34.146.23:2506" "8.34.146.23:8319" "8.34.146.17:8319" "8.34.146.28:1074" "8.34.146.29:3138" "8.34.146.24:1387" "8.34.146.29:3581" "8.34.146.24:7152" "8.34.146.25:1070" "8.34.146.24:8854" "8.34.146.26:3138" "8.34.146.27:7156" "8.34.146.24:5279" "8.34.146.32:1387" "8.34.146.25:1074" "8.34.146.25:3854" "8.34.146.26:8854" "8.34.146.25:7152" "8.34.146.28:3854" "8.34.146.30:8319" "8.34.146.26:3581" "8.34.146.28:1701" "8.34.146.29:3476" "8.34.146.32:4233" "8.34.146.33:7103" "8.34.146.33:7559" "8.34.146.34:2371" "8.34.146.34:3854" "8.34.146.33:8886" "8.34.146.34:8319" "8.34.146.34:7559" "8.34.146.34:7152" "8.34.146.33:2506" "8.34.146.37:4177" "8.34.146.37:7156" "8.34.146.38:4177" "8.34.146.38:7152" "8.34.146.38:8742" "8.34.146.39:1074" "8.34.146.39:3138" "188.114.99.244:3138" "188.114.99.254:1074" "188.114.99.255:8886" "8.34.146.40:8742" "8.34.146.40:4198" "8.34.146.43:1070" "8.34.146.40:8319" "8.34.146.44:1180" "8.34.146.42:2408" "8.34.146.41:5279" "8.34.146.43:1843" "8.34.146.43:3476" "8.34.146.42:7152" "8.34.146.44:4233" "8.34.146.44:5956")
    random_endpoint=${endpoints[$((RANDOM % ${#endpoints[@]}))]}
    echo ${random_endpoint}
    sed -i "s|^Endpoint *=.*|Endpoint = ${random_endpoint}|" /app/wg.conf
    echo 'If you have ping/speed problems, try to set: ENV WG_ENDPOINT="engage.cloudflareclient.com:2408"'
fi

echo "run wg..."
/app/wg --config /app/wg.conf > /dev/null 2>&1 &
sleep 3
echo "test w/o wg:"
curl -sS 2ip.me
echo "test with wg:"
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

ENTRYPOINT ["/entrypoint_with_wg.sh"]
CMD ["/bin/bash"]