//go:generate go-enum
package main

// ENUM(add, sub, mul, div)
type Operation string

// ENUM(op, x, y)
type Argument string
