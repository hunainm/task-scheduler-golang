package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"task-scheduler/internal/tasks"
	"task-scheduler/internal/users"

	"github.com/bnkamalesh/errors"
	"github.com/bnkamalesh/webgo/v6"
)

func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	u := new(users.User)
	err := json.NewDecoder(r.Body).Decode(u)

	if err != nil {
		errResponder(w, errors.InputBodyErr(err, "invalid JSON provided"))
		return
	}

	createdUser, err := h.api.Register(r.Context(), u)
	if err != nil {
		errResponder(w, err)
		return
	}
	tid, err := strconv.ParseInt(r.FormValue("tid"), 10, 64)
	task, err := h.api.GetTask(r.Context(), tid)
	if err != nil {
		errResponder(w, err)
		return
	}
	if task.AssignedTo == u.Email {
		task.UID = createdUser.UID
		editedTask, err := h.api.EditTask(r.Context(), tid, task)
		if err != nil {

			println(err.Error())
		}
		println(editedTask)
	}

	b, err := json.Marshal(createdUser)
	if err != nil {
		errResponder(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	u := new(users.User)
	err := json.NewDecoder(r.Body).Decode(u)

	if err != nil {
		errResponder(w, errors.InputBodyErr(err, "invalid JSON provided"))
		return
	}

	jwt, err := h.api.Login(r.Context(), u.Email, u.Password)
	if err != nil {
		errResponder(w, err)
		return
	}
	b, err := json.Marshal(jwt)

	if err != nil {
		errResponder(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (h *Handlers) AddTask(w http.ResponseWriter, r *http.Request) {
	t := new(tasks.Task)
	err := json.NewDecoder(r.Body).Decode(t)
	if err != nil {
		errResponder(w, errors.InputBodyErr(err, "Invalid JSON provided"))
		return
	}
	props, _ := r.Context().Value("props").(*users.Claims)
	t.UID, err = strconv.ParseInt(props.Id, 10, 64)
	createdTask, err := h.api.CreateTask(r.Context(), t)
	if err != nil {
		errResponder(w, err)
		return
	}

	b, err := json.Marshal(createdTask)
	if err != nil {
		errResponder(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	wctx := webgo.Context(r)
	tid, err := strconv.ParseInt(wctx.Params()["tid"], 10, 64)
	err = h.api.DeleteTask(r.Context(), tid)
	if err != nil {
		errResponder(w, err)
		return
	}

	webgo.R200(w, nil)
}

func (h *Handlers) EditTask(w http.ResponseWriter, r *http.Request) {
	t := new(tasks.Task)
	err := json.NewDecoder(r.Body).Decode(t)
	wctx := webgo.Context(r)
	tid, err := strconv.ParseInt(wctx.Params()["tid"], 10, 64)

	if err != nil {
		errResponder(w, errors.InputBodyErr(err, "Invalid JSON provided"))
		return
	}

	modifiedTask, err := h.api.EditTask(r.Context(), tid, t)
	if err != nil {
		errResponder(w, err)
		return
	}

	b, err := json.Marshal(modifiedTask)
	if err != nil {
		errResponder(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (h *Handlers) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(*users.Claims)
	uid, err := strconv.ParseInt(props.Id, 10, 64)
	tasks, err := h.api.GetAllTasks(r.Context(), uid)
	if err != nil {
		errResponder(w, err)
		return
	}

	b, err := json.Marshal(tasks)
	if err != nil {
		errResponder(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (h *Handlers) AssignTask(w http.ResponseWriter, r *http.Request) {
	t := new(tasks.Task)
	err := json.NewDecoder(r.Body).Decode(t)
	if err != nil {
		errResponder(w, errors.InputBodyErr(err, "Invalid JSON provided"))
		return
	}

	createdTask, err := h.api.AssignTask(r.Context(), t)
	if err != nil {
		errResponder(w, err)
		return
	}

	b, err := json.Marshal(createdTask)
	if err != nil {
		errResponder(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
