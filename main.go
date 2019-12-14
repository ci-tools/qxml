package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"os"
)

func init() {

}

// Node ...
type Node struct {
	XMLName    xml.Name
	Attributes []xml.Attr `xml:",any,attr"`
	Data       string     `xml:",chardata"`
	Nodes      []Node     `xml:",any"`
}

// JSONify ...
func (n Node) JSONify() interface{} {
	if len(n.Nodes) > 0 {
		nodes := map[string]interface{}{}
		for _, subn := range n.Nodes {
			nodes[subn.XMLName.Local] = subn.JSONify()
		}
		return nodes
	}
	return n.Data
}

func exec() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("Usage: cat file.xml | qxml 'template'")
	}
	v := &Node{}
	if err := xml.NewDecoder(os.Stdin).Decode(&v); err != nil {
		return err
	}
	// fmt.Printf("%v", v.XMLName)
	asJSON := map[string]interface{}{
		v.XMLName.Local: v.JSONify(),
	}
	type Inventory struct {
		Material string
		Count    uint
	}
	tmpl, err := template.New("result").Parse(os.Args[1])
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, asJSON)
	// if err != nil {
	// 	return err
	// }
	// pretty.JSON(os.Stdout, asJSON)
	// return nil
}

func main() {
	if err := exec(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
