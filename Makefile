.PHONY: dev deps test deploy

test:
	go test ./... -cover

# This command will start the pseudo-devserver on 8080.
# Note that the deployed version runs in Zeit and does not
# use an http server directly, so the following is a simulation
dev:
	go run cmd/dev/main.go

# `make deps` will fetch go deps, clean the modfile, and
# rebuild the generated classes based on the GraphQL schema
deps:
	go get -v
	go mod tidy
	gqlgen

deploy:
	now
