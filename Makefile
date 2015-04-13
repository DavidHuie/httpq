build: dep
	go build -o bin/httpq github.com/DavidHuie/httpq/cmd/httpq

test: dep
	go test ./...

dep:
	godep save -r ./...

install:
	go install github.com/DavidHuie/httpq/cmd/httpq
