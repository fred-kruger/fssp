package fssp

// Physical - данные для поиска физического лица.
type Physical struct {
	Region     int    `json:"region"`
	Firstname  string `json:"firstname"`
	Secondname string `json:"secondname"`
	Lastname   string `json:"lastname"`
	Birthdate  string `json:"birthdate"`
}
