package failureinjection

import "fmt"

type Journal interface {
	Write([]byte) error
	Sync() error
}

type Store struct {
	journal Journal
	values  map[string]string
}

func NewStore(journal Journal) *Store {
	return &Store{journal: journal, values: make(map[string]string)}
}

func (s *Store) Put(key, value string) error {
	if err := s.journal.Write([]byte(key + "=" + value)); err != nil {
		return fmt.Errorf("write journal: %w", err)
	}

	s.values[key] = value
	if err := s.journal.Sync(); err != nil {
		return fmt.Errorf("sync journal: %w", err)
	}
	return nil
}

func (s *Store) Get(key string) (string, bool) {
	value, ok := s.values[key]
	return value, ok
}
