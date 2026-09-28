package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
)

type Key struct {
	UIDL       string `json:"uidl"`
	Attachment string `json:"attachment"`
}

type State struct {
	Status   AttachStatus `json:"status"`
	Attempts int          `json:"attempts"`
}

type entry struct {
	Key   Key   `json:"key"`
	State State `json:"state"`
}

type Store struct {
	m           map[Key]State
	path        string
	maxAttempts int
}

func loadStore(path string, maxAttempts int) (*Store, error) {
	s := &Store{m: map[Key]State{}, path: path, maxAttempts: maxAttempts}
	raw, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []entry
	if err := json.Unmarshal(raw, &entries); err != nil {
		return nil, err
	}
	for _, e := range entries {
		s.m[e.Key] = e.State
	}
	return s, nil
}

func (s *Store) Uploadable(k Key) bool {
	st, ok := s.m[k]
	return !ok || (st.Status == FailedRetryable && st.Attempts < s.maxAttempts)
}

func (s *Store) Resolved(uidl string, names []string) bool {
	for _, n := range names {

		st, ok := s.m[Key{uidl, n}]
		if !ok {
			return false
		}

		if st.Status != Uploaded &&
			st.Status != SkippedTooSmall &&
			st.Status != FailedPermanent {
			return false
		}
	}
	return true
}

func (s *Store) Record(k Key, st AttachStatus) error {
	n := s.m[k].Attempts + 1
	if st == FailedRetryable && n >= s.maxAttempts {
		st = FailedPermanent
	}
	s.m[k] = State{Status: st, Attempts: n}
	return s.save()
}

func (s *Store) save() error {
	entries := make([]entry, 0, len(s.m))
	for k, v := range s.m {
		entries = append(entries, entry{k, v})
	}
	raw, err := json.Marshal(entries)
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
