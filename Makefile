
build:
	
	CC=arm-linux-gnueabihf-gcc GOARCH=arm CGO_ENABLED=1 go build -ldflags="-w -s" -a -o drawj2d-go

	