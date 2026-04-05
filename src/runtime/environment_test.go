package runtime

import (
	"fmt"
	"testing"
)

func TestNewEnvironment(t *testing.T) {
	env := NewEnvironment(nil)

	if env == nil {
		t.Fatal("NewEnvironment returned nil")
	}

	if env.vars == nil {
		t.Error("vars map is nil")
	}

	if env.constants == nil {
		t.Error("constants map is nil")
	}

	if env.outer != nil {
		t.Error("outer should be nil")
	}
}

func TestNewEnvironmentWithOuter(t *testing.T) {
	parent := NewEnvironment(nil)
	env := NewEnvironmentWithOuter(parent)

	if env == nil {
		t.Fatal("NewEnvironmentWithOuter returned nil")
	}

	if env.outer != parent {
		t.Error("outer reference not set correctly")
	}

	if env.vars == nil {
		t.Error("vars map is nil")
	}

	if env.constants == nil {
		t.Error("constants map is nil")
	}
}

func TestDeclareVariable(t *testing.T) {
	env := NewEnvironment(nil)

	t.Run("declare variable successfully", func(t *testing.T) {
		err := env.declareVariable("x", 42, false)
		if err != nil {
			t.Errorf("declareVariable returned error: %v", err)
		}

		val, exists := env.vars["x"]
		if !exists {
			t.Fatal("variable not stored in vars map")
		}

		if val != 42 {
			t.Errorf("expected 42, got %v", val)
		}
	})

	t.Run("declare constant successfully", func(t *testing.T) {
		env := NewEnvironment(nil)
		err := env.declareVariable("PI", 3.14, true)
		if err != nil {
			t.Errorf("declareVariable returned error: %v", err)
		}

		val, exists := env.constants["PI"]
		if !exists {
			t.Fatal("constant not stored in constants map")
		}

		if val != 3.14 {
			t.Errorf("expected 3.14, got %v", val)
		}

		// Constant should also be in vars
		val, exists = env.vars["PI"]
		if !exists {
			t.Fatal("constant not stored in vars map")
		}
	})

	t.Run("redeclare variable error", func(t *testing.T) {
		env := NewEnvironment(nil)
		env.declareVariable("x", 10, false)
		err := env.declareVariable("x", 20, false)

		if err == nil {
			t.Fatal("expected error when redeclaring variable, got nil")
		}

		if err.Error() != fmt.Sprintf("Variable already declared: %s", "x") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("redeclare constant error", func(t *testing.T) {
		env := NewEnvironment(nil)
		env.declareVariable("PI", 3.14, true)
		err := env.declareVariable("PI", 2.71, true)

		if err == nil {
			t.Fatal("expected error when redeclaring constant, got nil")
		}
	})

	t.Run("cannot declare variable with same name as constant", func(t *testing.T) {
		env := NewEnvironment(nil)
		env.declareVariable("PI", 3.14, true)
		err := env.declareVariable("PI", 10, false)

		if err == nil {
			t.Fatal("expected error when declaring variable with constant name, got nil")
		}
	})
}

func TestAssignVariable(t *testing.T) {
	t.Run("assign to existing variable", func(t *testing.T) {
		env := NewEnvironment(nil)
		env.declareVariable("x", 10, false)
		err := env.assignVariable("x", 20)

		if err != nil {
			t.Errorf("assignVariable returned error: %v", err)
		}

		val, _ := env.getVariable("x")
		if val != 20 {
			t.Errorf("expected 20, got %v", val)
		}
	})

	t.Run("cannot assign to constant", func(t *testing.T) {
		env := NewEnvironment(nil)
		env.declareVariable("PI", 3.14, true)
		err := env.assignVariable("PI", 2.71)

		if err == nil {
			t.Fatal("expected error when assigning to constant, got nil")
		}

		if err.Error() != fmt.Sprintf("Cannot assign to constant: %s", "PI") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("declare new variable if not found", func(t *testing.T) {
		env := NewEnvironment(nil)
		err := env.assignVariable("y", 42)

		if err != nil {
			t.Errorf("assignVariable returned error: %v", err)
		}

		val, _ := env.getVariable("y")
		if val != 42 {
			t.Errorf("expected 42, got %v", val)
		}
	})

	t.Run("assign different types", func(t *testing.T) {
		env := NewEnvironment(nil)
		env.assignVariable("a", "string")
		env.assignVariable("b", 3.14)
		env.assignVariable("c", true)

		val_a, _ := env.getVariable("a")
		val_b, _ := env.getVariable("b")
		val_c, _ := env.getVariable("c")

		if val_a != "string" {
			t.Errorf("expected 'string', got %v", val_a)
		}

		if val_b != 3.14 {
			t.Errorf("expected 3.14, got %v", val_b)
		}

		if val_c != true {
			t.Errorf("expected true, got %v", val_c)
		}
	})
}

func TestGetVariable(t *testing.T) {
	t.Run("get variable from current scope", func(t *testing.T) {
		env := NewEnvironment(nil)
		env.declareVariable("x", 42, false)

		val, err := env.getVariable("x")
		if err != nil {
			t.Errorf("getVariable returned error: %v", err)
		}

		if val != 42 {
			t.Errorf("expected 42, got %v", val)
		}
	})

	t.Run("get constant from current scope", func(t *testing.T) {
		env := NewEnvironment(nil)
		env.declareVariable("PI", 3.14, true)

		val, err := env.getVariable("PI")
		if err != nil {
			t.Errorf("getVariable returned error: %v", err)
		}

		if val != 3.14 {
			t.Errorf("expected 3.14, got %v", val)
		}
	})

	t.Run("get variable from outer scope", func(t *testing.T) {
		parent := NewEnvironment(nil)
		parent.declareVariable("x", 100, false)

		child := NewEnvironmentWithOuter(parent)

		val, err := child.getVariable("x")
		if err != nil {
			t.Errorf("getVariable returned error: %v", err)
		}

		if val != 100 {
			t.Errorf("expected 100, got %v", val)
		}
	})

	t.Run("variable not found error", func(t *testing.T) {
		env := NewEnvironment(nil)
		val, err := env.getVariable("nonexistent")

		if err == nil {
			t.Fatal("expected error when getting nonexistent variable, got nil")
		}

		if val != nil {
			t.Errorf("expected nil value, got %v", val)
		}

		if err.Error() != fmt.Sprintf("Variable not found: %s", "nonexistent") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("nested scope traversal", func(t *testing.T) {
		grandparent := NewEnvironment(nil)
		grandparent.declareVariable("a", 1, false)

		parent := NewEnvironmentWithOuter(grandparent)
		parent.declareVariable("b", 2, false)

		child := NewEnvironmentWithOuter(parent)
		child.declareVariable("c", 3, false)

		// Child can access all variables
		val_a, _ := child.getVariable("a")
		val_b, _ := child.getVariable("b")
		val_c, _ := child.getVariable("c")

		if val_a != 1 {
			t.Errorf("expected 1, got %v", val_a)
		}

		if val_b != 2 {
			t.Errorf("expected 2, got %v", val_b)
		}

		if val_c != 3 {
			t.Errorf("expected 3, got %v", val_c)
		}
	})

	t.Run("child overrides parent variable", func(t *testing.T) {
		parent := NewEnvironment(nil)
		parent.declareVariable("x", 10, false)

		child := NewEnvironmentWithOuter(parent)
		child.declareVariable("x", 20, false)

		parent_val, _ := parent.getVariable("x")
		child_val, _ := child.getVariable("x")

		if parent_val != 10 {
			t.Errorf("parent value should be 10, got %v", parent_val)
		}

		if child_val != 20 {
			t.Errorf("child value should be 20, got %v", child_val)
		}
	})
}

func TestEnvironmentIntegration(t *testing.T) {
	t.Run("complete workflow", func(t *testing.T) {
		// Create parent environment
		parent := NewEnvironment(nil)
		parent.declareVariable("global", "GLOBAL", false)
		parent.declareVariable("MAX", 100, true)

		// Create child environment
		child := NewEnvironmentWithOuter(parent)

		// Test declaring in child
		child.declareVariable("local", "LOCAL", false)
		child.declareVariable("localConst", 50, true)

		// Test assigning in child
		child.assignVariable("global", "MODIFIED")

		// Test getting from child
		local_val, _ := child.getVariable("local")
		global_val, _ := child.getVariable("global")
		max_val, _ := child.getVariable("MAX")

		if local_val != "LOCAL" {
			t.Errorf("expected 'LOCAL', got %v", local_val)
		}

		if global_val != "MODIFIED" {
			t.Errorf("expected 'MODIFIED', got %v", global_val)
		}

		if max_val != 100 {
			t.Errorf("expected 100, got %v", max_val)
		}

		// Parent's global should still be the original value
		parent_global_val, _ := parent.getVariable("global")
		if parent_global_val != "GLOBAL" {
			t.Errorf("parent variable should be 'GLOBAL', got %v", parent_global_val)
		}

		// Should not be able to assign to constant in current scope
		err := child.assignVariable("localConst", 200)
		if err == nil {
			t.Fatal("expected error when assigning to constant in current scope")
		}
	})
}
