# Gopt
[![GoDoc](https://godoc.org/github.com/go-andiamo/gopt?status.svg)](https://pkg.go.dev/github.com/go-andiamo/gopt)
[![Latest Version](https://img.shields.io/github/v/tag/go-andiamo/gopt.svg?sort=semver&style=flat&label=version&color=blue)](https://github.com/go-andiamo/gopt/releases)
[![codecov](https://codecov.io/gh/go-andiamo/gopt/branch/main/graph/badge.svg?token=igjnZdgh0e)](https://codecov.io/gh/go-andiamo/gopt)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-andiamo/gopt)](https://goreportcard.com/report/github.com/go-andiamo/gopt)

A very light Optional implementation in Golang

## Installation
To install Gopt, use go get:

    go get github.com/go-andiamo/gopt

To update Gopt to the latest version, run:

    go get -u github.com/go-andiamo/gopt

## Examples
```go
package main

import (
    . "github.com/go-andiamo/gopt"
)

func main() {
    optFlt := Of[float64](1.23)
    println(optFlt.IsPresent())
    println(optFlt.OrElse(-1))

    opt2 := Empty[float64]()
    println(opt2.IsPresent())
    println(opt2.OrElse(-1))

    opt2.OrElseSet(10)
    println(opt2.IsPresent())
    println(opt2.OrElse(-1))
}
```
[try on go-playground](https://go.dev/play/p/U0dKTrGlG-e)

Optionals can also be very useful when used in conjunction with JSON unmarshalling - to determine whether unmarshalled properties were actually present, with a valid value, in the JSON.  The following code demonstrates...
```go
package main

import (
    "encoding/json"
    "fmt"

    . "github.com/go-andiamo/gopt"
)

type NormalStruct struct {
    Foo string
    Bar int
    Baz float64
}

type OptsStruct struct {
    Foo Optional[string]
    Bar Optional[int]
    Baz Optional[float64]
}

func main() {
    jdata := `{"Foo": null, "Bar": 1}`

    normal := &NormalStruct{}
    err := json.Unmarshal([]byte(jdata), normal)
    if err == nil {
        // was property Foo set???
        fmt.Printf("'Foo' was set to \"%s\"???\n", normal.Foo) // was it really?
        // was property Bar actually set to 0???
        fmt.Printf("'Bar' was set to \"%d\"???\n", normal.Bar) // was it really?
        // was property Baz actually set to 0.000???
        fmt.Printf("'Baz' was set to \"%f\"???\n", normal.Baz) // was it really?
    } else {
        println(err.Error())
    }

    println()
    // now try with optionals...
    opts := &OptsStruct{}
    err = json.Unmarshal([]byte(jdata), opts)
    if err == nil {
        opts.Foo.IfSetOtherwise(
            func(v string) {
                fmt.Printf("'Foo' was set to \"%s\"\n", v)
            },
            func() {
                println("'Foo' was set but not to a valid value")
            },
            func() {
                println("'Foo' was not set at all")
            },
        )
        opts.Bar.IfSetOtherwise(
            func(v int) {
                fmt.Printf("'Bar' was set to %d\n", v)
            },
            func() {
                println("'Bar' was set but not to a valid value")
            },
            func() {
                println("'Bar' was not set at all")
            },
        )
        opts.Baz.IfSetOtherwise(
            func(v float64) {
                fmt.Printf("'Baz' was set to %f\n", v)
            },
            func() {
                println("'Baz' was set but not to a valid value")
            },
            func() {
                println("'Baz' was not set at all")
            },
        )
    } else {
        println(err.Error())
    }
}
```
[try on go-playground](https://go.dev/play/p/63eC1AJ3Qgn)

## Methods
<table>
    <tr>
        <th>Method and description</th>
        <th>Returns</th>
    </tr>
    <tr>
        <td>
            <code>AsEmpty()</code><br>
            returns a new empty optional of the same type
        </td>
        <td><code>*Optional[T]</code></td>
    </tr>
    <tr>
        <td>
            <code>Clear()</code><br>
            clears the optional<br>
            Clearing sets the present to false, the set flag to false and the value to an empty value
        </td>
        <td><code>*Optional[T]</code></td>
    </tr>
    <tr>
        <td>
            <code>Filter(f func(v T) bool)</code><br>
            if the value is present and calling the supplied filter function returns true, returns a new optional describing the value<br>
            Otherwise returns an empty optional
        </td>
        <td><code>*Optional[T]</code></td>
    </tr>
    <tr>
        <td>
            <code>Get()</code><br>
            returns the value and an error if the value is not present
        </td>
        <td><code>(T, error)</code></td>
    </tr>
    <tr>
        <td>
            <code>GetOk()</code><br>
            returns the value and true if the value is present<br>
            otherwise returns an empty value and false<br>
        </td>
        <td><code>(T, bool)</code></td>
    </tr>
    <tr>
        <td>
            <code>IfElse(condition bool, other T)</code><br>
            if the supplied condition is true and the value is present, returns the value<br>
            otherwise the other value is returned
        </td>
        <td><code>T</code></td>
    </tr>
    <tr>
        <td>
            <code>IfPresent(f func(v T))</code><br>
            if the value is present, calls the supplied function with the value, otherwise does nothing<br>
            <em>returns the original optional</em>
        </td>
        <td><code>*Optional[T]</code></td>
    </tr>
    <tr>
        <td>
            <code>IfPresentOtherwise(f func(v T), other func())</code><br>
            if the value is present, calls the supplied function with the value, otherwise calls the other function<br>
            <em>returns the original optional</em>
        </td>
        <td><code><code>*Optional[T]</code></code></td>
    </tr>
    <tr>
        <td>
            <code>IfSet(f func(v T), notPresent func())</code><br>
            if the value was set and is present, calls the supplied function with the value<br>
            if the value was set but is not present, calls the supplied notPresent function<br>
            otherwise, does nothing<br>
            <em>returns the original optional</em>
        </td>
        <td><code>*Optional[T]</code></td>
    </tr>
    <tr>
        <td>
            <code>IfSetOtherwise(f func(v T), notPresent func(), other func())</code><br>
            if the value was set and is present, calls the supplied function with the value<br>
            if the value was set but is not present, calls the supplied notPresent function<br>
            otherwise, calls the other func<br>
            <em>returns the original optional</em>
        </td>
        <td><code><code>*Optional[T]</code></code></td>
    </tr>
    <tr>
        <td>
            <code>IsPresent()</code><br>
            returns true if the value is present, otherwise false
        </td>
        <td><code>bool</code></td>
    </tr>
    <tr>
        <td>
            <code>Map(f func(v T) any)</code><br>
            if the value is present and the result of calling the supplied mapping function returns non-nil, returns
            an optional describing that returned value<br>
            Otherwise returns an empty optional
        </td>
        <td><code>*Optional[any]</code></td>
    </tr>
    <tr>
        <td>
            <code>MarshalJSON()</code><br>
            implements JSON marshal<br>
            if the value is present, returns the marshalled data for the value<br>
            Otherwise, returns the marshalled data for null
        </td>
        <td><code>([]byte, error)</code></td>
    </tr>
    <tr>
        <td>
            <code>OrElse(other T)</code><br>
            returns the value if present, otherwise returns other
        </td>
        <td><code>T</code></td>
    </tr>
    <tr>
        <td>
            <code>OrElseError(err error)</code><br>
            returns the supplied error if the value is not present, otherwise returns nil
        </td>
        <td><code>error</code></td>
    </tr>
    <tr>
        <td>
            <code>OrElseGet(f func() T)</code><br>
            returns the value if present, otherwise returns the result of calling the supplied function
        </td>
        <td><code>T</code></td>
    </tr>
    <tr>
        <td>
            <code>OrElsePanic(v any)</code><br>
            if the value is not present, panics with the supplied value, otherwise does nothing<br>
            <em>returns the original optional</em>
        </td>
        <td><code>*Optional[T]</code></td>
    </tr>
    <tr>
        <td>
            <code>OrElseSet(v T)</code><br>
            if the value is not present it is set to the supplied value
        </td>
        <td><code>*Optional[T]</code></td>
    </tr>
    <tr>
        <td>
            <code>Scan(value interface{})</code><br>
            implements sql.Scan
        </td>
        <td><code>error</code></td>
    </tr>
    <tr>
        <td>
            <code>UnSet()</code><br>
            clears the set flag (see <code>WasSet()</code>)<br>
            <em>returns the original optional</em>
        </td>
        <td><code>*Optional[T]</code></td>
    </tr>
    <tr>
        <td>
            <code>UnmarshalJSON(data []byte)</code><br>
            implements JSON unmarshal<br>
            if the supplied data is null representation, sets the present to false<br>
            Otherwise, unmarshal the data as the value and sets the optional to present (unless the result of
            unmarshalling the value returns an error - in which case the present is set to false)
        </td>
        <td><code>error</code></td>
    </tr>
    <tr>
        <td>
            <code>WasSet()</code><br>
            returns true if the last setting operation set the value, otherwise false<br>
            Setting operations are <code>UnmarshalJSON()</code>, <code>Scan()</code> and <code>OrElseSet()</code><br>
            Use method <code>UnSet()</code> to clear this flag alone
        </td>
        <td><code>bool</code></td>
    </tr>
    <tr>
        <td>
            <code>WasSetElse(other T)</code><br>
            returns the value if present and set, otherwise returns other
        </td>
        <td><code>T</code></td>
    </tr>
    <tr>
        <td>
            <code>WasSetElseError(err error)</code><br>
            returns the supplied error if the value is not present and set, otherwise returns nil<br>
            if the supplied error is nil and the value is not present and set, a <code>NotPresentError</code> is returned
        </td>
        <td><code>error</code></td>
    </tr>
    <tr>
        <td>
            <code>WasSetElseGet(f func() T)</code><br>
            returns the value if present and set, otherwise returns the result of calling the supplied function<br>
            if the supplied function is nil and the value is not present and set, returns a default empty value
        </td>
        <td><code>T</code></td>
    </tr>
    <tr>
        <td>
            <code>WasSetElsePanic(v any)</code><br>
            if the value is not present and set, panics with the supplied value, otherwise does nothing<br>
            <em>returns the original optional</em>
        </td>
        <td><code>*Optional[T]</code></td>
    </tr>
    <tr>
        <td>
            <code>WasSetElseSet(v T)</code><br>
            if the value is not present and set it is set to the value supplied
        </td>
        <td><code>*Optional[T]</code></td>
    </tr>
</table>

## Constructors
<table>
    <tr>
        <th>Constructor function and description</th>
    </tr>
    <tr>
        <td>
            <code>Of[T any](value T) *Optional[T]</code><br>
            Creates a new optional with the supplied value
        </td>
    </tr>
    <tr>
        <td>
            <code>OfNillable[T any](value T) *Optional[T]</code><br>
            Creates a new optional with the supplied value<br>
            If the supplied value is nil, an empty (not present) optional is returned
        </td>
    </tr>
    <tr>
        <td>
            <code>OfNillableString(value string) *Optional[string]</code><br>
            Creates a new string optional with the supplied value<br>
            If the supplied value is an empty string, an empty (not-present) optional is returned
        </td>
    </tr>
    <tr>
        <td>
            <code>Empty[T any]() *Optional[T]</code><br>
            Creates a new empty (not-present) optional of the specified type
        </td>
    </tr>
    <tr>
        <td>
            <code>EmptyString() *Optional[string]</code><br>
            returns an empty optional of type <code>string</code>
        </td>
    </tr>
    <tr>
        <td>
            <code>EmptyInterface() *Optional[interface{}]</code><br>
            returns an empty optional of type <code>interface{}</code>
        </td>
    </tr>
    <tr>
        <td>
            <code>EmptyInt() *Optional[int]</code><br>
            returns an empty optional of type <code>int</code>
        </td>
    </tr>
    <tr>
        <td>
            <code>EmptyInt8() *Optional[int8]</code><br>
7            returns an empty optional of type <code>int8</code>
        </td>
    </tr>
    <tr>
        <td>
            <code>EmptyInt16() *Optional[int16]</code><br>
            returns an empty optional of type <code>int16</code>
        </td>
    </tr>
    <tr>
        <td>
            <code>EmptyInt32() *Optional[int32]</code><br>
            returns an empty optional of type <code>int32</code>
        </td>
    </tr>
    <tr>
        <td>
            <code>EmptyInt64() *Optional[int64]</code><br>
            returns an empty optional of type <code>int64</code>
        </td>
    </tr>
    <tr>
        <td>
            <code>EmptyUint() *Optional[uint]</code><br>
            returns an empty optional of type <code>uint</code>
        </td>
    </tr>
    <tr>
        <td>
            <code>EmptyUint8() *Optional[uint8]</code><br>
            returns an empty optional of type <code>uint8</code>
        </td>
    </tr>
    <tr>
        <td>
            <code>EmptyUint16() *Optional[uint16]</code><br>
            returns an empty optional of type <code>uint16</code>
        </td>
    </tr>
    <tr>
        <td>
            <code>EmptyUint32() *Optional[uint32]</code><br>
            returns an empty optional of type <code>uint32</code>
        </td>
    </tr>
    <tr>
        <td>
            <code>func EmptyUint64() *Optional[uint64]</code><br>
            returns an empty optional of type <code>uint64</code>
        </td>
    </tr>
    <tr>
        <td>
            <code>EmptyBool() *Optional[bool]</code><br>
            returns an empty optional of type <code>bool</code>
        </td>
    </tr>
    <tr>
        <td>
            <code>EmptyFloat32() *Optional[float32]</code><br>
            returns an empty optional of type <code>float32</code>
        </td>
    </tr>
    <tr>
        <td>
            <code>EmptyFloat64() *Optional[float64]</code><br>
            returns an empty optional of type <code>float64</code>
        </td>
    </tr>
    <tr>
        <td>
            <code>EmptyByte() *Optional[byte]</code><br>
            returns an empty optional of type <code>byte</code>
        </td>
    </tr>
    <tr>
        <td>
            <code>EmptyRune() *Optional[rune]</code><br>
            returns an empty optional of type <code>rune</code>
        </td>
    </tr>
</table>
