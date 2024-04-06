package swagger

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserUpdate struct {
	ID         int         `json:"id"`
	Properties interface{} `json:"properties"`
}

type UserMe struct {
	Login                string      `json:"login"`
	CurrentTripStartDate interface{} `json:"current_trip_start_date"`
	Properties           interface{} `json:"properties"`
}
