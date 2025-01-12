package postgres

import (
	"SLON_tg_bot/src/domains/entities"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"strconv"
	"strings"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(connectionString string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) AddUser(user *entities.TargetUser) error {
	query := `INSERT INTO app.targets (creator_id, target, schedule, tags) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, user.CreatorID, user.Target, user.Schedule, user.Tags)
	return err
}

func (r *PostgresRepository) GetUsersByCreator(creator int64) ([]*entities.TargetUser, error) {
	query := `SELECT creator_id, target, schedule, tags FROM app.targets WHERE creator_id = $1`
	rows, err := r.db.Query(query, creator)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.TargetUser
	for rows.Next() {
		user := &entities.TargetUser{}
		if err := rows.Scan(&user.CreatorID, &user.Target, &user.Schedule, &user.Tags); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *PostgresRepository) GetSingleByCreatorAndTarget(creator int64, target string) (*entities.TargetUser, error) {
	query := `SELECT creator_id, target, schedule, tags FROM app.targets WHERE creator_id = $1 AND target = $2`
	user := &entities.TargetUser{}
	if err := r.db.QueryRow(query, creator, target).Scan(&user.CreatorID, &user.Target, &user.Schedule, &user.Tags); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *PostgresRepository) RemoveUser(creator int64, target string) error {
	query := `DELETE FROM app.targets WHERE creator_id = $1 AND target = $2`
	_, err := r.db.Exec(query, creator, target)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) UpdateUserTags(creator int64, target string, tags []entities.TagType) error {
	query := `UPDATE app.targets SET tags = $1 WHERE creator_id = $2 AND target = $3`
	intTags := "{"
	for _, tag := range tags {
		intTags += strconv.Itoa(int(tag)) + ","
	}
	intTags = strings.TrimSuffix(intTags, ",")
	_, err := r.db.Exec(query, intTags+"}", creator, target)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) UpdateUserSched(creator int64, target string, schedType entities.ScheduleType) error {
	query := `UPDATE app.targets SET schedule = $1 WHERE creator_id = $2 AND target = $3`
	_, err := r.db.Exec(query, schedType, creator, target)
	if err != nil {
		return err
	}
	return nil
}
