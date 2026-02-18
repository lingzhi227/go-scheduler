package agent

import (
	"context"
	"fmt"

	"github.com/lingzhi227/go-scheduler/pkg/vllm"
)

// CallableTool is the interface for tools that agents can invoke.
type CallableTool interface {
	// Declaration returns the tool's JSON schema definition for the LLM.
	Declaration() *ToolDeclaration
	// Call executes the tool with the given JSON arguments.
	Call(ctx context.Context, args []byte) (any, error)
}

// ToolDeclaration describes a tool for the LLM.
type ToolDeclaration struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"` // JSON Schema object
}

// ToolRegistry manages available tools by name.
type ToolRegistry struct {
	tools map[string]CallableTool
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{tools: make(map[string]CallableTool)}
}

func (r *ToolRegistry) Register(tool CallableTool) error {
	name := tool.Declaration().Name
	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("tool %q already registered", name)
	}
	r.tools[name] = tool
	return nil
}

func (r *ToolRegistry) Get(name string) (CallableTool, bool) {
	t, ok := r.tools[name]
	return t, ok
}

func (r *ToolRegistry) List() []CallableTool {
	out := make([]CallableTool, 0, len(r.tools))
	for _, t := range r.tools {
		out = append(out, t)
	}
	return out
}

// ToVLLMTools converts agent tools to vLLM tool definitions.
func ToVLLMTools(tools []CallableTool) []vllm.ToolDef {
	defs := make([]vllm.ToolDef, len(tools))
	for i, t := range tools {
		d := t.Declaration()
		defs[i] = vllm.ToolDef{
			Type: "function",
			Function: vllm.FunctionDef{
				Name:        d.Name,
				Description: d.Description,
				Parameters:  d.Parameters,
			},
		}
	}
	return defs
}

// FuncTool wraps a simple function as a CallableTool.
type FuncTool struct {
	decl *ToolDeclaration
	fn   func(ctx context.Context, args []byte) (any, error)
}

func NewFuncTool(name, description string, params any, fn func(ctx context.Context, args []byte) (any, error)) *FuncTool {
	return &FuncTool{
		decl: &ToolDeclaration{Name: name, Description: description, Parameters: params},
		fn:   fn,
	}
}

func (t *FuncTool) Declaration() *ToolDeclaration { return t.decl }
func (t *FuncTool) Call(ctx context.Context, args []byte) (any, error) {
	return t.fn(ctx, args)
}
