package days

import (
	"appengine/aetest"
	"appengine/datastore"
	"appengine/user"
	"reflect"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	t1 := Task{Summary: "task1", User: "test@example.com", Identifier: "123"}
	if t1.User != "test@example.com" {
		t.Error("Expected text@example.com, got: ", t1.User)
	}
	if t1.Identifier != "123" {
		t.Error("Expted <123>, got: ", t1.Identifier)
	}
	if t1.Summary != "task1" {
		t.Error("Exptected <task1>, got: ", t1.Summary)
	}
}

func TestTasks(t *testing.T) {
	t1 := Task{Summary: "task1", Content: "some content", Identifier: "123",
		Done: "Done"}
	t2 := Task{Content: "more content", Done: "Todo"}
	c, err := aetest.NewContext(nil)
	u := user.Current(c)
	if u != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	key := datastore.NewKey(c, "Task", "", 1, nil)
	if _, err := datastore.Put(c, key, &t1); err != nil {
		t.Fatal(err)
	}
	nkey := datastore.NewKey(c, "Task", "", 2, nil)
	if _, err := datastore.Put(c, nkey, &t2); err != nil {
		t.Fatal(err)
	}
	gt1 := Task{}
	gt2 := Task{}
	if err := datastore.Get(c, key, &gt1); err != nil {
		t.Fatal(err)

	}
	if err := datastore.Get(c, nkey, &gt2); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gt1, t1) {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gt2, t2) {
		t.Fatal(err)
	}
	// if reflect.TypeOf(t2) !=  {
	// 	t.Fatal(err)
	// }

	defer c.Close()
}

func TestParseTime(t *testing.T) {
	time1 := "01/02/2003"
	time2 := "02/03/2004"
	t1 := parseTime(time1)
	t2 := parseTime(time2)
	if t1.Month() != time.February {
		t.Error("exptected February, got: ", t1.Month())
	}
	if t2.Month() != time.March {
		t.Error("Expected March, got: ", t2.Month())
	}
}

func TestFormatDate(t *testing.T) {
	s1 := "01/02/2003"
	s2 := "03/04/2005"
	time1 := parseTime(s1)
	time2 := parseTime(s2)
	formatted := formatDate(time1)
	formatted2 := formatDate(time2)
	if formatted != s1 {
		t.Error("Expected <01/02/2003>, got: ", formatted)
	}
	if formatted2 != s2 {
		t.Error("Expected <01/02/2005>, got: ", formatted2)
	}
}
