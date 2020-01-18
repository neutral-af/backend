.PHONY: dev deps mod gen test deploy

test:
	go test ./... -cover

# This command will start the pseudo-devserver on 8080.
# Note that the deployed version runs in Zeit and does not
# use an http server directly, so the following is a simulation
dev:
	go run cmd/dev/main.go

mod:
	go get -v all
	go mod tidy

gen:
	gqlgen

deps: mod gen

deploy:
	now
