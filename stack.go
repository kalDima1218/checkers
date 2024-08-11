package main

type stack interface {
	push(string)
	top() string
	pop()
	empty() bool
}
