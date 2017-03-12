package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

const packageTemplate = `package kurento
{{ .Content }}
`

const strTemplate = `
{{ define "Arguments" }}{{ range $i, $e := .Params }}{{ if $i }} , {{ end }}{{ $e.name }} {{ $e.type | checkElement }}{{ end }}{{ end }}
{{ $name := .Name}}

{{/* Generator interface then struct */}}
{{ if ne .Name "MediaObject" }}
type I{{ .Name }} interface {
	{{ range .Methods }}
	{{ .Name | title }}({{ template "Arguments" .}})({{ if .Return.type }}{{ .Return.type }},{{ end }} error)
	{{ end }}
}
{{ end }}

{{ .Doc }}
type {{ .Name }} struct {
	{{ if eq .Name "MediaObject" }}connection *Connection
	{{ else }} {{ .Extends }}
	{{ end }}
	
	{{ range .Properties }}
	{{ .doc }}
	{{ .name | title }} {{ .type }}
	{{ end }}
}

// Return Constructor Params to be called by "Create".
func (elem *{{ .Name }}) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	{{ if len .Constructor.Params }}

	// Create basic constructor params
	ret := map[string]interface{} {
		{{ range .Constructor.Params }}{{ if eq .type "string" "float64" "boolean" "int" }}"{{ .name }}" : {{ .defaultValue }},
		{{ else }} "{{ .name }}" : fmt.Sprintf("%s", from),
		{{ end }}{{ end }}
	}
	
	mergeOptions(ret, options)
	return ret
	{{ else }}return options
	{{ end }}
}

{{ range .Methods }}
{{ .Doc }} 
{{ if .Return.doc }}
// Returns
{{ .Return.doc }}
{{ end }}
func (elem *{{$name}}) {{ .Name | title }}({{ template "Arguments" . }}) ({{ if .Return.type }}{{ .Return.type }}, {{ end}} error) {
	req := elem.getInvokeRequest()
	{{ if .Params }}
	params := make(map[string]interface{})
	{{ range  .Params }}
	setIfNotEmpty(params, "{{ .name }}", {{ .name }})
	{{end}}
	{{ end }}

	req["params"] = map[string]interface{} {
		"operation": "{{ .Name }}",
		"object": elem.Id, 
		{{ if .Params }}
		"operationParams": params,
		{{ end }}
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)
	{{ if .Return }}
	{{ .Return.doc }}
		{{ if eq .Return.type "string" "int" "float64" "boolean" }}
		return response.Result["value"], response.Error
		{{ else }}
		ret := {{ .Return.type }}{}
		return ret, response.Error
		{{ end }}
	{{ else }}
	return response.Error
	{{ end }}
}
{{ end }}
`

const complexTypeTemplate = `
{{ if eq .TypeFormat "ENUM" }}
{{ $name := .Name }}
{{ .Doc }}
type {{.Name}} string
// Implement fmt.Stringer interface
func (t {{.Name}}) String() string {
	return string(t)
}
const (
	{{ range .Values }}
		{{ $name | uppercase }}_{{ . | uppercase}} {{ $name }} = "{{ . }}" 
	{{ end}}
)
{{ else }}
type {{ .Name }} struct {
	{{ range .Properties}}{{ .name | title }} {{ .type }}
	{{ end }}
}
{{ end }}
`

const DOCLINELENGTH = 79

var re = regexp.MustCompile(`(.+)\[\]`)
var symbol = regexp.MustCompile(`<>`)

var CPXTYPES []string

type Core struct {
	RemoteClasses []Class
	ComplexTypes  []ComplexType
}

type Class struct {
	Name        string
	Doc         string
	Abstract    bool
	Properties  []map[string]interface{}
	Extends     string
	Methods     []Method
	Events      []string
	Constructor Constructor
}

type Constructor struct {
	Name   string
	Doc    string
	Params []map[string]interface{}
}

type Method struct {
	Constructor
	Return map[string]interface{}
}

type ComplexType struct {
	TypeFormat string
	Doc        string
	Values     []string
	Name       string
	Properties []map[string]interface{}
}

