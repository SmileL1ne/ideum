package mids

import "net/http"

/*
	Mids is the package that provides tools for 'chaining'
	middlewares into sequence to call them in particular order
	without making the code too long
*/

// Mid type is essentially a middleware
type Mid func(http.Handler) http.Handler

type Chain struct {
	mids []Mid
}

// New returns new chain of middlewares
func New(mids ...Mid) Chain {
	var temp []Mid
	return Chain{append(temp, mids...)}
}

// Then sends 'h' to every middleware starting from the last
// and ending with first in the chain
//
// It treats nil input as http.DefaultServerMux
func (c Chain) Then(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}

	for i := range c.mids {
		h = c.mids[len(c.mids)-1-i](h)
	}

	return h
}

// ThenFunc is working essentially the same as Then but
// accepting http.HandlerFunc as input
func (c Chain) ThenFunc(fn http.HandlerFunc) http.Handler {
	if fn == nil {
		return c.Then(nil)
	}
	return c.Then(fn)
}

// Append appends existing chain with inputed middlewares
func (c Chain) Append(mids ...Mid) Chain {
	newMids := make([]Mid, 0, len(c.mids)+len(mids))
	newMids = append(newMids, c.mids...)
	newMids = append(newMids, mids...)

	return Chain{newMids}
}
