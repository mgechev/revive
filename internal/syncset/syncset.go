// Package syncset provides a simple, mutex-protected set for strings.
package syncset

import "sync"

// Set is a concurrency-safe set of strings.
type Set struct {
	mu       sync.Mutex
	elements map[string]struct{}
}

// New returns an initialized, empty Set.
func New() *Set {
	return &Set{elements: map[string]struct{}{}}
}

// Has reports whether str is present in the set.
func (s *Set) Has(str string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, result := s.elements[str]
	return result
}

// Set adds str to the set.
func (s *Set) Set(str string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.elements[str] = struct{}{}
}
