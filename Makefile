.PHONY: build
build:
	go build -o bin/hugo-ai-translator cmd/hugo-ai-translator/main.go

.PHONY: translate-readme
translate-readme:
	./bin/hugo-ai-translator simple --target-languages en --target-languages ja --target-languages cn --target-languages es --target-languages fr --target-languages de --source-language ko --model gpt-4o-mini
