package entity

type Service struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ServiceDescribe struct {
	ID        int    `json:"id"`
	ServiceID int    `json:"id_service"`
	Describe  string `json:"describe"`

	Service Service `json:"service"`
}
