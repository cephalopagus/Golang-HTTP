package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"study/todo"
	"time"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	todoList *todo.List
}

func NewHTTPHandler(todoList *todo.List) *HTTPHandlers {
	return &HTTPHandlers{
		todoList: todoList,
	}
}

func (h *HTTPHandlers) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO TaskDTO
	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		errorDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		http.Error(w, errorDTO.ErrorToString(), http.StatusBadRequest)
		return
	}

	if err := taskDTO.Validate(); err != nil {
		errorDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		http.Error(w, errorDTO.ErrorToString(), http.StatusBadRequest)
		return
	}

	todoTask := todo.NewTask(taskDTO.Title, taskDTO.Description)
	if err := h.todoList.AddTask(todoTask); err != nil {

		errorDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, todo.ErrTaskAlreadyExists) {
			http.Error(w, errorDTO.ErrorToString(), http.StatusConflict)
		} else {
			http.Error(w, errorDTO.ErrorToString(), http.StatusInternalServerError)
		}
		return
	}

	b, err := json.MarshalIndent(todoTask, "", "    ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write response:", err)
		return
	}
}

func (h *HTTPHandlers) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]
	todoTask, err := h.todoList.GetTask(title)
	if err != nil {

		errorDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errorDTO.ErrorToString(), http.StatusNotFound)
		} else {
			http.Error(w, errorDTO.ErrorToString(), http.StatusInternalServerError)
		}
		return
	}
	b, err := json.MarshalIndent(todoTask, "", "   ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write response:", err)
		return
	}

}

func (h *HTTPHandlers) HandleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.todoList.ListTasks()
	b, err := json.MarshalIndent(tasks, "", "   ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write response:", err)
		return
	}

}

func (h *HTTPHandlers) HandleGetAllUncompletedTasks(w http.ResponseWriter, r *http.Request) {
	uncompletedTask := h.todoList.ListUncompletedTasks()
	b, err := json.MarshalIndent(uncompletedTask, "", "   ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write response:", err)
		return
	}

}

func (h *HTTPHandlers) HandleCompleteTask(w http.ResponseWriter, r *http.Request) {

	var completeDTO CompleteDTO
	if err := json.NewDecoder(r.Body).Decode(&completeDTO); err != nil {
		errorDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errorDTO.ErrorToString(), http.StatusNotFound)
		}
		return
	}
	title := mux.Vars(r)["title"]

	var (
		task todo.Task
		err  error
	)

	if completeDTO.Complete {
		task, err = h.todoList.CompleteTask(title)

	} else {
		task, err = h.todoList.UncompleteTask(title)
	}
	if err != nil {
		errorDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errorDTO.ErrorToString(), http.StatusNotFound)
		} else {
			http.Error(w, errorDTO.ErrorToString(), http.StatusInternalServerError)
		}
		return
	}

	b, err := json.MarshalIndent(task, "", "   ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write response:", err)
		return
	}

}
func (h *HTTPHandlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]
	if err := h.todoList.DeleteTask(title); err != nil {

		errorDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errorDTO.ErrorToString(), http.StatusNotFound)
		} else {
			http.Error(w, errorDTO.ErrorToString(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
