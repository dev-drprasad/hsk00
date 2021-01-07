.PHONY: generate

generate:
	pkger -include /assets -o cmd/
	pkger -include /assets -o gui/


build: generate
	@go build -o ./bin/hsk00 ./*.go
