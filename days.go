package days

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

// Task - struct for datastore table.
// Contains a summary and the contents of a task, the scheduled
// time for the task and whether it is done or not.
type Task struct {
	Summary    string
	Content    string
	Scheduled  string
	Done       bool
	Identifier string
}

func tasklistkey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Task", "default_tasklist", 0, nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTmpl := template.Must(template.New("home").ParseFiles("templates/home.tmpl",
		"templates/layout.tmpl"))
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Task").Ancestor(tasklistkey(c)).Order("Scheduled").Limit(10)
	tasks := make([]Task, 0, 10)
	if _, err := q.GetAll(c, &tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := homeTmpl.Execute(w, map[string]interface{}{"Pagetitle": "Tasks",
		"tasks": tasks}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func newtask(w http.ResponseWriter, r *http.Request) {
	newTmpl := template.Must(template.New("newtask").ParseFiles("templates/layout.tmpl",
		"templates/newtask.tmpl"))
	if err := newTmpl.Execute(w, map[string]interface{}{"Pagetitle": "New Task"}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func storetask(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	t := Task{Summary: r.FormValue("tinput"),
		Content:    r.FormValue("tarea"),
		Scheduled:  r.FormValue("scheduled"),
		Done:       false,
		Identifier: fmt.Sprintf("%d", time.Now().Unix())}
	key := datastore.NewIncompleteKey(c, "Task", tasklistkey(c))
	_, err := datastore.Put(c, key, &t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
func edittask(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := r.FormValue("taskid")
	var edittask []Task
	q := datastore.NewQuery("Task").Filter("Identifier =", id)
	q.GetAll(c, &edittask)
	tmpl := template.Must(template.New("edittask").ParseFiles("templates/layout.tmpl",
		"templates/edittask.tmpl"))
	tmpl.Execute(w, map[string]interface{}{"Pagetitle": "Edit Tasks",
		"Summary": edittask[0].Summary, "Content": edittask[0].Content,
		"Identifier": id, "Scheduled": edittask[0].Scheduled})
}

func updatetask(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := r.FormValue("taskid")
	scheduled := r.FormValue("scheduled")
	content := r.FormValue("tarea")
	q := datastore.NewQuery("Task").Filter("Identifier =", id)
	var task []Task
	q.GetAll(c, &task)
	//GetOrUpdate(c, content, scheduled)
	task[0].Scheduled = scheduled
	task[0].Content = content
	_, err := datastore.Put(c, tasklistkey(c), &task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func about(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("about").ParseFiles("templates/layout.tmpl",
		"templates/about.tmpl"))
	tmpl.Execute(w, map[string]interface{}{"Pagetitle": "About"})
}

func init() {
	http.HandleFunc("/", home)
	http.HandleFunc("/about", about)
	http.HandleFunc("/storetask", storetask)
	http.HandleFunc("/newtask", newtask)
	http.HandleFunc("/edittask", edittask)
	http.HandleFunc("/updatetask", updatetask)
}
func GetOrUpdate(c appengine.Context, content, scheduled string) error {
	return datastore.RunInTransaction(c, func(c appengine.Context) error {
		task := new(Task)
		err := datastore.Get(c, tasklistkey(c), task)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}
		task.Scheduled = scheduled
		task.Content = content
		_, err = datastore.Put(c, tasklistkey(c), &task)
		return err
	}, nil)
}
