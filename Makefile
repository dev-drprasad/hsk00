.PHONY: generate

generate:
	pkger -include /assets


build: generate
	@go build -o ./bin/hsk00 ./*.go
