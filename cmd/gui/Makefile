.PHONY: build

build:
	export PATH="$PATH:$(go env GOPATH)/bin"
	curl https://cdn.jsdelivr.net/npm/bulma@0.9.4/css/bulma.min.css --output frontend/src/assets/bulma.min.css
	wails build -platform "windows/amd64"
	wails build -platform "darwin/arm64"
	wails build -platform "darwin/amd64"