//go:build !darwin

package main

type systemBrowser struct{}

func (systemBrowser) Open(string) error {
	return newLocalError("browser_auth_unavailable", "browser authorization is available only on macOS; Linux requires the TASK_MANAGER_LOCAL_TOKEN runtime secret")
}