// template func that change Mediaxxx to IMediaxxx to be
// sure to work with interface.
// Set it global to be used by funcMap["paramValue"] above.
func tplCheckElement(p string) string {
	if len(p) > 5 && p[:5] == "Media" {
		if p[len(p)-4:] != "Type" {
			return "IMedia" + p[5:]
		}
	}
	return p
}

func isComplexType(t string) bool {
	for _, c := range CPXTYPES {
		if c == t {
			return true
		}
	}
	return false
}

var funcMap = template.FuncMap{
	"title":        strings.Title,
	"uppercase":    strings.ToUpper,
	"checkElement": tplCheckElement,
	"paramValue": func(p map[string]interface{}) string {
		name := p["name"].(string)
		t := p["type"].(string)
		t = tplCheckElement(t)

		ctype := isComplexType(t)
		switch t {
		case "float64", "int":
			return fmt.Sprintf("\"%s\" = %s", name, name)
		case "string", "boolean":
			return fmt.Sprintf("\"%s\" = %s", name, name)
		default:
			if !ctype && t[0] == 'I' {
				return fmt.Sprintf("\"%s\" = fmt.Sprintf(\"%%s\", %s)", name, name)
			}
		}
		return fmt.Sprintf("\"%s\" = %s", name, name)
	},
}

func parseComplexTypes(complexs []string, suffix string) {
	var paths []string

	// save kmds file to paths
	for _, path := range complexs {
		pathList, err := filepath.Glob(path)
		if err != nil {
			logFatal(err)
		}
		paths = append(paths, pathList...)
	}

	// parse
	var ret []string
	for _, path := range paths {

		ctypes := getModel(path).ComplexTypes

		if ctypes == nil {
			continue
		}

		for _, ctype := range ctypes {
			CPXTYPES = append(CPXTYPES, ctype.Name)

			ctype.Doc = formatDoc(ctype.Doc)

			for i, p := range ctype.Properties {
				ctype.Properties[i] = formatTypes(p)
			}

			buff := bytes.NewBufferString("")
			tpl, err := template.New("complexttypes").Funcs(funcMap).Parse(complexTypeTemplate)
			if err != nil {

			}

			tpl.Execute(buff, ctype)
			ret = append(ret, buff.String())
		}

		writeFile(createFile(path, suffix), ret)
	}
}

func parseRemotes(remotes []string) {

	var paths []string
	// save kmds file to paths
	for _, path := range remotes {
		pathList, err := filepath.Glob(path)
		if err != nil {
			logFatal(err)
		}
		paths = append(paths, pathList...)
	}

	for _, p := range paths {

		c := getModel(p).RemoteClasses
		ret := make([]string, len(c))

		for idx, cl := range c {

			fmt.Println("Generating ", cl.Name)

			for j, p := range cl.Properties {
				p = formatTypes(p)
				switch p["type"] {
				case "string", "float64", "bool", "[]string":
				default:
					if _, ok := p["type"].(string); ok {
						if p["type"].(string)[:2] == "[]" {
							t := p["type"].(string)[2:]
							if isComplexType(t) {
								p["type"] = "[]*" + t
							} else {
								p["type"] = "[]I" + t
							}
						} else {
							if isComplexType(p["type"].(string)) {
								p["type"] = "*" + p["type"].(string)
							} else {
								p["type"] = "I" + p["type"].(string)
							}
						}
					}
				}
				cl.Properties[j] = p
			}

			for j, m := range cl.Methods {

				for i, p := range m.Params {
					p := formatTypes(p)
					m.Params[i] = p
				}

				m.Doc = formatDoc(m.Doc)

				if m.Return["type"] != nil {
					m.Return = formatTypes(m.Return)
					m.Return["doc"] = formatDoc(m.Return["doc"].(string))
				}

				cl.Methods[j] = m
			}

			for j, p := range cl.Constructor.Params {
				p := formatTypes(p)
				cl.Constructor.Params[j] = p
			}

			tpl, err := template.New("structure").Funcs(funcMap).Parse(strTemplate)
			if err != nil {
				logFatal(err)
			}

			buff := bytes.NewBufferString("")
			cl.Doc = formatDoc(cl.Doc)

			err = tpl.Execute(buff, cl)
			if err != nil {
				logFatal(err)
			}

			ret[idx] = buff.String()
		}

		writeFile(createFile(p, ""), ret)
	}
}

