package storage

import (
    "encoding/json"
    "errors"
    "io/ioutil"
    "os"
    "sync"
)

type DataStore struct {
    filePath string
    data     map[string]interface{}
    mu       sync.Mutex
}

// Open initializes the DataStore with the given file path.
func Open(filePath string) (*DataStore, error) {
    ds := &DataStore{
        filePath: filePath,
        data:     make(map[string]interface{}),
    }

    // Check if the file exists
    if _, err := os.Stat(filePath); err == nil {
        // File exists, load data
        fileData, err := ioutil.ReadFile(filePath)
        if err != nil {
            return nil, err
        }
        if len(fileData) > 0 {
            if err := json.Unmarshal(fileData, &ds.data); err != nil {
                return nil, err
            }
        }
    } else if !os.IsNotExist(err) {
        return nil, err
    }

    return ds, nil
}

// Store saves a key-value pair to the data store and writes it to the file.
func (ds *DataStore) Store(key string, value interface{}) error {
    ds.mu.Lock()
    defer ds.mu.Unlock()

    ds.data[key] = value
    return ds.saveToFile()
}

// Get retrieves the value associated with the given key.
func (ds *DataStore) Get(key string) (interface{}, error) {
    ds.mu.Lock()
    defer ds.mu.Unlock()

    value, exists := ds.data[key]
    if !exists {
        return nil, errors.New("key not found")
    }
    return value, nil
}

// saveToFile writes the current data to the JSON file.
func (ds *DataStore) saveToFile() error {
    fileData, err := json.MarshalIndent(ds.data, "", "  ")
    if err != nil {
        return err
    }
    return ioutil.WriteFile(ds.filePath, fileData, 0644)
}