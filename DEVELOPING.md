# Developer's Guide

This document explains how to build, test, and develop features for revive.

## Installation

Clone the project:

```
git clone git@github.com:mgechev/revive.git
cd revive
```
## Build

In order to build the project run:
```bash
make build
```

The command will produce the `revive` binary in the root of the project.

## Development of rules

If you want to develop a new rule, follow as an example the already existing rules in the [rule package](https://github.com/mgechev/revive/tree/master/rule).

Each rule needs to implement the `lint.Rule` interface: 
```go
type Rule interface {
	Name() string
	Apply(*File, Arguments) []Failure
}
```
All rules with a configuration must implement `lint.ConfigurableRule` interface:
```go
type ConfigurableRule interface {
	Configure(Arguments) error
}
```

The `Arguments` type is an alias of the type `[]any`. The arguments of the rule are passed from the configuration file.

#### Example

Let's suppose we have developed a rule called `BanStructNameRule` which disallow us to name a structure with a given identifier. We can set the banned identifier by using the TOML configuration file:

```toml
[rule.ban-struct-name]
  arguments = ["Foo"]
```

With the snippet above we:

- Enable the rule with the name `ban-struct-name`. The `Name()` method of our rule should return a string that matches `ban-struct-name`.
- Configure the rule with the argument `Foo`. The list of arguments will be passed to `Apply(*File, Arguments)` together with the target file we're linting currently.

A sample rule implementation can be found [here](/rule/argument_limit.go).


## Development of formatters

If you want to develop a new formatter, follow as an example the already existing formatters in the [formatter package](https://github.com/mgechev/revive/tree/master/formatter).

All formatters should implement the following interface:

```go
type Formatter interface {
	Format(<-chan Failure, RulesConfig) (string, error)
	Name() string
}
```
