package runtime

import "fmt"

type Environment struct {
	vars      map[string]any
	constants map[string]any
	outer     *Environment
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		vars:      make(map[string]any),
		constants: make(map[string]any),
		outer:     nil,
	}
}

func NewEnvironmentWithOuter(outer *Environment) *Environment {
	return &Environment{
		vars:      make(map[string]any),
		constants: make(map[string]any),
		outer:     outer,
	}
}

func (e *Environment) declareVariable(name string, value any, isConstant bool) error {
	if e.vars[name] != nil || e.constants[name] != nil {
		return fmt.Errorf("Variable already declared: %s", name)
	}

	if isConstant {
		e.constants[name] = value
	}

	e.vars[name] = value
	return nil
}

func (e *Environment) assignVariable(name string, value any) error {
	if e.constants[name] != nil {
		return fmt.Errorf("Cannot assign to constant: %s", name)
	}

	if e.vars[name] != nil {
		e.vars[name] = value
		return nil
	}

	return e.declareVariable(name, value, true)
}

func (e *Environment) getVariable(name string) (any, error) {
	if val, ok := e.vars[name]; ok {
		return val, nil
	}

	if val, ok := e.constants[name]; ok {
		return val, nil
	}

	if e.outer != nil {
		return e.outer.getVariable(name)
	}

	return nil, fmt.Errorf("Variable not found: %s", name)
}
