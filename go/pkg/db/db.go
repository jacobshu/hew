package db

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

func (t *devDB) getTasks() ([]task, error) {
	tasksWithProjects, _ := t.db.Query(context.Background(),
		`select id, 
     description, 
     project_id, 
     project.name, 
     depends_on, 
     status, 
     created, 
     completed 
    from tasks 
    left join projects on tasks.project_id = projects.id`)

	results := []task{}
	for tasksWithProjects.Next() {
		var id int32
		var description string
		var project_name string
		var project_id int32
		var depends_on int32
		var status int32
		var created time.Time
		var completed time.Time
		err := tasksWithProjects.Scan(&id, &description, &project_id, &project_name, &depends_on, &status, &created, &completed)
		if err != nil {
			return results, err
		}

		results = append(results, task{
			ID:          id,
			Description: description,
			ProjectID:   project_id,
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
	if task.ProjectName != "" && task.ProjectID == 0 {
		row, _ := t.db.Query(context.Background(), "select id from projects where name = $1", task.ProjectName)
		if row.Next() {
			row.Scan(&project_id)
		} else {
			row, _ := t.db.Query(context.Background(), "insert into projects(name) values($1) returning id", task.ProjectName)
			row.Scan(&project_id)
		}
	}

	_, err := t.db.Exec(context.Background(),
		"insert into tasks(description, project_id, depends_on, created) values($1, $2, $3, $4)", task.Description, project_id, task.DependsOn, time.Now())
	return err

}

func (t *devDB) updateTask(itemNum int32, description string) error {
	_, err := t.db.Exec(context.Background(), "update tasks set description=$1 where id=$2", description, itemNum)
	return err
}

func (t *devDB) removeTask(taskID int32) error {
	_, err := t.db.Exec(context.Background(), "delete from tasks where id=$1", taskID)
	return err
}

func (t *devDB) Update(task task) error {
	_, err := t.db.Exec(context.Background(),
		`update tasks 
    set description=$1,
    project_id=$2,
    status=$3,
    depends_on=$4
    where id=$5`,
		task.Description, task.ProjectID, task.Status, task.DependsOn, task.ID)

	log.Printf("updated %+v documents in %+v")
	return err
}
