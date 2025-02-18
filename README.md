# Pick
[![CI Status](https://github.com/moukoublen/pick/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/moukoublen/pick/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/moukoublen/pick/graph/badge.svg?token=6X9MMYZJZ8)](https://codecov.io/gh/moukoublen/pick)
[![Go Report Card](https://goreportcard.com/badge/github.com/moukoublen/pick)](https://goreportcard.com/report/github.com/moukoublen/pick)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/moukoublen/pick)](https://pkg.go.dev/github.com/moukoublen/pick)

**Pick** is a go package to access (using dot and array notation) and cast any kind of data, using best effort performance and best effort cast. It is an alternative to [stretchr/objx](https://github.com/stretchr/objx) aiming to provide three main things:

1. Modular approach regarding the caster, traverser and selector format
2. Best effort performance using `reflect` as last resort.
3. Best effort cast aiming to cast and convert between types as much as possible.

## Install

```bash
go get -u github.com/moukoublen/pick
```

## Examples
```go
j := `{
    "item": {
        "one": 1,
        "two": "ok",
        "three": ["element 1", 2, "element 3"]
    },
    "float": 2.12,
    "floatDec": 2
}`

p1, _ := WrapJSON([]byte(j))
got, err := p1.String("item.three[1]")  // "2", nil
got, err := p1.Uint64("item.three[1]")  // uint64(2), nil
got, err := p1.String("item.three[-1]") // "element 3", nil | (access the last element)element
got, err := p1.Float32("float")         // float32(2.12), nil
got, err := p1.Int64("float")           // int64(2), ErrCastLostDecimals
got, err := p1.Int64("floatDec")        // int64(2), nil
got, err := p1.Int32("non-existing")    // 0, ErrFieldNotFound

// Relaxed API (no error in return)
got := p1.Relaxed().Int32("item.one")  // int32(1)
got := p1.Relaxed().Int32("none")      // int32(0)

// Relaxed with errors sink
sink := ErrorsSink{}
sm2 := p1.Relaxed(&sink)
got := sm2.String("item.three[1]") // "2"
got := sm2.Uint64("item.three[1]") // uint64(2)
got := sm2.Int32("item.one")       // int32(1)
got := sm2.Float32("float")        // float32(2.12)
got := sm2.Int64("float")          // int64(2) | error lost decimals
got := sm2.String("item.three")    // ""       | error field not found
err := sink.Outcome()              // joined error
```

#### Generics functions
```go
got, err := Get[int64](p1, "item.three[1]")  // (int64(2), nil)
got, err := Get[string](p1, "item.three[1]") // ("2", nil)

m := p1.Relaxed()
got := RelaxedGet[string](m, "item.three[1]") // "2"

got, err := Path[string](p1, Field("item"), Field("three"), Index(1)) // ("2", nil)
got := RelaxedPath[float32](m, Field("item"), Field("one")          // float32(1)
```

#### `Map` functions
```go
j2 := `{
    "items": [
        {"id": 34, "name": "test1", "array": [1,2,3]},
        {"id": 35, "name": "test2", "array": [4,5,6]},
        {"id": 36, "name": "test3", "array": [7,8,9]}
    ]
}`
p2, _ := WrapJSON([]byte(j2))

got, err := Map(p2, "items", func(p Picker) (int16, error) {
    n, _ := p.Int16("id")
    return n, nil
}) // ( []int16{34, 35, 36}, nil )

got, err := FlatMap(p2, "items", func(p Picker) ([]int16, error) {
    return p.Int16Slice("array")
}) // ( []int16{1, 2, 3, 4, 5, 6, 7, 8, 9}, nil )

got3, err3 := MapFilter(p2, "items", func(p Picker) (int32, bool, error) {
    i, err := p.Int32("id")
    return i, i%2 == 0, err
}) // ( []int32{34, 36}, nil )
```

#### `Each`/`EachField` functions
```go
j := `{
    "2023-01-01": [1,2,3],
    "2023-01-02": [1,2,3,4],
    "2023-01-03": [1,2,3,4,5]
}`
p, _ := WrapJSON([]byte(j))

lens := map[string]int{}
err := EachField(p, "", func(field string, value Picker, numOfFields int) error {
    l, _ := value.Len("")
    lens[field] = l
    return nil
}) // nil
// lens == map[string]int{"2023-01-01": 3, "2023-01-02": 4, "2023-01-03": 5})


sumEvenIndex := 0
err = Each(p, "2023-01-03", func(index int, item Picker, totalLength int) error {
    if index%2 == 0 {
        sumEvenIndex += item.Relaxed().Int("")
    }
    return nil
}) // nil
// sumEvenIndex == 6
```

  * [Each](root.go#L19) / [RelaxedEach](root.go#L143)
  * [EachField](root.go#L35) / [RelaxedEachField](root.go#L167)
  * [Map](root.go#L50) / [RelaxedMap](root.go#L201)
  * [FlatMap](root.go#L80) / [RelaxedFlatMap](root.go#L233)
  * [MapFilter](root.go#L65) / [RelaxedMapFilter](root.go#209)


### API
As an `API` we define a set of functions like this `Bool(T) Output` for all basic types. There are 2 different APIs for a picker.

  * Default API, the default one that is embedded in `Picker`. <br>E.g. `Picker.Bool(selector string) (bool, error)`
  * Relaxed API that can be accessed by calling `Picker.Relaxed()`. <br>E.g. `Picker.Relaxed().Bool(selector string) bool`. Relaxed API **does not** return error neither it panics in case of one. It just ignores any possible error and returns the default zero value in case of one. Errors can be optionally gathered using an `ErrorGatherer`.

**Supported Types in API**
  * `bool` / `[]bool`
  * `byte` / `[]byte`
  * `float{32,64}` / `[]float{32,64}`
  * `int{,8,16,32,64}` / `[]int{,8,16,32,64}`
  * `uint{,8,16,32,64}` / `[]uint{,8,16,32,64}`
  * `string` / `[]string`
  * `time.Time` / `[]time.Time`
  * `time.Duration` / `[]time.Duration`

### Relaxed API
```go
p1.Relaxed().String("item.three[1]") // == "2"
sm := p1.Relaxed()
sm.Uint64("item.three[1]") // == uint64(2)
```

Optionally an `ErrorGatherer` could be provided to `.Relaxed()` initializer to receive and handle each error produced by the operations.

A default implementation of `ErrorGatherer` is the `ErrorsSink`, which gathers all errors into a single one.

___
**Pick** is currently in a alpha stage, a lot of changes going to happen both to api and structure.


More technical details about pick internals, in [architecture](doc/architecture.md) documentation.

___
## Acknowledgements
Special thanks to **Konstantinos Pittas** ([@kostaspt](https://github.com/kostaspt)) for helping me ... **pick** the name of the library.
