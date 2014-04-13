package days

import (
	"appengine/aetest"
	"appengine/datastore"
	"testing"
)

func TestTasks(t *testing.T) {
	t1 := Task{Summary: "task1", Content: "some content", Identifier: "123",
		Done: false}
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	key := datastore.NewKey(c, "Task", "", 1, nil)
	if _, err := datastore.Put(c, key, &t1); err != nil {
		t.Fatal(err)
	}

	defer c.Close()
}
