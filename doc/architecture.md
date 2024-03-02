## Pick components

### 1) Notation
Notation is the text _format_ that is used in order to refer to a field. The main functionality of the notation interface is to parse a **selector** string into a path which is a slice `[]Key` and each `Key` could be of type **field** or **index**.

The default notation is the dot notation `DotNotation`. Example:

```golang
selectorString := "near_earth_objects[12].is_potentially_hazardous_asteroid"

DotNotation{}.Parse(selectorString)
    // will result to:
    []Key{
        Key{Name: "near_earth_objects",                Type: KeyTypeField},
        Key{Index: 12,                                 Type: KeyTypeIndex},
        Key{Name: "is_potentially_hazardous_asteroid", Type: KeyTypeField},
    }


// the Format function takes a []Key and formats it to the notation accordingly.
DotNotation{}.Format(Field("near_earth_objects"), Index(12), Field("is_potentially_hazardous_asteroid"))
    // will result to:
    "near_earth_objects[12].is_potentially_hazardous_asteroid"
```
The parse functionality aims to achieve the best possible performance with the least possible allocations. It iterates over the initial selector string, after converting it to rune slice, as much as possible without allocating new buffers.

Terminology:
  * **selector**: The `string` that describes a path (e.g. for dot notation `"near_earth_objects[12].is_potentially_hazardous_asteroid"`)
  * **path**: A slice of `[]Key`. The result of parsing a selector.
  * **Key**: A single descriptor of a path/selector that can be of type `Field`, indicating access to named fields, or `Index`, indicating access to arrays.
  * **Notation**: an implementation that specifies a format in which can parse selectors to path (`[]Key`) and format path (`[]Key`) back to selector.


### 2) Traverser
Traverser implements the functionality of accessing and retrieving a field on a data set (slice, map or struct) given a path (`[]Key`).

The default implementation `DefaultTraverser` aims to use reflect as last resort by attempting first to cast to most common types (`map[string]any` in case of `Field` and `[]any` in case of `Index`) and direct access to them. If the dataset is not one of those types, it attempts to access using reflect. This happens sequently for each `Key` of the path (`[]Key`).

Note: _Traverser needs a caster just in case it tries to traverse a map that the key of is of different type than `string` or `int`_

A simple example of traverser
```golang
data := map[string]any{
    "one": map[string]any{
        "two": map[string]any{
            "three": "value"
        },
    },
}

tr := NewDefaultTraverser(cast.NewCaster())
v, err := tr.Retrieve(data, []Key{ Field("one"), Field("two"), Field("three") })
// v == any("value")
// err == nil
```

### 3) Caster
Caster attempts to cast between types using reflect as a last resort. It also checks for overflows or lost decimals after casting and returns errors.

Example:
```golang
c := NewCaster()

got, err := c.AsInt8(int32(10))
// got == int8(10)
// err == nil

got, err := c.AsInt16("10")
// got == int8(10)
// err == nil

got, err := c.AsInt16(128)
// got == int8(-128)
// err is ErrCastOverFlow

got, err := c.AsInt8(10.12)
// got == int8(10)
// err is ErrCastLostDecimals

got, err := c.AsInt8(float64(10.00))
// got == int8(10)
// err == nil
```
