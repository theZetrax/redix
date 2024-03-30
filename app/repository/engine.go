// Description: This file contains the main storage engine of the application
// Author: Zablon Dawit
// Date Created: Mar-30-2024
package repository

import (
	"errors"
	"log"
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

type StorageEngine struct {
	timeout map[string]StoreTimeout
	store   map[string]string
	mutex   sync.Mutex
}

func NewStorageEngine() *StorageEngine {
	return &StorageEngine{
		store:   make(map[string]string),
		timeout: make(map[string]StoreTimeout),
	}
}

func (s *StorageEngine) hasTimeout(key string) (StoreTimeout, bool) {
	timeout, ok := s.timeout[key]
	return timeout, ok
}

func (s *StorageEngine) setTimeout(key string, timeout int) {
	s.timeout[key] = StoreTimeout{
		timeout:   timeout,
		createdAt: time.Now(),
	}
}

// ExpiredTimeout checks if the timeout for a key has expired
func (s *StorageEngine) ExpiredTimeout(key string) bool {
	timeout, ok := s.hasTimeout(key)
	if !ok {
		return false
	}

	expired := int(time.Since(timeout.createdAt).Milliseconds()) > timeout.timeout

	// if timed out, delete the key
	if s.ExpiredTimeout(key) {
		delete(s.store, key)
		delete(s.timeout, key)
	}

	return expired
}

func (s *StorageEngine) Set(key, value string, opts SetOptions) {
	log.Println("Setting key: ", key, " with value: ", value, " and options", opts)

	s.store[key] = value
	if opts.HasTimeout {
		s.setTimeout(key, opts.Timeout)
	}
}

func (s *StorageEngine) Get(key string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	val, ok := s.store[key]
	if !ok {
		return "", errors.New("Key not found")
	}
	return val, nil
}

func (s *StorageEngine) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.store, key)
	delete(s.timeout, key)
}
