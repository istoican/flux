GO := go

install:
	$(GO) install github.com/istoican/flux

deps:
	cd vendor && $(NPM) install

image: 
	docker build . -t "istoican/flux"

tests:
	docker-compose up --build ${TEST}

cmd/fluxd/fluxd:
	go build -o $@ ./$(@D)

cmd/flux/flux:
	go build -o $@ ./$(@D) 

clean:
	rm -f cmd/flux/flux
	rm -f cmd/fluxd/fluxd

stats:
	@echo "Number of printed pages: $(shell find ./ -type f \( -iname \*.go -o -iname \*.css -o -iname \*.js -o -iname \*.html \) -print0 | xargs -0 cat | wc -l) / 40"
	@echo "Go lines: \t\t$(shell find ./ -name '*.go' -print0 | xargs -0 cat | wc -l)"
	@echo "Javascript lines: \t$(shell find ./ -name '*.js' -print0 | xargs -0 cat | wc -l)"
	@echo "CSS lines: \t\t$(shell find ./ -name '*.css' -print0 | xargs -0 cat | wc -l)"

.PHONY: install deps stats
