package entity

type Resident struct {
	ID           int64  `json:"id"`
	Age          int8   `json:"age"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Patronymic   string `json:"patronymic"`
	Region       string `json:"region"`
	WorkActivity string `json:"work_activity"`
	CompanyName  string `json:"company_name"`
	Advantage    string `json:"advantage"`
	PhotoFileID  string `json:"photo_file_id"`

	ResidentHobie []ResidentHobie `json:"resident_hobie"`
	ResidentRole  ResidentRole    `json:"resident_role"`
}

type ResidentHobie struct {
	ID    int    `json:"id"`
	Hobie string `json:"hobie"`
}

type ResidentRole struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
}
