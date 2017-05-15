package tasks

import (
	"testing"

	"gopkg.in/mgo.v2"
)

func TestCRUD(t *testing.T) {
	sess, err := mgo.Dial("localhost:27017")
	if err != nil {
		t.Fatalf("Error dialing Mongo: %v", err)
	}
	defer sess.Close()

	store := &MongoStore{
		Session:        sess,
		DatabaseName:   "test",
		CollectionName: "tasks",
	}

	newtask := &NewTask{
		Title: "Learn MongoDB",
		Tags:  []string{"mongo", "info344"},
	}
	task, err := store.Insert(newtask)
	if err != nil {
		t.Errorf("error inserting new task: %v", err)
	}

	task2, err := store.Get(task.ID)
	if err != nil {
		t.Errorf("error fetching task: %v", err)
	}
	if task2.Title != task.Title {
		t.Errorf("task titles didn't match, expected %s but got %s", task.Title, task2.Title)
	}

	sess.DB(store.DatabaseName).C(store.CollectionName).RemoveAll(nil) //Nil -> delete all docs
}
