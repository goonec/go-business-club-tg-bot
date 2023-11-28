package entity

// v1
//type Resident struct {
//	ID           int    `json:"id"`
//	Age          int8   `json:"age"`
//	UsernameTG   string `json:"tg_username"`
//	PhoneNumber  string `json:"phone_number"`
//	Firstname    string `json:"firstname"`
//	Lastname     string `json:"lastname"`
//	Patronymic   string `json:"patronymic"`
//	Region       string `json:"region"`
//	WorkActivity string `json:"work_activity"`
//	CompanyName  string `json:"company_name"`
//	Advantage    string `json:"advantage"`
//	Hobie        string `json:"hobie"`
//	PhotoFileID  string `json:"photo_file_id"`
//}

// Resident v2
type Resident struct {
	ID           int    `json:"id"`
	UsernameTG   string `json:"tg_username"`
	ResidentData string `json:"resident_data"`
	PhotoFileID  string `json:"photo_file_id"`

	FIO FIO `json:"fio"`
}

type FIO struct {
	ID         int    `json:"id"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Patronymic string `json:"patronymic"`
}
