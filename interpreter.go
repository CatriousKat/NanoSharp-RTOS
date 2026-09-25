package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed script.ns
var scriptContent string

// --- Environment & State ---

type Environment struct {
	vars   map[string]any
	parent *Environment
}

func NewEnv(parent *Environment) *Environment {
	return &Environment{
		vars:   make(map[string]any),
		parent: parent,
	}
}

func (e *Environment) Set(name string, val any) {
	e.vars[name] = val
}

func (e *Environment) Get(name string) (any, bool) {
	if val, ok := e.vars[name]; ok {
		return val, true
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	return nil, false
}

type Interpreter struct {
	globalEnv *Environment
	files     map[string]string // In-Memory Virtual File System Storage
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		globalEnv: NewEnv(nil),
		files:     make(map[string]string),
	}
}

// --- Main Execution Engine ---

func (it *Interpreter) Run(source string) {
	lines := strings.Split(source, "\n")
	it.evalBlock(lines, it.globalEnv)
}

func (it *Interpreter) evalBlock(lines []string, env *Environment) any {
	i := 0
	for i < len(lines) {
		raw := lines[i]
		line := strings.TrimSpace(strings.ReplaceAll(raw, "\u00a0", " "))

		if line == "" || strings.HasPrefix(line, "//") {
			i++
			continue
		}

		if strings.HasPrefix(line, "if ") || line == "if" {
			condExpr := extractExpression(line, "if")
			newIdx, _ := it.handleIf(lines, i, condExpr, env)
			i = newIdx
			continue
		}

		it.evalLine(line, env)
		i++
	}
	return nil
}

func (it *Interpreter) evalLine(line string, env *Environment) any {
	if strings.Contains(line, "=") && !strings.Contains(line, "==") && !strings.HasPrefix(line, "if") {
		parts := strings.SplitN(line, "=", 2)
		varName := strings.TrimSpace(parts[0])
		rhsExpr := strings.TrimSpace(parts[1])
		varVal := it.evalExpr(rhsExpr, env)
		env.Set(varName, varVal)
		return varVal
	}
	return it.evalExpr(line, env)
}

// --- Control Flow (Simple If) ---

func (it *Interpreter) handleIf(lines []string, idx int, line string, env *Environment) (int, any) {
	condExpr := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "if"), "{"))
	condVal := it.evalExpr(condExpr, env)
	executed := it.isTruthy(condVal)

	subLines, nextIdx := collectBlock(lines, idx)
	if executed {
		it.evalBlock(subLines, NewEnv(env))
	}

	return nextIdx, nil
}

// --- Expression & Command Evaluation ---

func (it *Interpreter) evalExpr(expr string, env *Environment) any {
	expr = strings.TrimSpace(expr)
	if strings.HasPrefix(expr, "\"") && strings.HasSuffix(expr, "\"") {
		return expr[1 : len(expr)-1]
	}

	// --- Virtual File System (fs.*) Handlers ---
	if strings.HasPrefix(expr, "fs.write(") && strings.HasSuffix(expr, ")") {
		content := expr[9 : len(expr)-1]
		args := parseArgs(content)
		if len(args) == 2 {
			filename := fmt.Sprintf("%v", it.evalExpr(args[0], env))
			fileData := fmt.Sprintf("%v", it.evalExpr(args[1], env))
			it.files[filename] = fileData
			return true
		}
		return false
	}

	if strings.HasPrefix(expr, "fs.read(") && strings.HasSuffix(expr, ")") {
		content := expr[8 : len(expr)-1]
		args := parseArgs(content)
		if len(args) == 1 {
			filename := fmt.Sprintf("%v", it.evalExpr(args[0], env))
			if val, ok := it.files[filename]; ok {
				return val
			}
			return "" // Return empty string if file doesn't exist
		}
		return ""
	}

	if strings.HasPrefix(expr, "gui.") {
		return nil
	}

	if strings.HasPrefix(expr, "print(") && strings.HasSuffix(expr, ")") {
		content := expr[6 : len(expr)-1]
		args := parseArgs(content)
		var evaluated []string
		for _, arg := range args {
			val := it.evalExpr(arg, env)
			evaluated = append(evaluated, fmt.Sprintf("%v", val))
		}
		println(strings.Join(evaluated, " "))
		return nil
	}

	for _, op := range []string{">=", "<=", "==", "!=", ">", "<"} {
		if strings.Contains(expr, op) {
			parts := strings.SplitN(expr, op, 2)
			left := it.evalExpr(parts[0], env)
			right := it.evalExpr(parts[1], env)
			return compareValues(left, op, right)
		}
	}

	if val, ok := env.Get(expr); ok {
		return val
	}

	if n, err := strconv.Atoi(expr); err == nil {
		return n
	}

	if expr == "true" {
		return true
	}
	if expr == "false" {
		return false
	}

	return expr
}

// --- Utility Functions ---

func compareValues(left any, op string, right any) bool {
	lInt, lIsInt := toInt(left)
	rInt, rIsInt := toInt(right)

	if lIsInt && rIsInt {
		switch op {
		case ">":
			return lInt > rInt
		case "<":
			return lInt < rInt
		case ">=":
			return lInt >= rInt
		case "<=":
			return lInt <= rInt
		case "==":
			return lInt == rInt
		case "!=":
			return lInt != rInt
		}
	}
	return false
}

func toInt(val any) (int, bool) {
	switch v := val.(type) {
	case int:
		return v, true
	case string:
		if n, err := strconv.Atoi(v); err == nil {
			return n, true
		}
	}
	return 0, false
}

func (it *Interpreter) isTruthy(val any) bool {
	switch v := val.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case string:
		return v != "" && v != "false" && v != "0"
	default:
		return val != nil
	}
}

func extractExpression(line, keyword string) string {
	line = strings.TrimPrefix(line, keyword)
	line = strings.TrimSpace(line)
	line = strings.TrimSuffix(line, "{")
	return strings.TrimSpace(line)
}

func collectBlock(lines []string, startIdx int) ([]string, int) {
	var block []string
	depth := 0
	i := startIdx

	for i < len(lines) {
		line := strings.TrimSpace(strings.ReplaceAll(lines[i], "\u00a0", " "))
		if strings.Contains(line, "{") {
			depth++
			i++
			break
		}
		i++
	}

	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(strings.ReplaceAll(line, "\u00a0", " "))

		if strings.Contains(trimmed, "{") {
			depth++
		}
		if strings.Contains(trimmed, "}") {
			depth--
			if depth == 0 {
				i++
				break
			}
		}
		block = append(block, line)
		i++
	}
	return block, i
}

func parseArgs(s string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false

	for _, r := range s {
		if r == '"' {
			inQuotes = !inQuotes
			current.WriteRune(r)
		} else if r == ',' && !inQuotes {
			args = append(args, strings.TrimSpace(current.String()))
			current.Reset()
		} else {
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		args = append(args, strings.TrimSpace(current.String()))
	}
	return args
}

// --- Main Execution Point ---

func main() {
	println("Starting NanoSharp on Microcontroller...")
	interpreter := NewInterpreter()
	interpreter.Run(scriptContent)
	println("Execution finished.")
}