package tasks

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/bnkamalesh/errors"
	"github.com/jackc/pgx/v4/pgxpool"
)

type store interface {
	Create(ctx context.Context, t *Task) (int64, error)
	Delete(ctx context.Context, tid int64) error
	Edit(ctx context.Context, tid int64, t *Task) error
	Get(ctx context.Context, tid int64) (*Task, error)
	GetAll(ctx context.Context, uid int64) ([]Task, error)
}

type taskStore struct {
	qbuilder  squirrel.StatementBuilderType
	pqdriver  *pgxpool.Pool
	tableName string
}

func (ts *taskStore) Create(ctx context.Context, t *Task) (int64, error) {
	sqlStatement := `INSERT INTO tasks (uid, detail, assignedTo, completeBy, createdAt, updatedAt)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id`
	id := int64(0)
	err := ts.pqdriver.QueryRow(ctx, sqlStatement, t.UID, t.Detail, t.AssignedTo, t.CompleteBy, t.CreatedAt, t.UpdatedAt).Scan(&id)
	if err != nil {
		println(err.Error())
		return id, errors.InternalErr(err, errors.DefaultMessage)
	}
	return id, nil
}

func (ts *taskStore) Delete(ctx context.Context, tid int64) error {
	query, args, err := ts.qbuilder.Delete(ts.tableName).Where(
		squirrel.Eq{
			"id": tid,
		},
	).ToSql()
	if err != nil {
		return errors.InternalErr(err, errors.DefaultMessage)
	}

	_, err = ts.pqdriver.Query(ctx, query, args...)
	if err != nil {
		println(err.Error())
		return errors.InternalErr(err, errors.DefaultMessage)
	}
	return nil
}

func (ts *taskStore) Edit(ctx context.Context, tid int64, t *Task) error {
	query, args, err := ts.qbuilder.Update(ts.tableName).SetMap(map[string]interface{}{
		"uid":        t.UID,
		"detail":     t.Detail,
		"completeBy": t.CompleteBy,
		"updatedAt":  t.UpdatedAt,
	}).Where(squirrel.Eq{
		"id": tid,
	}).ToSql()
	if err != nil {
		return errors.InternalErr(err, errors.DefaultMessage)
	}

	_, err = ts.pqdriver.Exec(ctx, query, args...)
	if err != nil {
		println(err.Error())
		return errors.InternalErr(err, errors.DefaultMessage)
	}

	return nil
}

func (ts *taskStore) Get(ctx context.Context, tid int64) (*Task, error) {
	query, args, err := ts.qbuilder.Select(
		"uid",
		"detail",
		"completeBy",
		"assignedTo",
		"createdAt",
		"updatedAt",
	).From(
		ts.tableName,
	).Where(
		squirrel.Eq{
			"id": tid,
		},
	).ToSql()

	if err != nil {
		return nil, errors.InternalErr(err, errors.DefaultMessage)
	}

	task := new(Task)
	uid := new(sql.NullInt64)
	detail := new(sql.NullString)
	assignedTo := new(sql.NullString)
	completeBy := new(sql.NullTime)

	row := ts.pqdriver.QueryRow(ctx, query, args...)
	err = row.Scan(
		uid,
		detail,
		completeBy,
		assignedTo,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return nil, errors.InternalErr(err, err.Error())
	}

	task.TID = tid
	task.UID = uid.Int64
	task.Detail = detail.String
	task.CompleteBy = completeBy.Time
	task.AssignedTo = assignedTo.String

	return task, nil
}

func (ts *taskStore) GetAll(ctx context.Context, uid int64) ([]Task, error) {
	query, args, err := ts.qbuilder.Select(
		"id",
		"detail",
		"completeBy",
		"createdAt",
		"updatedAt",
	).From(
		ts.tableName,
	).Where(
		squirrel.Eq{
			"uid": uid,
		},
	).ToSql()

	if err != nil {
		return nil, errors.InternalErr(err, errors.DefaultMessage)
	}

	tasks := []Task{}
	id := new(sql.NullInt64)
	detail := new(sql.NullString)
	completeBy := new(sql.NullTime)

	rows, err := ts.pqdriver.Query(ctx, query, args...)
	defer rows.Close()
	for rows.Next() {
		task := new(Task)
		err = rows.Scan(
			id,
			detail,
			completeBy,
			&task.CreatedAt,
			&task.UpdatedAt,
		)

		if err != nil {
			return nil, errors.InternalErr(err, err.Error())
		}

		task.TID = id.Int64
		task.UID = uid
		task.Detail = detail.String
		task.CompleteBy = completeBy.Time
		tasks = append(tasks, *task)
	}

	return tasks, nil
}

func newStore(pqdriver *pgxpool.Pool) (*taskStore, error) {
	return &taskStore{
		pqdriver:  pqdriver,
		qbuilder:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		tableName: "tasks",
	}, nil
}
