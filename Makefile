
build:
	GOOS=linux CGO=1 CC=arm-linux-gnueabihf-gcc GOARCH=arm go build -ldflags="-w -s" -a

