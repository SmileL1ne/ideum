package mids

import "net/http"

type Mid func(http.Handler) http.Handler

type Chain struct {
	mids []Mid
}

func New(mids ...Mid) Chain {
	var temp []Mid
	return Chain{append(temp, mids...)}
}

func (c Chain) Then(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}

	for i := range c.mids {
		h = c.mids[len(c.mids)-1-i](h)
	}

	return h
}

func (c Chain) ThenFunc(fn http.HandlerFunc) http.Handler {
	if fn == nil {
		return c.Then(nil)
	}
	return c.Then(fn)
}

func (c Chain) Append(mids ...Mid) Chain {
	newMids := make([]Mid, 0, len(c.mids)+len(mids))
	newMids = append(newMids, c.mids...)
	newMids = append(newMids, mids...)

	return Chain{newMids}
}
