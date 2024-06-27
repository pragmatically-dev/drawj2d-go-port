
build:
	go build -a -gcflags=all="-l -B" -ldflags="-w -s"