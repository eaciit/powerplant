package main

import (
	"net/http"

	"github.com/eaciit/knot/knot.v1"
	_ "github.com/eaciit/powerplant/sec/webapp/webext"
)

func main() {
	app := knot.GetApp("sec")
	if app == nil {
		return
	}

	routes := make(map[string]knot.FnContent, 1)
	routes["/"] = func(r *knot.WebContext) interface{} {
		http.Redirect(r.Writer, r.Request, "/dashboard/default", http.StatusTemporaryRedirect)
		return true
	}
	knot.StartAppWithFn(app, "localhost:8009", routes)
}
