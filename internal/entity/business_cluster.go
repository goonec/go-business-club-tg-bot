package entity

type BusinessCluster struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type BusinessClusterResident struct {
	IDResident int    `json:"resident_id"`
	Name       string `json:"name"`
}
