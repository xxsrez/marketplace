package main

import "sync"

// volatileCredentialStore deliberately keeps OAuth client metadata and refresh
// tokens inside the companion process. Starting the MCP server and listing its
// tools therefore never touches macOS Keychain; the first actual upload starts
// OAuth, and a new process requires authorization again.
type volatileCredentialStore struct {
	mutex      sync.Mutex
	credential *oauthCredential
}

func (store *volatileCredentialStore) Load() (oauthCredential, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if store.credential == nil {
		return oauthCredential{}, errCredentialNotFound
	}
	return *store.credential, nil
}

func (store *volatileCredentialStore) Save(credential oauthCredential) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	copy := credential
	store.credential = &copy
	return nil
}

func (store *volatileCredentialStore) Delete() error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.credential = nil
	return nil
}
