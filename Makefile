ifndef $(tag)
	tag=latest
endif

build: cmd/*.go
	go build -o bin/mtoken ./cmd/main.go

test: pkg/*/*.go
	go test -v github.com/codemk8/mtoken/pkg/...

clean:
	-rm -rf bin/*

vendor:
	go mod vendor