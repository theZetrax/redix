// Description: This file contains the main storage engine of the application
// Author: Zablon Dawit
// Date Created: Mar-30-2024
package repository

import (
	"errors"
	"sync"
	"time"
)

type SetOptions struct {
	HasTimeout bool
	Timeout    int
}

type StoreTimeout struct {
	timeout   int
	createdAt time.Time
}

type Store struct {
	timeout map[string]StoreTimeout
	store   map[string]string
	mutex   sync.Mutex
}

func NewStore() *Store {
	return &Store{
		store:   make(map[string]string),
		timeout: make(map[string]StoreTimeout),
	}
}

func (s *Store) hasTimeout(key string) (StoreTimeout, bool) {
	timeout, ok := s.timeout[key]
	return timeout, ok
}

func (s *Store) setTimeout(key string, timeout int) {
	s.timeout[key] = StoreTimeout{
		timeout:   timeout,
		createdAt: time.Now(),
	}
}

// ExpiredTimeout checks if the timeout for a key has expired
func (s *Store) ExpiredTimeout(key string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	timeout, ok := s.hasTimeout(key)
	if !ok {
		return false
	}

	expired := int(time.Since(timeout.createdAt).Milliseconds()) > timeout.timeout

	// if timed out, delete the key
	if expired {
		delete(s.store, key)
		delete(s.timeout, key)
	}

	return expired
}

func (s *Store) Set(key, value string, opts SetOptions) {
	s.store[key] = value
	if opts.HasTimeout {
		s.setTimeout(key, opts.Timeout)
	}
}

func (s *Store) Get(key string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	val, ok := s.store[key]
	if !ok {
		return "", errors.New("Key not found: " + key)
	}
	return val, nil
}

func (s *Store) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.store, key)
	delete(s.timeout, key)
}
