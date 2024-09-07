/*
 The MIT License (MIT)

 Copyright © 2024 Zephyr <zephyr@coia.top>

 Permission is hereby granted, free of charge, to any person obtaining a copy
 of this software and associated documentation files (the "Software"), to deal
 in the Software without restriction, including without limitation the rights
 to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 copies of the Software, and to permit persons to whom the Software is
 furnished to do so, subject to the following conditions:

 The above copyright notice and this permission notice shall be included in
 all copies or substantial portions of the Software.

 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 THE SOFTWARE.
*/

package singleflight

import (
	"fmt"

	"golang.org/x/sync/singleflight"
)

type Result[T any] struct {
	Val    T
	Err    error
	Shared bool
}

var g singleflight.Group

func Do[T any](key string, fn func() (any, error)) (v T, err error, shared bool) {
	return doCall[T](key, fn)
}

func DoChan[T any](key string, fn func() (any, error)) <-chan Result[T] {
	ch := make(chan Result[T], 1)

	go func() {
		v, err, shared := doCall[T](key, fn)
		ch <- Result[T]{v, err, shared}
	}()

	return ch
}

func doCall[T any](key string, fn func() (any, error)) (v T, err error, shared bool) {
	res, err, shared := g.Do(key, fn)
	if err != nil {
		return
	}
	v, ok := res.(T)
	if !ok {
		err = fmt.Errorf("want type [%T], but actual [%T]", v, res)
	}
	return
}

func Forget(key string) {
	g.Forget(key)
}
