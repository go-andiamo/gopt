package gopt

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

type myStruct struct {
	Foo string
}

func TestOfNillable(t *testing.T) {
	opt := OfNillable[String](nil)
	require.False(t, opt.IsPresent())
	_, err := opt.Get()
	require.Error(t, err)

	x := &myStruct{
		Foo: "",
	}
	opt2 := Of(x).Empty()
	require.False(t, opt2.IsPresent())
	_, err = opt2.Get()
	require.Error(t, err)

	var v *myStruct
	v = opt2.OrElse(&myStruct{
		Foo: "bar",
	})
	require.Equal(t, "bar", v.Foo)

	opt3 := Of("").Empty()
	require.False(t, opt3.IsPresent())
	_, err = opt3.Get()
	require.Error(t, err)

	opt4 := OfNillable[String](nil)
	require.False(t, opt4.IsPresent())
}

func TestOf(t *testing.T) {
	opt := Of("aaa")
	require.True(t, opt.IsPresent())
	v, err := opt.Get()
	require.NoError(t, err)
	require.Equal(t, "aaa", v)

	str := "aaa"
	bigStr := String(&str)
	opt2 := OfNillable(bigStr)
	require.True(t, opt2.IsPresent())
	v2, err := opt2.Get()
	require.NoError(t, err)
	require.Equal(t, bigStr, v2)

	opt2a := OfNillable[String](nil)
	require.False(t, opt2a.IsPresent())

	type nillableTime *time.Time
	opt3 := OfNillable[nillableTime](nil)
	require.False(t, opt3.IsPresent())
}

func TestOptional_IsPresent(t *testing.T) {
	o := Of("aaa")
	require.True(t, o.IsPresent())
	o = Empty[string]()
	require.False(t, o.IsPresent())
}

func TestOptional_IfPresent(t *testing.T) {
	called := false
	collected := ""
	o := Of("aaa")
	f := func(v string) {
		called = true
		collected = v
	}
	o.IfPresent(f)
	require.True(t, called)
	require.Equal(t, "aaa", collected)

	o = Empty[string]()
	called = false
	o.IfPresent(f)
	require.False(t, called)
	require.Equal(t, "aaa", collected)
}

func TestOptional_OrElse(t *testing.T) {
	o := Empty[string]()
	v := o.OrElse("bbb")
	require.Equal(t, "bbb", v)

	o = Of("aaa")
	v = o.OrElse("bbb")
	require.Equal(t, "aaa", v)
}

func TestOptional_OrElseGet(t *testing.T) {
	o := Empty[string]()
	called := false
	f := func() string {
		called = true
		return "aaa"
	}
	v := o.OrElseGet(f)
	require.Equal(t, "aaa", v)
	require.True(t, called)

	called = false
	o = Of("bbb")
	v = o.OrElseGet(f)
	require.Equal(t, "bbb", v)
	require.False(t, called)
}

func TestOptional_OrElseError(t *testing.T) {
	o := Empty[string]()
	called := false
	f := func() error {
		called = true
		return errors.New("not there")
	}
	err := o.OrElseError(f)
	require.True(t, called)
	require.Error(t, err)
	require.Equal(t, "not there", err.Error())

	called = false
	o = Of("abc")
	err = o.OrElseError(f)
	require.False(t, called)
	require.NoError(t, err)
}

func TestOptional_Filter(t *testing.T) {
	o := Empty[string]()
	called := false
	filterOk := true
	f := func(v string) bool {
		called = true
		return filterOk
	}
	o2 := o.Filter(f)
	require.False(t, called)
	require.False(t, o2.IsPresent())

	o = Of("aaa")
	o2 = o.Filter(f)
	require.True(t, called)
	require.True(t, o2.IsPresent())

	called = false
	filterOk = false
	o2 = o.Filter(f)
	require.True(t, called)
	require.False(t, o2.IsPresent())
}

func TestOptional_Map(t *testing.T) {
	o := Of("123")
	called := false
	f := func(v string) any {
		called = true
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil
		}
		return int(i)
	}
	o2 := o.Map(f)
	require.True(t, called)
	require.True(t, o2.present)
	v, err := o2.Get()
	require.NoError(t, err)
	require.Equal(t, 123, v)

	o = Empty[string]()
	called = false
	o2 = o.Map(f)
	require.False(t, called)
	require.False(t, o2.present)
}

func TestOptional_MarshalUnmarshalJSON(t *testing.T) {
	type aStruct struct {
		Foo Optional[string]  `json:"foo"`
		Bar Optional[int]     `json:"bar"`
		Baz Optional[float64] `json:"baz"`
	}
	myA := &aStruct{
		Foo: Of("aaa"),
		Bar: Of(1),
		Baz: Of(1.2),
	}
	data, err := json.Marshal(myA)
	require.NoError(t, err)
	require.Equal(t, `{"foo":"aaa","bar":1,"baz":1.2}`, string(data[:]))

	myA2 := &aStruct{}
	err = json.Unmarshal(data, myA2)
	require.NoError(t, err)
	require.True(t, myA2.Foo.IsPresent())
	require.True(t, myA2.Bar.IsPresent())
	require.True(t, myA2.Baz.IsPresent())
	fooV, err := myA2.Foo.Get()
	require.NoError(t, err)
	require.Equal(t, "aaa", fooV)
	barV, err := myA2.Bar.Get()
	require.NoError(t, err)
	require.Equal(t, 1, barV)
	bazV, err := myA2.Baz.Get()
	require.NoError(t, err)
	require.Equal(t, 1.2, bazV)

	str := `{"foo":null,"bar":null,"baz":null}`
	myA3 := &aStruct{}
	err = json.Unmarshal([]byte(str), myA3)
	require.NoError(t, err)
	require.False(t, myA3.Foo.IsPresent())
	require.False(t, myA3.Bar.IsPresent())
	require.False(t, myA3.Baz.IsPresent())

	data, err = json.Marshal(myA3)
	require.NoError(t, err)
	require.Equal(t, `{"foo":null,"bar":null,"baz":null}`, string(data[:]))

	str = `{"foo":1.2,"bar":null,"baz":null}`
	err = json.Unmarshal([]byte(str), myA3)
	require.Error(t, err)
}
