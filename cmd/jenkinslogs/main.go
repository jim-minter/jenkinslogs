package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/jim-minter/jenkinslogs/pkg/detectors"
	"github.com/jim-minter/jenkinslogs/pkg/types"
)

var jenkinsRoot = flag.String("jenkinsRoot", "https://ci.openshift.redhat.com/jenkins/job", "jenkins root")
var addr = flag.String("l", "127.0.0.1:8080", "listen address")

func downloadAndSave(url string) (io.ReadCloser, error) {
	fn := strings.Replace(url, "/", "-", -1)
	if _, err := os.Stat(fn); err == nil {
		return os.Open(fn)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	f, err := os.Create(fn)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		f.Close()
		return nil, err
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		f.Close()
		return nil, err
	}

	return f, nil
}

func download(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	rx := regexp.MustCompile(`^/[a-z_]+/\d+/?$`)
	if !rx.MatchString(r.URL.Path) {
		http.NotFound(w, r)
		return
	}

	rc, err := downloadAndSave(*jenkinsRoot + strings.TrimRight(r.URL.Path, "/") + "/consoleText")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rc.Close()

	scanner := bufio.NewScanner(rc)
	s := []string{}
	for scanner.Scan() {
		if err = scanner.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s = append(s, scanner.Text())
	}
	g := types.NewGroup(0, len(s), nil, "jenkins", false)

	for _, detector := range detectors.Detectors {
		g.Walk(s, nil, nil, detector)
	}

	writeHeader(w)

	fmt.Fprintf(w, `<div style="font-family:monospace;">`)

	fmt.Fprintf(w, `[<a href="%s">jenkins</a>]`, *jenkinsRoot+strings.TrimRight(r.URL.Path, "/"))
	if r.URL.Query().Get("expanded") == "" {
		fmt.Fprintf(w, `<br>[<a href="%s?expanded=true">expanded</a>]`, r.URL.Path)
	} else {
		fmt.Fprintf(w, `<br>[<a href="%s">expanded</a>]`, r.URL.Path)
	}
	fmt.Fprintf(w, `<br><br></div><div style="font-family:monospace;">`)

	headlineTemplate := template.Must(template.New("headlineTemplate").Parse("<a href=\"#{{.ID}}\">{{.CollapsedText}}</a><br>\n"))

	g.Walk(s, func(g *types.Group) {
		if g.Expanded && g.ID != 1 {
			headlineTemplate.Execute(w, g)
		}
	}, nil, nil)

	fmt.Fprintf(w, `<br></div><div style="font-family:monospace;">`)

	g.Walk(s, nil, func(g *types.Group) {
		for _, child := range g.Children {
			if g.Expanded {
				break
			}
			g.Expanded = g.Expanded || child.Expanded
		}
	}, nil)

	if r.URL.Query().Get("expanded") != "" {
		var f func(g *types.Group)
		f = func(g *types.Group) {
			g.Expanded = true
			for _, child := range g.Children {
				f(child)
			}
		}
		f(g)
	}

	for _, child := range g.Children {
		writeTo(w, child, s)
	}
	writeFooter(w)
}

func writeHeader(w io.Writer) (int, error) {
	header := `<html>
  <head>
    <meta charset="UTF-8">
    <script src="http://ajax.googleapis.com/ajax/libs/jquery/1.9.1/jquery.min.js"></script>
  </head>
  <body>
    `

	return fmt.Fprint(w, header)
}

func writeFooter(w io.Writer) (int, error) {
	footer := `</div>
  </body>
  <script>
$(function() {
  $(".a").click(function() {
    $(this).parent().parent().children().toggle();
  });
});
  </script>
</html>`

	return fmt.Fprint(w, footer)
}

var groupTemplateHeader = template.Must(template.New("groupTemplateHeader").Parse(`
  <div>
  <div{{if not .Expanded}} style="display:none;"{{end}}>
  <span class="a">▾ <a name="{{.ID}}">{{.CollapsedText}}</a></span> [<a href="#{{.ID}}-end">end</a>]<br>
`))

var groupTemplateFooter = template.Must(template.New("groupTemplateFooter").Parse(`
  </div>
  <a name="{{.ID}}-end"></a>
  <div{{if .Expanded}} style="display:none;"{{end}}>
  <div class="a">▸ {{.CollapsedText}}</div>
  </div>
  </div>
`))

var lineTemplate = template.Must(template.New("lineTemplate").Parse("&nbsp;&nbsp;{{.}}<br>\n"))

func writeTo(w io.Writer, g *types.Group, s []string) {
	groupTemplateHeader.Execute(w, g)

	c := g.I
	for _, child := range g.Children {
		for ; c < child.I; c++ {
			lineTemplate.Execute(w, s[c])
		}
		writeTo(w, child, s)
		c = child.J
	}
	for ; c < g.J; c++ {
		lineTemplate.Execute(w, s[c])
	}

	groupTemplateFooter.Execute(w, g)
}

func main() {
	fmt.Println("running")
	flag.Parse()

	panic(http.ListenAndServe(*addr, handlers.CombinedLoggingHandler(os.Stderr, handlers.CompressHandler(http.HandlerFunc(handler)))))
}
