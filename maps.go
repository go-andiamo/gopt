package gopt

// OptMap can be used to cast an existing map for optional/functional methods
type OptMap[K comparable, V any] map[K]V

// Get returns an optional of the value in the map
//
// If the value is not present (or is nil) then an empty optional is returned
func (m OptMap[K, V]) Get(k K) *Optional[V] {
	return Map(m, k)
}

// Default returns the value (if present and non-nil) otherwise returns the default value
func (m OptMap[K, V]) Default(k K, def V) V {
	return Map(m, k).Default(def)
}

// ComputeIfAbsent if the specified key is not present (or the value is nil) sets the value according to the specified function
//
// returns either the existing value or the newly set value
func (m OptMap[K, V]) ComputeIfAbsent(k K, f func(k K) V) V {
	if v, ok := m[k]; ok && isPresent(v) {
		return v
	}
	var rv V
	if f != nil {
		rv = f(k)
		if isPresent(rv) {
			m[k] = rv
		}
	}
	return rv
}

// ComputeIfPresent if the specified key is present (and the value is non-nil) attempts to compute a new mapping using the supplied function
//
// If the supplied function is called but returns a nil value, the key is deleted
func (m OptMap[K, V]) ComputeIfPresent(k K, f func(k K, v V) V) V {
	var rv V
	if v, ok := m[k]; ok && isPresent(v) && f != nil {
		rv = f(k, v)
		if isPresent(rv) {
			m[k] = rv
		} else {
			delete(m, k)
		}
	}
	return rv
}

// PutIfAbsent if the specified key if absent (not present or nil-value) it is set to the specified value
//
// returns true if the value was set
func (m OptMap[K, V]) PutIfAbsent(k K, v V) bool {
	if ov, ok := m[k]; !ok || !isPresent(ov) {
		m[k] = v
		return true
	}
	return false
}

// ReplaceIfPresent is the specified key is present (and the value is non-nil) it is replaced with the specified value
//
// returns true if the value was replaced
//
// If the specified replacement value is nil and
func (m OptMap[K, V]) ReplaceIfPresent(k K, v V) bool {
	if isPresent(v) {
		if ov, ok := m[k]; ok && isPresent(ov) {
			m[k] = v
			return true
		}
	} else if _, ok := m[k]; ok {
		delete(m, k)
		return true
	}
	return false
}
