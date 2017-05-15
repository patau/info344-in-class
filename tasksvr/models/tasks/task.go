package tasks

import (
	"fmt"
	"time"
)

//NewTask represents a new task posted to the server
type NewTask struct {
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}

//Task represents a task stored in the database
type Task struct {
	ID         interface{} `json:"id" bson:"_id"` //interface{} - any data type accepted
	Title      string      `json:"title"`
	Tags       []string    `json:"tags"`
	CreatedAt  time.Time   `json:"createdAt"`
	ModifiedAt time.Time   `json:"modifiedAt"`
	Complete   bool        `json:"complete"`
}

//Validate will validate the NewTask
func (nt *NewTask) Validate() error {
	//Title field must be non-zero length
	if len(nt.Title) == 0 {
		return fmt.Errorf("title must be longer")
	}
	return nil
}

//ToTask converts a NewTask to a Task
func (nt *NewTask) ToTask() *Task {
	//Creating new task here
	t := &Task{
		Title:      nt.Title,
		Tags:       nt.Tags,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}
	return t
}
