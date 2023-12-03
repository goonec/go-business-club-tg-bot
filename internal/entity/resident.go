package entity

import (
	"regexp"
	"strconv"
	"strings"
)

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

var fioRegex = regexp.MustCompile(`^[a-zA-Zёа-яА-Я-]+$`)

func IsFIOValid(f, l, p string) string {
	var err []string

	if !fioRegex.MatchString(f) {
		err = append(err, "Некорректно введено имя. ")
	}
	if !fioRegex.MatchString(l) {
		err = append(err, "Некорректно введена фамилия. ")
	}
	if !fioRegex.MatchString(p) {
		err = append(err, "Некорректно введено отчество.")
	}

	validInfo := strings.Join(err, "")

	return validInfo
}
