package swagger

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserUpdate struct {
	ID         int         `json:"id"`
	Properties interface{} `json:"properties"`
}
