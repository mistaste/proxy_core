module segment

go 1.26.1

require (
	github.com/GFW-knocker/Xray-core v1.26.5-mahsa-r1
	github.com/Jigsaw-Code/outline-sdk v0.0.20
	github.com/things-go/go-socks5 v0.0.6
	github.com/tidwall/gjson v1.18.0
	github.com/tidwall/sjson v1.2.5
	github.com/xjasonlyu/tun2socks/v2 v2.6.0
	google.golang.org/grpc v1.79.3
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/GFW-knocker/wireguard v1.0.6 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/ajg/form v1.5.1 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/apernet/quic-go v0.59.1-0.20260217092621-db4786c77a22 // indirect
	github.com/cloudflare/circl v1.6.3 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/flynn/noise v1.0.0 // indirect
	github.com/ghodss/yaml v1.0.1-0.20220118164431-d8423dcdf344 // indirect
	github.com/go-chi/chi/v5 v5.2.1 // indirect
	github.com/go-chi/cors v1.2.1 // indirect
	github.com/go-chi/render v1.0.3 // indirect
	github.com/go-gost/relay v0.5.0 // indirect
	github.com/google/btree v1.1.3 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/schema v1.4.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/juju/ratelimit v1.0.2 // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/klauspost/reedsolomon v1.12.0 // indirect
	github.com/mahsanet/dnstt v0.0.0-20260208072530-0ad0213669c6 // indirect
	github.com/miekg/dns v1.1.72 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pires/go-proxyproto v0.11.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	github.com/refraction-networking/utls v1.8.2 // indirect
	github.com/sagernet/sing v0.5.1 // indirect
	github.com/sagernet/sing-shadowsocks v0.2.7 // indirect
	github.com/shadowsocks/go-shadowsocks2 v0.1.5 // indirect
	github.com/templexxx/cpu v0.1.1 // indirect
	github.com/templexxx/xorsimd v0.4.3 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	github.com/vishvananda/netlink v1.3.1 // indirect
	github.com/vishvananda/netns v0.0.5 // indirect
	github.com/xtaci/kcp-go/v5 v5.6.8 // indirect
	github.com/xtaci/smux v1.5.24 // indirect
	github.com/xtls/reality v0.0.0-20260322125925-9234c772ba8f // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	go4.org/netipx v0.0.0-20231129151722-fdeea329fbba // indirect
	golang.org/x/crypto v0.49.0 // indirect
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 // indirect
	golang.org/x/mod v0.33.0 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
	golang.org/x/time v0.12.0 // indirect
	golang.org/x/tools v0.42.0 // indirect
	golang.zx2c4.com/wintun v0.0.0-20230126152724-0fa3db229ce2 // indirect
	golang.zx2c4.com/wireguard v0.0.0-20250521234502-f333402bd9cb // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gvisor.dev/gvisor v0.0.0-20260122175437-89a5d21be8f0 // indirect
	lukechampine.com/blake3 v1.4.1 // indirect
)

replace github.com/xjasonlyu/tun2socks/v2 => ./patches/tun2socks
