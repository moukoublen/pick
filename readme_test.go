package pick

import (
	"errors"
	"testing"
	"time"

	"github.com/ifnotnil/x/tst"
	"github.com/moukoublen/pick/internal/testingx"
)

func TestReadme(t *testing.T) {
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

	t.Run("basic", func(t *testing.T) {
		assert2 := testingx.AssertEqualWithErrorFn(t)

		assert2(p1.String("item.three[1]"))("2", nil)
		assert2(p1.Uint64("item.three[1]"))(uint64(2), nil)
		assert2(p1.String("item.three[-1]"))("element 3", nil) // | (access the last element)
		assert2(p1.Float32("float"))(float32(2.12), nil)
		assert2(p1.Int64("float"))(int64(2), ErrConvertLostDecimals)
		assert2(p1.Int64("floatDec"))(int64(2), nil)
		assert2(p1.Int("non-existing"))(0, ErrFieldNotFound)
	})

	// Relaxed API
	t.Run("basic Relaxed api", func(t *testing.T) {
		assert := testingx.AssertEqualFn(t)

		assert(p1.Relaxed().Int32("item.one"), int32(1))
		assert(p1.Relaxed().Int32("non-existing"), int32(0))

		sink := &ErrorsSink{}
		m := p1.Relaxed(sink)
		assert(m.String("item.three"), "")
		assert(m.String("item.three[1]"), "2")
		assert(m.Uint64("item.three[1]"), uint64(2))
		assert(m.Int32("item.one"), int32(1))
		assert(m.Float32("float"), float32(2.12))
		assert(m.Int64("float"), int64(2))
		assert(sink.Outcome() != nil, true)
		// es.Outcome() = picker error with selector `item.three` ... invalid type | picker error with selector `float` missing decimals error
	})

	t.Run("generics", func(t *testing.T) {
		assert := testingx.AssertEqualFn(t)
		assert2 := testingx.AssertEqualWithErrorFn(t)

		assert2(Get[int64](p1, "item.three[1]"))(int64(2), nil)
		assert2(Get[string](p1, "item.three[1]"))("2", nil)

		m := p1.Relaxed()
		assert(RelaxedGet[string](m, "item.three[1]"), "2")

		assert2(Path[string](p1, Field("item"), Field("three"), Index(1)))("2", nil)
		assert(RelaxedPath[float32](m, Field("item"), Field("one")), float32(1))
	})

	t.Run("map examples", func(t *testing.T) {
		assert := testingx.AssertEqualFn(t)

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
		})
		assert(got, []int16{34, 35, 36})
		assert(err, nil)

		got2, err2 := FlatMap(p2, "items", func(p Picker) ([]int16, error) {
			return p.Int16Slice("array")
		})
		assert(got2, []int16{1, 2, 3, 4, 5, 6, 7, 8, 9})
		assert(err2, nil)

		got3, err3 := MapFilter(p2, "items", func(p Picker) (int32, bool, error) {
			i, err := p.Int32("id")
			return i, i%2 == 0, err
		})
		assert(got3, []int32{34, 36})
		assert(err3, nil)
	})

	t.Run("each and each field function", func(t *testing.T) {
		assert := testingx.AssertEqualFn(t)

		j := `{
				"2023-01-01": [0,1,2],
				"2023-01-02": [0,1,2,3],
				"2023-01-03": [0,1,2,3,4]
		}`
		p, _ := WrapJSON([]byte(j))

		lens := map[string]int{}
		err := EachField(p, "", func(field string, value Picker, numOfFields int) error {
			l, _ := value.Len("")
			lens[field] = l
			return nil
		})
		tst.NoError()(t, err)
		assert(lens, map[string]int{"2023-01-01": 3, "2023-01-02": 4, "2023-01-03": 5})

		sumEvenIndex := 0
		err = Each(p, "2023-01-03", func(index int, item Picker, totalLength int) error {
			if index%2 == 0 {
				sumEvenIndex += item.Relaxed().Int("")
			}
			return nil
		})
		assert(sumEvenIndex, 6)
		tst.NoError()(t, err)
	})

	// time API
	t.Run("time", func(t *testing.T) {
		assert := testingx.AssertEqualFn(t)

		dateData := map[string]any{
			"time1":     "1977-05-25T22:30:00Z",
			"time2":     "Wed, 25 May 1977 18:30:00 -0400",
			"timeSlice": []string{"1977-05-25T18:30:00Z", "1977-05-25T20:30:00Z", "1977-05-25T22:30:00Z"},
		}
		p3 := Wrap(dateData)
		{
			got, err := p3.Time("time1")
			assert(got, time.Date(1977, time.May, 25, 22, 30, 0, 0, time.UTC))
			assert(err, nil)
		}
		{
			loc, _ := time.LoadLocation("America/New_York")
			got, err := p3.TimeWithConfig(TimeConvertConfig{StringFormat: time.RFC1123Z}, "time2")
			assert(got, time.Date(1977, time.May, 25, 18, 30, 0, 0, loc))
			assert(err, nil)
		}
		{
			got, err := p3.TimeSlice("timeSlice")
			assert(
				got,
				[]time.Time{
					time.Date(1977, time.May, 25, 18, 30, 0, 0, time.UTC),
					time.Date(1977, time.May, 25, 20, 30, 0, 0, time.UTC),
					time.Date(1977, time.May, 25, 22, 30, 0, 0, time.UTC),
				},
			)
			assert(err, nil)
		}
	})
}

func TestReadmeConvert(t *testing.T) {
	eq := testingx.AssertEqualFn(t)

	c := NewDefaultConverter()

	{
		got, err := c.AsInt8(int32(10))
		eq(got, int8(10))
		eq(err, nil)
	}

	{
		got, err := c.AsInt8("10")
		eq(got, int8(10))
		eq(err, nil)
	}

	{
		got, err := c.AsInt8(128)
		eq(got, int8(-128))
		eq(errors.Is(err, ErrConvertOverFlow), true)
	}

	{
		got, err := c.AsInt8(10.12)
		eq(got, int8(10))
		eq(errors.Is(err, ErrConvertLostDecimals), true)
	}

	{
		got, err := c.AsInt8(float64(10.00))
		eq(got, int8(10))
		eq(err, nil)
	}
}
