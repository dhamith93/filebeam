
clean: 
	rm -rf build

build:
	mkdir -p build
	GOOS=windows GOARCH=amd64 go build -o build/main_win64.exe
	GOOS=linux GOARCH=amd64 go build -o build/main_linux64
	GOOS=darwin GOARCH=amd64 go build -o build/main_macos_amd64
	GOOS=darwin GOARCH=arm64 go build -o build/main_macos_arm64