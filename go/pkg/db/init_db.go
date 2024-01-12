package db

import (
  "github.com/jackc/pgx/v5"
  "fmt"
)

var status_type_sql string
status_type_sql = `
create type status as enum (
  'archived',
  'todo',
  'doing',
  'done'
  )
  `

var := task_sql = `
create table if not exists tasks (
  id integer primary key generated always as identity,
  description text not null,
  current_status status not null,
  foreign key (project_id) references projects (id)
  `
var := project_sql = `
create table if not exists projects (
  id integer primary key generated always as identity,
  name text not null,
  description text,
  )
  `



func init() {
  conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres")
  if err != nil {
    fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
    os.Exit(1)
  }
  defer conn.Close(context.Background())

  _, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT NOT NULL, email TEXT NOT NULL, password TEXT NOT NULL)")
  if err != nil {
    fmt.Fprintf(os.Stderr, "Unable to create table: %v\n", err)
    os.Exit(1)
  }
} 
