package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

type feature struct {
	XMLName     xml.Name `xml:"project"`
	ID          int      `xml:"id,attr"`
	Name        string   `xml:"name,attr"`
	Description string   `xml:",chardata"`
	Classes     []*class `xml:"class"`
}

type class struct {
	XMLName     xml.Name   `xml:"class"`
	ID          int        `xml:"id,attr"`
	Name        string     `xml:"name,attr"`
	Description string     `xml:",chardata"`
	Commands    []*command `xml:"cmd"`
}

type command struct {
	XMLName    xml.Name `xml:"cmd"`
	ID         int      `xml:"id,attr"`
	Name       string   `xml:"name,attr"`
	Deprecated bool     `xml:"deprecated,attr"`
	Content    string   `xml:"content,attr"`
	Timeout    string   `xml:"timeout,attr"`
	Comment    *comment `xml:"comment"`
	Args       []*arg   `xml:"arg"`
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
}

type enum struct {
	XMLName     xml.Name `xml:"enum"`
	Name        string   `xml:"name,attr"`
	Description string   `xml:",chardata"`
}

func main() {
	// Read
	const xmlFilePath = "/Users/kent/Downloads/common.xml"
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

	// Output
	fmt.Println("feature:")
	fmt.Printf("  id: %d\n", f.ID)
	fmt.Printf("  name: %s\n", f.Name)
	fmt.Printf("  description: %s\n", f.Description)
	if len(f.Classes) > 0 {
		fmt.Println("  classes:")
	}
	for _, c := range f.Classes {
		fmt.Printf("  - id: %d\n", c.ID)
		fmt.Printf("    name: %s\n", c.Name)
		fmt.Printf("    description: %s\n", c.Description)
		if len(c.Commands) > 0 {
			fmt.Println("    commands:")
		}
		for _, cmd := range c.Commands {
			fmt.Printf("    - id: %d\n", cmd.ID)
			fmt.Printf("      name: %s\n", cmd.Name)
			fmt.Printf("      deprecated: %t\n", cmd.Deprecated)
			if cmd.Comment != nil {
				fmt.Println("      comment:")
				fmt.Printf("        title: %s\n", cmd.Comment.Title)
				fmt.Printf("        description: %s\n", cmd.Comment.Description)
				fmt.Printf("        support: %s\n", cmd.Comment.Support)
				fmt.Printf("        triggered: %s\n", cmd.Comment.Triggered)
				fmt.Printf("        result: %s\n", cmd.Comment.Result)
			}
			if len(cmd.Args) > 0 {
				fmt.Println("      args:")
			}
			for _, arg := range cmd.Args {
				fmt.Printf("      - name: %s\n", arg.Name)
				fmt.Printf("        type: %s\n", arg.Type)
				fmt.Printf("        description: %s\n", arg.Description)
				if len(arg.Enums) > 0 {
					fmt.Println("        enums:")
				}
				for _, enum := range arg.Enums {
					fmt.Printf("        - name: %s\n", enum.Name)
					fmt.Printf("          description: %s\n", enum.Description)
				}
			}
		}
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
