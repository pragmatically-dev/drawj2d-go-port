
build:
	GOOS=linux GOARCH=arm go build -ldflags="-w -s" -a

