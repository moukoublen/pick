module github.com/moukoublen/pick

go 1.25

// Test dependencies. They will not be pushed downstream as indirect ones.
require (
	github.com/google/go-cmp v0.7.0
	github.com/ifnotnil/x/tst v0.0.2
	github.com/stretchr/testify v1.12.1
)

require (
	github.com/stretchr/objx v0.8.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
)
