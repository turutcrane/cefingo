#

template.qtpl.go: template.qtpl
	go generate *.go

parser/parse_string.go: parser/parse.go
	go generate parser/parse.go

.PHONY: vet
vet: template.qtpl.go *.go parser/parse_string.go
	@# adjust path produced in error meesage
	cd ../.. ; go vet tools/gen-cefingo

capi:
	@# output path is ../../capi
	#cd ../.. ; go run ./tools/gen-cefingo
	go generate github.com/turutcrane/cefingo/capi

.PHONY: fmt
fmt:
	go fmt *.go
	go fmt parser/*.go
