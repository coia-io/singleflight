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
	"errors"
	"fmt"
	"testing"
)

type test struct {
	name string
}

func (t *test) String() string { return t.name }

func TestDo(t *testing.T) {
	name := "zephyr"
	v1, err, _ := Do[string]("key", func() (any, error) { return name, nil })
	if got, want := fmt.Sprintf("%v (%T)", v1, v1), "zephyr (string)"; got != want {
		t.Errorf("Do = %v; want %v", got, want)
	}
	if err != nil {
		t.Errorf("Do error = %v", err)
	}

	v2, err, _ := Do[test]("key", func() (any, error) { return test{name}, nil })
	if got, want := fmt.Sprintf("%v (%T)", v2.String(), v2), "zephyr (singleflight.test)"; got != want {
		t.Errorf("Do = %v; want %v", got, want)
	}
	if err != nil {
		t.Errorf("Do error = %v", err)
	}

	_, err, _ = Do[test]("key", func() (any, error) { return nil, errors.New("custom error") })
	if got, want := err.Error(), "custom error"; got != want {
		t.Errorf("Do = %v; want %v", got, want)
	}

	_, err, _ = Do[test]("key", func() (any, error) { return name, nil })
	if got, want := err.Error(), "want type [singleflight.test], but actual [string]"; got != want {
		t.Errorf("Do = %v; want %v", got, want)
	}
}

func TestDoChan(t *testing.T) {
	name := "zephyr"
	ch1 := DoChan[string]("key", func() (any, error) { return name, nil })
	res1 := <-ch1
	v1 := res1.Val
	err := res1.Err
	if got, want := fmt.Sprintf("%v (%T)", v1, v1), "zephyr (string)"; got != want {
		t.Errorf("DoChan = %v; want %v", got, want)
	}
	if err != nil {
		t.Errorf("DoChan error = %v", err)
	}

	ch2 := DoChan[test]("key", func() (any, error) { return test{name}, nil })
	res2 := <-ch2
	v2 := res2.Val
	err = res2.Err
	if got, want := fmt.Sprintf("%v (%T)", v2.String(), v2), "zephyr (singleflight.test)"; got != want {
		t.Errorf("DoChan = %v; want %v", got, want)
	}
	if err != nil {
		t.Errorf("DoChan error = %v", err)
	}

	ch3 := DoChan[test]("key", func() (any, error) { return nil, errors.New("custom error") })
	res3 := <-ch3
	err = res3.Err
	if got, want := err.Error(), "custom error"; got != want {
		t.Errorf("DoChan = %v; want %v", got, want)
	}

	ch4 := DoChan[test]("key", func() (any, error) { return name, nil })
	res4 := <-ch4
	err = res4.Err
	if got, want := err.Error(), "want type [singleflight.test], but actual [string]"; got != want {
		t.Errorf("DoChan = %v; want %v", got, want)
	}
}
