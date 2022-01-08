package tasks

import (
	"context"
	"time"

	"task-scheduler/internal/platform/logger"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Task struct {
	TID        int64      `json:"tid,omitempty"`
	UID        int64      `json:"uid,omitempty"`
	Detail     string     `json:"detail,omitempty"`
	CompleteBy time.Time  `json:"completeBy,omitempty"`
	AssignedTo string     `json:"assignedTo,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

func (u *Task) init() {
	now := time.Now()
	if u.CreatedAt == nil {
		u.CreatedAt = &now
	}

	if u.UpdatedAt == nil {
		u.UpdatedAt = &now
	}
}

type Tasks struct {
	logHandler logger.Logger
	store      store
}

func (ts *Tasks) Create(ctx context.Context, t *Task) (*Task, error) {
	t.init()

	id, err := ts.store.Create(ctx, t)
	if err != nil {
		return nil, err
	}
	t.TID = id
	return t, nil
}

func (ts *Tasks) Delete(ctx context.Context, tid int64) error {
	err := ts.store.Delete(ctx, tid)
	if err != nil {
		return err
	}

	return nil
}

func (ts *Tasks) Edit(ctx context.Context, tid int64, t *Task) (*Task, error) {
	now := time.Now()
	t.UpdatedAt = &now
	err := ts.store.Edit(ctx, tid, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (ts *Tasks) Get(ctx context.Context, tid int64) (*Task, error) {
	task, err := ts.store.Get(ctx, tid)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (ts *Tasks) GetAll(ctx context.Context, uid int64) ([]Task, error) {
	tasks, err := ts.store.GetAll(ctx, uid)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func NewService(l logger.Logger, pqdriver *pgxpool.Pool) (*Tasks, error) {
	tstore, err := newStore(pqdriver)
	if err != nil {
		return nil, err
	}

	return &Tasks{
		logHandler: l,
		store:      tstore,
	}, nil
}
