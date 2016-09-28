package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// Store represents a bivouac issue tracker store.
// It contains the tracker's issues
type Store struct {
	Issues []*Issue `json:"issues"`
	Path   string   `json:"-"`

	NextID uint `json:"nextid"`
}

// CreateNewStore instanciates a new store, and write it
// to disk
func CreateNewStore(storePath string) (*Store, error) {
	store := &Store{}

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

// LoadStore reads a store content from file and returns
// an instancie of store
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

// GetOrCreateStore will try to load a store at provided path.
// If it does not exist, it will create it and load it.
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

// AddIssue adds an issue to the store
func (s *Store) AddIssue(issue Issue) {
	s.Issues = append(s.Issues, &issue)
}

// GetIssue retrieves issue with the provided id in the store
func (s *Store) GetIssue(id uint) (*Issue, error) {
	for _, issue := range s.Issues {
		if issue.ID == id {
			return issue, nil
		}
	}

	return nil, fmt.Errorf("no issue with id %d", id)
}

// ListIssues returns a list of the issues in store
func (s *Store) ListIssues() []*Issue {
	var issues []*Issue

	for i := len(s.Issues) - 1; i >= 0; i-- {
		issues = append(issues, s.Issues[i])
	}

	return issues
}

// FilterIssues lets you retrieve store issues with the provided status
func (s *Store) FilterIssues(status IssueStatus) []*Issue {
	var issues []*Issue

	for _, issue := range s.Issues {
		if issue.Status == status {
			issues = append(issues, issue)
		}
	}

	return issues
}

// HasIssues indicates if the store contains issues at all
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

func (s *Store) getNextID() uint {
	nextID := s.NextID
	s.NextID++
	return nextID
}

// NewStore creates a new Store instance
func NewStore(path string, issues []*Issue) *Store {
	return &Store{
		Issues: issues,
		Path:   path,
		NextID: 0,
	}
}
