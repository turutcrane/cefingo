#

template.qtpl.go: template.qtpl
	go generate *.go

parser/parse_string.go: parser/parse.go
	go generate parser/parse.go

.PHONY: vet
vet: template.qtpl.go *.go parser/parse_string.go
	@# エラーメッセージのパスを VSCode が解釈しやすいように cd する。
	cd .. ; go vet gen-cefingo/*.go

.PHONY: fmt
fmt:
	go fmt *.go
	go fmt parser/*.go