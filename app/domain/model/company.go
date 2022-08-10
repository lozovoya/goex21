package model

type Company struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Code     string `json:"code,omitempty"`
	Country  string `json:"country,omitempty"`
	Website  string `json:"website,omitempty"`
	Phone    string `json:"phone,omitempty"`
	IsActive bool   `json:"is_active,omitempty"`
}

type Conditions struct {
	Name struct {
		Value   string
		IsExist bool
	}
	Code struct {
		Value   string
		IsExist bool
	}
	Country struct {
		Value   string
		IsExist bool
	}
	Website struct {
		Value   string
		IsExist bool
	}
	Phone struct {
		Value   string
		IsExist bool
	}
	IsActive struct {
		Value   bool
		IsExist bool
	}
}
