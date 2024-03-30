// Description: This file contains the main storage engine of the application
// Author: Zablon Dawit
// Date Created: Mar-30-2024
package repository

import (
	"errors"
	"sync"
)

type StorageEngine struct {
	store map[string]string
	mutex sync.Mutex
}

func NewStorageEngine() *StorageEngine {
	return &StorageEngine{
		store: make(map[string]string),
	}
}

func (s *StorageEngine) Set(key, value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.store[key] = value
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
