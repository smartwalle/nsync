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

// This is a fork of https://github.com/golang/groupcache/blob/master/singleflight/singleflight.go

package singleflight

import (
	"runtime/debug"
	"sync"
)

type call[T any] struct {
	w     sync.WaitGroup
	value T
	err   error
	valid bool
}

func New[V any]() Group[string, V] {
	return NewGroup[string, V]()
}

func NewGroup[K Key, V any]() Group[K, V] {
	var nGroup = &group[K, V]{}
	nGroup.calls = make(map[K]*call[V])
	return nGroup
}

type group[K Key, V any] struct {
	mu    sync.Mutex
	calls map[K]*call[V]
}

func (g *group[K, V]) Do(key K, fn func(key K) (V, error)) (V, error) {
	g.mu.Lock()

	if c, ok := g.calls[key]; ok {
		g.mu.Unlock()
		c.w.Wait()

		if err, ok := c.err.(*stackError); ok {
			panic(err)
		}
		return c.value, c.err
	}

	var c = &call[V]{}
	c.valid = true
	c.w.Add(1)
	g.calls[key] = c
	g.mu.Unlock()

	g.do(key, c, fn)

	return c.value, c.err
}

func (g *group[K, V]) do(key K, c *call[V], fn func(key K) (V, error)) {
	defer func() {
		c.w.Done()

		g.mu.Lock()
		if c.valid {
			delete(g.calls, key)
		}
		g.mu.Unlock()

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

func (g *group[K, V]) Forget(key K) {
	g.mu.Lock()
	if c, ok := g.calls[key]; ok {
		c.valid = false
		delete(g.calls, key)
	}
	g.mu.Unlock()
}
