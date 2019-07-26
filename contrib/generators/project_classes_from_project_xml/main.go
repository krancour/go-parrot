package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var classTemplate = `
package {{ .PackageName }}

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// {{ .Description }}

// {{ .InterfaceName }} ...
// TODO: Document this
type {{ .InterfaceName }} interface{
	lock.ReadLockable
}

type {{ .StructName }} struct{
	sync.RWMutex
}

func ({{ .ShortVarName }} *{{ .StructName }}) ID() uint8 {
	return {{ .ID }}
}

func ({{ .ShortVarName }} *{{ .StructName }}) Name() string {
	return "{{ .Name }}"
}

func ({{ .ShortVarName }} *{{ .StructName }}) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
	{{- range .Commands }}
		arcommands.NewD2CCommand(
			{{ .ID }},
			"{{ .Name }}",
			[]interface{}{
			{{- range .Args }}
				{{ .GoType }}(0), // {{ .Name }},
			{{- end }}
			},
			{{ $.ShortVarName }}.{{ .FunctionName }},
		),
	{{- end }}
	}
}

{{- range .Commands }}

// TODO: Implement this
{{- if .Comment }}
// Title: {{ .Comment.Title }}
// Description: {{ .Comment.Description }}
// Support: {{ .Comment.Support }}
// Triggered: {{ .Comment.Triggered }}
// Result: {{ .Comment.Result }}
{{- end }}
{{- if .Deprecated }}
// WARNING: Deprecated
{{- end }}
func ({{ $.ShortVarName }} *{{ $.StructName }}) {{ .FunctionName }}(args []interface{}) error {
	{{ $.ShortVarName }}.Lock()
	defer {{ $.ShortVarName }}.Unlock()
	{{- range $i, $arg := .Args }}
	// {{ $arg.Name }} := ptr.To{{ $arg.UpperGoType }}(args[{{ $i }}].({{ $arg.GoType }}))
	//   {{ $arg.Description }}
	{{- range $j, $enum := $arg.Enums }}
	//   {{ $j }}: {{ $enum.Name }}: {{ $enum.Description }}
	{{- end }}
	{{- end }}
	log.Info("{{ $.PackageName }}.{{ .FunctionName }}() called")
	return nil
}
{{ end }}
`

type feature struct {
	XMLName     xml.Name `xml:"project"`
	ID          int      `xml:"id,attr"`
	Name        string   `xml:"name,attr"`
	Description string   `xml:",chardata"`
	Classes     []*class `xml:"class"`
}

type class struct {
	XMLName       xml.Name   `xml:"class"`
	ID            int        `xml:"id,attr"`
	Name          string     `xml:"name,attr"`
	Description   string     `xml:",chardata"`
	Commands      []*command `xml:"cmd"`
	PackageName   string
	InterfaceName string
	StructName    string
	ShortVarName  string
}

type command struct {
	XMLName      xml.Name `xml:"cmd"`
	ID           int      `xml:"id,attr"`
	Name         string   `xml:"name,attr"`
	Deprecated   bool     `xml:"deprecated,attr"`
	Content      string   `xml:"content,attr"`
	Timeout      string   `xml:"timeout,attr"`
	Comment      *comment `xml:"comment"`
	Args         []*arg   `xml:"arg"`
	FunctionName string
}

type comment struct {
	XMLName     xml.Name `xml:"comment"`
	Title       string   `xml:"title,attr"`
	Description string   `xml:"desc,attr"`
	Support     string   `xml:"support,attr"`
	Triggered   string   `xml:"triggered,attr"`
	Result      string   `xml:"result,attr"`
}

type arg struct {
	XMLName     xml.Name `xml:"arg"`
	Name        string   `xml:"name,attr"`
	Type        string   `xml:"type,attr"`
	Description string   `xml:",chardata"`
	Enums       []*enum  `xml:"enum"`
	GoType      string
	UpperGoType string
}

type enum struct {
	XMLName     xml.Name `xml:"enum"`
	Name        string   `xml:"name,attr"`
	Description string   `xml:",chardata"`
}

func main() {
	// Read
	const xmlFilePath = "/Users/kent/animation.xml"
	xmlFile, err := os.Open(xmlFilePath)
	if err != nil {
		log.Fatal(err)
	}
	xmlBytes, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Fatal(err)
	}

	// Parse
	f := &feature{}
	if err := xml.Unmarshal(xmlBytes, f); err != nil {
		log.Fatal(err)
	}

	// Normalize
	normalize(f)

	// Generate
	if err := generate(f); err != nil {
		log.Fatal(err)
	}
}

func normalize(f *feature) {
	f.Description = trim(f.Description)
	for _, c := range f.Classes {
		c.Description = trim(c.Description)
		for _, cmd := range c.Commands {
			if cmd.Comment != nil {
				cmd.Comment.Description = trim(cmd.Comment.Description)
				cmd.Comment.Result = trim(cmd.Comment.Result)
				cmd.Comment.Triggered = trim(cmd.Comment.Triggered)
			}
			for _, arg := range cmd.Args {
				arg.Description = trim(arg.Description)
				for _, enum := range arg.Enums {
					enum.Description = trim(enum.Description)
				}
			}
		}
	}
}

var space = regexp.MustCompile(`\s+`)

func trim(str string) string {
	str = strings.Replace(str, "\n", " ", -1)
	str = space.ReplaceAllString(str, " ")
	return strings.TrimSpace(str)
}

func generate(f *feature) error {
	if err := os.RemoveAll("generated"); err != nil {
		return err
	}
	if err := os.Mkdir("generated", 0755); err != nil {
		return err
	}
	stateRegex := regexp.MustCompile(`State`)
	eventRegex := regexp.MustCompile(`Event`)
	classTmpl, err := template.New("class").Parse(classTemplate)
	if err != nil {
		return err
	}
	for _, c := range f.Classes {
		if !stateRegex.MatchString(c.Name) && !eventRegex.MatchString(c.Name) {
			continue
		}
		c.PackageName = f.Name
		c.InterfaceName = c.Name
		c.StructName = fmt.Sprintf("%s%s", string(c.Name[0]+32), c.Name[1:len(c.Name)])
		c.ShortVarName = string(c.StructName[0])
		for _, cmd := range c.Commands {
			cmd.FunctionName = fmt.Sprintf("%s%s", string(cmd.Name[0]+32), cmd.Name[1:len(cmd.Name)])
			for _, arg := range cmd.Args {
				switch arg.Type {
				case "u8":
					arg.GoType = "uint8"
				case "i8":
					arg.GoType = "int8"
				case "u16":
					arg.GoType = "uint16"
				case "i16":
					arg.GoType = "int16"
				case "u32":
					arg.GoType = "uint32"
				case "i32":
					arg.GoType = "int32"
				case "u64":
					arg.GoType = "uint64"
				case "i64":
					arg.GoType = "int64"
				case "float":
					arg.GoType = "float32"
				case "double":
					arg.GoType = "float64"
				case "string":
					arg.GoType = "string"
				case "enum":
					arg.GoType = "int32"
				default:
					arg.GoType = "unknown"
				}
				arg.UpperGoType = fmt.Sprintf("%s%s", string(arg.GoType[0]+32), arg.GoType[1:len(arg.GoType)])
			}
		}
		f, err := os.Create(fmt.Sprintf("generated/%s.go", c.StructName))
		if err != nil {
			return err
		}
		if err := classTmpl.Execute(f, c); err != nil {
			return err
		}
		f.Close()
	}
	return nil
}
