package plugin

import (
	"context"
	"fmt"

	lua "github.com/Shopify/go-lua"
	"github.com/evcc-io/evcc/util"
)

// LuaPlugin implements Lua scripting for the plugin interface
type LuaPlugin struct {
	script string
	in     []inputTransformation
	out    []outputTransformation
	state  *lua.State
}

func init() {
	registry.AddCtx("lua", NewLuaPluginFromConfig)
}

// NewLuaPluginFromConfig creates a Lua provider
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

	l := lua.NewState()
	lua.OpenLibraries(l)

	p := &LuaPlugin{
		script: cc.Script,
		in:     in,
		out:    out,
		state:  l,
	}

	return p, nil
}

var _ FloatGetter = (*LuaPlugin)(nil)

// FloatGetter retrieves float from Lua script
func (p *LuaPlugin) FloatGetter() (func() (float64, error), error) {
	return func() (float64, error) {
		v, err := p.handleGetter()
		if err != nil {
			return 0, err
		}

		vv, ok := v.(float64)
		if !ok {
			return 0, fmt.Errorf("not a float: %v", v)
		}

		return vv, nil
	}, nil
}

var _ IntGetter = (*LuaPlugin)(nil)

// IntGetter retrieves int64 from Lua script
func (p *LuaPlugin) IntGetter() (func() (int64, error), error) {
	return func() (int64, error) {
		v, err := p.handleGetter()
		if err != nil {
			return 0, err
		}

		vv, ok := v.(int64)
		if !ok {
			return 0, fmt.Errorf("not an int: %v", v)
		}

		return vv, nil
	}, nil
}

var _ StringGetter = (*LuaPlugin)(nil)

// StringGetter retrieves string from Lua script
func (p *LuaPlugin) StringGetter() (func() (string, error), error) {
	return func() (string, error) {
		v, err := p.handleGetter()
		if err != nil {
			return "", err
		}

		vv, ok := v.(string)
		if !ok {
			return "", fmt.Errorf("not a string: %v", v)
		}

		return vv, nil
	}, nil
}

var _ BoolGetter = (*LuaPlugin)(nil)

// BoolGetter retrieves bool from Lua script
func (p *LuaPlugin) BoolGetter() (func() (bool, error), error) {
	return func() (bool, error) {
		v, err := p.handleGetter()
		if err != nil {
			return false, err
		}

		vv, ok := v.(bool)
		if !ok {
			return false, fmt.Errorf("not a bool: %v", v)
		}

		return vv, nil
	}, nil
}

var _ FloatSetter = (*LuaPlugin)(nil)

// FloatSetter sets float in Lua script
func (p *LuaPlugin) FloatSetter(param string) (func(float64) error, error) {
	return func(val float64) error {
		return p.setterEval(param, val)
	}, nil
}

var _ IntSetter = (*LuaPlugin)(nil)

// IntSetter sets int64 in Lua script
func (p *LuaPlugin) IntSetter(param string) (func(int64) error, error) {
	return func(val int64) error {
		return p.setterEval(param, val)
	}, nil
}

var _ StringSetter = (*LuaPlugin)(nil)

// StringSetter sets string in Lua script
func (p *LuaPlugin) StringSetter(param string) (func(string) error, error) {
	return func(val string) error {
		return p.setterEval(param, val)
	}, nil
}

var _ BoolSetter = (*LuaPlugin)(nil)

// BoolSetter sets bool in Lua script
func (p *LuaPlugin) BoolSetter(param string) (func(bool) error, error) {
	return func(val bool) error {
		return p.setterEval(param, val)
	}, nil
}

func (p *LuaPlugin) setParam() func(param string, val any) error {
	return func(param string, val any) error {
		fmt.Println("setParam", param, val)
		switch v := val.(type) {
		case int:
			p.state.PushInteger(v)
		case int64:
			p.state.PushInteger(int(v))
		case float64:
			p.state.PushNumber(v)
		case bool:
			p.state.PushBoolean(v)
		case string:
			p.state.PushString(v)
		default:
			return fmt.Errorf("unsupported value type for setter: %T", v)
		}
		p.state.SetGlobal(param)
		return nil
	}
}

func (p *LuaPlugin) handleGetter() (any, error) {
	if err := transformInputs(p.in, p.setParam()); err != nil {
		return nil, err
	}

	return p.evaluate()
}

// Eval evaluates the Lua script and returns the result
func (p *LuaPlugin) evaluate() (res any, err error) {
	// Create a new Lua state and open standard libraries
	if err := lua.DoString(p.state, p.script); err != nil {
		return nil, fmt.Errorf("lua error: %w", err)
	}
	if p.state.Top() == 0 {
		return nil, fmt.Errorf("lua script did not return a value")
	}
	result := p.state.ToValue(-1)
	p.state.Pop(1)
	return result, nil
}

func (p *LuaPlugin) handleSetter(param string, val any) error {
	setParam := p.setParam()

	if err := transformInputs(p.in, setParam); err != nil {
		return err
	}

	if err := setParam(param, val); err != nil {
		return err
	}

	vv := lua.NewState()

	// vv, err := p.evaluate(vm)
	// if err != nil {
	// 	return err
	// }

	return transformOutputs(p.out, vv)
}

// Eval evaluates the Lua script and returns the result
// func (p *LuaPlugin) handleSetter()(param string, val any) error {

// 	setParam := p.setParam()

// 	if err := transformInputs(p.in, setParam); err != nil {
// 		return err
// 	}

// 	l := lua.NewState()
// 	lua.OpenLibraries(l)
// 	// if err := p.setInputs(l, params); err != nil {
// 	// 	return nil, err
// 	// }
// 	// if err := lua.DoString(l, p.script); err != nil {
// 	// 	return nil, fmt.Errorf("lua error: %w", err)
// 	// }
// 	// if l.Top() == 0 {
// 	// 	return nil, fmt.Errorf("lua script did not return a value")
// 	// }
// 	l.ToValue(-1)
// 	l.Pop(1)
// 	return transformOutputs(p.out, vv)
// }

// setterEval executes the setter script with provided parameters and value
func (p *LuaPlugin) setterEval(param string, value any) error {
	l := lua.NewState()
	lua.OpenLibraries(l)
	if err := p.setInputs(l, param); err != nil {
		return err
	}
	switch v := value.(type) {
	case int:
		l.PushInteger(v)
	case int64:
		l.PushInteger(int(v))
	case float64:
		l.PushNumber(v)
	case bool:
		l.PushBoolean(v)
	case string:
		l.PushString(v)
	default:
		return fmt.Errorf("unsupported value type for setter: %T", v)
	}
	l.SetGlobal("value")
	if err := lua.DoString(l, p.script); err != nil {
		return fmt.Errorf("lua setter error: %w", err)
	}
	return nil
}

// setInputs sets plugin inputs as Lua globals
func (p *LuaPlugin) setInputs(l *lua.State, params string) error {

	// for name, value := range in {
	//  	switch v := value.(type) {
	//  	case int:
	//  		l.PushInteger(v)
	//  	case int64:
	//  		l.PushInteger(int(v))
	//  	case float64:
	//  		l.PushNumber(v)
	//  	case bool:
	//  		l.PushBoolean(v)
	//  	case string:
	//  		l.PushString(v)
	//  	default:
	//  		return fmt.Errorf("unsupported type for input %s: %T", name, v)
	// 	}
	//  	l.SetGlobal(name)
	// }
	return nil
}
