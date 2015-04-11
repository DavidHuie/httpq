build: dep
	go build -o bin/httpq github.com/DavidHuie/httpq/cmd

test: dep
	go test ./...

dep:
	godep save -r ./...
