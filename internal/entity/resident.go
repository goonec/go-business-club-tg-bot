package entity

import (
	"github.com/goonec/business-tg-bot/internal/boterror"
	"regexp"
	"strconv"
	"strings"
)

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

func FindID(data string) int {
	parts := strings.Split(data, "_")
	if len(parts) > 2 {
		return 0
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0
	}

	return id
}

var fioRegex = regexp.MustCompile(`^[a-zA-Zа-яА-Я]`)

func IsFIOValid(f, l, p string) error {
	var err *boterror.BotError

	//if !fioRegex.MatchString(f) {
	//	errF = errors.New("Некорректно введено имя.")
	//	fmt.Errorf()
	//}
	//if !fioRegex.MatchString(l) {
	//	errF = errors.New("Некорректно введена фамилия.")
	//}
	//if !fioRegex.MatchString(p) {
	//	errF = errors.New("Некорректно введено отчество.")
	//}

	return err
}
