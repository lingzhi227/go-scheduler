package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/lingzhi227/go-scheduler/pkg/vllm"
)

const defaultMaxTurns = 10

// RunAgentLoop executes the ReAct loop: call LLM -> handle tool calls -> repeat.
func RunAgentLoop(ctx context.Context, llmClient vllm.Client, ag Agent, task *Task, tools map[string]CallableTool) <-chan *Event {
	ch := make(chan *Event, 32)
	maxTurns := task.MaxTurns
	if maxTurns <= 0 {
		maxTurns = defaultMaxTurns
	}

	go func() {
		defer close(ch)

		// Build initial messages
		messages := []vllm.Message{
			{Role: "system", Content: ag.SystemPrompt()},
			{Role: "user", Content: task.Goal},
		}
		if task.Context != "" {
			messages = append(messages[:1], append([]vllm.Message{
				{Role: "user", Content: "Context: " + task.Context},
			}, messages[1:]...)...)
		}

		// Build tool definitions
		toolDefs := make([]vllm.ToolDef, 0, len(tools))
		for _, t := range tools {
			d := t.Declaration()
			toolDefs = append(toolDefs, vllm.ToolDef{
				Type: "function",
				Function: vllm.FunctionDef{
					Name:        d.Name,
					Description: d.Description,
					Parameters:  d.Parameters,
				},
			})
		}

		var usage TaskUsage

		for turn := 1; turn <= maxTurns; turn++ {
			if ctx.Err() != nil {
				emit(ch, NewEvent(EventTypeError, ag.ID(), task.ID, turn).WithError(ctx.Err().Error()))
				return
			}

			// Emit LLM call event
			emit(ch, NewEvent(EventTypeLLMCall, ag.ID(), task.ID, turn).WithData(&LLMCallData{
				MessageCount: len(messages),
				ToolCount:    len(toolDefs),
			}))

			req := &vllm.ChatCompletionRequest{
				Messages: messages,
			}
			if len(toolDefs) > 0 {
				req.Tools = toolDefs
				req.ToolChoice = "auto"
			}

			resp, err := llmClient.ChatCompletion(ctx, req)
			if err != nil {
				emit(ch, NewEvent(EventTypeError, ag.ID(), task.ID, turn).WithError(fmt.Sprintf("llm call failed: %v", err)))
				return
			}

			usage.LLMCalls++
			usage.PromptTokens += resp.Usage.PromptTokens
			usage.CompletionTokens += resp.Usage.CompletionTokens
			usage.TotalTokens += resp.Usage.TotalTokens
			usage.Turns = turn

			if len(resp.Choices) == 0 {
				emit(ch, NewEvent(EventTypeError, ag.ID(), task.ID, turn).WithError("no choices in response"))
				return
			}

			choice := resp.Choices[0]

			emit(ch, NewEvent(EventTypeLLMResult, ag.ID(), task.ID, turn).WithData(&LLMResultData{
				FinishReason: choice.FinishReason,
				ToolCalls:    len(choice.Message.ToolCalls),
				Tokens:       resp.Usage.TotalTokens,
			}))

			// Append assistant message to history
			messages = append(messages, choice.Message)

			// Check for tool calls
			if choice.FinishReason == "tool_calls" && len(choice.Message.ToolCalls) > 0 {
				for _, tc := range choice.Message.ToolCalls {
					emit(ch, NewEvent(EventTypeToolCall, ag.ID(), task.ID, turn).WithData(&ToolCallData{
						ToolName:  tc.Function.Name,
						Arguments: tc.Function.Arguments,
						CallID:    tc.ID,
					}))

					tool, ok := tools[tc.Function.Name]
					if !ok {
						toolMsg := vllm.Message{
							Role:       "tool",
							Content:    fmt.Sprintf("error: unknown tool %q", tc.Function.Name),
							ToolCallID: tc.ID,
						}
						messages = append(messages, toolMsg)
						emit(ch, NewEvent(EventTypeToolResult, ag.ID(), task.ID, turn).WithData(&ToolResultData{
							ToolName: tc.Function.Name,
							CallID:   tc.ID,
							Error:    fmt.Sprintf("unknown tool %q", tc.Function.Name),
						}))
						continue
					}

					result, err := tool.Call(ctx, []byte(tc.Function.Arguments))
					usage.ToolCalls++

					var resultStr string
					if err != nil {
						resultStr = fmt.Sprintf("error: %v", err)
						slog.Warn("tool call failed", "tool", tc.Function.Name, "error", err)
					} else {
						resultBytes, _ := json.Marshal(result)
						resultStr = string(resultBytes)
					}

					messages = append(messages, vllm.Message{
						Role:       "tool",
						Content:    resultStr,
						ToolCallID: tc.ID,
					})

					errStr := ""
					if err != nil {
						errStr = err.Error()
					}
					emit(ch, NewEvent(EventTypeToolResult, ag.ID(), task.ID, turn).WithData(&ToolResultData{
						ToolName: tc.Function.Name,
						CallID:   tc.ID,
						Result:   result,
						Error:    errStr,
					}))
				}
				continue // Next turn with tool results
			}

			// No tool calls -> agent is done
			emit(ch, NewEvent(EventTypeDone, ag.ID(), task.ID, turn).
				WithMessage(choice.Message.Content).
				WithResult(&TaskResult{
					TaskID:  task.ID,
					AgentID: ag.ID(),
					Output:  choice.Message.Content,
					Usage:   usage,
				}))
			return
		}

		// Max turns reached
		emit(ch, NewEvent(EventTypeError, ag.ID(), task.ID, maxTurns).
			WithError(fmt.Sprintf("max turns (%d) reached", maxTurns)))
	}()

	return ch
}

func emit(ch chan<- *Event, e *Event) {
	ch <- e
}
