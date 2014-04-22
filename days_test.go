package days

import (
	"appengine/aetest"
	"appengine/datastore"
	//	"reflect"
	"testing"
)

func TestUser(t *testing.T) {
	t1 := Task{Summary: "task1", Identifier: "123"}
	us := User{Email: "test@example.com", Task: t1}
	if us.Email != "test@example.com" {
		t.Error("Expected text@example.com, got: ", us.Email)
	}
	if us.Task.Identifier != "123" {
		t.Error("Expted <123>, got: ", us.Task.Identifier)
	}
	if us.Task.Summary != "task1" {
		t.Error("Exptected <task1>, got: ", us.Task.Summary)
	}
}

func TestTasks(t *testing.T) {
	t1 := Task{Summary: "task1", Content: "some content", Identifier: "123",
		Done: "Done"}
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	key := datastore.NewKey(c, "Task", "", 1, nil)
	if _, err := datastore.Put(c, key, &t1); err != nil {
		t.Fatal(err)
	}
	t2 := Task{}
	if err := datastore.Get(c, key, &t2); err != nil {
		t.Fatal(err)

	}
	// if reflect.TypeOf(t2) !=  {
	// 	t.Fatal(err)
	// }

	defer c.Close()
}
