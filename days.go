package days

import (
        "html/template"
        "net/http"
)

func init() {
        http.HandleFunc("/", root)
        http.HandleFunc("/sign", sign)
}

func root(w http.ResponseWriter, r *http.Request) {
        tm := make(map[string]interface{})
        tm["pagetitle"] = "Index"
        guestbookForm.Execute(w, tm)

}

var guestbookForm = template.Must(template.ParseFiles("templates/layout.tmpl"))

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
