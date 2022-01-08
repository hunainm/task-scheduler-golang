package api

import (
	"context"
	"strconv"
	"task-scheduler/internal/emailService"
	"task-scheduler/internal/tasks"
)

func (a *API) CreateTask(ctx context.Context, t *tasks.Task) (*tasks.Task, error) {
	t, err := a.tasks.Create(ctx, t)
	if err != nil {
		a.logger.Error(err)
		return nil, err
	}

	return t, nil
}

func (a *API) AssignTask(ctx context.Context, t *tasks.Task) (*tasks.Task, error) {
	u, err := a.users.GetUserByEmail(ctx, t.AssignedTo)
	if err == nil {
		t.UID = u.UID
	}

	t, err2 := a.tasks.Create(ctx, t)
	if err2 != nil {
		a.logger.Error(err)
		return nil, err
	}

	if err != nil {
		link := "http://localhost:8080/api/auth/register?tid=" + strconv.FormatInt(t.TID, 10)
		email := new(emailService.Email)
		email.Subject = "New Task Assigned"
		email.To = t.AssignedTo
		email.HtmlContent = `<p>Hello</p>
		<p>You've been assigned a new task. Please follow the link to register a new account and see your tasks.</p>
		<p><a href=` + link + `>Register (` + link + `)</a></p>`
		a.emailService.SendEmail(ctx, *email)
	}

	return t, nil
}

func (a *API) DeleteTask(ctx context.Context, tid int64) error {
	err := a.tasks.Delete(ctx, tid)
	if err != nil {
		a.logger.Error(err)
		return err
	}

	return nil
}

func (a *API) EditTask(ctx context.Context, tid int64, t *tasks.Task) (*tasks.Task, error) {
	t, err := a.tasks.Edit(ctx, tid, t)
	if err != nil {
		a.logger.Error(err)
		return nil, err
	}

	return t, nil
}

func (a *API) GetTask(ctx context.Context, tid int64) (*tasks.Task, error) {
	task, err := a.tasks.Get(ctx, tid)
	if err != nil {
		a.logger.Error(err)
		return nil, err
	}

	return task, nil
}

func (a *API) GetAllTasks(ctx context.Context, uid int64) ([]tasks.Task, error) {
	tasks, err := a.tasks.GetAll(ctx, uid)
	if err != nil {
		a.logger.Error(err)
		return nil, err
	}

	return tasks, nil
}
