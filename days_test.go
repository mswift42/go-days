package days

import (
	"appengine/aetest"
	"appengine/datastore"
	"testing"
)

func TestTasks(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	key := datastore.NewKey(c, "Task", "", 1, nil)
	if _, err := datastore.Put(c, key, &Task{Summary: "test1", Content: "some content"}); err != nil {
		t.Fatal(err)
	}
	defer c.Close()
}
