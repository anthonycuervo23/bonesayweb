package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	hapesay "github.com/anthonycuervo23/bonesay/v2"
)

type Bonesay struct {
	// Type of hape
	typ string
	// What the hape should say
	say string
}

type htmlVars struct {
	Bones    map[int]map[string]string
	AllBones map[int]string
	Selected string
}

func hapeList() []string {
	hapes, err := hapesay.Bones()
	if err != nil {
		return hapesay.BonesInBinary()
	}
	list := make([]string, 0)
	for _, hape := range hapes {
		list = append(list, hape.BoneFiles...)
	}
	return list
}

func serveTemplate(w http.ResponseWriter, req *http.Request) {
	html := &htmlVars{Bones: make(map[int]map[string]string, 0), AllBones: make(map[int]string, 0)}
	lp := filepath.Join("templates", "layout.html")

	hapes := hapeList()

	hape := &Bonesay{typ: "mobile", say: "Hello"}

	route := strings.Split(req.URL.Path, "/")

	if len(route) > 2 {
		hape.typ = route[1]
		hape.say = route[2]

		say, err := hapesay.Say(
			hape.say,
			hapesay.Type(hape.typ),
			hapesay.BallonWidth(15),
		)

		if err != nil {
			say, _ = hapesay.Say(
				"Error 404: Bone not found",
				hapesay.Type("mobile"),
				hapesay.BallonWidth(15),
			)
			w.WriteHeader(404)
		}

		html.Bones[0] = map[string]string{hape.typ: say}

		counter := 0
		for _, hapeFile := range hapes {
			html.AllBones[counter] = hapeFile
			counter++
		}

		html.Selected = hape.typ
	} else {
		if route[1] == "" {
			hape.say = "You can make me say anything, just type it in the 'Text to say' field above and press enter!"
		} else {
			hape.say = route[1]
		}

		say, _ := hapesay.Say(
			hape.say,
			hapesay.Type("mobile"),
			hapesay.BallonWidth(15),
		)
		html.Bones[0] = map[string]string{"mobile": say}
		html.AllBones[0] = "mobile"

		counter := 1
		for _, hapeFile := range hapes {
			if hapeFile == "mobile" {
				continue
			}
			say, _ := hapesay.Say(
				hape.say,
				hapesay.Type(hapeFile),
				hapesay.BallonWidth(15),
			)

			html.AllBones[counter] = hapeFile
			html.Bones[counter] = map[string]string{hapeFile: say}
			counter++
		}
	}

	// res := strings.Split(say, "\n")
	// fmt.Fprintf(w, "%v", res[2])
	tmpl, _ := template.ParseFiles(lp)
	tmpl.ExecuteTemplate(w, "layout", html)
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipAddress := r.RemoteAddr
		fwdAddress := r.Header.Get("X-Forwarded-For")
		if fwdAddress != "" {
			// Got X-Forwarded-For
			ipAddress = fwdAddress // If it's a single IP, then awesome!

			// If we got an array... grab the first IP
			ips := strings.Split(fwdAddress, ", ")
			if len(ips) > 1 {
				ipAddress = ips[0]
			}
		}
		fmt.Printf("%s %s %s\n", ipAddress, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", serveTemplate)

	log.Println("Listening on :8000...")
	err := http.ListenAndServe(":8000", Log(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}
