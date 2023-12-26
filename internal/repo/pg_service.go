package repo

import "github.com/goonec/business-tg-bot/pkg/postgres"

type serviceRepo struct {
	*postgres.Postgres
}

//func NewServiceRepo(pg *postgres.Postgres) Service {
//	return *serviceRepo{
//		pg,
//	}
//}