func formatDoc(doc string) string {
	doc = strings.Replace(doc, ":rom:cls:", "", -1)
	doc = strings.Replace(doc, ":term:", "", -1)
	doc = strings.Replace(doc, "``", `"`, -1)
	doc = strings.Replace(doc, "/*", "", -1)
	doc = strings.Replace(doc, "*/", "", -1)

	lines := strings.Split(doc, "\n")

	var part []string

	for _, line := range lines {

		part = append(part, line)
	}

	for i, p := range part {
		part[i] = "/*" + strings.TrimSpace(p) + "*/"
	}

	return strings.Join(part, "\n")
}

func formatTypes(p map[string]interface{}) map[string]interface{} {

	p["doc"] = formatDoc(p["doc"].(string))

	if p["type"] == "String[]" {
		p["type"] = "[]string"
	}

	if p["type"] == "String" || p["type"] == "String<>" {
		p["type"] = "string"
	}

	if p["type"] == "float" || p["type"] == "double" {
		p["type"] = "float64"
	}

	if p["type"] == "boolean" {
		p["type"] = "bool"
	}

	// type
	if re.MatchString(p["type"].(string)) {
		found := re.FindAllStringSubmatch(p["type"].(string), -1)
		p["type"] = "[]" + found[0][1]
	}

	// expect <>
	if symbol.MatchString(p["type"].(string)) {
		substr := p["type"].(string)
		p["type"] = substr[0 : len(substr)-2]
	}

	// default value
	if p["defaultValue"] == "" || p["defaultValue"] == nil {
		switch p["type"] {
		case "string":
			p["defaultValue"] = `""`
		case "bool":
			p["defaultValue"] = "false"
		case "int", "float64":
			p["defaultValue"] = "0"
		}
	}

	return p
}

func getModel(path string) Core {
	i := Core{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		logFatal(err)
	}

	str := strings.Replace(string(data), "\n", " ", -1)
	err = json.Unmarshal([]byte(str), &i)
	if err != nil {
		logFatal(err)
	}
	return i
}

func writeFile(path string, classess []string) {
	content := strings.Join(classess, "\n")
	tpl, _ := template.New("package").Parse(packageTemplate)
	buff := bytes.NewBufferString("")
	tpl.Execute(buff, map[string]string{
		"Content": content,
	})
	ioutil.WriteFile(path, buff.Bytes(), os.ModePerm)
}

func createFile(path string, suffix string) string {
	base := filepath.Base(path)
	base = strings.Replace(base, ".kmd.json", "", -1)
	base = strings.Replace(base, ".", "_", -1)
	base = "kurento/" + base
	if suffix != "" {
		base += "_" + suffix + ".go"
	} else {
		base += ".go"
	}
	return strings.ToLower(base)
}

func logFatal(err error) {
	log.Fatal(err)
}

func main() {

	// Write base data to dst
	data, err := ioutil.ReadFile("kurento_go_base/base.go")
	if err != nil {
		logFatal(err)
	}
	err = ioutil.WriteFile("kurento/base.go", data, os.ModePerm)
	if err != nil {
		logFatal(err)
	}

	// write ws data to dst
	data, err = ioutil.ReadFile("kurento_go_base/kurento_ws.go")
	if err != nil {
		logFatal(err)
	}
	err = ioutil.WriteFile("kurento/kurento_ws.go", data, os.ModePerm)
	if err != nil {
		logFatal(err)
	}

	// ComplexTypes list
	complexList := []string{
		"kms-core/src/server/interface/core.kmd.json",
		"kms-elements/src/server/interface/elements.*.kmd.json",
		"kms-filters/src/server/interface/filters.*.kmd.json",
	}
	parseComplexTypes(complexList, "complext_types")

	// RemoteClasses list
	remoteList := []string{
		"kms-core/src/server/interface/core.kmd.json",
		"kms-elements/src/server/interface/elements.*.kmd.json",
		"kms-filters/src/server/interface/filters.*.kmd.json",
	}
	parseRemotes(remoteList)
}
