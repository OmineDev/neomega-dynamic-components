

build:
	cd /workspace/component1 && CGO_CFLAGS= CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=/usr/bin/gcc  CXX=/usr/bin/gcc  /usr/local/go/bin/go build -buildmode=plugin -buildvcs=false -trimpath -ldflags "-s -w" -o /workspace/buildenv/component1.so .