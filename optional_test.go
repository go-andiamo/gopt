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
	opt2 := Of(x).AsEmpty()
	require.False(t, opt2.IsPresent())
	_, err = opt2.Get()
	require.Error(t, err)

	var v *myStruct
	v = opt2.OrElse(&myStruct{
		Foo: "bar",
	})
	require.Equal(t, "bar", v.Foo)

	opt3 := Of("").AsEmpty()
	require.False(t, opt3.IsPresent())
	_, err = opt3.Get()
	require.Error(t, err)

	opt4 := OfNillable[String](nil)
	require.False(t, opt4.IsPresent())
}

func TestOfNillableString(t *testing.T) {
	opt := OfNillableString("")
	require.False(t, opt.IsPresent())
	opt = OfNillableString("foo")
	require.True(t, opt.IsPresent())
	opt = Of("aaa")
	require.True(t, opt.IsPresent())
	opt = Of("")
	require.True(t, opt.IsPresent())
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

func TestOptional_AsEmpty(t *testing.T) {
	o := Of("abc")
	require.True(t, o.IsPresent())
	oe := o.AsEmpty()
	require.False(t, oe.IsPresent())

	o2 := Of(0)
	require.True(t, o2.IsPresent())
	oe2 := o.AsEmpty()
	require.False(t, oe2.IsPresent())
}

func TestOptional_Clear(t *testing.T) {
	o := EmptyString()
	require.False(t, o.IsPresent())
	require.False(t, o.WasSet())
	v, err := o.Get()
	require.Error(t, err)
	require.Equal(t, "", v)

	o.OrElseSet("foo")
	require.True(t, o.IsPresent())
	require.True(t, o.WasSet())
	v, err = o.Get()
	require.NoError(t, err)
	require.Equal(t, "foo", v)

	o.Clear()
	require.False(t, o.IsPresent())
	require.False(t, o.WasSet())
	v, err = o.Get()
	require.Error(t, err)
	require.Equal(t, "", v)
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

func TestOptional_Get(t *testing.T) {
	o := EmptyString()
	_, err := o.Get()
	require.Error(t, err)
	require.Equal(t, NotPresent, err)

	o = Of("aaa")
	v, err := o.Get()
	require.NoError(t, err)
	require.Equal(t, "aaa", v)
}

func TestOptional_GetOk(t *testing.T) {
	o := EmptyString()
	v, ok := o.GetOk()
	require.False(t, ok)
	require.Equal(t, "", v)

	o = Of("aaa")
	v, ok = o.GetOk()
	require.True(t, ok)
	require.Equal(t, "aaa", v)
}

func TestOptional_IfElse(t *testing.T) {
	o := Empty[string]()
	v := o.IfElse(true, "abc")
	require.Equal(t, "abc", v)

	o = Of("xyz")
	v = o.IfElse(false, "abc")
	require.Equal(t, "abc", v)
	v = o.IfElse(true, "abc")
	require.Equal(t, "xyz", v)
	v = o.IfElse(o.WasSet(), "abc")
	require.Equal(t, "abc", v)
	err := o.Scan("scanned")
	require.NoError(t, err)
	v = o.IfElse(o.WasSet(), "abc")
	require.Equal(t, "scanned", v)
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

func TestOptional_IfPresentOtherwise(t *testing.T) {
	calledPresent := false
	collected := ""
	calledOther := false
	f := func(v string) {
		calledPresent = true
		collected = v
	}
	oth := func() {
		calledOther = true
	}

	o := Of("aaa")
	o.IfPresentOtherwise(f, oth)
	require.True(t, calledPresent)
	require.Equal(t, "aaa", collected)
	require.False(t, calledOther)

	calledPresent = false
	calledOther = false
	o = o.AsEmpty()
	o.IfPresentOtherwise(f, oth)
	require.False(t, calledPresent)
	require.True(t, calledOther)
}

func TestOptional_IfSet(t *testing.T) {
	o := EmptyString()
	o.OrElseSet("aaa")
	setCalled := false
	notPresentCalled := false
	value := ""
	setFn := func(v string) {
		setCalled = true
		value = v
	}
	notPresentFn := func() {
		notPresentCalled = true
	}
	o.IfSet(setFn, notPresentFn)
	require.True(t, setCalled)
	require.Equal(t, "aaa", value)
	require.False(t, notPresentCalled)

	setCalled = false
	notPresentCalled = false
	o.Clear()
	err := o.UnmarshalJSON([]byte(`null`))
	require.NoError(t, err)
	o.IfSet(setFn, notPresentFn)
	require.False(t, setCalled)
	require.True(t, notPresentCalled)
}

func TestOptional_IfSetOtherwise(t *testing.T) {
	type oStruct struct {
		Foo Optional[string]
	}
	setCalled := false
	notPresentCalled := false
	otherCalled := false
	setFn := func(v string) {
		setCalled = true
	}
	notPresentFn := func() {
		notPresentCalled = true
	}
	otherFn := func() {
		otherCalled = true
	}

	strc := &oStruct{}
	err := json.Unmarshal([]byte(`{}`), strc)
	require.NoError(t, err)
	strc.Foo.IfSetOtherwise(setFn, notPresentFn, otherFn)
	require.False(t, setCalled)
	require.False(t, notPresentCalled)
	require.True(t, otherCalled)

	setCalled = false
	notPresentCalled = false
	otherCalled = false
	strc = &oStruct{}
	err = json.Unmarshal([]byte(`{"Foo":null}`), strc)
	require.NoError(t, err)
	strc.Foo.IfSetOtherwise(setFn, notPresentFn, otherFn)
	require.False(t, setCalled)
	require.True(t, notPresentCalled)
	require.False(t, otherCalled)

	setCalled = false
	notPresentCalled = false
	otherCalled = false
	strc = &oStruct{}
	err = json.Unmarshal([]byte(`{"Foo":"abc"}`), strc)
	require.NoError(t, err)
	strc.Foo.IfSetOtherwise(setFn, notPresentFn, otherFn)
	require.True(t, setCalled)
	require.False(t, notPresentCalled)
	require.False(t, otherCalled)
}

func TestOptional_IsPresent(t *testing.T) {
	o := Of("aaa")
	require.True(t, o.IsPresent())
	o = Empty[string]()
	require.False(t, o.IsPresent())
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
		Foo: *Of("aaa"),
		Bar: *Of(1),
		Baz: *Of(1.2),
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
	require.True(t, myA3.Foo.WasSet())
	require.False(t, myA3.Bar.IsPresent())
	require.True(t, myA3.Bar.WasSet())
	require.False(t, myA3.Baz.IsPresent())
	require.True(t, myA3.Baz.WasSet())

	data, err = json.Marshal(myA3)
	require.NoError(t, err)
	require.Equal(t, `{"foo":null,"bar":null,"baz":null}`, string(data[:]))

	str = `{"foo":null}`
	myA3 = &aStruct{}
	err = json.Unmarshal([]byte(str), myA3)
	require.NoError(t, err)
	require.False(t, myA3.Foo.IsPresent())
	require.True(t, myA3.Foo.WasSet())
	myA3.Foo.UnSet()
	require.False(t, myA3.Foo.WasSet())
	require.False(t, myA3.Bar.IsPresent())
	require.False(t, myA3.Bar.WasSet())
	require.False(t, myA3.Baz.IsPresent())
	require.False(t, myA3.Baz.WasSet())

	str = `{"foo":1.2,"bar":null,"baz":null}`
	err = json.Unmarshal([]byte(str), myA3)
	require.Error(t, err)
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

	o = Empty[string]()
	var ef func() string
	v = o.OrElseGet(ef)
	require.Equal(t, "", v)
}

func TestOptional_OrElseError(t *testing.T) {
	o := Empty[string]()
	err := o.OrElseError(errors.New("not there"))
	require.Error(t, err)
	require.Equal(t, "not there", err.Error())

	err = o.OrElseError(nil)
	require.Equal(t, NotPresent, err)

	o = Of("abc")
	err = o.OrElseError(errors.New("not there"))
	require.NoError(t, err)
}

func TestOptional_OrElsePanic(t *testing.T) {
	o := Of("str")
	o.OrElsePanic("whoops")
	o = Empty[string]()
	require.Panics(t, func() {
		o.OrElsePanic("whoops")
	})
}

func TestOptional_OrElseSet(t *testing.T) {
	o := Empty[map[string]interface{}]()
	require.False(t, o.IsPresent())

	o2 := o.OrElseSet(map[string]interface{}{})
	require.Equal(t, o, o2)
	require.True(t, o2.IsPresent())
	require.True(t, o.IsPresent())

	o = Empty[map[string]interface{}]()
	require.False(t, o.IsPresent())
	o2 = o.OrElseSet(nil)
	require.Equal(t, o, o2)
	require.False(t, o2.IsPresent())
	require.False(t, o.IsPresent())
}

func TestOptional_Scan(t *testing.T) {
	var o Optional[string]
	err := o.Scan("str")
	require.NoError(t, err)
	require.True(t, o.IsPresent())
	require.Equal(t, "str", o.OrElse("other"))
	err = o.Scan(nil)
	require.NoError(t, err)
	require.False(t, o.IsPresent())

	o2 := OfNillable(map[string]interface{}{}).AsEmpty()
	require.False(t, o2.IsPresent())
	err = o2.Scan([]byte(`{"foo":"bar"}`))
	require.NoError(t, err)
	require.True(t, o2.IsPresent())
	ov2, err := o2.Get()
	require.NoError(t, err)
	require.Equal(t, 1, len(ov2))

	o3 := OfNillable[map[string]interface{}](nil)
	require.False(t, o3.IsPresent())
	err = o3.Scan(nil)
	require.NoError(t, err)
	require.False(t, o3.IsPresent())
	err = o3.Scan("")
	require.NoError(t, err)
	require.False(t, o3.IsPresent())
	o3.UnSet()
	err = o3.Scan([]byte(`["foo","bar"]`))
	require.Error(t, err)
	require.False(t, o3.IsPresent())
	require.True(t, o3.WasSet())
	o3.UnSet()
	err = o3.Scan([]byte(`null`))
	require.NoError(t, err)
	require.False(t, o3.IsPresent())
	require.True(t, o3.WasSet())

	o4 := OfNillable[*scannable](nil)
	require.False(t, o4.present)
	err = o4.Scan("abc")
	require.NoError(t, err)
	o4v := o4.OrElse(nil)
	require.NotNil(t, o4v)
	require.True(t, o4v.called)
	require.Equal(t, "abc", o4v.value)

	o4 = Of(&scannable{err: errors.New("fooey")})
	require.True(t, o4.present)
	err = o4.Scan("abc")
	require.Error(t, err)

	o4.UnSet()
	err = o4.Scan(nil)
	require.NoError(t, err)
	require.False(t, o4.IsPresent())
	require.True(t, o4.WasSet())

	o5 := EmptyInterface()
	require.False(t, o5.IsPresent())
	require.False(t, o5.WasSet())
	err = o5.Scan(nil)
	require.NoError(t, err)
	require.False(t, o5.IsPresent())
	require.True(t, o5.WasSet())
}

func TestOptional_UnSet(t *testing.T) {
	o := Of("abc")
	require.True(t, o.IsPresent())
	require.False(t, o.WasSet())
	err := o.Scan("xyz")
	require.NoError(t, err)
	require.True(t, o.WasSet())
	o.UnSet()
	require.False(t, o.WasSet())
}

func TestOptional_WasSet(t *testing.T) {
	o := Of("abc")
	require.True(t, o.IsPresent())
	require.False(t, o.WasSet())

	err := o.Scan("xyz")
	require.NoError(t, err)
	require.True(t, o.WasSet())
}

func TestOptional_WasSetElse(t *testing.T) {
	o := Of("abc")
	require.True(t, o.IsPresent())
	require.False(t, o.WasSet())
	v := o.WasSetElse("def")
	require.Equal(t, "def", v)
	err := o.Scan("xyz")
	require.NoError(t, err)
	v = o.WasSetElse("def")
	require.Equal(t, "xyz", v)
}

func TestOptional_WasSetElseError(t *testing.T) {
	o := Of("abc")
	require.True(t, o.IsPresent())
	require.False(t, o.WasSet())
	err := o.WasSetElseError(errors.New("fooey"))
	require.Error(t, err)
	require.Equal(t, "fooey", err.Error())
	err = o.WasSetElseError(nil)
	require.Error(t, err)
	require.Equal(t, NotPresent, err)
	err = o.Scan("xyz")
	require.NoError(t, err)
	err = o.WasSetElseError(errors.New("fooey"))
	require.NoError(t, err)
}

func TestOptional_WasSetElseGet(t *testing.T) {
	o := Of("abc")
	require.True(t, o.IsPresent())
	require.False(t, o.WasSet())
	called := false
	f := func() string {
		called = true
		return "aaa"
	}
	v := o.WasSetElseGet(f)
	require.True(t, called)
	require.Equal(t, "aaa", v)

	called = false
	v = o.WasSetElseGet(nil)
	require.False(t, called)
	require.Equal(t, "", v)

	called = false
	err := o.Scan("xyz")
	require.NoError(t, err)
	v = o.WasSetElseGet(f)
	require.False(t, called)
	require.Equal(t, "xyz", v)
}

func TestOptional_WasSetElsePanic(t *testing.T) {
	o := Of("abc")
	require.True(t, o.IsPresent())
	require.False(t, o.WasSet())
	require.Panics(t, func() {
		o.WasSetElsePanic("fooey")
	})
	err := o.Scan("xyz")
	require.NoError(t, err)
	o = o.WasSetElsePanic("fooey")
	require.True(t, o.IsPresent())
	require.True(t, o.WasSet())
}

func TestOptional_WasSetElseSet(t *testing.T) {
	o := Of("abc")
	require.True(t, o.IsPresent())
	require.False(t, o.WasSet())
	v, err := o.Get()
	require.NoError(t, err)
	require.Equal(t, "abc", v)
	o2 := o.WasSetElseSet("xyz")
	require.True(t, o2.IsPresent())
	require.True(t, o2.WasSet())
	v, err = o2.Get()
	require.NoError(t, err)
	require.Equal(t, "xyz", v)

	my := &scannable{}
	o3 := Of(my)
	require.True(t, o3.IsPresent())
	require.False(t, o3.WasSet())
	o4 := o3.WasSetElseSet(nil)
	require.False(t, o4.IsPresent())
	require.True(t, o4.WasSet())
}

type scannable struct {
	called bool
	err    error
	value  any
}

func (s *scannable) Scan(src any) error {
	s.called = true
	s.value = src
	return s.err
}

func TestEmpties(t *testing.T) {
	require.False(t, EmptyString().IsPresent())
	require.False(t, EmptyInterface().IsPresent())
	require.False(t, EmptyInt().IsPresent())
	require.False(t, EmptyInt8().IsPresent())
	require.False(t, EmptyInt16().IsPresent())
	require.False(t, EmptyInt32().IsPresent())
	require.False(t, EmptyInt64().IsPresent())
	require.False(t, EmptyUint().IsPresent())
	require.False(t, EmptyUint8().IsPresent())
	require.False(t, EmptyUint16().IsPresent())
	require.False(t, EmptyUint32().IsPresent())
	require.False(t, EmptyUint64().IsPresent())
	require.False(t, EmptyBool().IsPresent())
	require.False(t, EmptyFloat32().IsPresent())
	require.False(t, EmptyFloat64().IsPresent())
	require.False(t, EmptyByte().IsPresent())
	require.False(t, EmptyRune().IsPresent())
}
