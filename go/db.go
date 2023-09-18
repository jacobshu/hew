package main

import (
  "context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/bson/primitive"
	//"os"
	//"reflect"
	"time"
)

type status int

const (
	todo status = iota
	inProgress
	done
)

func (s status) String() string {
	return [...]string{"todo", "in progress", "done"}[s]
}

type task struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
  Name      string
	Project   string
	Status    string
	Created   time.Time
  Completed time.Time `json:"completed,omitempty" bson:"completed,omitempty" optional:"yes"`
}

// implement list.Item & list.DefaultItem
func (t task) FilterValue() string {
	return t.Name
}

func (t task) Title() string {
	return t.Name
}

func (t task) Description() string {
	return t.Project
}

func (s status) Int() int {
	return int(s)
}

type devDB struct {
	db      *mongo.Client
  ctx     context.Context
  closeDb func()
}

func (t *devDB) insert(name, project string) error {
  newTask := task{
    Name: name, 
    Project: project, 
    Status: todo.String(), 
    Created: time.Now(),
  }
  result, err := t.db.Database("dev").Collection("tasks").InsertOne(context.TODO(), newTask)
  if err != nil {
    return err
  }

  fmt.Printf("insert result: %+v", result)

  t.closeDb()
  return nil
}

/*
func (t *taskDB) delete(id uint) error {
	_, err := t.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

// Update the task in the db. Provide new values for the fields you want to
// change, keep them empty if unchanged.
func (t *taskDB) update(task task) error {
	// Get the existing state of the task we want to update.
	orig, err := t.getTask(task.ID)
	if err != nil {
		return err
	}
	orig.merge(task)
	_, err = t.db.Exec(
		"UPDATE tasks SET name = ?, project = ?, status = ? WHERE id = ?",
		orig.Name,
		orig.Project,
		orig.Status,
		orig.ID)
	return err
}

// merge the changed fields to the original task
func (orig *task) merge(t task) {
	uValues := reflect.ValueOf(&t).Elem()
	oValues := reflect.ValueOf(orig).Elem()
	for i := 0; i < uValues.NumField(); i++ {
		uField := uValues.Field(i).Interface()
		if oValues.CanSet() {
			if v, ok := uField.(int64); ok && uField != 0 {
				oValues.Field(i).SetInt(v)
			}
			if v, ok := uField.(string); ok && uField != "" {
				oValues.Field(i).SetString(v)
			}
		}
	}
}

func (t *taskDB) getTasks() ([]task, error) {
	var tasks []task

  coll := client.Database("dev").Collection("tasks")
  filter := bson.D{{"name", "Bagels N Buns"}}

  var result Task
  err = coll.FindOne(context.TODO(), filter).Decode(&result)

  if err != nil {
    if err == mongo.ErrNoDocuments {
      // This error means your query did not match any documents.
      return tasks
    }
		return tasks, fmt.Errorf("unable to get values: %w", err)
  }

  //rows, err := t.db.Query("SELECT * FROM tasks")
	if err != nil {
	}
	for rows.Next() {
		var task task
		err = rows.Scan(
			&task.ID,
			&task.Name,
			&task.Project,
			&task.Status,
			&task.Created,
		)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}
	return tasks, err
}

func (t *taskDB) getTasksByStatus(status string) ([]task, error) {
	var tasks []task
	rows, err := t.db.Query("SELECT * FROM tasks WHERE status = ?", status)
	if err != nil {
		return tasks, fmt.Errorf("unable to get values: %w", err)
	}
	for rows.Next() {
		var task task
		err = rows.Scan(
			&task.ID,
			&task.Name,
			&task.Project,
			&task.Status,
			&task.Created,
		)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}
	return tasks, err
}

func (t *taskDB) getTask(id uint) (task, error) {
	var task task
	err := t.db.QueryRow("SELECT * FROM tasks WHERE id = ?", id).
		Scan(
			&task.ID,
			&task.Name,
			&task.Project,
			&task.Status,
			&task.Created,
		)
	return task, err
}
*/
