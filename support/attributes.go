package support

import "encoding/json"

type Attributes struct {
  values map[string]interface{}
}

// Get gets the first value associated with the given key.
func (a *Attributes) Get(key string) interface{} {
	if a.values == nil {
		a.values = map[string]interface{}{}
	}
	vs := a.values[key]
	if a.Has(key) == false {
		return nil
	}
	return vs
}

// Add adds the value to key. It appends to any existing
// attributes associated with key.
func (a *Attributes) Add(key string, value  interface{}) {
	if a.values == nil {
		a.values = map[string]interface{}{}
	}
	a.values[key] = value
}

// Del deletes the values associated with key.
func (a *Attributes) Del(key  string) {
	if a.values == nil {
		a.values = map[string]interface{}{}
	}
	delete(a.values, key)
}

// Has checks whether a given key is set.
func (a *Attributes) Has(key string) bool {
	if a.values == nil {
		a.values = map[string]interface{}{}
	}
	_, ok := a.values[key]
	return ok
}
// MarshalJSON custom parse JSON
func (a *Attributes) MarshalJSON() ([]byte, error) {
	values :=a.values
	return json.Marshal(&values)
}