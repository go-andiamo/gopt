package gopt

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestExtract(t *testing.T) {
	m := map[string]interface{}{
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

	m2 := map[int]interface{}{
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

func TestExtractJson(t *testing.T) {
	m := map[string]interface{}{
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

func TestGet(t *testing.T) {
	m := map[string]interface{}{
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
	m := map[string]interface{}{
		"foo": "foo value",
		"bar": nil,
	}
	om := OptMap[string, interface{}](m)
	ov := om.Get("foo")
	require.True(t, ov.IsPresent())
	ov = om.Get("bar")
	require.False(t, ov.IsPresent())
	ov = om.Get("baz")
	require.False(t, ov.IsPresent())
}

func TestOptMap_IfPresent(t *testing.T) {
	m := map[string]interface{}{
		"foo": "foo value",
		"bar": nil,
	}
	om := OptMap[string, interface{}](m)
	called := false
	f := func(key string, v interface{}) {
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
	m := map[string]interface{}{
		"foo": "foo value",
		"bar": nil,
	}
	om := OptMap[string, interface{}](m)
	called := false
	otherCalled := false
	f := func(key string, v interface{}) {
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
	m := map[string]interface{}{
		"foo": "foo value",
		"bar": nil,
	}
	om := OptMap[string, interface{}](m)
	v := om.Default("foo", "defaulted")
	require.Equal(t, "foo value", v)
	v = om.Default("bar", "defaulted")
	require.Equal(t, "defaulted", v)
	v = om.Default("baz", "defaulted")
	require.Equal(t, "defaulted", v)
}

func TestOptMap_ComputeIfAbsent(t *testing.T) {
	m := map[string]interface{}{
		"foo": "foo value",
		"bar": nil,
	}
	om := OptMap[string, interface{}](m)
	f := func(k string) interface{} {
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
	m := map[string]interface{}{
		"foo": "foo value",
		"bar": nil,
		"baz": "baz value",
	}
	om := OptMap[string, interface{}](m)
	f := func(k string, v interface{}) interface{} {
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
	m := map[string]interface{}{
		"foo": "foo value",
		"bar": nil,
	}
	om := OptMap[string, interface{}](m)
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
	m := map[string]interface{}{
		"foo": "foo value",
		"bar": nil,
	}
	om := OptMap[string, interface{}](m)
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
