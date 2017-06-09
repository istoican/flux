GO := go
NPM := npm 
TSC := vendor/node_modules/.bin/tsc
ELM := vendor/node_modules/.bin/elm

install:
	$(GO) install github.com/istoican/flux

deps:
	cd vendor && $(NPM) install

cmd/flux/flux.go.js:
	$(TSC) -p cmd/flux --outFile $@

.PHONY: install deps