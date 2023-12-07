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
    "two": "ok",
    "three": ["element 1", 2, "element 3"]
},
"float": 2.12
}`
p, _ := WrapJSON([]byte(j))

{
    returned, err := p.String("item.three[1]")
    assert(t, "2", returned, nil, err)
}
{
    returned, err := p.Uint64("item.three[1]")
    assert(t, uint64(2), returned, nil, err)
}
{
    returned, err := p.Int32("item.one")
    assert(t, int32(1), returned, nil, err)
}
{
    returned, err := p.Float32("float")
    assert(t, float32(2.12), returned, nil, err)
}
{
    returned, err := p.Int64("float")
    assert(t, int64(2), returned, cast.ErrCastLostDecimals, err)
}
```

**`Map` function**
```golang
j2 := `{
    "items": [
        {"id": 34, "name": "test1"},
        {"id": 35, "name": "test2"},
        {"id": 36, "name": "test3"}
    ]
}`
p2, _ := WrapJSON([]byte(j2))

type Foo struct{ ID int16 }

slice, err := Map(p2, "items", func(p *Picker) (Foo, error) {
    f := Foo{}
    f.ID, _ = p.Int16("id")
    return f, nil
})
assert(t, []Foo{{ID: 34}, {ID: 35}, {ID: 36}}, slice, nil, err)
```

**Pick** is currently in a pre-alpha stage, a lot of changes going to happen both to api and structure.


More technical details about pick internals, in [architecture](doc/architecture.md) documentation.

___
## Special Mentions
Special thanks to **Konstantinos Pittas** ([@kostaspt](https://github.com/kostaspt)) for helping me ... **pick** the name of the library.
