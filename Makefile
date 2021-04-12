NAME := vscale-backup
VERSION?=$(git version > /dev/null 2>&1 && git describe --dirty=-dirty --always 2>/dev/null || echo NO_VERSION)
LDFLAGS=-ldflags "-X=main.version=$(VERSION)"

fmt:
	@goimports -local ${NAME} -l -w .
build-linux:
	@CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) $(GOFLAGS) -o ${NAME} $(MAIN)
build-darwin:
	@CGO_ENABLE=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) $(GOFLAGS) -o ${NAME} $(MAIN)
build-windows:
	@CGO_ENABLE=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) $(GOFLAGS) -o ${NAME} $(MAIN)
