package server

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"restServer/taskstore"
	"strconv"
	"time"
)

type TaskServer struct {
	store *taskstore.TaskStore
}

func NewTaskServer() *TaskServer {
	store := taskstore.New()
	return &TaskServer{store: store}
}

//-------------------------------------------- Controladores ----------------------------------------//

// CreateTaskHandler godoc
// @Summary Crear una tarea
// @Description Crea una nueva tarea
// @Tags task
// @Accept json
// @Produce json
// @Param task body object true "Nueva tarea"
// @Success 200 {object} map[string]int
// @Failure 400 {string} string
// @Failure 415 {string} string
// @Router /task/ [post]
func (ts *TaskServer) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling task create at %s\n", r.URL.Path)

	type RequestTask struct {
		Text string   `json:"text"`
		Tags []string `json:"tags"`
		Due  string   `json:"due"`
	}

	type ResponseId struct {
		Id int `json:"id"`
	}

	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediaType != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var req RequestTask

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	due, err := time.Parse("2006-01-02 15:04:05", req.Due)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateTask(req.Text, req.Tags, due)
	js, err := json.Marshal(ResponseId{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)

}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
	w.WriteHeader(http.StatusOK)
}

// GetTaskHandler godoc
// @Summary Obtener una tarea
// @Description Obtiene una tarea por ID
// @Tags task
// @Produce json
// @Param id path int true "ID de la tarea"
// @Success 200 {object} taskstore.Task
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /task/{id}/ [get]
func (ts *TaskServer) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling task get at %s\n", r.URL.Path)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := ts.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderJSON(w, task)
}

// DeleteTaskHandler godoc
// @Summary Eliminar una tarea
// @Description Elimina una tarea por ID
// @Tags task
// @Param id path int true "ID de la tarea"
// @Success 204
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /task/{id}/ [delete]
// @Security BasicAuth
func (ts *TaskServer) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling task delete at %s\n", r.URL.Path)

	idTask, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ts.store.DeleteTask(idTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

// GetAllTasksHandler godoc
// @Summary Obtener todas las tareas
// @Description Devuelve todas las tareas
// @Tags task
// @Produce json
// @Success 200 {array} taskstore.Task
// @Failure 500 {string} string
// @Router /task/ [get]
// @Security BasicAuth
func (ts *TaskServer) GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("handling task get all at %s\n", r.URL.Path)
	tasks, err := ts.store.GetAllTasks()
	js, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}

// DeleteAllTasksHandler godoc
// @Summary Eliminar todas las tareas
// @Description Borra todas las tareas
// @Tags task
// @Success 204
// @Failure 500 {string} string
// @Router /task/ [delete]
// @Security BasicAuth
func (ts *TaskServer) DeleteAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling task delete all at %s\n", r.URL.Path)

	err := ts.store.DeleteAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// TagHandler godoc
// @Summary Obtener tareas por tag
// @Description Devuelve tareas que contienen un tag
// @Tags task
// @Produce json
// @Param tag path string true "Tag"
// @Success 200 {array} taskstore.Task
// @Failure 404 {string} string
// @Router /tag/{tag}/ [get]
func (ts *TaskServer) TagHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling task tag at %s\n", r.URL.Path)

	tag := r.PathValue("tag")

	tasks, err := ts.store.GetTasksByTag(tag)

	if err != nil {
		_, _ = w.Write([]byte("Not found Tasks by Tag"))
		return
	}

	js, err := json.Marshal(tasks)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
}

// DueHandler godoc
// @Summary Obtener tareas por fecha
// @Description Devuelve tareas por fecha límite
// @Tags task
// @Produce json
// @Param year path int true "Año"
// @Param month path int true "Mes"
// @Param day path int true "Día"
// @Success 200 {array} taskstore.Task
// @Failure 400 {string} string
// @Router /due/{year}/{month}/{day}/ [get]
func (ts *TaskServer) DueHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling tasks by due at %s\n", req.URL.Path)

	badRequestError := func() {
		http.Error(w, fmt.Sprintf("expect /due/<year>/<month>/<day>, got %v", req.URL.Path), http.StatusBadRequest)
	}

	year, errYear := strconv.Atoi(req.PathValue("year"))
	month, errMonth := strconv.Atoi(req.PathValue("month"))
	day, errDay := strconv.Atoi(req.PathValue("day"))
	if errYear != nil || errMonth != nil || errDay != nil || month < int(time.January) || month > int(time.December) {
		badRequestError()
		return
	}

	tasks, _ := ts.store.GetTasksByDue(year, time.Month(month), day)
	js, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
}
