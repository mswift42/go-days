package days

import (
        "github.com/codegangsta/martini"
        "github.com/martini-contrib/render"
        "net/http"
        "time"
)

type Task struct {
        Summary   string
        Scheduled time.Time
        Done      bool
}

func init() {
        m := martini.Classic()
        m.Use(render.Renderer(render.Options{Layout: "layout",
                Directory: "templates"}))
        m.Get("/", func(r render.Render) {
                r.HTML(200, "home", map[string]interface{}{"Pagetitle": "Tasks"})
        })
        m.Get("/newtask", func(r render.Render) {
                r.HTML(200, "newtask", map[string]interface{}{"Pagetitle": "New Task"})
        })
        m.Get("/about", func(r render.Render) {
                r.HTML(200, "about", map[string]interface{}{"Pagetitle": "About"})
        })
        //      http.HandleFunc("/", root)
        http.Handle("/", m)
        http.Handle("/newtask", m)
        http.Handle("/about", m)

}
