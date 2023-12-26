package entity

type Service struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ServiceDescribe struct {
	ID          int    `json:"id"`
	ServiceID   int    `json:"id_service"`
	Name        string `json:"name"`
	Describe    string `json:"describe"`
	PhotoFileID string `json:"photo_file_id"`

	Service Service `json:"service"`
}
