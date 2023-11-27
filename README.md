# Pick
[![Go Report Card](https://goreportcard.com/badge/github.com/moukoublen/pick)](https://goreportcard.com/report/github.com/moukoublen/pick)
[![CI Status](https://github.com/moukoublen/pick/actions/workflows/ci.yml/badge.svg)](https://github.com/moukoublen/pick/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/moukoublen/pick/graph/badge.svg?token=6X9MMYZJZ8)](https://codecov.io/gh/moukoublen/pick)


**Pick** is a go package to access (using dot and array notation) and cast any kind of data, using best effort performance and best effort cast. It is an alternative to [stretchr/objx](https://github.com/stretchr/objx) aiming to provide three main things:

1. Modular approach regarding the caster, traverser and selector format
2. Best effort performance using `reflect` as last resort.
3. Best effort cast aiming to cast and convert between types as much as possible.

### Examples
```golang
j := `{
    "item": {
        "one": 1,
        "two": "ok"
        "three": ["element 1", 2, "element 3"]
    }
    "float": 2.12
}`
p, _ := pick.WrapJSON([]byte(j))

returned, found, err := p.String("item.three[1]") // returned: string("2")    err: nil
returned, found, err := p.Uint64("item.three[1]") // returned: uint64("2")    err: nil
returned, found, err := p.Int32("item.one")       // returned: int32(1)       err: nil
returned, found, err := p.Float32("float")        // returned: float32(2.12)  err: nil
returned, found, err := p.Int64("float")          // returned: int64(2)       err: ErrCastLostDecimals
```

**Pick** is still in alpha version, a lot of changes going to happen both to api and structure.


Special thanks to **Konstantinos Pittas** ([@kostaspt](https://github.com/daydroidmuchiri)) for helping me ... **pick** the name of the library.
