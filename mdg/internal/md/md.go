package md

import (
	"bytes"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"html/template"
	"os/exec"
	"runtime"
	"time"
)

const defaultTemplate = `<!DOCTYPE html>
<html>
<head>
<meta http-equiv="content-type" content="text/html; charset=utf-8">
<title>{{ .Title }}</title>
</head>
<body>
{{ .Body }}
</body>
</html>
`

type content struct {
	Title string
	Body  template.HTML
}

func ParseContent(input []byte, tFname string) ([]byte, error) {
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	t := template.New("mdg")
	var err error
	if tFname != "" {
		t, err = template.ParseFiles(tFname)
	} else {
		t, err = t.Parse(defaultTemplate)
	}
	if err != nil {
		return nil, err
	}

	c := content{
		Title: "Markdown Preview Glancer",
		Body:  template.HTML(body),
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, c); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Preview(fname string) error {
	var cmdName string
	var cmdArgs []string

	switch runtime.GOOS {
	case "linux":
		cmdName = "xdg-open"
	case "windows":
		cmdName = "cmd.exe"
		cmdArgs = []string{"/C", "start"}
	case "darwin":
		cmdName = "open"
	default:
		return fmt.Errorf("unsupported OS")
	}

	cmdArgs = append(cmdArgs, fname)

	path, err := exec.LookPath(cmdName)
	if err != nil {
		return err
	}

	err = exec.Command(path, cmdArgs...).Run()
	time.Sleep(2 * time.Second)

	return err
}
