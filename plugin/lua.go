package plugin

import (
	"context"
	"fmt"
	"strconv"

	lua "github.com/Shopify/go-lua"
	"github.com/evcc-io/evcc/util"
	"github.com/traefik/yaegi/interp"
)

// LuaPlugin implements scripting using Lua for evcc plugin interface
type LuaPlugin struct {
	script string
	in     []inputTransformation
	out    []outputTransformation
}

// NewLuaPluginWithSetter allows specifying a setter script and output variable
func NewLuaPluginFromConfig(ctx context.Context, other map[string]interface{}) (Plugin, error) {
	var cc struct {
		Script string
		In     []transformationConfig
		Out    []transformationConfig
	}

	if err := util.DecodeOther(other, &cc); err != nil {
		return nil, err
	}

	in, err := configureInputs(ctx, cc.In)
	if err != nil {
		return nil, err
	}

	out, err := configureOutputs(ctx, cc.Out)
	if err != nil {
		return nil, err
	}

	p := &LuaPlugin{
		script: cc.Script,
		in:     in,
		out:    out,
	}

	return p, nil
}

func (p *LuaPlugin) Inputs() []string {
	return p.in
}

// registerBit32 registers bit32.band for scripts using it
func registerBit32(l *lua.State) {
	l.NewTable()
	l.PushGoFunction(func(l *lua.State) int {
		a := int(lua.CheckInteger(l, 1))
		b := int(lua.CheckInteger(l, 2))
		l.PushInteger(a & b)
		return 1
	})
	l.SetField(-2, "band")
	l.SetGlobal("bit32")
}

// setInputs sets plugin inputs as Lua globals
func (p *LuaPlugin) setInputs(l *lua.State, params map[string]interface{}) error {
	for _, name := range p.inputs {
		val, ok := params[name]
		if !ok {
			return fmt.Errorf("missing input: %s", name)
		}
		switch v := val.(type) {
		case int:
			l.PushInteger(v)
		case int64:
			l.PushInteger(int(v))
		case float64:
			l.PushNumber(float64(v))
		case float32:
			l.PushNumber(float64(v))
		case bool:
			l.PushBoolean(v)
		case string:
			l.PushString(v)
		default:
			return fmt.Errorf("unsupported type for input %s: %T", name, v)
		}
		l.SetGlobal(name)
	}
	return nil
}

// Eval evaluates the Lua script and returns the result as interface{}
func (p *LuaPlugin) Eval(params map[string]interface{}) (interface{}, error) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	registerBit32(l)
	if err := p.setInputs(l, params); err != nil {
		return nil, err
	}
	if err := lua.DoString(l, p.script); err != nil {
		return nil, fmt.Errorf("lua error: %w", err)
	}
	if l.Top() == 0 {
		return nil, fmt.Errorf("lua script did not return a value")
	}
	ret := l.ToValue(-1)
	l.Pop(1)
	switch v := ret.(type) {
	case float64:
		return v, nil
	case int64:
		return v, nil
	case float32:
		return float64(v), nil
	case bool:
		return v, nil
	case string:
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported return type: %T", ret)
	}
}

func (p *Go) setParam(vm *interp.Interpreter) func(param string, val any) error {

}

func (p *LuaPlugin) handleGetter() (any, error) {
	if err := transformInputs(p.in, p.setParam(vm)); err != nil {
		return nil, err
	}

	return p.evaluate(vm)
}

// === Getter interfaces ===
var _ FloatGetter = (*Go)(nil)

func (p *LuaPlugin) FloatGetter() (func() (float64, error), error) {

	return func() (float64, error) {
		v, err := p.Eval(nil)
		if err != nil {
			return 0, err
		}
		vv, ok := v.(float64)
		if !ok {
			return 0, fmt.Errorf("not a float: %v", v)
		}
		return vv, nil
	}, nil
	// val, err := p.Eval(params)
	// if err != nil {
	// 	return 0, err
	// }
	// switch v := val.(type) {
	// case float64:
	// 	return v, nil
	// case int:
	// 	return float64(v), nil
	// case int64:
	// 	return float64(v), nil
	// case float32:
	// 	return float64(v), nil
	// case string:
	// 	return strconv.ParseFloat(v, 64)
	// case bool:
	// 	if v {
	// 		return 1, nil
	// 	}
	// 	return 0, nil
	// default:
	// 	return 0, fmt.Errorf("cannot convert return value (%T) to float64", v)
	// }
}

var _ IntGetter = (*Go)(nil)

func (p *LuaPlugin) IntGetter(params map[string]interface{}) (int64, error) {
	val, err := p.Eval(params)
	if err != nil {
		return 0, err
	}
	switch v := val.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert return value (%T) to int64", v)
	}
}

var _ BoolGetter = (*Go)(nil)

func (p *LuaPlugin) BoolGetter(params map[string]interface{}) (bool, error) {
	val, err := p.Eval(params)
	if err != nil {
		return false, err
	}
	switch v := val.(type) {
	case bool:
		return v, nil
	case int:
		return v != 0, nil
	case int64:
		return v != 0, nil
	case float64:
		return v != 0, nil
	case string:
		return v != "" && v != "0" && v != "false", nil
	default:
		return false, fmt.Errorf("cannot convert return value (%T) to bool", v)
	}
}

var _ StringGetter = (*Go)(nil)

func (p *LuaPlugin) StringGetter(params map[string]interface{}) (string, error) {
	val, err := p.Eval(params)
	if err != nil {
		return "", err
	}
	switch v := val.(type) {
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case bool:
		if v {
			return "true", nil
		}
		return "false", nil
	default:
		return "", fmt.Errorf("cannot convert return value (%T) to string", v)
	}
}

// === Setter interfaces ===

var _ FloatSetter = (*Go)(nil)

func (p *LuaPlugin) FloatSetter(params map[string]interface{}, v float64) error {
	return p.setterEval(params, v)
}

var _ IntSetter = (*Go)(nil)

func (p *LuaPlugin) IntSetter(params map[string]interface{}, v int64) error {
	return p.setterEval(params, v)
}

var _ BoolSetter = (*Go)(nil)

func (p *LuaPlugin) BoolSetter(params map[string]interface{}, v bool) error {
	return p.setterEval(params, v)
}

var _ StringSetter = (*Go)(nil)

func (p *LuaPlugin) StringSetter(params map[string]interface{}, v string) error {
	return p.setterEval(params, v)
}

// setterEval runs the setterScript with params and value as "value"
func (p *LuaPlugin) setterEval(params map[string]interface{}, value interface{}) error {
	if p.setterScript == "" {
		return fmt.Errorf("no setter script defined")
	}
	l := lua.NewState()
	lua.OpenLibraries(l)
	registerBit32(l)
	// Set all params as globals
	if err := p.setInputs(l, params); err != nil {
		return err
	}
	// Set "value" for the new value
	switch v := value.(type) {
	case int:
		l.PushInteger(v)
	case int64:
		l.PushInteger(int(v))
	case float64:
		l.PushNumber(v)
	case float32:
		l.PushNumber(float64(v))
	case bool:
		l.PushBoolean(v)
	case string:
		l.PushString(v)
	default:
		return fmt.Errorf("unsupported value type for setter: %T", v)
	}
	l.SetGlobal("value")
	if err := lua.DoString(l, p.setterScript); err != nil {
		return fmt.Errorf("lua setter error: %w", err)
	}
	return nil
}
