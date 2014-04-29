package days

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

// Task - struct for datastore table.
// Contains a summary and the contents of a task, the scheduled
// time for the task and whether it is done or not.
type Task struct {
	User       string
	Summary    string
	Content    string
	Scheduled  string
	Done       string
	Identifier string
}

func parseTime(s string) time.Time {
	layout := "02/01/2006"
	t, _ := time.Parse(layout, s)
	return t
}

// withLayout - take a template name and a templatefile
// and return it combined with layout.tmpl.
func withLayout(name, templ string) *template.Template {
	return template.Must(template.New(name).ParseFiles(templ, "templates/layout.tmpl"))
}

func tasklistkey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Task", "default_tasklist", 0, nil)

}

func home(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	NotSignedIn := ``

	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}
	url, _ := user.LogoutURL(c, "/signout")
	q := datastore.NewQuery("Task").Ancestor(tasklistkey(c)).Filter("User =", fmt.Sprintf("%s", u)).Order("Scheduled").Limit(10)
	tasks := make([]Task, 0, 10)
	if _, err := q.GetAll(c, &tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := withLayout("home", "templates/home.tmpl").Execute(w, map[string]interface{}{"Pagetitle": "Tasks",
		"tasks": tasks, "User": u, "NotSignedIn": NotSignedIn, "Logout": url}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func newtask(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if err := withLayout("newtask", "templates/newtask.tmpl").Execute(w,
		map[string]interface{}{"Pagetitle": "New Task", "User": u}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func storetask(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	user := user.Current(c)
	t := Task{User: fmt.Sprintf("%s", user),
		Summary:    r.FormValue("tinput"),
		Content:    r.FormValue("tarea"),
		Scheduled:  r.FormValue("scheduled"),
		Done:       "Todo",
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
	u := user.Current(c)
	id := r.FormValue("taskid")
	var edittask []Task
	q := datastore.NewQuery("Task").Filter("Identifier =", id)
	q.GetAll(c, &edittask)
	done := edittask[0].Done
	check1, check2 := "", ""
	if done == "Todo" {
		check1, check2 = "checked", ""
	} else {
		check1, check2 = "", "checked"
	}
	withLayout("edittask", "templates/edittask.tmpl").Execute(w,
		map[string]interface{}{"Pagetitle": "Edit Tasks", "User": u,
			"Summary": edittask[0].Summary, "Content": edittask[0].Content,
			"Identifier": id, "Scheduled": edittask[0].Scheduled,
			"Check1": check1, "Check2": check2})
}

func updatetask(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := r.FormValue("taskid")
	scheduled := r.FormValue("scheduled")
	content := r.FormValue("tarea")
	done := r.FormValue("Done")
	q := datastore.NewQuery("Task").Filter("Identifier =", id)
	var edittask []Task
	key, err := q.GetAll(c, &edittask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if btn := r.FormValue("taskidbutton"); btn == "delete" {
		datastore.Delete(c, key[0])
	} else {
		edittask[0].Scheduled = scheduled
		edittask[0].Content = content
		edittask[0].Done = done
		_, nerr := datastore.Put(c, key[0], &edittask[0])
		if nerr != nil {
			http.Error(w, nerr.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)

}

func about(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	url, _ := user.LogoutURL(c, "/signout")
	withLayout("about", "templates/about.tmpl").Execute(w,
		map[string]interface{}{"Pagetitle": "About",
			"Logout": url, "User": u})
}
func signout(w http.ResponseWriter, r *http.Request) {
	withLayout("signout", "templates/signout.tmpl").Execute(w,
		map[string]interface{}{"Pagetitle": "signout"})

}

func init() {
	http.HandleFunc("/", home)
	http.HandleFunc("/about", about)
	http.HandleFunc("/storetask", storetask)
	http.HandleFunc("/newtask", newtask)
	http.HandleFunc("/edittask", edittask)
	http.HandleFunc("/updatetask", updatetask)
	http.HandleFunc("/signout", signout)
}
