package days

import (
        "github.com/codegangsta/martini"
        "github.com/martini-contrib/render"
        "net/http"
)

func init() {
        m := martini.Classic()
        m.Use(render.Renderer(render.Options{Layout: "layout",
                Directory: "templates"}))
        m.Get("/", func(r render.Render) {
                r.HTML(200, "home", map[string]interface{}{"Pagetitle": "Tasks"})
        })
        //      http.HandleFunc("/", root)
        http.Handle("/", m)
}
