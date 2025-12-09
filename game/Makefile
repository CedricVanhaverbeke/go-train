mock:
	go run main.go -workout "108-600;216-600;270-900;162-600;270-900;203-300;108-300" -mock true
build:
	GOOS=windows GOARCH=amd64 go build -o ../bin/overlay-windows-amd64 
	GOOS=windows GOARCH=arm go build -o ../bin/overlay-windows-arm 
	GOOS=windows GOARCH=arm64 go build -o ../bin/overlay-windows-arm64 
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o ../bin/overlay-darwin-amd64 
	GOOS=darwin GOARCH=arm64 go build -o ../bin/overlay-darwin-arm64 

build-linux:
	CC=x86_64-linux-gnu-gcc CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ../bin/overlay-linux-amd64

