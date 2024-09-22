package models

type UserInfoExpectedFormat struct {
	Following bool                        `json:"following"`
	Profile   PublicProfileExpectedFormat `json:"profile"`
}

type PublicProfileExpectedFormat struct {
	ID        string `json:"id"`
	FisrtName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Location  string `json:"location"`
	Following int    `json:"following"`
	Followers int    `json:"followers"`
}

type Pagination struct {
	Next_Offset int `json:"next_offset,omitempty"`
	Limit       int `json:"limit"`
}

type PrivateProfileExpectedFormat struct {
	ID        string `json:"id"`
	FisrtName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Location  string `json:"location"`
	Email    string `json:"email"`
	Interests []string `json:"interests"`
	Following int    `json:"following"`
	Followers int    `json:"followers"`
}