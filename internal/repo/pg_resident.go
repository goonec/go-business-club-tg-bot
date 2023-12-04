package repo

import (
	"context"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/pkg/postgres"
	"github.com/jackc/pgx/v5"
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

	residents := make([]entity.Resident, 0, 256)

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

	fio := make([]entity.FIO, 0, 256)

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

func (r *residentRepository) Create(ctx context.Context, resident *entity.Resident) (int, error) {
	query := `insert into resident (tg_username,firstname,lastname,patronymic,resident_data,photo_file_id) 
				values ($1,$2,$3,$4,$5,$6) returning id`
	var id int

	err := r.Pool.QueryRow(ctx, query, resident.UsernameTG,
		resident.FIO.Firstname,
		resident.FIO.Lastname,
		resident.FIO.Patronymic,
		resident.ResidentData,
		resident.PhotoFileID).Scan(&id)

	return id, err
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
		if err == pgx.ErrNoRows {
			return nil, boterror.ErrNotFound
		}
		return nil, err
	}

	return &resident, nil
}

func (r *residentRepository) DeleteByID(ctx context.Context, id int) error {
	query := `delete from resident where id = $1`

	_, err := r.Pool.Exec(ctx, query, id)
	return err
}

func (r *residentRepository) GetAllByClusterID(ctx context.Context, id int) ([]entity.FIO, error) {
	query := `select r.id, r.firstname, substring(r.lastname,1,1), substring(r.patronymic,1,1) from resident r
				join business_cluster_resident bcr on bcr.id_resident = r.id
				join business_cluster bc on bc.id = bcr.id_business_cluster
				where bc.id = $1`

	rows, err := r.Pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	fio := make([]entity.FIO, 0, 128)

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
