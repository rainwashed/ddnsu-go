main:
	GOOS=linux GOARCH=amd64 go build -o build/ddnsu.x64 src/*.go
	GOOS=linux GOARCH=arm64 go build -o build/ddnsu.arm64 src/*.go
	GOOS=linux GOARCH=386 go build -o build/ddnsu.x32 src/*.go

