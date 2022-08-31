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
// This is a fork of https://github.com/golang/groupcache/blob/master/singleflight/singleflight.go written in a
// more extendable way
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

func New() *Group {
	var nGroup = &Group{}
	nGroup.mu = &sync.Mutex{}
	nGroup.calls = make(map[string]*call)
	return nGroup
}

type Group struct {
	mu    *sync.Mutex
	calls map[string]*call
}

func (this *Group) Do(key string, fn func(key string) (interface{}, error)) (interface{}, error) {
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

func (this *Group) do(key string, c *call, fn func(key string) (interface{}, error)) {
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

func (this *Group) Forget(key string) {
	this.mu.Lock()
	if c, ok := this.calls[key]; ok {
		c.valid = false
		delete(this.calls, key)
	}
	this.mu.Unlock()
}
