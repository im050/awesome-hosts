package main

type Parameters struct {
	data map[string]interface{}
}

func NewParameters() *Parameters {
	m := new(Parameters)
	return m
}

func (p *Parameters) From(data map[string]interface{}) {
	p.data = data
}

func (p *Parameters) GetString(name string, defaultValue string) string {
	v := p.Get(name, defaultValue)
	if v == nil {
		return ""
	}
	return v.(string)
}

func (p *Parameters) GetInt(name string, defaultValue int) int {
	v := p.Get(name, defaultValue)
	if v == nil {
		return 0
	}
	switch v.(type) {
	case float64:
		return int(v.(float64))
	}
	return v.(int)
}

func (p *Parameters) GetBool(name string, defaultValue bool) bool {
	v := p.Get(name, defaultValue)
	if v == nil {
		return false
	}
	return v.(bool)
}

func (p *Parameters) GetFloat(name string, defaultValue float64) float64 {
	v := p.Get(name, defaultValue)
	if v == nil {
		return 0
	}
	return v.(float64)
}

func (p *Parameters) Get(name string, defaultValue interface{}) interface{} {
	v, ok := p.data[name]
	if !ok {
		if defaultValue != nil {
			v = defaultValue
		} else {
			return nil
		}
	}
	return v
}