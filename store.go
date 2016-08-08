package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type Store struct {
	Issues []*Issue `json:"issues"`
	Path   string   `json:"-"`

	NextId uint `json:"nextid"`
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
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path.Join(cwd, bivouacFile)); os.IsNotExist(err) {
		_, err := CreateNewStore(bivouacFile)
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
	var issues []*Issue

	for i := len(s.Issues) - 1; i >= 0; i-- {
		issues = append(issues, s.Issues[i])
	}

	return issues
}

func (s *Store) FilterIssues(status IssueStatus) []*Issue {
	var issues []*Issue

	for _, issue := range s.Issues {
		if issue.Status == status {
			issues = append(issues, issue)
		}
	}

	return issues
}

func (s *Store) HasIssues() bool {
	return len(s.Issues) > 0
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
	nextId := s.NextId
	s.NextId += 1
	return nextId
}

func NewStore(path string, issues []*Issue) *Store {
	return &Store{
		Issues: issues,
		Path:   path,
		NextId: 0,
	}
}
