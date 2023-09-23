package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"reflect"
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
	Name      string             `json:"name" bson:"name"`
	Project   string             `json:"project" bson:"project"`
	Status    string             `json:"status" bson:"status"`
	Created   time.Time          `json:"created" bson:"created"`
	Completed time.Time          `json:"completed,omitempty" bson:"completed,omitempty" optional:"yes"`
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
	tasks   *mongo.Collection
	links   *mongo.Collection
}

func (t *devDB) ObjectIdFromString(str string) (primitive.ObjectID, error) {
	_id, err := primitive.ObjectIDFromHex(str)
	if err != nil {
		log.Printf("error getting object id: %+v", err)
		return primitive.ObjectIDFromHex("00000000000000000000")
	}
	return _id, nil
}

func (t *devDB) insertTask(name, project string) error {
	newTask := task{
		Name:    name,
		Project: project,
		Status:  todo.String(),
		Created: time.Now(),
	}

	result, err := t.tasks.InsertOne(context.TODO(), newTask)
	if err != nil {
		return err
	}

	log.Printf("insert task: %+v", result)

	t.closeDb()
	return nil
}

func (t *devDB) deleteTaskById(strId string) error {
	id, err := t.ObjectIdFromString(strId)
	if err != nil {
		return err
	}
	result, err := t.tasks.DeleteOne(context.TODO(), bson.D{{"_id", id}})
	log.Printf("deleted: %+v", result)
	t.closeDb()
	return err
}

// Update the task in the db. Provide new values for the fields you want to
// change, keep them empty if unchanged.
func (t *devDB) updateTask(strId string, task task) error {
	id, err := t.ObjectIdFromString(strId)
	if err != nil {
		return err
	}

	orig, err := t.getTask(id)
	if err != nil {
		return err
	}
	orig.merge(task)

	filter := bson.D{{"_id", orig.ID}}
	update := bson.D{{"$set", bson.D{
		{"name", orig.Name},
		{"project", orig.Project},
		{"status", orig.Status},
	}}}
	result, err := t.tasks.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	log.Printf("updated: %+v", result)
	t.closeDb()
	return nil
}

// merge the changed fields to the original task
func (orig *task) merge(t task) {
	//log.Printf("merging, \n%+v \nwith \n%+v", orig, t)
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

func (t *devDB) getTasks() ([]task, error) {
	var tasks []task

	var results []bson.M
	cursor, err := t.tasks.Find(context.TODO(), bson.D{{}})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return tasks, nil
		}
		return tasks, fmt.Errorf("unable to get values: %w", err)
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	for _, result := range results {
		//cursor.Decode(&result)
		var task = task{
			ID:      result["_id"].(primitive.ObjectID),
			Name:    result["name"].(string),
			Project: result["project"].(string),
			Status:  result["status"].(string),
			Created: result["created"].(primitive.DateTime).Time(),
			//Completed: result["completed"].(primitive.DateTime).Time(),
		}

		tasks = append(tasks, task)
	}

	t.closeDb()
	return tasks, err
}

/*
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
*/

func (t *devDB) getTask(id primitive.ObjectID) (task, error) {
	var task task
	err := t.tasks.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&task)
	return task, err
}
