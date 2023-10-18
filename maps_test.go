package gopt

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestExtract(t *testing.T) {
	m := map[string]any{
		"str":  "Str",
		"int":  16,
		"time": time.Now(),
	}

	s := Extract[string, string](m, "str")
	require.True(t, s.IsPresent())
	require.Equal(t, "Str", s.OrElse(""))

	s = Extract[string, string](m, "int")
	require.False(t, s.IsPresent())

	i := Extract[string, int](m, "str")
	require.False(t, i.IsPresent())

	i = Extract[string, int](m, "int")
	require.True(t, i.IsPresent())
	require.Equal(t, 16, i.OrElse(0))

	dt := Extract[string, time.Time](m, "time")
	require.True(t, dt.IsPresent())

	dt = Extract[string, time.Time](m, "str")
	require.False(t, dt.IsPresent())

	m2 := map[int]any{
		1: "Str",
		2: 16,
	}
	s = Extract[int, string](m2, 1)
	require.True(t, s.IsPresent())
	require.Equal(t, "Str", s.OrElse(""))
	s = Extract[int, string](m2, 2)
	require.False(t, s.IsPresent())
	i = Extract[int, int](m2, 2)
	require.True(t, i.IsPresent())
	require.Equal(t, 16, i.OrElse(0))
	i = Extract[int, int](m2, 1)
	require.False(t, i.IsPresent())
}

func TestExtract_WithConverters(t *testing.T) {
	m := map[string]any{
		"int":  json.Number("16"),
		"time": time.Now(),
		"str":  "2023-10-18T13:00:00Z",
	}

	jnToInt := func(v any) (int, bool) {
		switch vt := v.(type) {
		case json.Number:
			if v, err := vt.Int64(); err == nil {
				return int(v), true
			}
		}
		return 0, false
	}
	i := Extract[string, int](m, "int")
	require.False(t, i.IsPresent())
	i = Extract[string, int](m, "int", jnToInt)
	require.True(t, i.IsPresent())

	timeToStr := func(v any) (string, bool) {
		switch vt := v.(type) {
		case time.Time:
			return vt.Format(time.RFC3339), true
		}
		return "", false
	}
	s := Extract[string, string](m, "time")
	require.False(t, s.IsPresent())
	s = Extract[string](m, "time", timeToStr)
	require.True(t, s.IsPresent())

	strToTime := func(v any) (time.Time, bool) {
		switch vt := v.(type) {
		case string:
			if tv, err := time.Parse(time.RFC3339, vt); err == nil {
				return tv, true
			}
		}
		return time.Time{}, false
	}
	tm := Extract[string, time.Time](m, "str")
	require.False(t, tm.IsPresent())
	tm = Extract[string, time.Time](m, "time")
	require.True(t, tm.IsPresent())
	tm = Extract[string, time.Time](m, "str", strToTime)
	require.True(t, tm.IsPresent())
	tm = Extract[string, time.Time](m, "time", strToTime)
	require.True(t, tm.IsPresent())
}

func TestExtractJson(t *testing.T) {
	m := map[string]any{
		"str":  "Str",
		"int":  16,
		"time": time.Now(),
	}

	s := ExtractJson[string](m, "str")
	require.True(t, s.IsPresent())
	require.Equal(t, "Str", s.OrElse(""))

	s = ExtractJson[string](m, "int")
	require.False(t, s.IsPresent())

	i := ExtractJson[int](m, "str")
	require.False(t, i.IsPresent())

	i = ExtractJson[int](m, "int")
	require.True(t, i.IsPresent())
	require.Equal(t, 16, i.OrElse(0))

	dt := ExtractJson[time.Time](m, "time")
	require.True(t, dt.IsPresent())

	dt = ExtractJson[time.Time](m, "str")
	require.False(t, dt.IsPresent())
}

