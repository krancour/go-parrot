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

var featTemplate = `
package {{ .Name }}

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// {{ .Description }}

type Feature interface {
	arcommands.D2CFeature
	lock.ReadLockable
}

type feature struct {
	sync.RWMutex
}

func NewFeature(c2dCommandClient arcommands.C2DCommandClient) Feature {
	return &feature{}
}

func (f *feature) FeatureID() uint8 {
	return {{ .ID }}
}

func (f *feature) FeatureName() string {
	return "{{ .Name }}"
}

func (f *feature) D2CClasses() []arcommands.D2CClass {
	return []arcommands.D2CClass{f}
}

func (f *feature) ClassID() uint8 {
	return 0
}

func (f *feature) ClassName() string {
	return ""
}

func (f *feature) D2CCommands(log *log.Entry) []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
	{{- range .Messages.Events }}
		arcommands.NewD2CCommand(
			{{ .ID }},
			"{{ .Name }}",
			[]interface{}{
			{{- range .Args }}
				{{ .GoType }}(0), // {{ .Name }},
			{{- end }}
			},
			f.{{ .FunctionName }},
			log,
		),
	{{- end }}
	}
}

{{- range .Messages.Events }}

// TODO: Implement this
{{- if .Comment }}
// Title: {{ .Comment.Title }}
// Description: {{ .Comment.Description }}
// Support: {{ .Comment.Support }}
// Triggered: {{ .Comment.Triggered }}
// Result: {{ .Comment.Result }}
{{- end }}
func (f *feature) {{ .FunctionName }}(args []interface{}) error {
	f.Lock()
	defer f.Unlock()
	{{- range $i, $arg := .Args }}
	// {{ $arg.Name }} := ptr.To{{ $arg.UpperGoType }}(args[{{ $i }}].({{ $arg.GoType }}))
	//   {{ $arg.Description }}
	{{- range $j, $enum := $arg.Enums }}
	//   {{ $j }}: {{ $enum.Name }}: {{ $enum.Description }}
	{{- end }}
	{{- end }}
	log.Info("{{ $.Name }}.{{ .FunctionName }}() called")
	return nil
}
{{ end }}
`

type feature struct {
	XMLName      xml.Name  `xml:"feature"`
	ID           int       `xml:"id,attr"`
	Name         string    `xml:"name,attr"`
	Description  string    `xml:",chardata"`
	Enums        *enums    `xml:"enums"`
	Messages     *messages `xml:"msgs"`
	ShortVarName string
}

type enums struct {
	XMLName xml.Name `xml:"enums"`
	Enums   []*enum  `xml:"enum"`
}

type enum struct {
	XMLName     xml.Name     `xml:"enum"`
	Name        string       `xml:"name,attr"`
	Description string       `xml:",chardata"`
	Values      []*enumValue `xml:"value"`
}

type enumValue struct {
	XMLName     xml.Name `xml:"value"`
	Name        string   `xml:"name,attr"`
	Description string   `xml:",chardata"`
}

type messages struct {
	XMLName  xml.Name   `xml:"msgs"`
	Events   []*event   `xml:"evt"`
	Commands []*command `xml:"cmd"`
}

type event struct {
	XMLName      xml.Name `xml:"evt"`
	ID           int      `xml:"id,attr"`
	Name         string   `xml:"name,attr"`
	Comment      *comment `xml:"comment"`
	Args         []*arg   `xml:"arg"`
	FunctionName string
}

type command struct {
	XMLName      xml.Name `xml:"cmd"`
	ID           int      `xml:"id,attr"`
	Name         string   `xml:"name,attr"`
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

func main() {
	// Read
	const xmlFilePath = "/Users/kent/Code/go/src/github.com/krancour/go-parrot/contrib/generators/ref/animation.xml"
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

	// spew.Dump(f)

	// Generate
	if err := generate(f); err != nil {
		log.Fatal(err)
	}
}

func normalize(f *feature) {
	for _, enum := range f.Enums.Enums {
		enum.Description = trim(enum.Description)
	}
	for _, evt := range f.Messages.Events {
		if evt.Comment != nil {
			evt.Comment.Description = trim(evt.Comment.Description)
			evt.Comment.Result = trim(evt.Comment.Result)
			evt.Comment.Triggered = trim(evt.Comment.Triggered)
		}
		for _, arg := range evt.Args {
			arg.Description = trim(arg.Description)
			for _, enum := range arg.Enums {
				enum.Description = trim(enum.Description)
			}
		}
	}
	for _, cmd := range f.Messages.Commands {
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

var space = regexp.MustCompile(`\s+`)

func trim(str string) string {
	str = strings.Replace(str, "\n", " ", -1)
	str = space.ReplaceAllString(str, " ")
	return strings.TrimSpace(str)
}

func generate(feat *feature) error {
	if err := os.RemoveAll("generated"); err != nil {
		return err
	}
	if err := os.Mkdir("generated", 0755); err != nil {
		return err
	}
	featTmpl, err := template.New("feature").Parse(featTemplate)
	if err != nil {
		return err
	}
	for _, e := range feat.Messages.Events {
		e.FunctionName = fmt.Sprintf("%sEvent", e.Name)
		for _, arg := range e.Args {
			if strings.HasPrefix(arg.Type, "bitfield:") || strings.HasPrefix(arg.Type, "enum:") {
				arg.GoType = "int32"
			} else {
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
			}
			arg.UpperGoType = fmt.Sprintf("%s%s", string(arg.GoType[0]+32), arg.GoType[1:len(arg.GoType)])
		}
		f, err := os.Create(fmt.Sprintf("generated/%s.go", feat.Name))
		if err != nil {
			return err
		}
		if err := featTmpl.Execute(f, feat); err != nil {
			return err
		}
		f.Close()
	}
	return nil
}
