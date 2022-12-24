package gopt

// Extract extracts an optional value, of the specified type, from a map
//
// If the key is present (and the value is non-nil and of the specified type) then an optional with the value is returned, otherwise an empty optional is returned
func Extract[K comparable, T any](m map[K]any, key K) *Optional[T] {
	result := Empty[T]()
	if rv, ok := m[key]; ok && isPresent(rv) {
		if v, ok := rv.(T); ok {
			result = Of[T](v)
		}
	}
	return result
}

// ExtractJson extracts an optional value, of the specified type, from a map[string]interface{}
//
// If the key is present (and the value is non-nil and of the specified type) then an optional with the value is returned, otherwise an empty optional is returned
func ExtractJson[T any](m map[string]interface{}, key string) *Optional[T] {
	result := Empty[T]()
	if rv, ok := m[key]; ok && isPresent(rv) {
		if v, ok := rv.(T); ok {
			result = Of[T](v)
		}
	}
	return result
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
