
build:
	go fmt . ./lib/...
	go vet . ./lib/...
	go test . ./lib/...
	go install .
	cyberspace

js:
	./node_modules/.bin/webpack -d -w --progress

compile:
	./node_modules/.bin/webpack -p --progress
	go install .
