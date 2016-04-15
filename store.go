package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Store struct {
	Issues []*Issue `json:"issues"`
	Path   string   `json:"-"`

	nextId uint `json:"nextid"`
}

func CreateNewStore(storePath string) (*Store, error) {
	var store *Store = &Store{}

	jsonData, err := json.Marshal(store)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(storePath, jsonData, os.FileMode(0600))
	if err != nil {
		return nil, err
	}

	store.Path, err = filepath.Abs(storePath)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func LoadStore(storePath string) (*Store, error) {
	var store *Store
	var storeContent []byte
	var err error

	storeContent, err = ioutil.ReadFile(storePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(storeContent, &store)
	if err != nil {
		return nil, err
	}

	store.Path, err = filepath.Abs(storePath)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func GetOrCreateStore(storePath string) (*Store, error) {
	if _, err := os.Stat(".groundcontrol"); os.IsNotExist(err) {
		_, err := CreateNewStore(".groundcontrol")
		if err != nil {
			return nil, fmt.Errorf("Creating store file failed with error: %s", err.Error())
		}
	}

	return LoadStore(storePath)
}

func (s *Store) AddIssue(issue Issue) {
	s.Issues = append(s.Issues, &issue)
}

func (s *Store) GetIssue(id uint) (*Issue, error) {
	for _, issue := range s.Issues {
		if issue.Id == id {
			return issue, nil
		}
	}

	return nil, fmt.Errorf("no issue with id %d", id)
}

func (s *Store) ListIssues() []*Issue {
	return s.Issues
}

func (s *Store) Write() error {
	jsonData, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(s.Path, jsonData, os.FileMode(0600))
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) getNextId() uint {
	nextId := s.nextId
	s.nextId += 1
	return nextId
}

func NewStore(path string, issues []*Issue) *Store {
	return &Store{
		Issues: issues,
		Path:   path,
		nextId: 0,
	}
}
