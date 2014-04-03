package days

import (
        "appengine"
        "appengine/datastore"
        "html/template"
        "net/http"
        "time"
)

type Task struct {
        Summary   string
        Content   string
        Scheduled time.Time
        Done      bool
}

// guestbookKey returns the key used for all guestbook entries.
func tasklistkey(c appengine.Context) *datastore.Key {
        // The string "default_guestbook" here could be varied to have multiple guestbooks.
        return datastore.NewKey(c, "Task", "default_tasklist", 0, nil)
}

func home(w http.ResponseWriter, r *http.Request) {
        // home := make(map[string]*template.Template)
        // homeTmpl["home.tmpl"]
        homeTmpl := template.Must(template.New("home").ParseFiles("templates/home.tmpl", "templates/layout.tmpl"))
        c := appengine.NewContext(r)
        q := datastore.NewQuery("Task").Ancestor(tasklistkey(c)).Order("-Date")
        tasks := make([]Task, 0, 10)
        if _, err := q.GetAll(c, &tasks); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        if err := homeTmpl.Execute(w, map[string]interface{}{"Pagetitle": "Tasks",
                "Tasks": tasks}); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        }
}
func newtask(w http.ResponseWriter, r *http.Request) {
        newTmpl := template.Must(template.New("newtask").ParseFiles("templates/layout.tmpl", "templates/newtask.tmpl"))
        newTmpl.Execute(w, map[string]interface{}{"Pagetitle": "New Task"})
}

func storetask(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        t := Task{Summary: r.FormValue("summary"),
                Content:   r.FormValue("content"),
                Scheduled: time.Now()}
        key := datastore.NewIncompleteKey(c, "Task", tasklistkey(c))
        _, err := datastore.Put(c, key, &t)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        http.Redirect(w, r, "/", http.StatusFound)
}

func about(w http.ResponseWriter, r *http.Request) {
        tmpl := template.Must(template.New("about").ParseFiles("templates/layout.tmpl", "templates/about.tmpl"))
        tmpl.Execute(w, map[string]interface{}{"Pagetitle": "About"})
}

func init() {
        http.HandleFunc("/", home)
        http.HandleFunc("/about", about)
        http.HandleFunc("/storetask", storetask)
        http.HandleFunc("/newtask", newtask)
}