func TestExtractJson_WithConverters(t *testing.T) {
	m := map[string]any{
		"int":  json.Number("16"),
		"time": time.Now(),
		"str":  "2023-10-18T13:00:00Z",
	}

	jnToInt := func(v any) (int, bool) {
		switch vt := v.(type) {
		case json.Number:
			if v, err := vt.Int64(); err == nil {
				return int(v), true
			}
		}
		return 0, false
	}
	i := ExtractJson[int](m, "int")
	require.False(t, i.IsPresent())
	i = ExtractJson[int](m, "int", jnToInt)
	require.True(t, i.IsPresent())

	timeToStr := func(v any) (string, bool) {
		switch vt := v.(type) {
		case time.Time:
			return vt.Format(time.RFC3339), true
		}
		return "", false
	}
	s := ExtractJson[string](m, "time")
	require.False(t, s.IsPresent())
	s = ExtractJson[string](m, "time", timeToStr)
	require.True(t, s.IsPresent())

	strToTime := func(v any) (time.Time, bool) {
		switch vt := v.(type) {
		case string:
			if tv, err := time.Parse(time.RFC3339, vt); err == nil {
				return tv, true
			}
		}
		return time.Time{}, false
	}
	tm := ExtractJson[time.Time](m, "str")
	require.False(t, tm.IsPresent())
	tm = ExtractJson[time.Time](m, "time")
	require.True(t, tm.IsPresent())
	tm = ExtractJson[time.Time](m, "str", strToTime)
	require.True(t, tm.IsPresent())
	tm = ExtractJson[time.Time](m, "time", strToTime)
	require.True(t, tm.IsPresent())
}

func TestExtractJsonPath(t *testing.T) {
	m := map[string]any{
		"arr": []any{
			"first",
			2,
			map[string]any{
				"foo": "bar",
			},
		},
		"foo": map[string]any{
			"bar": map[string]any{
				"baz": []any{
					"first",
					2,
					map[string]any{
						"foo": "bar",
					},
				},
			},
		},
	}
	sl, p := ExtractJsonPath[[]any](m, "arr")
	require.True(t, sl.IsPresent())
	require.Equal(t, 0, len(p))
	sm, p := ExtractJsonPath[map[string]any](m, "foo")
	require.True(t, sm.IsPresent())
	require.Equal(t, 0, len(p))

	first, p := ExtractJsonPath[string](m, "arr[0]")
	require.True(t, first.IsPresent())
	require.Equal(t, 2, len(p))
	require.True(t, p[0])
	require.True(t, p[1])

	last, p := ExtractJsonPath[map[string]any](m, "arr[-1]")
	require.True(t, last.IsPresent())
	require.Equal(t, 2, len(p))
	require.True(t, p[0])
	require.True(t, p[1])

	foo, p := ExtractJsonPath[string](m, "arr[-1].foo")
	require.True(t, foo.IsPresent())
	require.Equal(t, 2, len(p))
	require.True(t, p[0])
	require.True(t, p[1])
	require.Equal(t, "bar", foo.Default(""))

	x, p := ExtractJsonPath[any](m, "arr[-1].foo.xxx")
	require.False(t, x.IsPresent())
	require.Equal(t, 3, len(p))
	require.True(t, p[0])
	require.True(t, p[1])
	require.False(t, p[2])

	x, p = ExtractJsonPath[any](m, "arr[0].foo")
	require.False(t, x.IsPresent())
	require.Equal(t, 2, len(p))
	require.True(t, p[0])
	require.False(t, p[1])

	x, p = ExtractJsonPath[any](m, "arr[-4]")
	require.False(t, x.IsPresent())
	require.Equal(t, 2, len(p))
	require.True(t, p[0])
	require.False(t, p[1])

	x, p = ExtractJsonPath[any](m, "foo[0]")
	require.False(t, x.IsPresent())
	require.Equal(t, 2, len(p))
	require.True(t, p[0])
	require.False(t, p[1])

	x, p = ExtractJsonPath[any](m, "xxx[0]")
	require.False(t, x.IsPresent())
	require.Equal(t, 1, len(p))
	require.False(t, p[0])

	x, p = ExtractJsonPath[any](m, "xxx.yyy")
	require.False(t, x.IsPresent())
	require.Equal(t, 1, len(p))
	require.False(t, p[0])

	foo, p = ExtractJsonPath[string](m, "foo.bar.baz[-1].foo")
	require.True(t, foo.IsPresent())
	require.Equal(t, 4, len(p))
	require.True(t, p[0])
	require.True(t, p[1])
	require.True(t, p[2])
	require.True(t, p[3])

	foo, p = ExtractJsonPath[string](m, "foo.bar.baz[-4].foo")
	require.False(t, foo.IsPresent())
	require.Equal(t, 4, len(p))
	require.True(t, p[0])
	require.True(t, p[1])
	require.True(t, p[2])
	require.False(t, p[3])

	first, p = ExtractJsonPath[string](m, "foo.bar.baz[0]")
	require.True(t, first.IsPresent())
	require.Equal(t, 4, len(p))
	require.True(t, p[0])
	require.True(t, p[1])
	require.True(t, p[2])
	require.True(t, p[3])

	foo, p = ExtractJsonPath[string](m, "foo.bar.baz.foo")
	require.False(t, foo.IsPresent())
	require.Equal(t, 3, len(p))
	require.True(t, p[0])
	require.True(t, p[1])
	require.False(t, p[2])
}

