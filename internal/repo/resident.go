package repo

import (
	"context"
	"database/sql"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/pkg/postgres"
)

type residentRepository struct {
	*postgres.Postgres
}

func NewResidentRepository(postgres *postgres.Postgres) Resident {
	return &residentRepository{
		postgres,
	}
}

func (r *residentRepository) GetAll(ctx context.Context) ([]entity.Resident, error) {
	query := `select id, tg_username ,resident_data,photo_file_id from resident`

	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	residents := make([]entity.Resident, 0, 512)

	for rows.Next() {
		var r entity.Resident

		err := rows.Scan(&r.ID,
			&r.UsernameTG,
			&r.ResidentData,
			&r.PhotoFileID,
		)
		if err != nil {
			return nil, err
		}

		residents = append(residents, r)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return residents, nil
}

func (r *residentRepository) GetAllFIO(ctx context.Context) ([]entity.FIO, error) {
	query := `select id, firstname,substring(lastname,1,1),substring(patronymic,1,1) from resident`

	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	fio := make([]entity.FIO, 0, 512)

	for rows.Next() {
		var f entity.FIO

		err := rows.Scan(&f.ID,
			&f.Firstname,
			&f.Lastname,
			&f.Patronymic,
		)
		if err != nil {
			return nil, err
		}

		fio = append(fio, f)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return fio, nil
}

func (r *residentRepository) Create(ctx context.Context, resident *entity.Resident) error {
	query := `insert into resident (tg_username,firstname,lastname,patronymic,resident_data,photo_file_id) 
				values ($1,$2,$3,$4,$5,$6)`

	_, err := r.Pool.Exec(ctx, query, resident.UsernameTG,
		resident.FIO.Firstname,
		resident.FIO.Lastname,
		resident.FIO.Patronymic,
		resident.ResidentData,
		resident.PhotoFileID)

	return err
}

func (r *residentRepository) GetByID(ctx context.Context, id int) (*entity.Resident, error) {
	query := `select tg_username,firstname,lastname,patronymic,resident_data,photo_file_id from resident where id =$1`
	var resident entity.Resident

	err := r.Pool.QueryRow(ctx, query, id).Scan(&resident.UsernameTG,
		&resident.FIO.Firstname,
		&resident.FIO.Lastname,
		&resident.FIO.Patronymic,
		&resident.ResidentData,
		&resident.PhotoFileID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, boterror.ErrNotFound
		}
		return nil, err
	}

	return &resident, nil
}
