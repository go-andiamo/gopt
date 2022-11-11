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

## Methods
<table>
    <tr>
        <th>Method and description</th>
        <th>Returns</th>
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
            <code>AsEmpty()</code><br>
            returns a new empty optional of the same type
        </td>
        <td><code>Optional[T]</code></td>
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
            <code>IfPresent(f func(v T))</code><br>
            if the value is present, calls the supplied function with the value, otherwise does nothing<br>
            <em>returns the original optional</em>
        </td>
        <td><code>Optional[T]</code></td>
    </tr>
    <tr>
        <td>
            <code>IfPresentOtherwise(f func(v T), other func())</code><br>
            if the value is present, calls the supplied function with the value, otherwise calls the other function<br>
            <em>returns the original optional</em>
        </td>
        <td><code><code>Optional[T]</code></code></td>
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
            <code>OrElseGet(f func() T)</code><br>
            returns the value if present, otherwise returns the result of calling the supplied function
        </td>
        <td><code>T</code></td>
    </tr>
    <tr>
        <td>
            <code>OrElseSet(v T)</code><br>
            if the value is not present it is set to the supplied value
        </td>
        <td><code>Optional[T]</code></td>
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
            <code>OrElsePanic(v any)</code><br>
            if the value is not present, panics with the supplied value, otherwise does nothing
        </td>
        <td><em>nothing</em></td>
    </tr>
    <tr>
        <td>
            <code>DoWith(f func(v T))</code><br>
            if the value is present, calls the supplied function with the value<br>
            <em>returns the original optional</em>
        </td>
        <td><code>Optional[T]</code></td>
    </tr>
    <tr>
        <td>
            <code>Filter(f func(v T) bool)</code><br>
            if the value is present and calling the supplied filter function returns true, returns a new optional describing the value<br>
            Otherwise returns an empty optional
        </td>
        <td><code>Optional[T]</code></td>
    </tr>
    <tr>
        <td>
            <code>Map(f func(v T) any)</code><br>
            if the value is present and the result of calling the supplied mapping function returns non-nil, returns
            an optional describing that returned value<br>
            Otherwise returns an empty optional
        </td>
        <td><code>Optional[any]</code></td>
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
            <code>Scan(value interface{})</code><br>
            implements sql.Scan
        </td>
        <td><code>error</code></td>
    </tr>
</table>

## Constructors
<table>
    <tr>
        <th>Constructor function and description</th>
    </tr>
    <tr>
        <td>
            <code>Of[T any](value T) Optional[T]</code><br>
            Creates a new optional with the supplied value
        </td>
    </tr>
    <tr>
        <td>
            <code>OfNillable[T any](value T) Optional[T]</code><br>
            Creates a new optional with the supplied value<br>
            If the supplied value is nil, an empty (not present) optional is returned
        </td>
    </tr>
    <tr>
        <td>
            <code>OfNillableString(value string) Optional[string]</code><br>
            Creates a new string optional with the supplied value<br>
            If the supplied value is an empty string, an empty (not-present) optional is returned
        </td>
    </tr>
    <tr>
        <td>
            <code>Empty[T any]() Optional[T]</code><br>
            Creates a new empty (not-present) optional of the specified type
        </td>
    </tr>
</table>
