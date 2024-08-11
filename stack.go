package main

type stack interface {
	push(string) string
	top() string
	pop()
	empty() bool
}
