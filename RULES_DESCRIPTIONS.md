# Description of available rules

List of all available rules.

<!-- toc -->

- [add-constant](#add-constant)
- [argument-limit](#argument-limit)
- [atomic](#atomic)
- [banned-characters](#banned-characters)
- [bare-return](#bare-return)
- [blank-imports](#blank-imports)
- [bool-literal-in-expr](#bool-literal-in-expr)
- [call-to-gc](#call-to-gc)
- [cognitive-complexity](#cognitive-complexity)
- [comment-spacings](#comment-spacings)
- [comments-density](#comments-density)
- [confusing-naming](#confusing-naming)
- [confusing-results](#confusing-results)
- [constant-logical-expr](#constant-logical-expr)
- [context-as-argument](#context-as-argument)
- [context-keys-type](#context-keys-type)
- [cyclomatic](#cyclomatic)
- [datarace](#datarace)
- [deep-exit](#deep-exit)
- [defer](#defer)
- [dot-imports](#dot-imports)
- [duplicated-imports](#duplicated-imports)
- [early-return](#early-return)
- [empty-block](#empty-block)
- [empty-lines](#empty-lines)
- [enforce-map-style](#enforce-map-style)
- [enforce-repeated-arg-type-style](#enforce-repeated-arg-type-style)
- [enforce-slice-style](#enforce-slice-style)
- [enforce-switch-style](#enforce-switch-style)
- [error-naming](#error-naming)
- [error-return](#error-return)
- [error-strings](#error-strings)
- [errorf](#errorf)
- [exported](#exported)
- [file-header](#file-header)
- [file-length-limit](#file-length-limit)
- [filename-format](#filename-format)
- [flag-parameter](#flag-parameter)
- [forbidden-call-in-wg-go](#forbidden-call-in-wg-go)
- [function-length](#function-length)
- [function-result-limit](#function-result-limit)
- [get-return](#get-return)
- [identical-branches](#identical-branches)
- [identical-ifelseif-branches](#identical-ifelseif-branches)
- [identical-ifelseif-conditions](#identical-ifelseif-conditions)
- [identical-switch-branches](#identical-switch-branches)
- [identical-switch-conditions](#identical-switch-conditions)
- [if-return](#if-return)
- [import-alias-naming](#import-alias-naming)
- [import-shadowing](#import-shadowing)
- [imports-blocklist](#imports-blocklist)
- [increment-decrement](#increment-decrement)
- [indent-error-flow](#indent-error-flow)
- [inefficient-map-lookup](#inefficient-map-lookup)
- [line-length-limit](#line-length-limit)
- [max-control-nesting](#max-control-nesting)
- [max-public-structs](#max-public-structs)
- [modifies-parameter](#modifies-parameter)
- [modifies-value-receiver](#modifies-value-receiver)
- [nested-structs](#nested-structs)
- [optimize-operands-order](#optimize-operands-order)
- [package-comments](#package-comments)
- [package-directory-mismatch](#package-directory-mismatch)
- [range-val-address](#range-val-address)
- [range-val-in-closure](#range-val-in-closure)
- [range](#range)
- [receiver-naming](#receiver-naming)
- [redefines-builtin-id](#redefines-builtin-id)
- [redundant-build-tag](#redundant-build-tag)
- [redundant-import-alias](#redundant-import-alias)
- [redundant-test-main-exit](#redundant-test-main-exit)
- [string-format](#string-format)
- [string-of-int](#string-of-int)
- [struct-tag](#struct-tag)
- [superfluous-else](#superfluous-else)
- [time-date](#time-date)
- [time-equal](#time-equal)
- [time-naming](#time-naming)
- [unchecked-type-assertion](#unchecked-type-assertion)
- [unconditional-recursion](#unconditional-recursion)
- [unexported-naming](#unexported-naming)
- [unexported-return](#unexported-return)
- [unhandled-error](#unhandled-error)
- [unnecessary-if](#unnecessary-if)
- [unnecessary-format](#unnecessary-format)
- [unnecessary-stmt](#unnecessary-stmt)
- [unreachable-code](#unreachable-code)
- [unsecure-url-scheme](#unsecure-url-scheme)
- [unused-parameter](#unused-parameter)
- [unused-receiver](#unused-receiver)
- [use-any](#use-any)
- [use-errors-new](#use-errors-new)
- [use-fmt-print](#use-fmt-print)
- [use-waitgroup-go](#use-waitgroup-go)
- [useless-break](#useless-break)
- [useless-fallthrough](#useless-fallthrough)
- [var-declaration](#var-declaration)
- [var-naming](#var-naming)
- [waitgroup-by-value](#waitgroup-by-value)

<!-- tocstop -->

## add-constant

_Description_: Suggests using constant for [magic numbers](https://en.wikipedia.org/wiki/Magic_number_(programming)#Unnamed_numerical_constants)
and string literals.

_Configuration_:

- `maxLitCount` (`maxlitcount`, `max-lit-count`): (string) maximum number of instances of a string literal that are tolerated before warn.
- `allowStrs` (`allowstrs`, `allow-strs`): (string) comma-separated list of allowed string literals
- `allowInts` (`allowints`, `allow-ints`): (string) comma-separated list of allowed integers
- `allowFloats` (`allowfloats`, `allow-floats`): (string) comma-separated list of allowed floats
- `ignoreFuncs` (`ignorefuncs`, `ignore-funcs`): (string) comma-separated list of function names regexp patterns to exclude

Configuration examples:

```toml
[rule.add-constant]
arguments = [
  { maxLitCount = "3", allowStrs = "\"\"", allowInts = "0,1,2", allowFloats = "0.0,0.,1.0,1.,2.0,2.", ignoreFuncs = "os\\.*,fmt\\.Println,make" },
]
```

```toml
[rule.add-constant]
arguments = [
  { max-lit-count = "3", allow-strs = "\"\"", allow-ints = "0,1,2", allow-floats = "0.0,0.,1.0,1.,2.0,2.", ignore-funcs = "os\\.*,fmt\\.Println,make" },
]
```

## argument-limit

_Description_: Warns when a function receives more parameters than the maximum set by the rule's configuration.
Enforcing a maximum number of parameters helps to keep the code readable and maintainable.

_Configuration_: (int) the maximum number of parameters allowed per function.

Configuration example:

```toml
[rule.argument-limit]
arguments = [4]
```

## atomic

_Description_: Check for commonly mistaken usages of the `sync/atomic` package

_Configuration_: N/A

## banned-characters

_Description_: Checks given banned characters in identifiers(func, var, const). Comments are not checked.

_Configuration_: This rule requires a slice of strings, the characters to ban.

Configuration example:

```toml
[rule.banned-characters]
arguments = ["Ω", "Σ", "σ"]
```

## bare-return

_Description_: Warns on bare (a.k.a. naked) returns

_Configuration_: N/A

## blank-imports

_Description_: Blank import should be only in a main or test package, or have a comment justifying it.

_Configuration_: N/A

## bool-literal-in-expr

_Description_: Using Boolean literals (`true`, `false`) in logic expressions may make the code less readable.
This rule suggests removing Boolean literals from logic expressions.

### Examples (bool-literal-in-expr)

Before (violation):

```go
if attachRequired == true {
  // do something
}

if mustReply == false {
  // do something
}
```

After (fixed):

```go
if attachRequired {
  // do something
}

if !mustReply {
  // do something
}
```

_Configuration_: N/A

## call-to-gc

_Description_: Explicitly invoking the garbage collector is, except for specific uses in benchmarking, very dubious.

The garbage collector can be configured through environment variables as [described here](https://golang.org/pkg/runtime/).

_Configuration_: N/A

## cognitive-complexity

_Description_: [Cognitive complexity](https://www.sonarsource.com/docs/CognitiveComplexity.pdf) is a measure of how hard code is to understand.
While cyclomatic complexity is good to measure "testability" of the code,
cognitive complexity aims to provide a more precise measure of the difficulty of understanding the code.
Enforcing a maximum complexity per function helps to keep code readable and maintainable.

_Configuration_: (int) the maximum function complexity

Configuration example:

```toml
[rule.cognitive-complexity]
arguments = [7]
```

## comment-spacings

_Description_: Spots comments of the form:

```go
//This is a malformed comment: no space between // and the start of the sentence
```

_Configuration_: ([]string) list of exceptions. For example, to accept comments of the form

```go
//mypragma: activate something
//+optional
```

You need to add both `"mypragma:"` and `"+optional"` in the configuration

Configuration example:

```toml
[rule.comment-spacings]
arguments = ["mypragma:", "+optional"]
```

## comments-density

_Description_: Spots files not respecting a minimum value for the [_comments lines density_](https://docs.sonarsource.com/sonarqube/latest/user-guide/metric-definitions/)
metric = _comment lines / (lines of code + comment lines) * 100_

_Configuration_: (int) the minimum expected comments lines density.

Configuration example:

```toml
[rule.comments-density]
arguments = [15]
```

## confusing-naming

_Description_: Methods or fields of `struct` that have names different only by capitalization could be confusing.

_Configuration_: N/A

## confusing-results

_Description_: Function or methods that return multiple, no named, values of the same type could induce error.

### Examples (confusing-results)

Before (violation):

```go
// getPos yields the geographical position of this tag.
func (t tag) getPos() (float32, float32)
```

After (fixed):

```go
// getPos yields the geographical position of this tag.
func (t tag) getPos() (longitude float32, latitude float32)
```

_Configuration_: N/A

## constant-logical-expr

_Description_: The rule spots logical expressions that evaluate always to the same value.

_Configuration_: N/A

## context-as-argument

_Description_: By [convention](https://go.dev/wiki/CodeReviewComments#contexts), `context.Context` should be the first parameter of a function.
This rule spots function declarations that do not follow the convention.

_Configuration_:

- `allowTypesBefore` (`allowtypesbefore`, `allow-types-before`): (string) comma-separated list of types that may be before 'context.Context'

Configuration examples:

```toml
[rule.context-as-argument]
arguments = [
  { allowTypesBefore = "*testing.T,*github.com/user/repo/testing.Harness" },
]
```

```toml
[rule.context-as-argument]
arguments = [
  { allow-types-before = "*testing.T,*github.com/user/repo/testing.Harness" },
]
```

## context-keys-type

_Description_: Basic types should not be used as a key in `context.WithValue`.

_Configuration_: N/A

## cyclomatic

_Description_: [Cyclomatic complexity](https://en.wikipedia.org/wiki/Cyclomatic_complexity) is a measure of code complexity.
Enforcing a maximum complexity per function helps to keep code readable and maintainable.

_Configuration_: (int) the maximum function complexity

Configuration example:

```toml
[rule.cyclomatic]
arguments = [3]
```

## datarace

_Description_: This rule spots potential dataraces caused by goroutines capturing (by-reference) particular identifiers of the function from
which goroutines are created.
The rule is able to spot two of such cases: go-routines capturing named return values, and capturing `for-range` values.

_Configuration_: N/A

## deep-exit

_Description_: Packages exposing functions that can stop program execution by exiting are hard to reuse.
This rule looks for program exits in functions other than `main()` or `init()`.

_Configuration_: N/A

## defer

_Description_: This rule warns on some common mistakes when using `defer` statement. It currently alerts on the following situations:

<!-- markdownlint-disable MD013 -->

| name              | description                                                                                                                                                                                     |
| ----------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| call-chain (callChain, callchain)        | even if deferring call-chains of the form `foo()()` is valid, it does not helps code understanding (only the last call is deferred)                                                             |
| loop              | deferring inside loops can be misleading (deferred functions are not executed at the end of the loop iteration but of the current function) and it could lead to exhausting the execution stack |
| method-call (methodCall, methodcall)       | deferring a call to a method can lead to subtle bugs if the method does not have a pointer receiver                                                                                             |
| recover           | calling `recover` outside a deferred function has no effect                                                                                                                                     |
| immediate-recover (immediateRecover, immediaterecover) | calling `recover` at the time a defer is registered, rather than as part of the deferred callback.  e.g. `defer recover()` or equivalent.                                                       |
| return            | returning values form a deferred function has no effect                                                                                                                                         |

<!-- markdownlint-enable MD013 -->

These gotchas are [described here](https://blog.learngoprogramming.com/gotchas-of-defer-in-go-1-8d070894cb01).

_Configuration_: By default, all warnings are enabled but it is possible selectively enable them through configuration.
For example to enable only `call-chain` and `loop`:

Configuration examples:

```toml
[rule.defer]
arguments = [["callChain", "loop"]]
```

```toml
[rule.defer]
arguments = [["call-chain", "loop"]]
```

## dot-imports

_Description_: Importing with `.` makes the programs much harder to understand because it is unclear whether names belong to the current package or
to an imported package.

More [information here](https://go.dev/wiki/CodeReviewComments#import-dot).

_Configuration_:

- `allowedPackages` (`allowedpackages`, `allowed-packages`): (list of strings) comma-separated list of allowed dot import packages

Configuration examples:

```toml
[rule.dot-imports]
arguments = [
  { allowedPackages = [
    "github.com/onsi/ginkgo/v2",
    "github.com/onsi/gomega",
  ] },
]
```

```toml
[rule.dot-imports]
arguments = [
  { allowed-packages = [
    "github.com/onsi/ginkgo/v2",
    "github.com/onsi/gomega",
  ] },
]
```

## duplicated-imports

_Description_: It is possible to unintentionally import the same package twice. This rule looks for packages that are imported two or more times.

_Configuration_: N/A

## early-return

_Description_: In Go it is idiomatic to minimize nesting statements, a typical example is to avoid if-then-else constructions.
This rule spots constructions like

```golang
if cond {
	// do something
} else {
	// do other thing
	return ...
}
```

where the `if` condition may be inverted in order to reduce nesting:

```golang
if !cond {
	// do other thing
	return ...
}

// do something
```

_Configuration_: ([]string) rule flags. Available flags are:

- `preserveScope` (`preservescope`, `preserve-scope`): do not suggest refactorings that would increase variable scope
- `allowJump` (`allowjump`, `allow-jump`): suggest a new jump (`return`, `continue` or `break` statement) if it could unnest multiple statements.
By default, only relocation of _existing_ jumps (i.e. from the `else` clause) are suggested.

Configuration examples:

```toml
[rule.early-return]
arguments = ["preserveScope", "allowJump"]
```

```toml
[rule.early-return]
arguments = ["preserve-scope", "allow-jump"]
```

## empty-block

_Description_: Empty blocks make code less readable and could be a symptom of a bug or unfinished refactoring.

_Configuration_: N/A

## empty-lines

_Description_: Sometimes `gofmt` is not enough to enforce a common formatting of a code-base;
this rule warns when there are heading or trailing newlines in code blocks.

_Configuration_: N/A

## enforce-map-style

_Description_: This rule enforces consistent usage of `make(map[type]type)` or `map[type]type{}` for map initialization.
It does not affect `make(map[type]type, size)` constructions as well as `map[type]type{k1: v1}`.

_Configuration_: (string) Specifies the enforced style for map initialization. The options are:

- "any": No enforcement (default).
- "make": Enforces the usage of `make(map[type]type)`.
- "literal": Enforces the usage of `map[type]type{}`.

Configuration example:

```toml
[rule.enforce-map-style]
arguments = ["make"]
```

## enforce-repeated-arg-type-style

_Description_: This rule is designed to maintain consistency in the declaration of repeated argument and return value types in Go functions.
It supports three styles: 'any', 'short', and 'full'.
The 'any' style is lenient and allows any form of type declaration.
The 'short' style encourages omitting repeated types for conciseness,
whereas the 'full' style mandates explicitly stating the type for each argument
and return value, even if they are repeated, promoting clarity.

_Configuration (1)_: (string) as a single string, it configures both argument
and return value styles. Accepts 'any', 'short', or 'full' (default: 'any').

_Configuration (2)_: (map[string]any) as a map, allows separate configuration
for function arguments and return values. Valid keys are `funcArgStyle` (`funcargstyle`, `func-arg-style`) and
`funcRetValStyle` (`funcretvalstyle`, `func-ret-val-style`), each accepting 'any', 'short', or 'full'. If a key is not
specified, the default value of 'any' is used.

_Note_: The rule applies checks based on the specified styles. For 'full' style,
it flags instances where types are omitted in repeated arguments or return values.
For 'short' style, it highlights opportunities to omit repeated types for brevity.
Incorrect or unknown configuration values will result in an error.

Example (1):

```toml
[rule.enforce-repeated-arg-type-style]
arguments = ["short"]
```

Examples (2):

```toml
[rule.enforce-repeated-arg-type-style]
arguments = [{ funcArgStyle = "full", funcRetValStyle = "short" }]
```

```toml
[rule.enforce-repeated-arg-type-style]
arguments = [{ func-arg-style = "full", func-ret-val-style = "short" }]
```

## enforce-slice-style

_Description_: This rule enforces consistent usage of `make([]type, 0)`, `[]type{}`, or `var []type` for slice initialization.
It does not affect `make([]type, non_zero_len, or_non_zero_cap)` constructions as well as `[]type{v1}`.
Nil slices are always permitted.

_Configuration_: (string) Specifies the enforced style for slice initialization. The options are:

- "any": No enforcement (default).
- "make": Enforces the usage of `make([]type, 0)`.
- "literal": Enforces the usage of `[]type{}`.
- "nil": Enforces the usage of `var []type`.

Configuration example:

```toml
[rule.enforce-slice-style]
arguments = ["make"]
```

## enforce-switch-style

_Description_: This rule enforces consistent usage of `default` on `switch` statements.
It can check for `default` case clause occurrence and/or position in the list of case clauses.

_Configuration_: ([]string) Specifies what to enforced: occurrence and/or position. The, non-mutually exclusive, options are:

- "allowNoDefault": allows `switch` without `default` case clause.
- "allowDefaultNotLast": allows `default` case clause to be not the last clause of the `switch`.

Configuration examples:

To enforce that all `switch` statements have a `default` clause as its the last case clause:

```toml
[rule.enforce-switch-style]
```

To enforce that all `switch` statements have a `default` clause but its position is unimportant:

```toml
[rule.enforce-switch-style]
arguments = ["allowDefaultNotLast"]
```

To enforce that in all `switch` statements with a `default` clause, the `default` is the last case clause:

```toml
[rule.enforce-switch-style]
arguments = ["allowNoDefault"]
```

Notice that a configuration including both options will effectively deactivate the whole rule.

## error-naming

_Description_: By convention, for the sake of readability, variables of type `error` must be named with the prefix `err`.

_Configuration_: N/A

## error-return

_Description_: By convention, for the sake of readability, the errors should be last in the list of returned values by a function.

_Configuration_: N/A

## error-strings

_Description_: By convention, for better readability, error messages should not be capitalized or end with punctuation or a newline.
By default, the rule analyzes functions for creating errors from `fmt`, `errors`, and `github.com/pkg/errors`.
Optionally, the rule can be configured to analyze user functions that create errors.

More [information here](https://go.dev/wiki/CodeReviewComments#error-strings).

_Configuration_: ([]string) the list of additional error functions to check.
The format of values is `package.FunctionName`.

Configuration example:

```toml
[rule.error-strings]
arguments = ["xerrors.Errorf"]
```

## errorf

_Description_: It is possible to get a simpler program by replacing `errors.New(fmt.Sprintf())` with `fmt.Errorf()`.
This rule spots that kind of simplification opportunities.

_Configuration_: N/A

## exported

_Description_: Exported function and methods should have comments. This warns on undocumented exported functions and methods.

More [information here](https://go.dev/wiki/CodeReviewComments#doc-comments).

_Configuration_: ([]string) rule flags.
Please notice that without configuration, the default behavior of the rule is that of its `golint` counterpart.
Available flags are:

- `checkPrivateReceivers` (`checkprivatereceivers`, `check-private-receivers`) enables checking public methods of private types
- `disableStutteringCheck` (`disablestutteringcheck`, `disable-stuttering-check`) disables checking for method names that stutter with the package name
(i.e. avoid failure messages of the form _type name will be used as x.XY by other packages, and that stutters; consider calling this Y_)
- `sayRepetitiveInsteadOfStutters` (`sayrepetitiveinsteadofstutters`, `say-repetitive-instead-of-stutters`) replaces the use of the term _stutters_
by _repetitive_ in failure messages
- `checkPublicInterface` (`checkpublicinterface`, `check-public-interface`) enabled checking public method definitions in public interface types
- `disableChecksOnConstants` (`disablechecksonconstants`, `disable-checks-on-constants`) disable all checks on constant declarations
- `disableChecksOnFunctions` (`disablechecksonfunctions`, `disable-checks-on-functions`) disable all checks on function declarations
- `disableChecksOnMethods` (`disablechecksonmethods`, `disable-checks-on-methods`) disable all checks on method declarations
- `disableChecksOnTypes` (`disablechecksontypes`, `disable-checks-on-types`) disable all checks on type declarations
- `disableChecksOnVariables` (`disablechecksonvariables`, `disable-checks-on-variables`) disable all checks on variable declarations

Configuration examples:

```toml
[rule.exported]
arguments = [
  "checkPrivateReceivers",
  "disableStutteringCheck",
  "checkPublicInterface",
  "disableChecksOnFunctions",
]
```

```toml
[rule.exported]
arguments = [
  "check-private-receivers",
  "disable-stuttering-check",
  "check-public-interface",
  "disable-checks-on-functions",
]
```

## file-header

_Description_: This rule helps to enforce a common header for all source files in a project by spotting those files that do not have the specified header.

_Configuration_: (string) the header to look for in source files.

Configuration example:

```toml
[rule.file-header]
arguments = ["This is the text that must appear at the top of source files."]
```

## file-length-limit

_Description_: This rule enforces a maximum number of lines per file, in order to aid in maintainability and reduce complexity.

_Configuration_:

- `max`: (int) a maximum number of lines in a file. Must be non-negative integers. 0 means the rule is disabled (default `0`);
- `skipComments` (`skipcomments`, `skip-comments`): (bool) if true ignore and do not count lines containing just comments (default `false`);
- `skipBlankLines` (`skipblanklines`, `skip-blank-lines`): (bool) if true ignore and do not count lines made up purely of whitespace (default `false`).

Configuration examples:

```toml
[rule.file-length-limit]
arguments = [{ max = 100, skipComments = true, skipBlankLines = true }]
```

```toml
[rule.file-length-limit]
arguments = [{ max = 100, skip-comments = true, skip-blank-lines = true }]
```

## filename-format

_Description_: enforces conventions on source file names. By default, the rule enforces filenames of the form `^[_A-Za-z0-9][_A-Za-z0-9-]*\.go$`.
Optionally, the rule can be configured to enforce other forms.

_Configuration_: (string) regular expression for source filenames.

Configuration example:

```toml
[rule.filename-format]
arguments = ["^[_a-z][_a-z0-9]*\\.go$"]
```

## flag-parameter

_Description_: If a function controls the flow of another by passing it information on what to do, both functions are said to be [control-coupled](https://en.wikipedia.org/wiki/Coupling_(computer_programming)#Procedural_programming).
Coupling among functions must be minimized for better maintainability of the code.
This rule warns on boolean parameters that create a control coupling.

_Configuration_: N/A

## forbidden-call-in-wg-go

_Description_: Since Go 1.25, it is possible to create goroutines with the method `waitgroup.Go`.
The `Go` method calls a function in a new goroutine and adds (`Add`) that task to the WaitGroup.
When the function returns, the task is removed (`Done`) from the WaitGroup.

This rule ensures that functions don't panic as is specified
in the [documentation of `WaitGroup.Go`](https://pkg.go.dev/sync#WaitGroup.Go).

The rule also warns against a common mistake when refactoring legacy code:
accidentally leaving behind a call to `WaitGroup.Done`, which can cause subtle bugs or panics.

### Examples (forbidden-call-in-wg-go)

Legacy code with a call to `wg.Done`:

```go
wg := sync.WaitGroup{}

wg.Add(1)
go func() {
  doSomething()
  wg.Done()
}()

wg.Wait
```

Refactored, incorrect, code:

```go
wg := sync.WaitGroup{}

wg.Go(func() {
  doSomething()
  wg.Done()
})

wg.Wait
```

Fixed code:

```go
wg := sync.WaitGroup{}

wg.Go(func() {
  doSomething()
})

wg.Wait
```

_Configuration_: N/A

## function-length

_Description_: Functions too long (with many statements and/or lines) can be hard to understand.

_Configuration_: (int,int) the maximum allowed statements and lines. Must be non-negative integers. Set to 0 to disable the check

Configuration example:

```toml
[rule.function-length]
arguments = [10, 0]
```

Will check for functions exceeding 10 statements and will not check the number of lines of functions

## function-result-limit

_Description_: Functions returning too many results can be hard to understand/use.

_Configuration_: (int) the maximum allowed return values

Configuration example:

```toml
[rule.function-result-limit]
arguments = [3]
```

## get-return

_Description_: Typically, functions with names prefixed with _Get_ are supposed to return a value.

_Configuration_: N/A

## identical-branches

_Description_: An `if-then-else` conditional with identical implementations in both branches is an error.

_Configuration_: N/A

## identical-ifelseif-branches

_Description_: an `if ... else if` chain with identical branches makes maintenance harder
and might be a source of bugs. Duplicated branches should be consolidated in one.

_Configuration_: N/A

## identical-ifelseif-conditions

_Description_: an `if ... else if` chain  with identical conditions can lead to
unreachable code and is a potential source of bugs while making the code harder to read and maintain.

_Configuration_: N/A

## identical-switch-branches

_Description_: a `switch` with identical branches makes maintenance harder
and might be a source of bugs. Duplicated branches should be consolidated
in one case clause.

## identical-switch-conditions

_Description_: a `switch` statement with cases with the same condition can lead to
unreachable code and is a potential source of bugs while making the code harder to read and maintain.

_Configuration_: N/A

## if-return

_Description_: Checking if an error is _nil_ to just after return the error or nil is redundant.

_Configuration_: N/A

## import-alias-naming

_Description_: Aligns with Go's naming conventions, as outlined in the official
[blog post](https://go.dev/blog/package-names). It enforces clear and lowercase import alias names, echoing
the principles of good package naming. Users can follow these guidelines by default or define a custom regex rule.
Importantly, aliases with underscores ("_") are always allowed.

_Configuration_ (1): (`string`) as plain string accepts allow regexp pattern for aliases (default: `^[a-z][a-z0-9]{0,}$`).

_Configuration_ (2): (`map[string]string`) as a map accepts two values:

- for a key `allowRegex` (`allowregex`, `allow-regex`) accepts allow regexp pattern
- for a key `denyRegex` (`denyregex`, `deny-regex`) deny regexp pattern

_Note_: If both `allowRegex` and `denyRegex` are provided, the alias must comply with both of them.
If none are given (i.e. an empty map), the default value `^[a-z][a-z0-9]{0,}$` for allowRegex is used.
Unknown keys will result in an error.

Configuration example (1):

```toml
[rule.import-alias-naming]
arguments = ["^[a-z][a-z0-9]{0,}$"]
```

Configuration examples (2):

```toml
[rule.import-alias-naming]
arguments = [{ allowRegex = "^[a-z][a-z0-9]{0,}$", denyRegex = '^v\d+$' }]
```

```toml
[rule.import-alias-naming]
arguments = [{ allow-regex = "^[a-z][a-z0-9]{0,}$", deny-regex = '^v\d+$' }]
```

## import-shadowing

_Description_: In Go it is possible to declare identifiers (packages, structs,
interfaces, parameters, receivers, variables, constants...) that conflict with the
name of an imported package. This rule spots identifiers that shadow an import.

The rule ignores versioned import paths such as `k8s.io/api/core/v1` when `v1` is the package name,
which allows identifiers like `v1`. This is a deliberate trade-off to keep the rule simple.

_Configuration_: N/A

## imports-blocklist

_Description_: Warns when importing block-listed packages.

_Configuration_: block-list of package names (or regular expression package names).

Configuration example:

```toml
[rule.imports-blocklist]
arguments = ["crypto/md5", "crypto/sha1", "crypto/**/pkix"]
```

## increment-decrement

_Description_: By convention, for better readability, incrementing an integer variable by 1 is recommended to be done using the `++` operator.
This rule spots expressions like `i += 1` and `i -= 1` and proposes to change them into `i++` and `i--`.

_Configuration_: N/A

## indent-error-flow

_Description_: To improve the readability of code, it is recommended to reduce the indentation as much as possible.
This rule highlights redundant _else-blocks_ that can be eliminated from the code.

More [information here](https://go.dev/wiki/CodeReviewComments#indent-error-flow).

_Configuration_: ([]string) rule flags. Available flags are:

- `preserveScope` (`preservescope`, `preserve-scope`): do not suggest refactorings that would increase variable scope

Configuration examples:

```toml
[rule.indent-error-flow]
arguments = ["preserveScope"]
```

```toml
[rule.indent-error-flow]
arguments = ["preserve-scope"]
```

## inefficient-map-lookup

_Description_: This rule identifies code that iteratively searches for a key in a map.

This inefficiency is usually introduced when refactoring code from using a slice to a map.
For example if during refactoring the `elements` slice is transformed into a map.

```diff
-       elements             []string
+       elements             map[string]float64
```

and then a loop over `elements` is changed in an obvious but inefficient way:

```diff
-       for _, e := range elements {
+       for e := range elements {
                if e == someStaticValue {
                        // do something
                }
        }
```

Configuration example:

```go
aMap := map[string]bool{}{}
aValue := false

// Inefficient map lookup
for k := range aMap {
  if k == aValue {
    // do something
  }
}

// Simpler and more efficient version
if _, ok := aMap[aValue]; ok {
  // do something
}
```

_Configuration_: N/A

## line-length-limit

_Description_: Warns in the presence of code lines longer than a configured maximum.

_Configuration_: (int) maximum line length in characters.

Configuration example:

```toml
[rule.line-length-limit]
arguments = [80]
```

## max-control-nesting

_Description_: Warns if nesting level of control structures (`if-then-else`, `for`, `switch`) exceeds a given maximum.

_Configuration_: (int) maximum accepted nesting level of control structures (defaults to 5)

Configuration example:

```toml
[rule.max-control-nesting]
arguments = [3]
```

## max-public-structs

_Description_: Packages declaring too many public structs can be hard to understand/use,
and could be a symptom of bad design.

This rule warns on files declaring more than a configured, maximum number of public structs.

_Configuration_: (int) the maximum allowed public structs

Configuration example:

```toml
[rule.max-public-structs]
arguments = [3]
```

## modifies-parameter

_Description_: A function that modifies its parameters can be hard to understand.
It can also be misleading if the arguments are passed by value by the caller.
This rule warns when a function modifies one or more of its parameters or when
parameters are passed to functions that modify them (e.g. `slices.Delete`).

_Configuration_: N/A

## modifies-value-receiver

_Description_: A method that modifies its receiver value can have undesired behavior.
The modification can be also the root of a bug because the actual value receiver could be a copy of that used at the calling site.
This rule warns when a method modifies its receiver.

_Configuration_: N/A

## nested-structs

_Description_: Packages declaring structs that contain other inline struct definitions can be hard to understand/read for other developers.

_Configuration_: N/A

## optimize-operands-order

_Description_: Conditional expressions can be written to take advantage of short circuit evaluation and speed up its average evaluation time
by forcing the evaluation of less time-consuming terms before more costly ones.
This rule spots logical expressions where the order of evaluation of terms seems non optimal.
Please notice that confidence of this rule is low and is up to the user to decide if the suggested rewrite of the expression
keeps the semantics of the original one.

_Configuration_: N/A

Example:

```golang
if isGenerated(content) && !config.IgnoreGeneratedHeader {
```

Swap left and right side :

```golang
if !config.IgnoreGeneratedHeader && isGenerated(content) {
```

## package-comments

_Description_: Packages should have comments. This rule warns on undocumented packages and when packages comments are detached to the `package` keyword.

More [information here](https://go.dev/wiki/CodeReviewComments#package-comments).

_Configuration_: N/A

## package-directory-mismatch

_Description_: It is considered a good practice to name a package after the directory containing it.
This rule warns when the package name declared in the file does not match the name of the directory containing the file.

The following cases are excluded from this check:

- Package `main` (executable packages)
- Files in `testdata` directories (at any level) - by default
- Files directly in `internal` directories (but files in subdirectories of `internal` are checked)

For test files (files with `_test` suffix), package name additionally checked if it matches directory name with  `_test` suffix appended.

The rule normalizes both directory and package names before comparison by removing hyphens (`-`),
underscores (`_`), and dots (`.`). This allows package `foo_barbuz` be equal with directory `foo-bar.buz`.

For files in version directories (`v1`, `v2`, etc.), package name is checked if it matches either the version directory or its parent directory.

_Configuration_: Named arguments for directory exclusions.

Configuration examples:

Default behavior excludes paths containing `testdata`

```toml
[rule.package-directory-mismatch]
```

Ignore specific directories with `ignore-directories`

```toml
[rule.package-directory-mismatch]
arguments = [{ ignore-directories = ["testcases", "testinfo"] }]
```

Include all directories (`testdata` also)

```toml
[rule.package-directory-mismatch]
arguments = [{ ignoreDirectories = [] }]
```

## range-val-address

_Description_: Range variables in a loop are reused at each iteration.
This rule warns when assigning the address of the variable, passing the address to append() or using it in a map.

_Configuration_: N/A

_Note_: This rule is irrelevant for Go 1.22+.

## range-val-in-closure

_Description_: Range variables in a loop are reused at each iteration; therefore a goroutine created in a loop will point to the range variable
with from the upper scope. This way, the goroutine could use the variable with an undesired value.
This rule warns when a range value (or index) is used inside a closure.

_Configuration_: N/A

_Note_: This rule is irrelevant for Go 1.22+.

## range

_Description_: This rule suggests a shorter way of writing ranges that do not use the second value.

_Configuration_: N/A

## receiver-naming

_Description_: By convention, receiver names in a method should reflect their identity.
For example, if the receiver is of type `Parts`, `p` is an adequate name for it.
Contrary to other languages, it is not idiomatic to name receivers as `this` or `self`.

_Configuration_: (optional) list of key-value-pair-map (`[]map[string]any`).

- `maxLength` (`maxlength`, `max-length`): (int) max length of receiver name

Configuration examples:

```toml
[rule.receiver-naming]
arguments = [{ maxLength = 2 }]
```

```toml
[rule.receiver-naming]
arguments = [{ max-length = 2 }]
```

## redefines-builtin-id

_Description_: Constant names like `false`, `true`, `nil`, function names like `append`, `make`,
and basic type names like `bool`, and `byte` are not reserved words of the language; therefore the can be redefined.
Even if possible, redefining these built in names can lead to bugs very difficult to detect.

_Configuration_: N/A

## redundant-build-tag

_Description_: This rule warns about redundant build tag comments `// +build` when `//go:build` is present.
`gofmt` in Go 1.17+ automatically adds the `//go:build` constraint, making the `// +build` comment unnecessary.

_Configuration_: N/A

## redundant-import-alias

_Description_: This rule warns on redundant import aliases. This happens when the alias used on the import statement matches the imported package name.

_Configuration_: N/A

## redundant-test-main-exit

_Description_: This rule warns about redundant `Exit` calls in the `TestMain` function,
as the Go test runner automatically handles program termination starting from Go 1.15.

_Configuration_: N/A

_Note_: This rule is irrelevant for Go versions below 1.15.

## string-format

_Description_: This rule allows you to configure a list of regular expressions that string literals in certain function calls are checked against.
This is geared towards user facing applications where string literals are often used for messages that will be presented to users,
so it may be desirable to enforce consistent formatting.

_Configuration_: Each argument is a slice containing 2-3 strings: a scope, a regex, and an optional error message.

1. The first string defines a **scope**. This controls which string literals the regex will apply to, and is defined as a function argument.
It must contain at least a function name (`core.WriteError`).
Scopes may optionally contain a number specifying which argument in the function to check (`core.WriteError[1]`),
as well as a struct field (`core.WriteError[1].Message`, only works for top level fields).
Function arguments are counted starting at 0, so `[0]` would refer to the first argument, `[1]` would refer to the second, etc.
If no argument number is provided, the first argument will be used (same as `[0]`).
You can use multiple scopes to one regex. Split them by `,` (`core.WriteError,fmt.Errorf`).

2. The second string is a **regular expression** (beginning and ending with a `/` character), which will be used to check the string literals in the scope.
The default semantics is "_strings matching the regular expression are OK_".
If you need to inverse the semantics you can add a `!` just before the first `/`. Examples:

    - with `"/^[A-Z].*$/"` the rule will **accept** strings starting with capital letters
    - with `"!/^[A-Z].*$/"` the rule will a **fail** on strings starting with capital letters

3. The third string (optional) is a **message** containing the purpose for the regex, which will be used in lint errors.

Configuration example:

```toml
[rule.string-format]
arguments = [
  [
    "core.WriteError[1].Message",
    "/^([^A-Z]|$)/",
    "must not start with a capital letter",
  ],
  [
    "fmt.Errorf[0]",
    "/(^|[^\\.!?])$/",
    "must not end in punctuation",
  ],
  [
    "panic",
    "/^[^\\n]*$/",
    "must not contain line breaks",
  ],
  [
    "fmt.Errorf[0],core.WriteError[1].Message",
    "!/^.*%w.*$/",
    "must not contain '%w'",
  ],
]
```

## string-of-int

_Description_: Explicit type conversion `string(i)` where `i` has an integer type other than `rune` might behave not as expected by the developer
(e.g. `string(42)` is not `"42"`). This rule spot that kind of suspicious conversions.

_Configuration_: N/A

## struct-tag

_Description_: The rule spots errors in struct tags.
This is useful because struct tags are not checked at compile time.

The list of [supported tags](https://go.dev/wiki/Well-known-struct-tags):

| Tag           | Documentation                                                            |
| ------------- | ------------------------------------------------------------------------ |
| `asn1`         | <https://pkg.go.dev/encoding/asn1>                                      |
| `bson`         | <https://pkg.go.dev/go.mongodb.org/mongo-driver/bson>                   |
| `cbor`         | <https://pkg.go.dev/github.com/fxamacker/cbor/v2>                   |
| `datastore`    | <https://pkg.go.dev/cloud.google.com/go/datastore>                      |
| `default`      | The type of "default" must match the type of the field.                 |
| `json`         | <https://pkg.go.dev/encoding/json>                                      |
| `mapstructure` | <https://pkg.go.dev/github.com/mitchellh/mapstructure>                  |
| `properties`   | <https://pkg.go.dev/github.com/magiconair/properties#Properties.Decode> |
| `protobuf`     | <https://github.com/golang/protobuf>                                    |
| `required`     | Should be only "true" or "false".                                       |
| `spanner`      | <https://pkg.go.dev/cloud.google.com/go/spanner>                        |
| `toml`         | <https://pkg.go.dev/github.com/pelletier/go-toml/v2>                    |
| `url`          | <https://github.com/google/go-querystring>                              |
| `validate`     | <https://github.com/go-playground/validator>                            |
| `xml`          | <https://pkg.go.dev/encoding/xml>                                       |
| `yaml`         | <https://pkg.go.dev/gopkg.in/yaml.v2>                                   |

_Configuration_: (optional) The list of struct tags that can be accepted by the rule additionally to the supported tags.

Configuration example:

To accept the `inline` option in JSON tags (and `outline` and `gnu` in BSON tags) you must provide the following configuration

```toml
[rule.struct-tag]
arguments = ["json,inline", "bson,outline,gnu"]
```

To prevent a tag from being checked, simply add a `!` before its name.
For example, to instruct the rule not to check `validate` tags
(and accept `outline` and `gnu` in BSON tags) you can provide the following configuration

```toml
[rule.struct-tag]
arguments = ["!validate", "bson,outline,gnu"]
```

## superfluous-else

_Description_: To improve the readability of code, it is recommended to reduce the indentation as much as possible.
This rule highlights redundant _else-blocks_ that can be eliminated from the code.

_Configuration_: ([]string) rule flags. Available flags are:

- `preserveScope` (`preservescope`, `preserve-scope`): (string) do not suggest refactorings that would increase variable scope

Configuration examples:

```toml
[rule.superfluous-else]
arguments = ["preserveScope"]
```

```toml
[rule.superfluous-else]
arguments = ["preserve-scope"]
```

## time-date

_Description_: Reports bad usage of `time.Date`.

_Configuration_: N/A

_Examples_:

- Invalid dates reporting:

  - 0 for the month or day argument
  - out of bounds argument for the month (12), day (31), hour (23), minute (59), or seconds (59)
  - an invalid date: 31st of June, 29th of February in 2023, ...

- Non-decimal integers are used as arguments

  This includes:

  - leading zero notation like using 00 for hours, minutes, and seconds.
  - octal notation 0o1, 0o0 that are often caused by using gofumpt on leading zero notation.
  - padding zeros such as 00123456 that are source of bugs.
  - ... and some other use cases.

<details>
<summary>More information about what is detected and reported</summary>

```go
import "time"

var (
	// here we can imagine zeroes were used for padding purpose
	a = time.Date(2023, 1, 2, 3, 4, 0, 00000000, time.UTC) // 00000000 is octal and equals 0 in decimal
	b = time.Date(2023, 1, 2, 3, 4, 0, 00000006, time.UTC) // 00000006 is octal and equals 6 in decimal
	c = time.Date(2023, 1, 2, 3, 4, 0, 00123456, time.UTC) // 00123456 is octal and equals 42798 in decimal
)
```

But here, what was expected 123456 or 42798 ? It's a source of bugs.

So the rule will report this

> octal notation with padding zeroes found: choose between 123456 and 42798 (decimal value of 123456 octal value)

This rule also reports strange notations used with time.Date.

Example:

```go
import "time"

var _ = time.Date(
	0x7e7,    // hexadecimal notation: use 2023 instead of 0x7e7/
	0b1,      // binary notation: use 1 instead of 0b1/
	0x_2,     // hexadecimal notation: use 2 instead of 0x_2/
	1_3,      // alternative notation: use 13 instead of 1_3/
	1e1,      // exponential notation: use 10 instead of 1e1/
	0.,       // float literal: use 0 instead of 0./
	0x1.Fp+6, // float literal: use 124 instead of 0x1.Fp+6/
	time.UTC)
```

All these are considered to be an uncommon usage of time.Date, are reported with a 0.8 confidence.

Note: even if 00, 01, 02, 03, 04, 05, 06, 07 are octal numbers, they can be considered as valid, and reported with 0.5 confidence.

```go
import "time"

var _ = time.Date(2023, 01, 02, 03, 04, 00, 0, time.UTC)
```

</details>

## time-equal

_Description_: This rule warns when using `==` and `!=` for equality check `time.Time` and suggest to `time.time.Equal` method,
for about information follow [this link](https://pkg.go.dev/time#Time)

_Configuration_: N/A

## time-naming

_Description_: Using unit-specific suffix like "Secs", "Mins", ... when naming variables of type `time.Duration` can be misleading,
this rule highlights those cases.

_Configuration_: N/A

## unchecked-type-assertion

_Description_: This rule checks whether a type assertion result is checked (the `ok` value), preventing unexpected `panic`s.

_Configuration_: list of key-value-pair-map (`[]map[string]any`).

- `acceptIgnoredAssertionResult` (`acceptignoredassertionresult`, `accept-ignored-assertion-result`): (bool) default `false`,
set it to `true` to accept ignored type assertion results like this:

```golang
foo, _ := bar(.*Baz).
//   ^
```

Configuration examples:

```toml
[rule.unchecked-type-assertion]
arguments = [{ acceptIgnoredAssertionResult = true }]
```

```toml
[rule.unchecked-type-assertion]
arguments = [{ accept-ignored-assertion-result = true }]
```

## unconditional-recursion

_Description_: Unconditional recursive calls will produce infinite recursion, thus program stack overflow.
This rule detects and warns about unconditional (direct) recursive calls.

_Configuration_: N/A

## unexported-naming

_Description_: this rule warns on wrongly named un-exported symbols, i.e. un-exported symbols whose name start with a capital letter.

_Configuration_: N/A

## unexported-return

_Description_: This rule warns when an exported function or method returns a value of an un-exported type.

_Configuration_: N/A

## unhandled-error

_Description_: This rule warns when errors returned by a function are not explicitly handled on the caller side.

_Configuration_: function names regexp patterns to ignore

Configuration example:

```toml
[rule.unhandled-error]
arguments = [
  '^os\.(CreateTemp|WriteFile|Chmod)$',
  '^fmt\.Print',
  'myFunction',
  '^net\.',
  '^(bytes\.Buffer|string\.Writer)\.Write(Byte|Rune|String)?$',
]
```

## unnecessary-if

_Description_: Detects unnecessary `if-else` statements that return or assign a boolean value
based on a condition and suggests a simplified, direct return or assignment.
The `if-else` block is redundant because the condition itself is already a boolean expression.
The simplified version is immediately clearer, more idiomatic, and reduces cognitive load for the reader.

### Examples (unnecessary-if)

```go
if y <= 0 {
  z = true
} else {
  z = false
}

if x > 10 {
  return false
} else {
  return true
}
```

Fixed code:

```go
z = y <= 0

return x <= 10
```

_Configuration_: N/A

## unnecessary-format

_Description_: This rule identifies calls to formatting functions where the format string does not contain any formatting verbs
and recommends switching to the non-formatting, more efficient alternative.

_Configuration_: N/A

## unnecessary-stmt

_Description_: This rule suggests to remove redundant statements like a `break` at the end of a case block, for improving the code's readability.

_Configuration_: N/A

## unreachable-code

_Description_: This rule spots and proposes to remove [unreachable code](https://en.wikipedia.org/wiki/Unreachable_code).

_Configuration_: N/A

## unsecure-url-scheme

_Description_: Checks for usage of potentially unsecure URL schemes (`http`, `ws`) in string literals.
Using unencrypted URL schemes can expose sensitive data during transmission and
make applications vulnerable to man-in-the-middle attacks.
Secure alternatives like `https` should be preferred when possible.

_Configuration_: N/A

The rule will not warn on local URLs (`localhost`, `127.0.0.1`).

## unused-parameter

_Description_: This rule warns on unused parameters. Functions or methods with unused parameters can be a symptom of an unfinished refactoring or a bug.

_Configuration_: Supports arguments with single of `map[string]any` with option `allowRegex` (`allowregex`, `allow-regex`) to provide additional
to `_` mask to allowed unused parameter names.

Configuration examples:

This allows any names starting with `_`, not just `_` itself:

```go
func SomeFunc(_someObj *MyStruct) {} // matches rule
```

```toml
[rule.unused-parameter]
arguments = [{ allowRegex = "^_" }]
```

```toml
[rule.unused-parameter]
arguments = [{ allow-regex = "^_" }]
```

## unused-receiver

_Description_: This rule warns on unused method receivers. Methods with unused receivers can be a symptom of an unfinished refactoring or a bug.

_Configuration_:
Supports arguments with single of `map[string]any` with option `allowRegex` to provide additional to `_` mask to allowed unused receiver names.

Configuration examples:

This allows any names starting with `_`, not just `_` itself:

```go
func (_my *MyStruct) SomeMethod() {} // matches rule
```

```toml
[rule.unused-receiver]
arguments = [{ allowRegex = "^_" }]
```

```toml
[rule.unused-receiver]
arguments = [{ allow-regex = "^_" }]
```

## use-any

_Description_: Since Go 1.18, `interface{}` has an alias: `any`. This rule proposes to replace instances of `interface{}` with `any`.

_Configuration_: N/A

## use-errors-new

_Description_: This rule identifies calls to `fmt.Errorf` that can be safely replaced by, the more efficient, `errors.New`.

_Configuration_: N/A

## use-fmt-print

_Description_: This rule proposes to replace calls to built-in `print` and `println` with their equivalents from `fmt` standard package.

`print` and `println` built-in functions are not recommended for use-cases other than
[language bootstrapping and are not guaranteed to stay in the language](https://go.dev/ref/spec#Bootstrapping).

_Configuration_: N/A

## use-waitgroup-go

_Description_: Since Go 1.25 the `sync` package proposes the [`WaitGroup.Go`](https://pkg.go.dev/sync#WaitGroup.Go) method.
This method is a shorter and safer replacement for the idiom `wg.Add ... go { ... wg.Done ... }`.
The rule proposes to replace these legacy idioms with calls to the new method.

_Limitations_: The rule doesn't rely on type information but on variable names to identify waitgroups.
This means the rule search for `wg` (the defacto standard name for wait groups);
if the waitgroup variable is named differently than `wg` the rule will skip it.

_Configuration_: N/A

## useless-break

_Description_: This rule warns on useless `break` statements in case clauses of switch and select statements. Go,
unlike other programming languages like C, only executes statements of the selected case while ignoring the subsequent case clauses.
Therefore, inserting a `break` at the end of a case clause has no effect.

Because `break` statements are rarely used in case clauses, when switch or select statements are inside a for-loop,
the programmer might wrongly assume that a `break` in a case clause will take the control out of the loop.
The rule emits a specific warning for such cases.

_Configuration_: N/A

## useless-fallthrough

_Description_: This rule warns on useless `fallthrough` statements in case clauses of switch statements.
A `fallthrough` is considered _useless_ if it's the single statement of a case clause block.

Go allows `switch` statements with clauses that group multiple cases.
Thus, for example:

```go
switch category {
  case "Lu":
    fallthrough
  case "Ll":
    fallthrough    
  case "Lt":
    fallthrough
  case "Lm": 
    fallthrough
  case "Lo":
    return true
  default:
    return false
}
```

can be written as

```go
switch category {
  case "Lu", "Ll", "Lt", "Lm", "Lo":
      return true
  default:
    return false
}
```

_Configuration_: N/A

## var-declaration

_Description_: This rule proposes simplifications of variable declarations.

_Configuration_: N/A

## var-naming

_Description_: This rule warns when [initialism](https://go.dev/wiki/CodeReviewComments#initialisms), [variable](https://go.dev/wiki/CodeReviewComments#variable-names)
or [package](https://go.dev/wiki/CodeReviewComments#package-names) naming conventions are not followed.
It ignores functions starting with `Example`, `Test`, `Benchmark`, and `Fuzz` in test files, preserving `golint` original behavior.

_Configuration_: This rule accepts two slices of strings and one optional slice containing a single map with named parameters.
(This is because TOML does not support "slice of any," and we maintain backward compatibility with the previous configuration version).
The first slice is an allowlist, and the second one is a blocklist of initialisms.
You can add a boolean parameter `skipInitialismNameChecks` (`skipinitialismnamechecks` or `skip-initialism-name-checks`) to control how names
of functions, variables, consts, and structs handle known initialisms (e.g., JSON, HTTP, etc.) when written in `camelCase`.
When `skipInitialismNameChecks` is set to true, the rule allows names like `readJson`, `HttpMethod` etc.
In the map, you can add a boolean `upperCaseConst` (`uppercaseconst`, `upper-case-const`) parameter to allow `UPPER_CASE` for `const`.
You can also add a boolean `skipPackageNameChecks` (`skippackagenamechecks`, `skip-package-name-checks`) to skip package name checks.
When `skipPackageNameChecks` is false (the default), you can configure
`extraBadPackageNames` (`extrabadpackagenames`, `extra-bad-package-names`)
to forbid using the values from the list as package names additionally
to the standard meaningless ones: "common", "interfaces", "misc",
"types", "util", "utils".
When `skipPackageNameCollisionWithGoStd`
(`skippackagenamecollisionwithgostd`, `skip-package-name-collision-with-go-std`)
is set to true, the rule disables checks on package names that collide
with Go standard library packages.

By default, the rule behaves exactly as the alternative in `golint` but optionally, you can relax it (see [golint/lint/issues/89](https://github.com/golang/lint/issues/89)).

Configuration examples:

```toml
[rule.var-naming]
arguments = [[], [], [{ skipInitialismNameChecks = true }]]
```

```toml
[rule.var-naming]
arguments = [["ID"], ["VM"], [{ upperCaseConst = true }]]
```

```toml
[rule.var-naming]
arguments = [[], [], [{ skipPackageNameChecks = true }]]
```

```toml
[rule.var-naming]
arguments = [[], [], [{ extraBadPackageNames = ["helpers", "models"] }]]
```

```toml
[rule.var-naming]
arguments = [[], [], [{ skip-initialism-name-checks = true }]]
```

```toml
[rule.var-naming]
arguments = [["ID"], ["VM"], [{ upper-case-const = true }]]
```

```toml
[rule.var-naming]
arguments = [[], [], [{ skip-package-name-checks = true }]]
```

```toml
[rule.var-naming]
arguments = [[], [], [{ extra-bad-package-names = ["helpers", "models"] }]]
```

```toml
[rule.var-naming]
arguments = [[], [], [{ skip-package-name-collision-with-go-std = true }]]
```

## waitgroup-by-value

_Description_: Function parameters that are passed by value, are in fact a copy of the original argument.
Passing a copy of a `sync.WaitGroup` is usually not what the developer wants to do.
This rule warns when a `sync.WaitGroup` expected as a by-value parameter in a function or method.

_Configuration_: N/A
