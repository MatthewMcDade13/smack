package interp

import "fmt"

type Env struct {
	outer *Env
	data  SmackMap
}

func NewEnv(outer *Env) *Env {
	data := make(SmackMap)
	return &Env{
		outer, data,
	}
}

func (e *Env) Set(key_sym string, val Value) {
	e.data[key_sym] = val
}

// NOTE :: Check returned value for v.Type() != VAL_NONE, as result may be nil
func (e *Env) Find(key_sym string) Value {
	v := e.data[key_sym]
	if !v.IsNone() {
		return v
	}

	if e.outer == nil {
		return NoValue()
	}

	return e.outer.Find(key_sym)
}

func (e *Env) Get(key_sym string) (Value, error) {
	v := e.Find(key_sym)
	if v.IsNone() {
		return NoValue(), fmt.Errorf("Value not found in environment for given Symbol string(name): %s", key_sym)
	} else {
		return v, nil
	}
}
