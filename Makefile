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

stats:
	@echo "Number of printed pages: $(shell find ./ -type f \( -iname \*.go -o -iname \*.css -o -iname \*.js -o -iname \*.html \) -print0 | xargs -0 cat | wc -l) / 40"
	@echo "Go lines: \t\t$(shell find ./ -name '*.go' -print0 | xargs -0 cat | wc -l)"
	@echo "Javascript lines: \t$(shell find ./ -name '*.js' -print0 | xargs -0 cat | wc -l)"
	@echo "CSS lines: \t\t$(shell find ./ -name '*.css' -print0 | xargs -0 cat | wc -l)"

pb/rpc.pb.go: pb/rpc.proto
	protoc -I $(@D) $< --go_out=plugins=grpc:$(@D)

.PHONY: install deps stats