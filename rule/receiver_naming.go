package rule

import (
	"fmt"
	"go/ast"
	"sync"

	"github.com/mgechev/revive/internal/typeparams"
	"github.com/mgechev/revive/lint"
)

// ReceiverNamingRule lints a receiver name.
type ReceiverNamingRule struct {
	receiverNameMaxLength int

	configureOnce sync.Once
}

const defaultReceiverNameMaxLength = -1 // thus will not check

func (r *ReceiverNamingRule) configure(arguments lint.Arguments) {
	r.receiverNameMaxLength = defaultReceiverNameMaxLength
	if len(arguments) < 1 {
		return
	}

	args, ok := arguments[0].(map[string]any)
	if !ok {
		panic(fmt.Sprintf("Unable to get arguments for rule %s. Expected object of key-value-pairs.", r.Name()))
	}

	for k, v := range args {
		switch k {
		case "maxLength":
			value, ok := v.(int64)
			if !ok {
				panic(fmt.Sprintf("Invalid value %v for argument %s of rule %s, expected integer value got %T", v, k, r.Name(), v))
			}
			r.receiverNameMaxLength = int(value)
		default:
			panic(fmt.Sprintf("Unknown argument %s for %s rule.", k, r.Name()))
		}
	}
}

// Apply applies the rule to given file.
func (r *ReceiverNamingRule) Apply(file *lint.File, args lint.Arguments) []lint.Failure {
	r.configureOnce.Do(func() { r.configure(args) })

	typeReceiver := map[string]string{}
	var failures []lint.Failure
	for _, decl := range file.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
			continue
		}

		names := fn.Recv.List[0].Names
		if len(names) < 1 {
			continue
		}
		name := names[0].Name

		if name == "_" {
			failures = append(failures, lint.Failure{
				Node:       decl,
				Confidence: 1,
				Category:   "naming",
				Failure:    "receiver name should not be an underscore, omit the name if it is unused",
			})
			continue
		}

		if name == "this" || name == "self" {
			failures = append(failures, lint.Failure{
				Node:       decl,
				Confidence: 1,
				Category:   "naming",
				Failure:    `receiver name should be a reflection of its identity; don't use generic names such as "this" or "self"`,
			})
			continue
		}

		if r.receiverNameMaxLength > 0 && len([]rune(name)) > r.receiverNameMaxLength {
			failures = append(failures, lint.Failure{
				Node:       decl,
				Confidence: 1,
				Category:   "naming",
				Failure:    fmt.Sprintf("receiver name %s is longer than %d characters", name, r.receiverNameMaxLength),
			})
			continue
		}

		recv := typeparams.ReceiverType(fn)
		if prev, ok := typeReceiver[recv]; ok && prev != name {
			failures = append(failures, lint.Failure{
				Node:       decl,
				Confidence: 1,
				Category:   "naming",
				Failure:    fmt.Sprintf("receiver name %s should be consistent with previous receiver name %s for %s", name, prev, recv),
			})
			continue
		}

		typeReceiver[recv] = name
	}

	return failures
}

// Name returns the rule name.
func (*ReceiverNamingRule) Name() string {
	return "receiver-naming"
}
