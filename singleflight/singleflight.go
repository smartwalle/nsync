/*
Copyright 2012 Google Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package singleflight provides a duplicate function call suppression
// mechanism.
//
// This is a fork of https://github.com/golang/groupcache/blob/master/singleflight/singleflight.go
package singleflight

import (
	"runtime/debug"
	"sync"
)

type call struct {
	w     sync.WaitGroup
	value interface{}
	err   error
	valid bool
}

func New() Group[string] {
	return NewGroup[string]()
}

func NewGroup[K Key]() Group[K] {
	var nGroup = &group[K]{}
	nGroup.calls = make(map[K]*call)
	return nGroup
}

type group[K Key] struct {
	mu    sync.Mutex
	calls map[K]*call
}

func (this *group[K]) Do(key K, fn func(key K) (interface{}, error)) (interface{}, error) {
	this.mu.Lock()

	if c, ok := this.calls[key]; ok {
		this.mu.Unlock()
		c.w.Wait()

		if err, ok := c.err.(*stackError); ok {
			panic(err)
		}
		return c.value, c.err
	}

	var c = &call{}
	c.valid = true
	c.w.Add(1)
	this.calls[key] = c
	this.mu.Unlock()

	this.do(key, c, fn)

	return c.value, c.err
}

func (this *group[K]) do(key K, c *call, fn func(key K) (interface{}, error)) {
	defer func() {
		c.w.Done()

		this.mu.Lock()
		if c.valid {
			delete(this.calls, key)
		}
		this.mu.Unlock()

		if err, ok := c.err.(*stackError); ok {
			panic(err)
		}
	}()

	defer func() {
		if v := recover(); v != nil {
			c.err = newStackError(v, debug.Stack())
		}
	}()

	c.value, c.err = fn(key)
}

func (this *group[K]) Forget(key K) {
	this.mu.Lock()
	if c, ok := this.calls[key]; ok {
		c.valid = false
		delete(this.calls, key)
	}
	this.mu.Unlock()
}
