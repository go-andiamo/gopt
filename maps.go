package gopt

import (
	"strconv"
	"strings"
)

// Converter is function(s) that can be passed to Extract, ExtractJson and ExtractJsonPath to convert the value found to the required type
//
// Each converter passed is called successively until one returns true as the second return arg
type Converter[T any] func(value any) (T, bool)

// Extract extracts an optional value, of the specified type, from a map
//
// If the key is present (and the value is non-nil and of the specified type) then an optional with the value is returned, otherwise an empty optional is returned
func Extract[K comparable, T any](m map[K]any, key K, converters ...Converter[T]) *Optional[T] {
	result := Empty[T]()
	if rv, ok := m[key]; ok && isPresent(rv) {
		if v, ok := rv.(T); ok {
			result = Of[T](v)
		} else if v, ok := runConverters(rv, converters...); ok {
			result = Of[T](v)
		}
	}
	return result
}

// ExtractJson extracts an optional value, of the specified type, from a map[string]any
//
// If the key is present (and the value is non-nil and of the specified type) then an optional with the value is returned, otherwise an empty optional is returned
func ExtractJson[T any](m map[string]any, key string, converters ...Converter[T]) *Optional[T] {
	result := Empty[T]()
	if rv, ok := m[key]; ok && isPresent(rv) {
		if v, ok := rv.(T); ok {
			result = Of[T](v)
		} else if v, ok := runConverters(rv, converters...); ok {
			result = Of[T](v)
		}
	}
	return result
}

// ExtractJsonPath extracts an optional value, of the specified type, from a map[string]any by traversing the supplied JSON path
//
// If the key is present (and the value is non-nil and of the specified type) then an optional with the value is returned, otherwise an empty optional is returned
//
// The supplied JSON path is a string path with parts separated by "." - where array properties can be indexed using notation "property[n]"
//
// The index n may be positive (or zero) or negative (indicating relative index to the end).  For example, given a map of:
//   m := map[string]any{
//     "foo": map[string]any{
//       "bar": []any{
//         map[string]any{
//           "baz": "X",
//         },
//         map[string]any{
//           "baz": "Y",
//         },
//       }
//     }
//   }
// then using:
//  o, _ := ExtractJsonPath[string](m, "foo.bar[0].baz")
// would yield a present Optional with value "X"
//
// or using:
//  o, _ := ExtractJsonPath[string](m, "foo.bar[-1].baz")
// would yield a present Optional with value "Y"
//
//
// The second result arg of ExtractJsonPath is a slice of bools indicating the whether the path items existed - if the returned
// Optional is not present it may be that the final property of the path was not found (or the incorrect type) or that the path was not found...
// if the path was not found, then the final element in that second slice return arg will be false
func ExtractJsonPath[T any](m map[string]any, path string, converters ...Converter[T]) (*Optional[T], []bool) {
	parts := strings.Split(path, ".")
	l := len(parts)
	if l == 1 && !(strings.HasSuffix(parts[0], "]") && strings.Contains(parts[0], "[")) {
		return ExtractJson[T](m, parts[0]), nil
	}
	result := Empty[T]()
	pathPresent := make([]bool, 0, l-1)
	curr := m
	ok := false
	for i, part := range parts {
		if pty, isIndexed, idx := getProperty(part); isIndexed {
			if curr, result, pathPresent, ok = extractPathIndexed[T](pty, idx, i == l-1, curr, result, pathPresent, converters...); !ok {
				break
			}
		} else if curr, result, pathPresent, ok = extractPathProperty[T](pty, i == l-1, curr, result, pathPresent, converters...); !ok {
			break
		}
	}
	return result, pathPresent
}

func extractPathProperty[T any](pty string, last bool, curr map[string]any, result *Optional[T], pathPresent []bool, converters ...Converter[T]) (map[string]any, *Optional[T], []bool, bool) {
	isOk := false
	if rv, ok := curr[pty]; ok {
		if last {
			if av, ok := rv.(T); ok && isPresent(rv) {
				isOk = true
				result = Of[T](av)
			} else if cv, ok := runConverters[T](rv, converters...); ok {
				isOk = true
				result = Of[T](cv)
			}
		} else if curr, ok = rv.(map[string]any); ok {
			isOk = true
			pathPresent = append(pathPresent, true)
		} else {
			pathPresent = append(pathPresent, false)
		}
	} else {
		pathPresent = append(pathPresent, false)
	}
	return curr, result, pathPresent, isOk
}

