package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/ivanjabrony/personApi/internal/model"
	"github.com/jmoiron/sqlx"
)

type PgPersonRepository struct {
	db *sqlx.DB
}

func NewPgPersonRepository(db *sqlx.DB) *PgPersonRepository {
	return &PgPersonRepository{db}
}

func (r *PgPersonRepository) Create(ctx context.Context, person *model.Person) (int, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return -1, err
	}

	defer func() {
		var e error
		if err == nil {
			e = tx.Commit()
		} else {
			e = tx.Rollback()
		}

		if err == nil && e != nil {
			err = fmt.Errorf("finishing transaction: %w", e)
		}
	}()

	query, args, err := squirrel.
		Insert("persons").
		Columns("name", "surname", "patronymic", "age", "gender", "nationality").
		Values(
			person.Name,
			person.Surname,
			person.Patronymic,
			person.Age,
			person.Gender,
			person.Nationality).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return -1, fmt.Errorf("failed to build query: %w", err)
	}

	err = tx.QueryRowxContext(ctx, query, args...).Scan(&person.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, fmt.Errorf("person not inserted: %w", err)
		}
		return -1, fmt.Errorf("failed to execute query: %w", err)
	}

	return person.Id, nil
}

func (r *PgPersonRepository) GetById(ctx context.Context, id int) (*model.Person, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		var e error
		if err == nil {
			e = tx.Commit()
		} else {
			e = tx.Rollback()
		}

		if err == nil && e != nil {
			err = fmt.Errorf("finishing transaction: %w", e)
		}
	}()

	query, args, err := squirrel.
		Select("id, name", "surname", "patronymic", "age", "gender", "nationality").
		From("persons").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var person model.Person

	err = tx.QueryRowxContext(ctx, query, args...).Scan(
		&person.Id,
		&person.Name,
		&person.Surname,
		&person.Patronymic,
		&person.Age,
		&person.Gender,
		&person.Nationality,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("person not inserted: %w", err)
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &person, nil
}

func (r *PgPersonRepository) GetFiltered(ctx context.Context, filter *model.PersonFilter) ([]model.Person, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		var e error
		if err == nil {
			e = tx.Commit()
		} else {
			e = tx.Rollback()
		}

		if err == nil && e != nil {
			err = fmt.Errorf("finishing transaction: %w", e)
		}
	}()

	queryString := squirrel.
		Select("id, name", "surname", "patronymic", "age", "gender", "nationality").
		From("persons")

	if filter.Name != nil {
		queryString = queryString.Where(squirrel.Eq{"name": *filter.Name})
	}
	if filter.Surname != nil {
		queryString = queryString.Where(squirrel.Eq{"surname": *filter.Surname})
	}
	if filter.Patronymic != nil {
		queryString = queryString.Where(squirrel.Eq{"patronymic": *filter.Patronymic})
	}
	if len(filter.Nationalities) != 0 {
		queryString = queryString.Where(squirrel.Eq{"nationality": filter.Nationalities})
	}
	if len(filter.Genders) != 0 {
		queryString = queryString.Where(squirrel.Eq{"gender": filter.Genders})
	}

	if filter.NameLike != nil {
		queryString = queryString.Where(squirrel.Like{"name": "%" + *filter.NameLike + "%"})
	}
	if filter.SurnameLike != nil {
		queryString = queryString.Where(squirrel.Like{"surname": "%" + *filter.SurnameLike + "%"})
	}
	if filter.PatronymicLike != nil {
		queryString = queryString.Where(squirrel.Like{"patronymic": "%" + *filter.PatronymicLike + "%"})
	}

	if filter.AgeMax != nil {
		queryString = queryString.Where(squirrel.LtOrEq{"age": *filter.AgeMax})
	}
	if filter.AgeMin != nil {
		queryString = queryString.Where(squirrel.GtOrEq{"age": *filter.AgeMin})
	}

	query, args, err := queryString.
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var persons []model.Person

	err = tx.SelectContext(ctx, &persons, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return persons, nil
}

func (r *PgPersonRepository) GetAll(ctx context.Context) ([]model.Person, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		var e error
		if err == nil {
			e = tx.Commit()
		} else {
			e = tx.Rollback()
		}

		if err == nil && e != nil {
			err = fmt.Errorf("finishing transaction: %w", e)
		}
	}()

	query, args, err := squirrel.
		Select("id, name", "surname", "patronymic", "age", "gender", "nationality").
		From("persons").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var persons []model.Person

	err = tx.SelectContext(ctx, &persons, query, args...)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return persons, nil
}

func (r *PgPersonRepository) Update(ctx context.Context, person *model.Person) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		var e error
		if err == nil {
			e = tx.Commit()
		} else {
			e = tx.Rollback()
		}

		if err == nil && e != nil {
			err = fmt.Errorf("finishing transaction: %w", e)
		}
	}()

	query, args, err := squirrel.
		Update("persons").
		Set("name", person.Name).
		Set("surname", person.Surname).
		Set("patronymic", person.Patronymic).
		Where(squirrel.Eq{"id": person.Id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("person with id %d not found", person.Id)
	}
	return nil
}

func (r *PgPersonRepository) DeleteById(ctx context.Context, id int) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		var e error
		if err == nil {
			e = tx.Commit()
		} else {
			e = tx.Rollback()
		}

		if err == nil && e != nil {
			err = fmt.Errorf("finishing transaction: %w", e)
		}
	}()

	query, args, err := squirrel.
		Delete("persons").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("person with id %d not found", id)
	}

	return nil
}
