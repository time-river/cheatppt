all:
	go build


upgrade:
	go get -t -u ./...
