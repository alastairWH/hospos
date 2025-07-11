package adminapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type User struct {
	Name string `json:"name"`
	Role string `json:"role"`
	Pin  string `json:"pin"`
}

type Role struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

func apiBaseURL() string {
	url := os.Getenv("HOSPOS_API_URL")
	if url == "" {
		url = "http://localhost:8080/api"
	}
	return url
}

func AddUser(u User) error {
	b, err := json.Marshal(u)
	if err != nil {
		return err
	}
	resp, err := http.Post(apiBaseURL()+"/users", "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return errors.New("failed to add user: " + resp.Status)
	}
	return nil
}

func GetUsers() ([]User, error) {
	resp, err := http.Get(apiBaseURL() + "/users")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get users: " + resp.Status)
	}
	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}
	return users, nil
}

func GetRoles() ([]Role, error) {
	resp, err := http.Get(apiBaseURL() + "/roles")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get roles: " + resp.Status)
	}
	var roles []Role
	if err := json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		return nil, err
	}
	return roles, nil
}

func AddRole(role string) error {
	b, err := json.Marshal(Role{Role: role})
	if err != nil {
		return err
	}
	resp, err := http.Post(apiBaseURL()+"/roles", "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return errors.New("failed to add role: " + resp.Status)
	}
	return nil
}

func CheckAPIStatus() bool {
	resp, err := http.Get(apiBaseURL() + "/users")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func InitDB() error {
	resp, err := http.Post(apiBaseURL()+"/dbinit", "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to init db: " + resp.Status)
	}
	return nil
}

func SeedTestData() error {
	resp, err := http.Post(apiBaseURL()+"/devtools/seed", "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errors.New("failed to seed test data: " + resp.Status)
	}
	return nil
}

func ClearTestData() error {
	resp, err := http.Post(apiBaseURL()+"/devtools/clear", "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errors.New("failed to clear test data: " + resp.Status)
	}
	return nil
}
