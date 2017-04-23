SOURCES := $(shell find . -type f -name '*.go' -not -path './vendor/*')

dev:
	@goimports -l -w ${SOURCES}
	@go vet . ./lib/...
	@go test . ./lib/...
	go install

test:
	go test . ./lib/...

build:
	go fmt . ./lib/...
	go vet . ./lib/...
	staticcheck . ./lib/...
	gosimple . ./lib/...
	go test . ./lib/...
	go install -race .
	cyberspace

cpu:
	go tool pprof /Users/stojg/Code/golang/bin/cyberspace cpu_profile.out

mem:
	go tool pprof /Users/stojg/Code/golang/bin/cyberspace mem_profile.out

alloc:
	go tool pprof -alloc_objects /Users/stojg/Code/golang/bin/cyberspace mem_profile.out

js:
	./node_modules/.bin/webpack -d -w --progress

compile:
	./node_modules/.bin/webpack -p --progress
	go install .
