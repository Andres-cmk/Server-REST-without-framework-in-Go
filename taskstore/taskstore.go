package taskstore

import (
	"errors"
	"restServer/graph/model"
	"strconv"
	"sync"
	"time"
)

// simple memoria dentro del servidor
type TaskStore struct {
	// Mutex para evitar condiciones de carrera
	sync.Mutex
	tasks  map[string]model.Task
	nextId int
}

// Funcion para declarar una nueva memoria de Tasks
func New() *TaskStore {
	ts := &TaskStore{}
	ts.tasks = make(map[string]model.Task)
	ts.nextId = 0
	return ts
}

// ------------------------------- Creacion de metodos para la memoria --------------------------------------------------//

// Creacion de una nueva tarea
func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time, attachments []*model.Attachment) string {
	ts.Lock()
	defer ts.Unlock()

	idStr := strconv.Itoa(ts.nextId)

	// creamos una nueva variable de tipo Task
	newTask := model.Task{
		ID:          idStr,
		Text:        text,
		Tags:        tags,
		Due:         due,
		Attachments: attachments,
	}

	// En la memoria guardamos la nueva tarea con el Id asignado
	ts.tasks[idStr] = newTask
	ts.nextId++
	return newTask.ID
}

// O(1) obtenemos la tarea por Id
func (ts *TaskStore) GetTask(id string) (model.Task, error) {
	ts.Lock()
	defer ts.Unlock()

	task, ok := ts.tasks[id]
	if ok {
		return task, nil
	} else {
		return model.Task{}, errors.New("task not found")
	}
}

// O(1) eliminamos la tarea por Id
func (ts *TaskStore) DeleteTask(id string) error {
	ts.Lock()
	defer ts.Unlock()

	if _, ok := ts.tasks[id]; !ok {
		return errors.New("task not found")
	}

	delete(ts.tasks, id)
	return nil
}

// Eliminar las tareas de la memoria, creamos un nuevo mapa vacio
func (ts *TaskStore) DeleteAllTasks() error {
	ts.Lock()
	defer ts.Unlock()

	ts.tasks = make(map[string]model.Task)
	return nil
}

// Obtener todas las tareas de la memoria O(n)
func (ts *TaskStore) GetAllTasks() ([]model.Task, error) {
	ts.Lock()
	defer ts.Unlock()

	allTasks := make([]model.Task, 0, len(ts.tasks)) // declaracion slice con capacidad inicial hasta la cantidad de tareas
	for _, task := range ts.tasks {
		allTasks = append(allTasks, task)
	}

	return allTasks, nil
}

// Obtener tareas por tag O(n) (se guarda en memoria sin orden)
func (ts *TaskStore) GetTasksByTag(tag string) ([]model.Task, error) {

	ts.Lock()
	defer ts.Unlock()

	tasks := make([]model.Task, 0, len(ts.tasks))

	// Debug: ver cuántas tareas hay
	println("DEBUG: Total tareas en store:", len(ts.tasks))
	println("DEBUG: Buscando tag:", tag)

	// Bucle etiquetado para continuar con el siguiente elemento del bucle externo
taskloop: //label asociada a un bucle
	for _, task := range ts.tasks {
		println("DEBUG: Tarea ID:", task.ID, "Text:", task.Text, "Tags:", len(task.Tags))
		for i, taskTag := range task.Tags {
			println("  DEBUG: Tag[", i, "]:", taskTag)
			if taskTag == tag {
				tasks = append(tasks, task)
				continue taskloop
			}
		}
	}

	println("DEBUG: Tareas encontradas:", len(tasks))
	return tasks, nil

}

// Obtener tareas por fecha de vencimiento O(n) (se guarda en memoria sin orden)
func (ts *TaskStore) GetTasksByDue(year int, month time.Month, day int) ([]model.Task, error) {
	ts.Lock()
	defer ts.Unlock()

	tasksMatch := make([]model.Task, 0, len(ts.tasks))

	for _, task := range ts.tasks {
		y, m, d := task.Due.Date()
		if y == year && m == month && d == day {
			tasksMatch = append(tasksMatch, task)
		}
	}

	return tasksMatch, nil
}
