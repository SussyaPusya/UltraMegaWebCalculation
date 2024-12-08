package pkg

import "errors"

// Errors for calc
var (
	ErrorEmptyExpr      = errors.New("empty expression")
	ErrorMismatchedParh = errors.New("mismatched parentheses")
	ErrInvalidExpr      = errors.New("invalid expression")
	ErrDivisonByZero    = errors.New("division by zero")
	ErrIvalidOperat     = errors.New("invalid operator")
)

//Errors for main app

var (
	ErrIvalJson = errors.New("invalid json")
)
