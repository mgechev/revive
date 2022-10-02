# Developer's Guide

This document explains how to build, test, and develop features for revive.

## Installation

Clone the project:

```
git clone git@github.com:mgechev/revive.git
cd revive
```

In order to fetch all the dependencies run:

```bash
make install
```

## Build

In order to build the project run:

```bash
make build
```

The command will produce the `revive` binary in the root of the project.

## Development of rules

If you want to develop a new rule, follow as an example the already existing rules in the [rule package](https://github.com/mgechev/revive/tree/master/rule).

All rules should implement the following interface:

```go
type Rule interface {
	Name() string
	Apply(*File, Arguments) []Failure
}
```

## Development of formatters

If you want to develop a new formatter, follow as an example the already existing formatters in the [formatter package](https://github.com/mgechev/revive/tree/master/formatter).

All formatters should implement the following interface:

```go
type Formatter interface {
	Format(<-chan Failure, RulesConfig) (string, error)
	Name() string
}
```
