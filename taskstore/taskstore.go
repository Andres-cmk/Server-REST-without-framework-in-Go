package taskstore

import (
	"errors"
	"sync"
	"time"
)

// Modelo principal de una tarea para la API REST
type Task struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

// simple memoria dentro del servidor
type TaskStore struct {
	// Mutex para evitar condiciones de carrera
	sync.Mutex
	tasks  map[int]Task
	nextId int
}

// Funcion para declarar una nueva memoria de Tasks
func New() *TaskStore {
	ts := &TaskStore{}
	ts.tasks = make(map[int]Task)
	ts.nextId = 0
	return ts
}

// ------------------------------- Creacion de metodos para la memoria --------------------------------------------------//

// Creacion de una nueva tarea
func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time) int {
	ts.Lock()
	defer ts.Unlock()

	// creamos una nueva variable de tipo Task
	newTask := Task{
		Id:   ts.nextId,
		Text: text,
		Tags: tags,
		Due:  due,
	}

	// En la memoria guardamos la nueva tarea con el Id asignado
	ts.tasks[newTask.Id] = newTask
	ts.nextId++
	return newTask.Id
}

// O(1) obtenemos la tarea por Id
func (ts *TaskStore) GetTask(id int) (Task, error) {
	ts.Lock()
	defer ts.Unlock()

	task, ok := ts.tasks[id]
	if ok {
		return task, nil
	} else {
		return Task{}, errors.New("task not found")
	}
}

// O(1) eliminamos la tarea por Id
func (ts *TaskStore) DeleteTask(id int) error {
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

	ts.tasks = make(map[int]Task)
	return nil
}

// Obtener todas las tareas de la memoria O(n)
func (ts *TaskStore) GetAllTasks() ([]Task, error) {
	ts.Lock()
	defer ts.Unlock()

	allTasks := make([]Task, 0, len(ts.tasks)) // declaracion slice con capacidad inicial hasta la cantidad de tareas
	for _, task := range ts.tasks {
		allTasks = append(allTasks, task)
	}

	return allTasks, nil
}

// Obtener tareas por tag O(n) (se guarda en memoria sin orden)
func (ts *TaskStore) GetTasksByTag(tag string) ([]Task, error) {

	ts.Lock()
	defer ts.Unlock()

	tasks := make([]Task, 0, len(ts.tasks))

	// Bucle etiquetado para continuar con el siguiente elemento del bucle externo
taskloop: //label asociada a un bucle
	for _, task := range ts.tasks {
		for _, taskTag := range task.Tags {
			if taskTag == tag {
				tasks = append(tasks, task)
				continue taskloop
			}
		}
	}
	return tasks, nil

}

// Obtener tareas por fecha de vencimiento O(n) (se guarda en memoria sin orden)
func (ts *TaskStore) GetTasksByDue(year int, month time.Month, day int) ([]Task, error) {
	ts.Lock()
	defer ts.Unlock()

	tasksMatch := make([]Task, 0, len(ts.tasks))

	for _, task := range ts.tasks {
		y, m, d := task.Due.Date()
		if y == year && m == month && d == day {
			tasksMatch = append(tasksMatch, task)
		}
	}

	return tasksMatch, nil
}
