package days

import (
        "appengine"
        "appengine/datastore"
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
        return datastore.NewKey(c, "Tasklist", "default_tasklist", 0, nil)
}
