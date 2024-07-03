# Pick
[![Go Report Card](https://goreportcard.com/badge/github.com/moukoublen/pick)](https://goreportcard.com/report/github.com/moukoublen/pick)
[![CI Status](https://github.com/moukoublen/pick/actions/workflows/ci.yml/badge.svg)](https://github.com/moukoublen/pick/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/moukoublen/pick/graph/badge.svg?token=6X9MMYZJZ8)](https://codecov.io/gh/moukoublen/pick)


**Pick** is a go package to access (using dot and array notation) and cast any kind of data, using best effort performance and best effort cast. It is an alternative to [stretchr/objx](https://github.com/stretchr/objx) aiming to provide three main things:

1. Modular approach regarding the caster, traverser and selector format
2. Best effort performance using `reflect` as last resort.
3. Best effort cast aiming to cast and convert between types as much as possible.

### Examples
```go
j := `{
    "item": {
        "one": 1,
        "two": "ok",
        "three": ["element 1", 2, "element 3"]
    },
    "float": 2.12
}`

p1, _ := WrapJSON([]byte(j))

got, err := p1.String("item.three[1]")
// got == "2"

got, err := p1.Uint64("item.three[1]")
// got == uint64(2)

got, err := p1.String("item.three[-1]") // access the last element
// got == "element 3"

got, err := p1.Float32("float")
// got == float32(2.12)

got, err := p1.Int64("float")
// got == int64(2)
// err is cast.ErrCastLostDecimals

got, err := p1.Int32("non-existing")
// got == int32(0)
// err is ErrFieldNotFound

m := p1.Must()

got := m.Int32("item.one")
// got == int32(1)

got := m.Int32("non-existing")
// got == int32(0)
```

**`Map` function**
```go
j2 := `{
    "items": [
        {"id": 34, "name": "test1"},
        {"id": 35, "name": "test2"},
        {"id": 36, "name": "test3"}
    ]
}`
p2, _ := WrapJSON([]byte(j2))

got, err := Map(p2, "items", func(p *Picker) (int16, error) {
    n, _ := p.Int16("id")
    return n, nil
})
// got == []int16{34, 35, 36}
// err == nil
```

**Time functions**
```go
dateData := map[string]any{
    "time1":     "1977-05-25T22:30:00Z",
    "time2":     "Wed, 25 May 1977 18:30:00 -0400",
    "timeSlice": []string{"1977-05-25T18:30:00Z", "1977-05-25T20:30:00Z", "1977-05-25T22:30:00Z"},
}

p3 := Wrap(dateData)

got, err := p3.Time("time1")
// got == time.Date(1977, time.May, 25, 22, 30, 0, 0, time.UTC)
// err == nil

loc, _ := time.LoadLocation("America/New_York")
got, err := p3.TimeWithConfig(cast.TimeCastConfig{StringFormat: time.RFC1123Z}, "time2")
// got == time.Date(1977, time.May, 25, 18, 30, 0, 0, loc)
// err == nil

got, err := p3.TimeSlice("timeSlice")
// got == []time.Time{
//     time.Date(1977, time.May, 25, 18, 30, 0, 0, time.UTC),
//     time.Date(1977, time.May, 25, 20, 30, 0, 0, time.UTC),
//     time.Date(1977, time.May, 25, 22, 30, 0, 0, time.UTC),
// },
// err == nil
```


### API
As an `API` we define a set of functions like this `Bool(T) Output` for all basic types. There are 4 different APIs for a picker.

  * Selector API, the default one that is embedded in `Picker`. <br>E.g. `Picker.Bool(selector string) (bool, error)`
  * Selector Must API that can be accessed by calling `Picker.Must()`. <br>E.g. `Picker.Must().Bool(selector string) bool`

**Supported Types in API**
  * `bool` / `[]bool`
  * `byte` / `[]byte`
  * `float{32,64}` / `[]float{32,64}`
  * `int{,8,16,32,64}` / `[]int{,8,16,32,64}`
  * `uint{,8,16,32,64}` / `[]uint{,8,16,32,64}`
  * `string` / `[]string`
  * `time.Time` / `[]time.Time`
  * `time.Duration` / `[]time.Duration`

Examples:

**Selector Must API**
```go
p1.Must().String("item.three[1]") // == "2"
sm := p1.Must()
sm.Uint64("item.three[1]") // == uint64(2)
sm.Int32("item.one")       // == int32(1)
sm.Float32("float")        // == float32(2.12)
sm.Int64("float")          // == int64(2)
```

**Pick** is currently in a pre-alpha stage, a lot of changes going to happen both to api and structure.


More technical details about pick internals, in [architecture](doc/architecture.md) documentation.

___
## Acknowledgements
Special thanks to **Konstantinos Pittas** ([@kostaspt](https://github.com/kostaspt)) for helping me ... **pick** the name of the library.
