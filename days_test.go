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
func TestElapsedDays(t *testing.T) {
	day1 := parseTime("01/01/2000")
	day2 := parseTime("02/01/2000")
	day3 := parseTime("05/01/2000")
	dur := elapsedDays(day1, day2)
	dur2 := elapsedDays(day1, day3)
	dur3 := elapsedDays(day3, day2)
	if dur != 1 {
		t.Error("Exptected <1>, got: ", dur)
	}
	if dur2 != 4 {
		t.Error("Exptected <4>, got: ", dur2)
	}
	if dur3 != -3 {
		t.Error("Expected <-3>, got: ", dur3)
	}
}

func TestAddDay(t *testing.T) {
	startday := parseTime("01/01/2000")
	fday := formatDate(addDay(startday, 1))
	if fday != "02/01/2000" {
		t.Error("Expected <02/01/2000>, got: ", fday)
	}
}

func TestWeekDates(t *testing.T) {
	startday := "01/01/2000"
	week := weekDates(startday)
	sndday := week[1]
	if sndday.Day() != 2 {
		t.Error("Expected <2>, got: ", sndday.Day())
	}
	thirddate := week[2]
	thirddatemonth := thirddate.Month()
	if thirddatemonth != time.January {
		t.Error("Expected <January>, got: ", thirddatemonth)
	}
	fourthdate := week[3]
	formateddate := formatDate(fourthdate)
	if formateddate != "04/01/2000" {
		t.Error("Expected <04/01/2000>, got: ", formateddate)
	}
}
