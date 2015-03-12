package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template/parse"
)

var funcs = map[string]interface{}{
	"gen_markdown": dummy,
}

func main() {
	w := os.Stdout
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	files , err := filepath.Glob(filepath.Join(cwd, "*.html"))
	if err != nil {
		log.Fatal(err)
	}
	w.WriteString("package main\n")
	for _, f := range files {
		w.WriteString("\n")
		Parse(f, w)
	}
}

func Parse(filename string, w io.Writer) {
	content, _ := ioutil.ReadFile(filename)
	treeSet, _ := parse.Parse("index", string(content), "{{", "}}", funcs)
	p := &Printer{1, w}
	funcName := strings.Replace(filepath.Base(filename), filepath.Ext(filename), "", -1)
	fmt.Println("func " + funcName + "(w io.Writer, arg *View) {")
	for _, node := range treeSet["index"].Root.Nodes {
		p.Dispatch(node)
	}
	fmt.Println("}")
}

type Printer struct {
	indent int
	writer io.Writer
}

func (p *Printer) Write(s string) {
	io.WriteString(p.writer, strings.Repeat("\t", p.indent))
	io.WriteString(p.writer, s)
	io.WriteString(p.writer, "\n")
}

func (p *Printer) Dispatch(node parse.Node) {
	switch node := node.(type) {
	case *parse.TextNode:
		body := node.String()
		body = strings.Replace(body, "\n", "\\n", -1)
		body = strings.Replace(body, "\"", "\\\"", -1)
		p.Write(`io.WriteString(w, "` + body + `")`)
	case *parse.IfNode:
		action := p.ProcessArgs(node.Pipe.Cmds[0].Args, nil)
		p.Write(`if ` + action + " {")
		p.indent++
		for _, subnode := range node.List.Nodes {
			p.Dispatch(subnode)
		}
		if node.ElseList != nil {
			p.indent--
			p.Write("} else {")
			p.indent++
			for _, subnode := range node.ElseList.Nodes {
				p.Dispatch(subnode)
			}
		}
		p.indent--
		p.Write("}")
	case *parse.ActionNode:
		p.Write(`io.WriteString(w, string(` + p.ProcessCommands(node.Pipe.Cmds) + `))`)
	case *parse.RangeNode:
		action := p.ProcessArgs(node.Pipe.Cmds[0].Args, nil)
		p.Write(`if len(` + action + `) > 0 {`)
		p.indent++
		// TODO: change this var assign
		p.Write(`for _, arg := range ` + action + " {")
		p.indent++
		for _, subnode := range node.List.Nodes {
			p.Dispatch(subnode)
		}
		p.indent--
		p.Write(`}`)
		if node.ElseList != nil {
			p.indent--
			p.Write("} else {")
			p.indent++
			for _, subnode := range node.ElseList.Nodes {
				p.Dispatch(subnode)
			}
		}
		p.indent--
		p.Write("}")
	default:
		log.Fatalf("Unknown node: %T\n", node)
	}
}

func (p *Printer) ProcessCommands(cmds []*parse.CommandNode) string {
	result := ""
	for i, cmd := range cmds {
		if i == 0 {
			result = p.ProcessArgs(cmd.Args, nil)
		} else {
			result = p.ProcessArgs(cmd.Args, &result)
		}
	}
	return result
}

func (p *Printer) ProcessArgs(args []parse.Node, lastArg *string) string {
	argStrings := []string{}
	for _, arg := range args {
		argString := ""
		switch arg := arg.(type) {
		case *parse.FieldNode:
			argString = "arg" + arg.String()
		case *parse.IdentifierNode:
			argString = arg.String()
		default:
			fmt.Printf("Unknown: %T\n", arg)
		}
		argStrings = append(argStrings, argString)
	}
	if len(argStrings) == 1 {
		return argStrings[0]
	} else {
		return argStrings[0] + "(" + strings.Join(argStrings[1:len(argStrings)], ",") + ")"
	}
}

func dummy() {
}