func extractPathIndexed[T any](pty string, idx int, last bool, curr map[string]any, result *Optional[T], pathPresent []bool, converters ...Converter[T]) (map[string]any, *Optional[T], []bool, bool) {
	isOk := false
	if rv, ok := curr[pty]; ok {
		pathPresent = append(pathPresent, true)
		if sv, ok := rv.([]any); ok {
			if idx < 0 {
				idx = len(sv) + idx
			}
			if idx >= 0 && idx < len(sv) {
				v := sv[idx]
				if last {
					pathPresent = append(pathPresent, true)
					if av, ok := v.(T); ok && isPresent(v) {
						result = Of[T](av)
						isOk = true
					} else if cv, ok := runConverters[T](v, converters...); ok {
						result = Of[T](cv)
						isOk = true
					}
				} else if curr, ok = v.(map[string]any); ok {
					pathPresent = append(pathPresent, true)
					isOk = true
				} else {
					pathPresent = append(pathPresent, false)
				}
			} else {
				pathPresent = append(pathPresent, false)
			}
		} else {
			pathPresent = append(pathPresent, false)
		}
	} else {
		pathPresent = append(pathPresent, false)
	}
	return curr, result, pathPresent, isOk
}

func getProperty(pathPart string) (string, bool, int) {
	if oat := strings.LastIndexByte(pathPart, '['); oat != -1 && strings.HasSuffix(pathPart, "]") {
		pty := pathPart[:oat]
		idxs := pathPart[oat+1 : len(pathPart)-1]
		if idx, err := strconv.ParseInt(idxs, 10, 64); err == nil {
			return pty, true, int(idx)
		}
	}
	return pathPart, false, -1
}

func runConverters[T any](value any, converters ...Converter[T]) (result T, ok bool) {
	for _, converter := range converters {
		if converter != nil {
			if result, ok = converter(value); ok {
				return
			}
		}
	}
	return
}

// Get obtains an optional from a map
//
// If the key is present (and the value is non-nil) then an optional with the value is returned, otherwise an empty optional is returned
func Get[K comparable, T any](m map[K]T, key K) *Optional[T] {
	if v, ok := m[key]; ok && isPresent(v) {
		return Of(v)
	}
	return Empty[T]()
}

// OptMap can be used to cast an existing map for optional/functional methods
type OptMap[K comparable, V any] map[K]V

// Get returns an optional of the value in the map
//
// If the value is not present (or is nil) then an empty optional is returned
func (m OptMap[K, V]) Get(key K) *Optional[V] {
	return Get(m, key)
}

// Default returns the value (if present and non-nil) otherwise returns the default value
func (m OptMap[K, V]) Default(key K, def V) V {
	return Get(m, key).Default(def)
}

// IfPresent if the key is present (and the value is non-nil) calls the supplied function with the key and value
//
// otherwise does nothing
func (m OptMap[K, V]) IfPresent(key K, f func(key K, v V)) OptMap[K, V] {
	if f != nil {
		if v, ok := m[key]; ok && isPresent(v) {
			f(key, v)
		}
	}
	return m
}

// IfPresentOtherwise if the key is present (and the value is non-nil) calls the supplied function with the key and value
//
// otherwise calls the other function with they key
func (m OptMap[K, V]) IfPresentOtherwise(key K, f func(key K, v V), other func(key K)) OptMap[K, V] {
	if v, ok := m[key]; ok && isPresent(v) {
		if f != nil {
			f(key, v)
		}
	} else if other != nil {
		other(key)
	}
	return m
}

// ComputeIfAbsent if the specified key is not present (or the value is nil) sets the value according to the specified function
//
// returns either the existing value or the newly set value
func (m OptMap[K, V]) ComputeIfAbsent(key K, f func(key K) V) V {
	if v, ok := m[key]; ok && isPresent(v) {
		return v
	}
	var rv V
	if f != nil {
		rv = f(key)
		if isPresent(rv) {
			m[key] = rv
		}
	}
	return rv
}

// ComputeIfPresent if the specified key is present (and the value is non-nil) attempts to compute a new mapping using the supplied function
//
// If the supplied function is called but returns a nil value, the key is deleted
func (m OptMap[K, V]) ComputeIfPresent(key K, f func(key K, v V) V) V {
	var rv V
	if v, ok := m[key]; ok && isPresent(v) && f != nil {
		rv = f(key, v)
		if isPresent(rv) {
			m[key] = rv
		} else {
			delete(m, key)
		}
	}
	return rv
}

// PutIfAbsent if the specified key if absent (not present or nil-value) it is set to the specified value
//
// returns true if the value was set
func (m OptMap[K, V]) PutIfAbsent(key K, v V) bool {
	if ov, ok := m[key]; !ok || !isPresent(ov) {
		m[key] = v
		return true
	}
	return false
}

// ReplaceIfPresent is the specified key is present (and the value is non-nil) it is replaced with the specified value
//
// returns true if the value was replaced
//
// If the specified replacement value is nil and
func (m OptMap[K, V]) ReplaceIfPresent(key K, v V) bool {
	if isPresent(v) {
		if ov, ok := m[key]; ok && isPresent(ov) {
			m[key] = v
			return true
		}
	} else if _, ok := m[key]; ok {
		delete(m, key)
		return true
	}
	return false
}
