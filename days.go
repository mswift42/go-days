package days

import (
        "github.com/codegangsta/martini"
        "github.com/martini-contrib/render"
        "html/template"
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
        http.HandleFunc("/sign", sign)
}

func sign(w http.ResponseWriter, r *http.Request) {
        err := signTemplate.Execute(w, r.FormValue("content"))
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        }
}

var signTemplate = template.Must(template.New("sign").Parse(signTemplateHtml))

const signTemplateHtml = `
<html>
  <body>
    <p>You wrote:</p>
    <pre>{{.}}</pre>
  </body>
</html>
`