func TestExtractJsonPath_WithConverters(t *testing.T) {
	m := map[string]any{
		"arr": []any{
			"first",
			2,
		},
		"foo": map[string]any{
			"bar": map[string]any{
				"baz": 2,
			},
		},
	}
	intToStr := func(v any) (string, bool) {
		switch vt := v.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			return fmt.Sprintf("%d", vt), true
		}
		return "", false
	}

	s, _ := ExtractJsonPath[string](m, "arr[-1]")
	require.False(t, s.IsPresent())
	s, _ = ExtractJsonPath[string](m, "arr[-1]", intToStr)
	require.True(t, s.IsPresent())
	require.Equal(t, "2", s.Default(""))
	s, _ = ExtractJsonPath[string](m, "arr[0]")
	require.True(t, s.IsPresent())
	require.Equal(t, "first", s.Default(""))
	s, _ = ExtractJsonPath[string](m, "arr[0]", intToStr)
	require.True(t, s.IsPresent())
	require.Equal(t, "first", s.Default(""))

	s, _ = ExtractJsonPath[string](m, "foo.bar.baz")
	require.False(t, s.IsPresent())
	s, _ = ExtractJsonPath[string](m, "foo.bar.baz", intToStr)
	require.True(t, s.IsPresent())
	require.Equal(t, "2", s.Default(""))
}

func TestGet(t *testing.T) {
	m := map[string]any{
		"foo": "foo value",
		"bar": 1,
	}
	o := Get(m, "foo")
	require.True(t, o.IsPresent())
	require.Equal(t, "foo value", o.OrElse(nil))
	o = Get(m, "bar")
	require.True(t, o.IsPresent())
	require.Equal(t, 1, o.OrElse(nil))
	o = Get(m, "baz")
	require.False(t, o.IsPresent())
	require.Equal(t, "", o.OrElse(""))

	m2 := map[int]*myStruct{
		1: {},
		2: nil,
	}
	o2 := Get(m2, 1)
	require.True(t, o2.IsPresent())
	o2 = Get(m2, 2)
	require.False(t, o2.IsPresent())
	o2 = Get(m2, 3)
	require.False(t, o2.IsPresent())
}

func TestOptMap_Get(t *testing.T) {
	m := map[string]any{
		"foo": "foo value",
		"bar": nil,
	}
	om := OptMap[string, any](m)
	ov := om.Get("foo")
	require.True(t, ov.IsPresent())
	ov = om.Get("bar")
	require.False(t, ov.IsPresent())
	ov = om.Get("baz")
	require.False(t, ov.IsPresent())
}

func TestOptMap_IfPresent(t *testing.T) {
	m := map[string]any{
		"foo": "foo value",
		"bar": nil,
	}
	om := OptMap[string, any](m)
	called := false
	f := func(key string, v any) {
		called = true
	}
	om.IfPresent("foo", f)
	require.True(t, called)
	called = false
	om.IfPresent("bar", f)
	require.False(t, called)
	om.IfPresent("baz", f)
	require.False(t, called)

	om.IfPresent("bar", f).IfPresent("foo", f)
	require.True(t, called)
}

