package parameters

type Parameters struct {
	data map[string]interface{}
}

func New() *Parameters {
	p := new(Parameters)
	return p
}

func (p *Parameters) From(data map[string]interface{}) {
	p.data = data
}

func (p *Parameters) GetString(name string, args ...interface{}) (string, bool) {
	v := p.Get(name, args)
	if v == nil {
		return "", false
	}
	return v.(string), true
}

func (p *Parameters) GetArray(name string, args ...interface{}) ([]interface{}, bool) {
	v := p.Get(name, args).([]interface{})
	if v == nil {
		return v, false
	}
	return v, true
}

func (p *Parameters) GetInt(name string, args ...interface{}) (int, bool) {
	v := p.Get(name, args)
	if v == nil {
		return 0, false
	}
	switch v.(type) {
	case float64:
		return int(v.(float64)), true
	}
	return v.(int), true
}

func (p *Parameters) GetBool(name string, args ...interface{}) (bool, bool) {
	v := p.Get(name, args)
	if v == nil {
		return false, false
	}
	return v.(bool), true
}

func (p *Parameters) GetFloat(name string, args ...interface{}) (float64, bool) {
	v := p.Get(name, args)
	if v == nil {
		return 0, false
	}
	return v.(float64), true
}

func (p *Parameters) Get(name string, args []interface{}) interface{} {
	v, ok := p.data[name]
	if !ok {
		if len(args) > 0 {
			v = args[0]
		} else {
			return nil
		}
	}
	return v
}