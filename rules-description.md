# Description of available rules

List of all available rules. The rules ported from `golint` are left unchanged and indicated in the `golint` column.

| Name                                                  | Short description                                                | `golint` | 
| ----------------------------------------------------- | ---------------------------------------------------------------- | :------: |
| [`add-constant`](#add-constant)                       | Suggests using constant for magic numbers and string literals    |    no    |
| [`atomic`](#atomic)                                   | Check for common mistaken usages of the `sync/atomic` package    |    no    |
| [`argument-limit`](#argument-limit)                   | Specifies the maximum number of arguments a function can receive |    no    |
| [`blank-imports`](#blank-imports)                     | Disallows blank imports                                          |   yes    |
| [`bool-literal-in-expr`](#bool-literal-in-expr)       | Suggests removing Boolean literals from logic expressions        |    no    |
| [`confusing-naming`](#confusing-naming)               | Warns on methods with names that differ only by capitalization   |    no    |
| [`confusing-results`](#confusing-results)             | Suggests to name potentially confusing function results          |    no    |
| [`constant-logical-expr`](#constant-logical-expr)     | Warns on constant logical expressions                            |    no    |
| [`context-as-argument`](#context-as-argument)         | `context.Context` should be the first argument of a function.    |   yes    |
| [`context-keys-type`](#context-keys-type)             | Disallows the usage of basic types in `context.WithValue`.       |   yes    |
| [`cyclomatic`](#cyclomatic)                           | Sets restriction for maximum Cyclomatic complexity.              |    no    |
| [`deep-exit`](#deep-exit)                             | Looks for program exits in funcs other than `main()` or `init()` |    no    |
| [`dot-imports`](#dot-imports)                         | Forbids `.` imports.                                             |   yes    |
| [`empty-block`](#empty-block)                         | Warns on empty code blocks                                       |    no    |
| [`error-naming`](#error-naming)                       | Naming of error variables.                                       |   yes    |
| [`error-return`](#error-return)                       | The error return parameter should be last.                       |   yes    |
| [`error-strings`](#error-strings)                     | Conventions around error strings.                                |   yes    |
| [`errorf`](#errorf)                                   | Should replace `errors.New(fmt.Sprintf())` with `fmt.Errorf()`   |   yes    |
| [`exported`](#exported)                               | Naming and commenting conventions on exported symbols.           |   yes    |
| [`file-header`](#file-header)                         |  Header which each file should have.                             |    no    |
| [`flag-parameter`](#flag-parameter)                   | Warns on boolean parameters that create a control coupling       |    no    |
| [`function-result-limit`](#function-result-limit)     | Specifies the maximum number of results a function can return    |    no    |
| [`get-return`](#get-return)                           | Warns on getters that do not yield any result                    |    no    |
| [`if-return`](#if-return)                             | Redundant if when returning an error.                            |   yes    |
| [`increment-decrement`](#increment-decrement)         | Use `i++` and `i--` instead of `i += 1` and `i -= 1`.            |   yes    |
| [`indent-error-flow`](#indent-error-flow)             | Prevents redundant else statements.                              |   yes    |
| [`imports-blacklist`](#imports-blacklist)             | Disallows importing the specified packages                       |    no    |
| [`max-public-structs`](#max-public-structs)           | The maximum number of public structs in a file.                  |    no    |
| [`modifies-parameter`](#modifies-parameter)           | Warns on assignments to function parameters                      |    no    |
| [`modifies-value-receiver`](#modifies-value-receiver) | Warns on assignments to value-passed method receivers            |    no    |
| [`package-comments`](#package-comments)               | Package commenting conventions.                                  |   yes    |
| [`range`](#range)                                     | Prevents redundant variables when iterating over a collection.   |   yes    |
| [`range-val-in-closure`](#range-val-in-closure)       | Warns if range value is used in a closure dispatched as goroutine|    no    |
| [`receiver-naming`](#receiver-naming)                 | Conventions around the naming of receivers.                      |   yes    |
| [`redefines-builtin-id`](#redefines-builtin-id)       | Warns on redefinitions of builtin identifiers                    |    no    |
| [`struct-tag`](#struct-tag)                           | Checks common struct tags like `json`,`xml`,`yaml`               |    no    |
| [`superfluous-else`](#superfluous-else)               | Prevents redundant else statements (extends `indent-error-flow`) |    no    |
| [`time-naming`](#time-naming)                         | Conventions around the naming of time variables.                 |   yes    |
| [`var-naming`](#var-naming)                           | Naming rules.                                                    |   yes    |
| [`var-declaration`](#var-declaration)                 | Reduces redundancies around variable declaration.                |   yes    |
| [`unexported-return`](#unexported-return)             | Warns when a public return is from unexported type.              |   yes    |
| [`unnecessary-stmt`](#unnecessary-stmt)               | Suggests removing or simplifying unnecessary statements          |    no    |
| [`unreachable-code`](#unreachable-code)               | Warns on unreachable code                                        |    no    |
| [`unused-parameter`](#unused-parameter)               | Suggests to rename or remove unused function parameters          |    no    |
| [`waitgroup-by-value  `](#waitgroup-by-value  )       | Warns on functions taking sync.WaitGroup as a by-value parameter |    no    |

## add-constant

_Description_: Suggests using constant for [magic numbers](https://en.wikipedia.org/wiki/Magic_number_(programming)#Unnamed_numerical_constants) and string literals.

_Configuration_:

* `maxLitCount` : (string) maximum number of instances of a string literal that are tolerated before warn.
* `allowStr`: (string) comma separated list of allowed string literals
* `allowInts`: (string) comma separated list of allowed integers
* `allowFloats`: (string) comma separated list of allowed floats

Example:

```toml
[rule.add-constant]
  arguments = [{maxLitCount = "3",allowStrs ="\"\"",allowInts="0,1,2",allowFloats="0.0,0.,1.0,1.,2.0,2."}]
```

## atomic

_Description_: Check for common mistaken usages of the `sync/atomic` package

_Configuration_: N/A

## argument-limit

_Description_: Warns when a function receives more parameters than the maximum set by the rule's configuration.
Enforcing a maximum number of parameters helps keeping code readable and maintainable.

_Configuration_: (int) the maximum number of parameters allowed per function.

Example:

```toml
[argument-limit]
  arguments =[4]
```

## blank-imports

_Description_: Blank import should be only in a main or test package, or have a comment justifying it; this rule warns if that is not the case.

_Configuration_: N/A

## bool-literal-in-expr

_Description_: Using Boolean literals (`true`, `false`) in logic expressions may make the code less readable. This rule
suggests removing Boolean literals from logic expressions.
The rule will also try to spot logical expressions that evaluate always to the same value.

_Configuration_: N/A

## confusing-naming

_Description_: Methods or fields of `struct`s that have similar names may be misleading.
Warns on methods and fields with names that differ only by capitalization.

_Configuration_: N/A

## confusing-results

_Description_: Function or methods that return multiple, no named, values of the same type may induce to error. This rule warns on such cases.

_Configuration_: N/A

## constant-logical-expr

_Description_: Using Boolean literals (`true`, `false`) in logic expressions may make the code less readable. This rule
suggests removing Boolean literals from logic expressions.
The rule will also try to spot logical expressions that evaluate always to the same value.

_Configuration_: N/A

## context-as-argument

_Description_: Is a language [convention](https://github.com/golang/go/wiki/CodeReviewComments#contexts) that `context.Context` should be the first parameter of a function. This rule spots function declarations that do not follow the convention.

_Configuration_: N/A

## context-keys-type

_Description_: Basic types should not be used as key in `context.WithValue`.
This rule warns when that is the case.

_Configuration_: N/A

## cyclomatic

_Description_: [Cyclomatic complexity](https://en.wikipedia.org/wiki/Cyclomatic_complexity) is a measure of code complexity.
Enforcing a maximum complexity per function helps keeping code readable and maintainable.

_Configuration_: (int) the maximum function complexity

Example:
```toml
[cyclomatic]
  arguments =[3]
```

## deep-exit

_Description_: Packages exposing functions that can stop program execution by exiting are hard to reuse. This rule looks for program exits in functions other than `main()` or `init()`.

_Configuration_: N/A

## dot-imports
 
 _Description_: Importing with `.` makes the programs much harder to understand because it is unclear where names belongs to the current package or to an imported package.
 This rule warns when using imports `.`

 More information [here](https://github.com/golang/go/wiki/CodeReviewComments#import-dot)

_Configuration_: N/A

## empty-block

_Description_: Empty blocks make code less readable and may be the symptom of a bug (unfinished refactoring?). This rule spots empty blocks in the code.

_Configuration_: N/A

## error-naming

_Description_: By convention, for the sake of readability, variables holding errors must be named with the prefix `err`. This rule warns when this convention is not followed.

_Configuration_: N/A

## error-return

_Description_: By convention, for the sake of readability, when a function returns a value of type `error`, the value must be the last one of the list of return values. This rule warns when this convention is not followed.

_Configuration_: N/A

## error-strings

_Description_: By convention, for better readability, error strings should not be capitalized or end with punctuation or a newline. This rule warns when this convention is not followed.

More information [here](https://github.com/golang/go/wiki/CodeReviewComments#error-strings)

_Configuration_: N/A

## errorf

_Description_: It is possible get simpler code by replacing `errors.New(fmt.Sprintf())` with `fmt.Errorf()`. This rule spots that kind of simplification opportunities.

_Configuration_: N/A

## exported

_Description_: Exported function and methods should have comments. 
This warns on undocumented exported functions and methods.

More information [here](https://github.com/golang/go/wiki/CodeReviewComments#doc-comments)

_Configuration_: N/A

## file-header

_Description_: This rule helps enforcing a common header for all source files in a project
by spotting those files that do not have the specified header.

_Configuration_: (string) the header to look for in source files.

Example:
```toml
[file-header]
  arguments =["This is the text that must appear at the top of source files."]
```

## flag-parameter

_Description_: If a function controls the flow of another by passing it information on what to do, both functions are said to be [control-coupled](https://en.wikipedia.org/wiki/Coupling_(computer_programming)#Procedural_programming).
Coupling among functions must be minimized for a better maintainability of the code.
This rule warns on boolean parameters that create a control coupling.

_Configuration_: N/A

## function-result-limit
## get-return
## if-return
## increment-decrement
## indent-error-flow
## imports-blacklist
## max-public-structs
## modifies-parameter
## modifies-value-receiver
## package-comments
## range    
## range-val-in-closure
## receiver-naming
## redefines-builtin-id
## struct-tag
## superfluous-else
## time-naming
## var-naming
## var-declaration
## unexported-return
## unnecessary-stmt
## unreachable-code
## unused-parameter
## waitgroup-by-value  