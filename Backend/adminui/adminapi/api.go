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

type BusinessInfo struct {
	CompanyName      string   `json:"companyName"`
	CompanyAddress   string   `json:"companyAddress"`
	FinanceEmail     string   `json:"financeEmail"`
	VATID            string   `json:"vatId"`
	CompanyRegNumber string   `json:"companyRegNumber"`
	Phone            string   `json:"phone"`
	Website          string   `json:"website"`
	LogoURL          string   `json:"logoUrl"`
	SalesIDPrefix    string   `json:"salesIdPrefix"`
	Currency         string   `json:"currency"`
	DefaultTaxRate   float64  `json:"defaultTaxRate"`
	BankDetails      string   `json:"bankDetails"`
	LegalFooter      string   `json:"legalFooter"`
	OpeningHours     string   `json:"openingHours"`
	SocialLinks      []string `json:"socialLinks"`
	CustomReceiptMsg string   `json:"customReceiptMsg"`
	InvoiceFormat    string   `json:"invoiceFormat"`
	Country          string   `json:"country"`
	LastSalesNumber  int      `json:"lastSalesNumber"`
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

func GetBusinessInfo() (*BusinessInfo, error) {
	resp, err := http.Get(apiBaseURL() + "/business")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get business info: " + resp.Status)
	}
	var info BusinessInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

func SetBusinessInfo(info *BusinessInfo) error {
	b, err := json.Marshal(info)
	if err != nil {
		return err
	}
	resp, err := http.Post(apiBaseURL()+"/business", "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errors.New("failed to set business info: " + resp.Status)
	}
	return nil
}
