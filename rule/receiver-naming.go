package rule

import (
	"fmt"
	"go/ast"
	"strconv"
	"strings"
	"sync"

	"github.com/mgechev/revive/internal/typeparams"
	"github.com/mgechev/revive/lint"
)

// ReceiverNamingRule lints given else constructs.
type ReceiverNamingRule struct {
	receiverNameMaxLength int
	sync.Mutex
}

const defaultReceiverNameMaxLength = -1 // thus will not check

func (r *ReceiverNamingRule) configure(arguments lint.Arguments) {
	r.Lock()
	defer r.Unlock()
	if r.receiverNameMaxLength == 0 {
		if len(arguments) < 1 {
			r.receiverNameMaxLength = defaultReceiverNameMaxLength
			return
		}
		arg := arguments[0]
		argStr, ok := arg.(string)
		if !ok {
			panic(fmt.Sprintf("Invalid argument for %s rule. Expecting an string, got %T", r.Name(), arg))
		}

		parts := strings.Split(argStr, "=")
		if len(parts) != 2 {
			panic(fmt.Sprintf("Invalid argument for %s rule. Expecting an string of the form 'key=value', got %s", r.Name(), argStr))
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "max-length":
			var err error
			r.receiverNameMaxLength, err = strconv.Atoi(value)
			if err != nil {
				panic(fmt.Sprintf("Invalid value %s for the configuration key max-length, expected integer value: %v", value, err))
			}
		default:
			panic(fmt.Sprintf("Unknown configuration key %s for %s rule.", key, r.Name()))
		}
	}
}

// Apply applies the rule to given file.
func (r *ReceiverNamingRule) Apply(file *lint.File, args lint.Arguments) []lint.Failure {
	r.configure(args)

	var failures []lint.Failure

	fileAst := file.AST
	walker := lintReceiverName{
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
		typeReceiver:          map[string]string{},
		receiverNameMaxLength: r.receiverNameMaxLength,
	}

	ast.Walk(walker, fileAst)

	return failures
}

// Name returns the rule name.
func (*ReceiverNamingRule) Name() string {
	return "receiver-naming"
}

type lintReceiverName struct {
	onFailure             func(lint.Failure)
	typeReceiver          map[string]string
	receiverNameMaxLength int
}

func (w lintReceiverName) Visit(n ast.Node) ast.Visitor {
	fn, ok := n.(*ast.FuncDecl)
	if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
		return w
	}
	names := fn.Recv.List[0].Names
	if len(names) < 1 {
		return w
	}
	name := names[0].Name
	if name == "_" {
		w.onFailure(lint.Failure{
			Node:       n,
			Confidence: 1,
			Category:   "naming",
			Failure:    "receiver name should not be an underscore, omit the name if it is unused",
		})
		return w
	}
	if name == "this" || name == "self" {
		w.onFailure(lint.Failure{
			Node:       n,
			Confidence: 1,
			Category:   "naming",
			Failure:    `receiver name should be a reflection of its identity; don't use generic names such as "this" or "self"`,
		})
		return w
	}

	if w.receiverNameMaxLength > 0 && len(name) > w.receiverNameMaxLength {
		w.onFailure(lint.Failure{
			Node:       n,
			Confidence: 1,
			Category:   "naming",
			Failure:    fmt.Sprintf("receiver name %s is longer than %d characters", name, w.receiverNameMaxLength),
		})
		return w
	}

	recv := typeparams.ReceiverType(fn)
	if prev, ok := w.typeReceiver[recv]; ok && prev != name {
		w.onFailure(lint.Failure{
			Node:       n,
			Confidence: 1,
			Category:   "naming",
			Failure:    fmt.Sprintf("receiver name %s should be consistent with previous receiver name %s for %s", name, prev, recv),
		})
		return w
	}
	w.typeReceiver[recv] = name
	return w
}
