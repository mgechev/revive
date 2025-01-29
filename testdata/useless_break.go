package fixtures

import (
	ast "go/ast"
	"reflect"
)

func UselessBreaks() {

	switch {
	case true:
		break // MATCH /useless break in case clause/
	case false:
		if false {
			break
		}
	}

	select {
	case c:
		break // MATCH /useless break in case clause/
	case n:
		if true {
			if false {
				break
			}
			break
		}
	}

	for {
		switch {
		case c1:
			break // MATCH /useless break in case clause (WARN: this break statement affects this switch or select statement and not the loop enclosing it)/
		}
	}

	for _, node := range desc.Args {
		switch node := node.(type) {
		case *ast.FuncLit:
			found = true
			funcLit = node
			break // MATCH /useless break in case clause (WARN: this break statement affects this switch or select statement and not the loop enclosing it)/
		}
	}

	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		if val.Len() == 0 {
			break
		}
		for i := 0; i < val.Len(); i++ {
			oneIteration(reflect.ValueOf(i), val.Index(i))
		}
		return
	case reflect.Map:
		if val.Len() == 0 {
			break
		}
		om := fmtsort.Sort(val)
		for i, key := range om.Key {
			oneIteration(key, om.Value[i])
		}
		return
	case reflect.Chan:
		if val.IsNil() {
			break
		}
		if val.Type().ChanDir() == reflect.SendDir {
			s.errorf("range over send-only channel %v", val)
			break
		}
		i := 0
		for ; ; i++ {
			elem, ok := val.Recv()
			if !ok {
				break
			}
			oneIteration(reflect.ValueOf(i), elem)
		}
		if i == 0 {
			break
		}
		return
	case reflect.Invalid:
		break // MATCH /useless break in case clause/
	default:
		s.errorf("range can't iterate over %v", val)
	}
}
