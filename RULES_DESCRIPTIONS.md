# Description of available rules

List of all available rules.


- [Description of available rules](#description-of-available-rules)
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
  - [comments-density](#comment-spacings)
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
  - [error-naming](#error-naming)
  - [error-return](#error-return)
  - [error-strings](#error-strings)
  - [errorf](#errorf)
  - [exported](#exported)
  - [file-header](#file-header)
  - [flag-parameter](#flag-parameter)
  - [function-length](#function-length)
  - [function-result-limit](#function-result-limit)
  - [get-return](#get-return)
  - [identical-branches](#identical-branches)
  - [if-return](#if-return)
  - [import-alias-naming](#import-alias-naming)
  - [import-shadowing](#import-shadowing)
  - [imports-blocklist](#imports-blocklist)
  - [increment-decrement](#increment-decrement)
  - [indent-error-flow](#indent-error-flow)
  - [line-length-limit](#line-length-limit)
  - [max-control-nesting](#max-control-nesting)
  - [max-public-structs](#max-public-structs)
  - [modifies-parameter](#modifies-parameter)
  - [modifies-value-receiver](#modifies-value-receiver)
  - [nested-structs](#nested-structs)
  - [optimize-operands-order](#optimize-operands-order)
  - [package-comments](#package-comments)
  - [range-val-address](#range-val-address)
  - [range-val-in-closure](#range-val-in-closure)
  - [range](#range)
  - [receiver-naming](#receiver-naming)
  - [redefines-builtin-id](#redefines-builtin-id)
  - [redundant-import-alias](#redundant-import-alias)
  - [string-format](#string-format)
  - [string-of-int](#string-of-int)
  - [struct-tag](#struct-tag)
  - [superfluous-else](#superfluous-else)
  - [time-equal](#time-equal)
  - [time-naming](#time-naming)
  - [unchecked-type-assertion](#unchecked-type-assertion)
  - [unconditional-recursion](#unconditional-recursion)
  - [unexported-naming](#unexported-naming)
  - [unexported-return](#unexported-return)
  - [unhandled-error](#unhandled-error)
  - [unnecessary-stmt](#unnecessary-stmt)
  - [unreachable-code](#unreachable-code)
  - [unused-parameter](#unused-parameter)
  - [unused-receiver](#unused-receiver)
  - [use-any](#use-any)
  - [useless-break](#useless-break)
  - [var-declaration](#var-declaration)
  - [var-naming](#var-naming)
  - [waitgroup-by-value](#waitgroup-by-value)

## add-constant

_Description_: Suggests using constant for [magic numbers](https://en.wikipedia.org/wiki/Magic_number_(programming)#Unnamed_numerical_constants) and string literals.

_Configuration_:

* `maxLitCount` : (string) maximum number of instances of a string literal that are tolerated before warn.
* `allowStr`: (string) comma-separated list of allowed string literals
* `allowInts`: (string) comma-separated list of allowed integers
* `allowFloats`: (string) comma-separated list of allowed floats
* `ignoreFuncs`: (string) comma-separated list of function names regexp patterns to exclude

Example:

```toml
[rule.add-constant]
  arguments = [{ maxLitCount = "3", allowStrs = "\"\"", allowInts = "0,1,2", allowFloats = "0.0,0.,1.0,1.,2.0,2.", ignoreFuncs = "os\\.*,fmt\\.Println,make" }]
```

## argument-limit

_Description_: Warns when a function receives more parameters than the maximum set by the rule's configuration.
Enforcing a maximum number of parameters helps to keep the code readable and maintainable.

_Configuration_: (int) the maximum number of parameters allowed per function.

Example:

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

Example:

```toml
[rule.banned-characters]
  arguments = ["Ω","Σ","σ"]
```

## bare-return

_Description_: Warns on bare (a.k.a. naked) returns

_Configuration_: N/A

## blank-imports

_Description_: Blank import should be only in a main or test package, or have a comment justifying it.

_Configuration_: N/A

## bool-literal-in-expr

_Description_: Using Boolean literals (`true`, `false`) in logic expressions may make the code less readable. This rule suggests removing Boolean literals from logic expressions.

_Configuration_: N/A

## call-to-gc

_Description_:  Explicitly invoking the garbage collector is, except for specific uses in benchmarking, very dubious.

The garbage collector can be configured through environment variables as described [here](https://golang.org/pkg/runtime/).

_Configuration_: N/A

## cognitive-complexity

_Description_: [Cognitive complexity](https://www.sonarsource.com/docs/CognitiveComplexity.pdf) is a measure of how hard code is to understand.
While cyclomatic complexity is good to measure "testability" of the code, cognitive complexity aims to provide a more precise measure of the difficulty of understanding the code.
Enforcing a maximum complexity per function helps to keep code readable and maintainable.

_Configuration_: (int) the maximum function complexity

Example:

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

Example:

```toml
[rule.comment-spacings]
  arguments = ["mypragma:", "+optional"]
```

## comments-density

_Description_: Spots files not respecting a minimum value for the [_comments lines density_](https://docs.sonarsource.com/sonarqube/latest/user-guide/metric-definitions/) metric = _comment lines / (lines of code + comment lines) * 100_

_Configuration_: (int) the minimum expected comments lines density.

Example:

```toml
[rule.comments-density]
  arguments = [15]
```

## confusing-naming

_Description_: Methods or fields of `struct` that have names different only by capitalization could be confusing.

_Configuration_: N/A

## confusing-results

_Description_: Function or methods that return multiple, no named, values of the same type could induce error.

_Configuration_: N/A

## constant-logical-expr

_Description_: The rule spots logical expressions that evaluate always to the same value.

_Configuration_: N/A

## context-as-argument

_Description_: By [convention](https://github.com/golang/go/wiki/CodeReviewComments#contexts), `context.Context` should be the first parameter of a function. This rule spots function declarations that do not follow the convention.

_Configuration_:

* `allowTypesBefore` : (string) comma-separated list of types that may be before 'context.Context'

Example:

```toml
[rule.context-as-argument]
  arguments = [{allowTypesBefore = "*testing.T,*github.com/user/repo/testing.Harness"}]
```

## context-keys-type

_Description_: Basic types should not be used as a key in `context.WithValue`.

_Configuration_: N/A

## cyclomatic

_Description_: [Cyclomatic complexity](https://en.wikipedia.org/wiki/Cyclomatic_complexity) is a measure of code complexity. Enforcing a maximum complexity per function helps to keep code readable and maintainable.

_Configuration_: (int) the maximum function complexity

Example:

```toml
[rule.cyclomatic]
  arguments = [3]
```

## datarace

_Description_: This rule spots potential dataraces caused by go-routines capturing (by-reference) particular identifiers of the function from which go-routines are created. The rule is able to spot two of such cases: go-routines capturing named return values, and capturing `for-range` values.

_Configuration_: N/A

## deep-exit

_Description_: Packages exposing functions that can stop program execution by exiting are hard to reuse. This rule looks for program exits in functions other than `main()` or `init()`.

_Configuration_: N/A

## defer

_Description_: This rule warns on some common mistakes when using `defer` statement. It currently alerts on the following situations:

| name              | description                                                                                                                                                                                     |
|-------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| call-chain        | even if deferring call-chains of the form `foo()()` is valid, it does not helps code understanding (only the last call is deferred)                                                             |
| loop              | deferring inside loops can be misleading (deferred functions are not executed at the end of the loop iteration but of the current function) and it could lead to exhausting the execution stack |
| method-call       | deferring a call to a method can lead to subtle bugs if the method does not have a pointer receiver                                                                                             |
| recover           | calling `recover` outside a deferred function has no effect                                                                                                                                     |
| immediate-recover | calling `recover` at the time a defer is registered, rather than as part of the deferred callback.  e.g. `defer recover()` or equivalent.                                                       |
| return            | returning values form a deferred function has no effect                                                                                                                                         |

These gotchas are described [here](https://blog.learngoprogramming.com/gotchas-of-defer-in-go-1-8d070894cb01)

_Configuration_: by default all warnings are enabled but it is possible selectively enable them through configuration. For example to enable only `call-chain` and `loop`:

```toml
[rule.defer]
  arguments = [["call-chain","loop"]]
```

## dot-imports

_Description_: Importing with `.` makes the programs much harder to understand because it is unclear whether names belong to the current package or to an imported package.

More information [here](https://github.com/golang/go/wiki/CodeReviewComments#import-dot)

_Configuration_:

* `allowedPackages`: (list of strings) comma-separated list of allowed dot import packages

Example:

```toml
[rule.dot-imports]
  arguments = [{ allowedPackages = ["github.com/onsi/ginkgo/v2","github.com/onsi/gomega"] }]
```

## duplicated-imports

_Description_: It is possible to unintentionally import the same package twice. This rule looks for packages that are imported two or more times.

_Configuration_: N/A

## early-return

_Description_: In Go it is idiomatic to minimize nesting statements, a typical example is to avoid if-then-else constructions. This rule spots constructions like
```go
if cond {
  // do something
} else {
  // do other thing
  return ...
}
```
where the `if` condition may be inverted in order to reduce nesting:
```go
if !cond {
  // do other thing
  return ...
}

// do something
```

_Configuration_: ([]string) rule flags. Available flags are:

* _preserveScope_: do not suggest refactorings that would increase variable scope

Example:

```toml
[rule.early-return]
  arguments = ["preserveScope"]
```

## empty-block

_Description_: Empty blocks make code less readable and could be a symptom of a bug or unfinished refactoring.

_Configuration_: N/A

## empty-lines

_Description_: Sometimes `gofmt` is not enough to enforce a common formatting of a code-base; this rule warns when there are heading or trailing newlines in code blocks.

_Configuration_: N/A

## enforce-map-style

_Description_: This rule enforces consistent usage of `make(map[type]type)` or `map[type]type{}` for map initialization. It does not affect `make(map[type]type, size)` constructions as well as `map[type]type{k1: v1}`.

_Configuration_: (string) Specifies the enforced style for map initialization. The options are:
- "any": No enforcement (default).
- "make": Enforces the usage of `make(map[type]type)`.
- "literal": Enforces the usage of `map[type]type{}`.

Example:

```toml
[rule.enforce-map-style]
  arguments = ["make"]
```

## enforce-repeated-arg-type-style

_Description_: This rule is designed to maintain consistency in the declaration
of repeated argument and return value types in Go functions. It supports three styles:
'any', 'short', and 'full'. The 'any' style is lenient and allows any form of type
declaration. The 'short' style encourages omitting repeated types for conciseness,
whereas the 'full' style mandates explicitly stating the type for each argument
and return value, even if they are repeated, promoting clarity.

_Configuration (1)_: (string) as a single string, it configures both argument
and return value styles. Accepts 'any', 'short', or 'full' (default: 'any').

_Configuration (2)_: (map[string]any) as a map, allows separate configuration
for function arguments and return values. Valid keys are "funcArgStyle" and
"funcRetValStyle", each accepting 'any', 'short', or 'full'. If a key is not
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

Example (2):

```toml
[rule.enforce-repeated-arg-type-style]
  arguments = [{ funcArgStyle = "full", funcRetValStyle = "short" }]
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

Example:

```toml
[rule.enforce-slice-style]
  arguments = ["make"]
```

## error-naming

_Description_: By convention, for the sake of readability, variables of type `error` must be named with the prefix `err`.

_Configuration_: N/A

## error-return

_Description_: By convention, for the sake of readability, the errors should be last in the list of returned values by a function.

_Configuration_: N/A

## error-strings

_Description_: By convention, for better readability, error messages should not be capitalized or end with punctuation or a newline.

More information [here](https://github.com/golang/go/wiki/CodeReviewComments#error-strings)

_Configuration_: N/A

## errorf

_Description_: It is possible to get a simpler program by replacing `errors.New(fmt.Sprintf())` with `fmt.Errorf()`. This rule spots that kind of simplification opportunities.

_Configuration_: N/A

## exported

_Description_: Exported function and methods should have comments. This warns on undocumented exported functions and methods.

More information [here](https://github.com/golang/go/wiki/CodeReviewComments#doc-comments)

_Configuration_: ([]string) rule flags.
Please notice that without configuration, the default behavior of the rule is that of its `golint` counterpart.
Available flags are:

* _checkPrivateReceivers_ enables checking public methods of private types
* _disableStutteringCheck_ disables checking for method names that stutter with the package name (i.e. avoid failure messages of the form _type name will be used as x.XY by other packages, and that stutters; consider calling this Y_)
* _sayRepetitiveInsteadOfStutters_ replaces the use of the term _stutters_ by _repetitive_ in failure messages

Example:

```toml
[rule.exported]
  arguments = ["checkPrivateReceivers", "disableStutteringCheck"]
```

## file-header

_Description_: This rule helps to enforce a common header for all source files in a project by spotting those files that do not have the specified header.

_Configuration_: (string) the header to look for in source files.

Example:

```toml
[rule.file-header]
  arguments = ["This is the text that must appear at the top of source files."]
```

## flag-parameter

_Description_: If a function controls the flow of another by passing it information on what to do, both functions are said to be [control-coupled](https://en.wikipedia.org/wiki/Coupling_(computer_programming)#Procedural_programming).
Coupling among functions must be minimized for better maintainability of the code.
This rule warns on boolean parameters that create a control coupling.

_Configuration_: N/A

## function-length

_Description_: Functions too long (with many statements and/or lines) can be hard to understand.

_Configuration_: (int,int) the maximum allowed statements and lines. Must be non-negative integers. Set to 0 to disable the check

Example:

```toml
[rule.function-length]
  arguments = [10, 0]
```
Will check for functions exceeding 10 statements and will not check the number of lines of functions

## function-result-limit

_Description_: Functions returning too many results can be hard to understand/use.

_Configuration_: (int) the maximum allowed return values

Example:

```toml
[rule.function-result-limit]
  arguments = [3]
```

## get-return

_Description_: Typically, functions with names prefixed with _Get_ are supposed to return a value.

_Configuration_: N/A

## identical-branches

_Description_: an `if-then-else` conditional with identical implementations in both branches is an error.

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
* for a key "allowRegex" accepts allow regexp pattern
* for a key "denyRegex deny regexp pattern

_Note_: If both `allowRegex` and `denyRegex` are provided, the alias must comply with both of them.
If none are given (i.e. an empty map), the default value `^[a-z][a-z0-9]{0,}$` for allowRegex is used.
Unknown keys will result in an error.

Example (1):

```toml
[rule.import-alias-naming]
  arguments = ["^[a-z][a-z0-9]{0,}$"]
```

Example (2):

```toml
[rule.import-alias-naming]
  arguments = [{ allowRegex = "^[a-z][a-z0-9]{0,}$", denyRegex = '^v\d+$' }]
```

## import-shadowing

_Description_: In Go it is possible to declare identifiers (packages, structs,
interfaces, parameters, receivers, variables, constants...) that conflict with the
name of an imported package. This rule spots identifiers that shadow an import.

_Configuration_: N/A

## imports-blocklist

_Description_: Warns when importing block-listed packages.

_Configuration_: block-list of package names (or regular expression package names).

Example:

```toml
[rule.imports-blocklist]
  arguments = ["crypto/md5", "crypto/sha1", "crypto/**/pkix"]
```

## increment-decrement

_Description_: By convention, for better readability, incrementing an integer variable by 1 is recommended to be done using the `++` operator.
This rule spots expressions  like `i += 1` and `i -= 1` and proposes to change them into `i++` and `i--`.

_Configuration_: N/A

## indent-error-flow

_Description_: To improve the readability of code, it is recommended to reduce the indentation as much as possible.
This rule highlights redundant _else-blocks_ that can be eliminated from the code.

More information [here](https://github.com/golang/go/wiki/CodeReviewComments#indent-error-flow)

_Configuration_: ([]string) rule flags. Available flags are:

* _preserveScope_: do not suggest refactorings that would increase variable scope

Example:

```toml
[rule.indent-error-flow]
  arguments = ["preserveScope"]
```

## line-length-limit

_Description_: Warns in the presence of code lines longer than a configured maximum.

_Configuration_: (int) maximum line length in characters.

Example:

```toml
[rule.line-length-limit]
  arguments = [80]
```

## max-control-nesting

_Description_: Warns if nesting level of control structures (`if-then-else`, `for`, `switch`) exceeds a given maximum.

_Configuration_: (int) maximum accepted nesting level of control structures (defaults to 5)

Example:

```toml
[rule.max-control-nesting]
  arguments = [3]
```

## max-public-structs

_Description_: Packages declaring too many public structs can be hard to understand/use,
and could be a symptom of bad design.

This rule warns on files declaring more than a configured, maximum number of public structs.

_Configuration_: (int) the maximum allowed public structs

Example:

```toml
[rule.max-public-structs]
  arguments = [3]
```

## modifies-parameter

_Description_: A function that modifies its parameters can be hard to understand. It can also be misleading if the arguments are passed by value by the caller.
This rule warns when a function modifies one or more of its parameters.

_Configuration_: N/A

## modifies-value-receiver

_Description_: A method that modifies its receiver value can have undesired behavior. The modification can be also the root of a bug because the actual value receiver could be a copy of that used at the calling site.
This rule warns when a method modifies its receiver.

_Configuration_: N/A

## nested-structs

_Description_: Packages declaring structs that contain other inline struct definitions can be hard to understand/read for other developers.

_Configuration_: N/A

## optimize-operands-order

_Description_: conditional expressions can be written to take advantage of short circuit evaluation and speed up its average evaluation time by forcing the evaluation of less time-consuming terms before more costly ones. This rule spots logical expressions where the order of evaluation of terms seems non optimal. Please notice that confidence of this rule is low and is up to the user to decide if the suggested rewrite of the expression keeps the semantics of the original one.

_Configuration_: N/A

Example:

```go
if isGenerated(content) && !config.IgnoreGeneratedHeader {
```

Swap left and right side :

```go
if !config.IgnoreGeneratedHeader && isGenerated(content) {
```

## package-comments

_Description_: Packages should have comments. This rule warns on undocumented packages and when packages comments are detached to the `package` keyword.

More information [here](https://github.com/golang/go/wiki/CodeReviewComments#package-comments)

_Configuration_: N/A

## range-val-address

_Description_: Range variables in a loop are reused at each iteration. This rule warns when assigning the address of the variable, passing the address to append() or using it in a map.

_Configuration_: N/A

## range-val-in-closure

_Description_: Range variables in a loop are reused at each iteration; therefore a goroutine created in a loop will point to the range variable with from the upper scope. This way, the goroutine could use the variable with an undesired value.
This rule warns when a range value (or index) is used inside a closure

_Configuration_: N/A

## range

_Description_: This rule suggests a shorter way of writing ranges that do not use the second value.

_Configuration_: N/A

## receiver-naming

_Description_: By convention, receiver names in a method should reflect their identity. For example, if the receiver is of type `Parts`, `p` is an adequate name for it. Contrary to other languages, it is not idiomatic to name receivers as `this` or `self`.

_Configuration_: N/A

## redefines-builtin-id

_Description_: Constant names like `false`, `true`, `nil`, function names like `append`, `make`, and basic type names like `bool`, and `byte` are not reserved words of the language; therefore the can be redefined.
Even if possible, redefining these built in names can lead to bugs very difficult to detect.

_Configuration_: N/A

## redundant-import-alias

_Description_: This rule warns on redundant import aliases. This happens when the alias used on the import statement matches the imported package name.

_Configuration_: N/A

## string-format

_Description_: This rule allows you to configure a list of regular expressions that string literals in certain function calls are checked against.
This is geared towards user facing applications where string literals are often used for messages that will be presented to users, so it may be desirable to enforce consistent formatting.

_Configuration_: Each argument is a slice containing 2-3 strings: a scope, a regex, and an optional error message.

1. The first string defines a **scope**. This controls which string literals the regex will apply to, and is defined as a function argument. It must contain at least a function name (`core.WriteError`). Scopes may optionally contain a number specifying which argument in the function to check (`core.WriteError[1]`), as well as a struct field (`core.WriteError[1].Message`, only works for top level fields). Function arguments are counted starting at 0, so `[0]` would refer to the first argument, `[1]` would refer to the second, etc. If no argument number is provided, the first argument will be used (same as `[0]`).

2. The second string is a **regular expression** (beginning and ending with a `/` character), which will be used to check the string literals in the scope. The default semantics is "_strings matching the regular expression are OK_". If you need to inverse the semantics you can add a `!` just before the first `/`. Examples:

    * with `"/^[A-Z].*$/"` the rule will **accept** strings starting with capital letters
    * with `"!/^[A-Z].*$/"` the rule will a **fail** on strings starting with capital letters

3. The third string (optional) is a **message** containing the purpose for the regex, which will be used in lint errors.

Example:

```toml
[rule.string-format]
  arguments = [
    ["core.WriteError[1].Message", "/^([^A-Z]|$)/", "must not start with a capital letter"],
    ["fmt.Errorf[0]", "/(^|[^\\.!?])$/", "must not end in punctuation"],
    ["panic", "/^[^\\n]*$/", "must not contain line breaks"],
  ]
```

## string-of-int

_Description_:  explicit type conversion `string(i)` where `i` has an integer type other than `rune` might behave not as expected by the developer (e.g. `string(42)` is not `"42"`). This rule spot that kind of suspicious conversions.

_Configuration_: N/A

## struct-tag

_Description_: Struct tags are not checked at compile time.
This rule, checks and warns if it finds errors in common struct tags types like: asn1, default, json, protobuf, xml, yaml.

_Configuration_: (optional) list of user defined options.

Example:
To accept the `inline` option in JSON tags (and `outline` and `gnu` in BSON tags) you must provide the following configuration

```toml
[rule.struct-tag]
  arguments = ["json,inline", "bson,outline,gnu"]
```

## superfluous-else

_Description_: To improve the readability of code, it is recommended to reduce the indentation as much as possible.
This rule highlights redundant _else-blocks_ that can be eliminated from the code.

_Configuration_: ([]string) rule flags. Available flags are:

* _preserveScope_: do not suggest refactorings that would increase variable scope

Example:

```toml
[rule.superfluous-else]
  arguments = ["preserveScope"]
```

## time-equal

_Description_: This rule warns when using `==` and `!=` for equality check `time.Time` and suggest to `time.time.Equal` method, for about information follow this [link](https://pkg.go.dev/time#Time)

_Configuration_: N/A

## time-naming

_Description_: Using unit-specific suffix like "Secs", "Mins", ... when naming variables of type `time.Duration` can be misleading, this rule highlights those cases.

_Configuration_: N/A

## unchecked-type-assertion

_Description_: This rule checks whether a type assertion result is checked (the `ok` value), preventing unexpected `panic`s.

_Configuration_: list of key-value-pair-map (`[]map[string]any`).

- `acceptIgnoredAssertionResult` : (bool) default `false`, set it to `true` to accept ignored type assertion results like this:

```go
foo, _ := bar(.*Baz).
//   ^
```

Example:

```toml
[rule.unchecked-type-assertion]
  arguments = [{acceptIgnoredAssertionResult=true}]
```

## unconditional-recursion

_Description_: Unconditional recursive calls will produce infinite recursion, thus program stack overflow. This rule detects and warns about unconditional (direct) recursive calls.

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

Example:

```toml
[rule.unhandled-error]
  arguments = ["os\.(Create|WriteFile|Chmod)", "fmt\.Print", "myFunction", "net\..*", "bytes\.Buffer\.Write"]
```

## unnecessary-stmt

_Description_: This rule suggests to remove redundant statements like a `break` at the end of a case block, for improving the code's readability.

_Configuration_: N/A

## unreachable-code

_Description_: This rule spots and proposes to remove [unreachable code](https://en.wikipedia.org/wiki/Unreachable_code).

_Configuration_: N/A

## unused-parameter

_Description_: This rule warns on unused parameters. Functions or methods with unused parameters can be a symptom of an unfinished refactoring or a bug.

_Configuration_: Supports arguments with single of `map[string]any` with option `allowRegex` to provide additional to `_` mask to allowed unused parameter names, for example:

```toml
[rule.unused-parameter]
    Arguments = [{ allowRegex = "^_" }]
```

allows any names started with `_`, not just `_` itself:

```go
func SomeFunc(_someObj *MyStruct) {} // matches rule
```

## unused-receiver

_Description_: This rule warns on unused method receivers. Methods with unused receivers can be a symptom of an unfinished refactoring or a bug.

_Configuration_: Supports arguments with single of `map[string]any` with option `allowRegex` to provide additional to `_` mask to allowed unused receiver names, for example:

```toml
[rule.unused-receiver]
    Arguments = [{ allowRegex = "^_" }]
```

allows any names started with `_`, not just `_` itself:

```go
func (_my *MyStruct) SomeMethod() {} // matches rule
```

## use-any

_Description_: Since Go 1.18, `interface{}` has an alias: `any`. This rule proposes to replace instances of `interface{}` with `any`.

_Configuration_: N/A

## useless-break

_Description_: This rule warns on useless `break` statements in case clauses of switch and select statements. Go, unlike other programming languages like C, only executes statements of the selected case while ignoring the subsequent case clauses.
Therefore, inserting a `break` at the end of a case clause has no effect.

Because `break` statements are rarely used in case clauses, when switch or select statements are inside a for-loop, the programmer might wrongly assume that a `break` in a case clause will take the control out of the loop.
The rule emits a specific warning for such cases.

_Configuration_: N/A

## var-declaration

_Description_: This rule proposes simplifications of variable declarations.

_Configuration_: N/A

## var-naming

_Description_: This rule warns when [initialism](https://github.com/golang/go/wiki/CodeReviewComments#initialisms), [variable](https://github.com/golang/go/wiki/CodeReviewComments#variable-names) or [package](https://github.com/golang/go/wiki/CodeReviewComments#package-names) naming conventions are not followed.

_Configuration_: This rule accepts two slices of strings and one optional slice with single map with named parameters.
(it's due to TOML hasn't "slice of any" and we keep backward compatibility with previous config version)
First slice is an allowlist and second one is a blocklist of initialisms.
In map, you can add "upperCaseConst=true" parameter to allow `UPPER_CASE` for `const`
By default, the rule behaves exactly as the alternative in `golint` but optionally, you can relax it (see [golint/lint/issues/89](https://github.com/golang/lint/issues/89))

Example:

```toml
[rule.var-naming]
  arguments = [["ID"], ["VM"], [{upperCaseConst=true}]]
```

You can also add "skipPackageNameChecks=true" to skip package name checks.

Example:


```toml
[rule.var-naming]
  arguments = [[], [], [{skipPackageNameChecks=true}]]
```

## waitgroup-by-value

_Description_: Function parameters that are passed by value, are in fact a copy of the original argument. Passing a copy of a `sync.WaitGroup` is usually not what the developer wants to do.
This rule warns when a `sync.WaitGroup` expected as a by-value parameter in a function or method.

_Configuration_: N/A