func TestOptMap_IfPresentOtherwise(t *testing.T) {
	m := map[string]any{
		"foo": "foo value",
		"bar": nil,
	}
	om := OptMap[string, any](m)
	called := false
	otherCalled := false
	f := func(key string, v any) {
		called = true
	}
	other := func(key string) {
		otherCalled = true
	}
	om.IfPresentOtherwise("foo", f, other)
	require.True(t, called)
	require.False(t, otherCalled)
	called = false
	otherCalled = false
	om.IfPresentOtherwise("bar", f, other)
	require.False(t, called)
	require.True(t, otherCalled)
	called = false
	otherCalled = false
	om.IfPresentOtherwise("baz", f, other)
	require.False(t, called)
	require.True(t, otherCalled)

	called = false
	otherCalled = false
	om.IfPresentOtherwise("bar", f, other).IfPresentOtherwise("foo", f, other)
	require.True(t, called)
	require.True(t, otherCalled)
}

func TestOptMap_Default(t *testing.T) {
	m := map[string]any{
		"foo": "foo value",
		"bar": nil,
	}
	om := OptMap[string, any](m)
	v := om.Default("foo", "defaulted")
	require.Equal(t, "foo value", v)
	v = om.Default("bar", "defaulted")
	require.Equal(t, "defaulted", v)
	v = om.Default("baz", "defaulted")
	require.Equal(t, "defaulted", v)
}

func TestOptMap_ComputeIfAbsent(t *testing.T) {
	m := map[string]any{
		"foo": "foo value",
		"bar": nil,
	}
	om := OptMap[string, any](m)
	f := func(k string) any {
		if k == "baz" {
			return nil
		}
		return "computed"
	}
	v := om.ComputeIfAbsent("foo", f)
	require.Equal(t, "foo value", v)
	orgV, ok := m["foo"]
	require.True(t, ok)
	require.Equal(t, "foo value", orgV)
	v = om.ComputeIfAbsent("bar", f)
	require.Equal(t, "computed", v)
	orgV, ok = m["bar"]
	require.True(t, ok)
	require.Equal(t, "computed", orgV)
	v = om.ComputeIfAbsent("baz", f)
	require.Nil(t, v)
	_, ok = m["baz"]
	require.False(t, ok)
}

func TestOptMap_ComputeIfPresent(t *testing.T) {
	m := map[string]any{
		"foo": "foo value",
		"bar": nil,
		"baz": "baz value",
	}
	om := OptMap[string, any](m)
	f := func(k string, v any) any {
		if k == "baz" {
			return nil
		}
		return "computed"
	}
	v := om.ComputeIfPresent("foo", f)
	require.Equal(t, "computed", v)
	orgV, ok := m["foo"]
	require.True(t, ok)
	require.Equal(t, "computed", orgV)
	v = om.ComputeIfPresent("bar", f)
	require.Nil(t, v)
	orgV, ok = m["bar"]
	require.True(t, ok)
	require.Nil(t, orgV)
	v = om.ComputeIfPresent("baz", f)
	require.Nil(t, v)
	_, ok = m["baz"]
	require.False(t, ok)
}

func TestOptMap_PutIfAbsent(t *testing.T) {
	m := map[string]any{
		"foo": "foo value",
		"bar": nil,
	}
	om := OptMap[string, any](m)
	r := om.PutIfAbsent("foo", "absent value")
	require.False(t, r)
	orgV, ok := m["foo"]
	require.True(t, ok)
	require.Equal(t, "foo value", orgV)
	r = om.PutIfAbsent("bar", "absent value")
	require.True(t, r)
	orgV, ok = m["bar"]
	require.True(t, ok)
	require.Equal(t, "absent value", orgV)
	r = om.PutIfAbsent("baz", "absent value")
	require.True(t, r)
	orgV, ok = m["baz"]
	require.True(t, ok)
	require.Equal(t, "absent value", orgV)
}

func TestOptMap_ReplaceIfPresent(t *testing.T) {
	m := map[string]any{
		"foo": "foo value",
		"bar": nil,
	}
	om := OptMap[string, any](m)
	r := om.ReplaceIfPresent("foo", "replacement value")
	require.True(t, r)
	orgV, ok := m["foo"]
	require.True(t, ok)
	require.Equal(t, "replacement value", orgV)
	r = om.ReplaceIfPresent("bar", "replacement value")
	require.False(t, r)
	orgV, ok = m["bar"]
	require.True(t, ok)
	require.Nil(t, orgV)
	r = om.ReplaceIfPresent("baz", "replacement value")
	require.False(t, r)

	r = om.ReplaceIfPresent("foo", nil)
	require.True(t, r)
	_, ok = m["foo"]
	require.False(t, ok)
}
