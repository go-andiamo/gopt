package gopt

import (
	"github.com/stretchr/testify/require"
	"testing"
)

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

func TestOptMap_GetOrDefault(t *testing.T) {
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
