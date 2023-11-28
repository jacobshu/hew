package main

import (
	"context"
	"errors"
	"log"
	"reflect"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/jackc/pgx/v5"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	ID          int32
	Description string
	ProjectID   int32
	ProjectName string
	DependsOn   int32
	Status      int32
	Created     time.Time
	Completed   time.Time
}

func (s status) Int() int {
	return int(s)
}

type devDB struct {
	db      *pgx.Conn
	ctx     context.Context
	closeDb func()
}

func (t *devDB) listTasks() ([]task, error) {
	tasksWithProjects, _ := t.db.Query(context.Background(),
		"select * from tasks left join projects on tasks.project_id = projects.id")

	results := []task{}
	for tasksWithProjects.Next() {
		var id int32
		var description string
		var project_name string
		var project int32
		var depends_on int32
		var status int32
		var created time.Time
		var completed time.Time
		err := tasksWithProjects.Scan(&id, &description, &project, &project_name, &depends_on, &status, &created, &completed)
		if err != nil {
			return results, err
		}

		results = append(results, task{
			ID:          id,
			Description: description,
			ProjectID:   project,
			ProjectName: project_name,
			DependsOn:   depends_on,
			Status:      status,
			Created:     created,
			Completed:   completed,
		})
	}

	return results, nil
}

func (t *devDB) addTask(task task) error {
	var project_id int32
  // project name is provided but there's no project id, insert a new project
	if task.ProjectName != "" && task.ProjectID == 0 {
		row, _ := t.db.Query(context.Background(), "select id from projects where name = $1", task.ProjectName)
		row.Scan(&project_id)
	}

	tx, err := t.db.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), "insert into tasks(description, depends_on) values($1, $2)", task.Description, task.DependsOn)
	_, err = tx.Exec(context.Background(), "insert into tasks(description, depends_on) values($1, $2)", task.Description, task.DependsOn)
	return err
}

func (t *devDB) updateTask(itemNum int32, description string) error {
	_, err := t.db.Exec(context.Background(), "update tasks set description=$1 where id=$2", description, itemNum)
	return err
}

func (t *devDB) removeTask(itemNum int32) error {
	_, err := t.db.Exec(context.Background(), "delete from tasks where id=$1", itemNum)
	return err
}

func (t *devDB) ObjectIdFromString(str string) (primitive.ObjectID, error) {
	_id, err := primitive.ObjectIDFromHex(str)
	if err != nil {
		log.Printf("error getting object id: %+v", err)
		return primitive.ObjectIDFromHex("00000000000000000000")
	}
	return _id, nil
}

func (t *devDB) Find(collection string, filter bson.D, opts options.FindOptions) ([]bson.M, error) {
	col := t.db.Database("dev").Collection(collection)
	cursor, err := col.Find(t.ctx, filter)
	if err != nil {
		log.Fatalf("failed to find in %s: %+v", collection, err)
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatalf("failed to iterate on find in %s, %+v", collection, err)
	}

	log.Printf("find in %s: %+v", collection, results)
	return results, err
}

func (t *devDB) Update(collection string, update bson.D, filter bson.D, opts options.UpdateOptions) error {
	col := t.db.Database("dev").Collection("tasks")
	result, err := col.UpdateMany(context.TODO(), filter, update, &opts)
	if err != nil {
		return err
	}

	log.Printf("updated %+v documents in %+v", result.ModifiedCount, collection)
	return nil
}

func (t *devDB) Delete(collection string, filter bson.D, opts *options.DeleteOptions) (int, error) {
	col := t.db.Database("dev").Collection(collection)

	empty := bson.D{{}}
	if cmp.Equal(filter, empty) {
		return 0, errors.New("you don't want to do that")
	}

	result, err := col.DeleteMany(context.TODO(), filter, opts)
	if err != nil {
		return 0, err
	}

	return int(result.DeletedCount), nil
}

func (t *devDB) insertTask(name, project string) error {
	newTask := task{
		Name:    name,
		Project: project,
		Status:  todo.String(),
		Created: time.Now(),
	}

	err := t.InsertOne("task", newTask, *options.InsertOne().SetBypassDocumentValidation(false))
	if err != nil {
		return err
	}

	return nil
}

func (t *devDB) deleteTaskById(strId string) error {
	id, err := t.ObjectIdFromString(strId)
	if err != nil {
		return err
	}
	result, err := t.Delete("tasks", bson.D{{Key: "_id", Value: id}}, &options.DeleteOptions{})
	if err != nil {
		log.Printf("error deleting task with id %+v", strId)
	} else if result > 1 {
		log.Printf("oops, deleted more than 1")
	} else {
		log.Printf("deleted id: %+v", strId)
	}

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

	filter := bson.D{{Key: "_id", Value: orig.ID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "name", Value: orig.Name},
		{Key: "project", Value: orig.Project},
		{Key: "status", Value: orig.Status},
	}}}
	err = t.Update("tasks", update, filter, options.UpdateOptions{})
	if err != nil {
		return err
	}

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
	results, err := t.Find("tasks", bson.D{{}}, options.FindOptions{})
	if err != nil {
		return tasks, err
	}

	for _, result := range results {
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
	var resTask task
	result, err := t.Find("tasks", bson.D{{Key: "_id", Value: id}}, options.FindOptions{})
	if err != nil {
		return resTask, err
	}
	// err := t.tasks.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&task)
	resTask = task{
		ID:      result[0]["_id"].(primitive.ObjectID),
		Name:    result[0]["name"].(string),
		Project: result[0]["project"].(string),
		Status:  result[0]["status"].(string),
		Created: result[0]["created"].(primitive.DateTime).Time(),
		//Completed: result["completed"].(primitive.DateTime).Time(),
	}
	return resTask, err
}
