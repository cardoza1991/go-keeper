package main

import (
    "encoding/json"
    "errors"
    "io/ioutil"
    "os"
    "sync"
)

type PasswordEntry struct {
    Service  string `json:"service"`
    Password string `json:"password"`
}

type PasswordStore struct {
    Entries []PasswordEntry `json:"entries"`
    mutex   sync.Mutex
}

var store PasswordStore

const storageFile = "passwords.json"

// LoadStore loads the password store from the storage file
func LoadStore() error {
    store.mutex.Lock()
    defer store.mutex.Unlock()

    if _, err := os.Stat(storageFile); os.IsNotExist(err) {
        store.Entries = []PasswordEntry{}
        return nil
    }

    data, err := ioutil.ReadFile(storageFile)
    if err != nil {
        return err
    }

    // Check if the file is empty
    if len(data) == 0 {
        store.Entries = []PasswordEntry{}
        return SaveStore() // Initialize the file with an empty array
    }

    if err := json.Unmarshal(data, &store.Entries); err != nil {
        return err
    }

    return nil
}

// SaveStore saves the password store to the storage file
func SaveStore() error {
    // Removed mutex locking here
    data, err := json.MarshalIndent(store.Entries, "", "  ")
    if err != nil {
        return err
    }

    return ioutil.WriteFile(storageFile, data, 0600) // 0600 for read/write by owner only
}

// AddEntry adds a new password entry
func AddEntry(entry PasswordEntry) error {
    store.mutex.Lock()
    defer store.mutex.Unlock()

    // Check for duplicates
    for _, e := range store.Entries {
        if e.Service == entry.Service {
            return errors.New("service already exists")
        }
    }

    store.Entries = append(store.Entries, entry)
    return SaveStore()
}

// GetEntry retrieves a password for a given service
func GetEntry(service string) (PasswordEntry, error) {
    store.mutex.Lock()
    defer store.mutex.Unlock()

    for _, e := range store.Entries {
        if e.Service == service {
            return e, nil
        }
    }

    return PasswordEntry{}, errors.New("service not found")
}

// DeleteEntry deletes a password entry for a given service
func DeleteEntry(service string) error {
    store.mutex.Lock()
    defer store.mutex.Unlock()

    for i, e := range store.Entries {
        if e.Service == service {
            store.Entries = append(store.Entries[:i], store.Entries[i+1:]...)
            return SaveStore()
        }
    }

    return errors.New("service not found")
}

// ListEntries lists all services
func ListEntries() []PasswordEntry {
    store.mutex.Lock()
    defer store.mutex.Unlock()

    return store.Entries
}
