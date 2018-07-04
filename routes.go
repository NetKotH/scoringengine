package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

const wwwroot = "www"

func handleRoot(w http.ResponseWriter, r *http.Request) {
	logger := log.New(os.Stderr, "[https] ", log.LstdFlags)
	filepath := r.URL.Path
	if strings.HasSuffix(filepath, "/") {
		filepath = path.Join(filepath, "index.html")
	}
	// Figure out template version of filepath
	templatePath := path.Join(wwwroot, filepath+".tmpl")
	// bail if things get weird
	if !strings.HasPrefix(templatePath, wwwroot) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	// Try to open the template file
	templateFile, err := os.Open(templatePath)
	// If the request is for the index.html file, and it doesn't exist, try to create it for reading later
	// This lets the user modify the template without recompiling
	if os.IsNotExist(err) && filepath == "/index.html" {
		err = os.Mkdir(wwwroot, 0755)
		if err != nil {
			logger.Println("Couldn't create www root: %s. Using builtin template.", wwwroot)
		} else {
			logger.Println("Writing out template file")
			ioutil.WriteFile(templatePath, []byte(indexTemplate), 0644)
			templateFile, err = os.Open(templatePath)
		}
	}
	// having an error here means that the template doesn't exist so the code after this block shouldn't execute
	if err != nil {
		// maybe we have a raw file?
		http.FileServer(http.Dir(wwwroot)).ServeHTTP(w, r)
		return
	}

	templateData, err := ioutil.ReadAll(templateFile)
	t, err := template.New("foo").Parse(string(templateData))
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	w.Header().Add("Content-Type", "text/html")

	teams := make(map[string]int)
	for _, e := range Events {
		teams[e.Team] += int(e.Type)
	}
	var teamsSorted []TeamScore
	for t, v := range teams {
		teamsSorted = append(teamsSorted, TeamScore{
			Name:  t,
			Score: v,
		})
	}
	sort.Sort(ByScore(teamsSorted))

	var v []victim
	for _, VM := range victims {
		v = append(v, *VM)
	}
	for i := range v {
		v[i].LastSeenRel = time.Now().Sub(v[i].LastSeen).Round(time.Second)
	}
	sort.Sort(ByIP(v))

	err = t.Execute(w, templateVars{
		Victims: v,
		Teams:   teamsSorted,
	}) // merge.
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	return
}

func handleRootRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://scoreboard.netkoth.org", http.StatusFound)
}
