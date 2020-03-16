package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"github.com/k0kubun/pp"
)

type nippoRenderer struct {
	level  int
	format string
}

func main() {
	md, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic("Read error")
	}
	extensions := parser.CommonExtensions | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	node := markdown.Parse(md, p)
	pp.Print(node)
	textBytes := markdown.Render(node, &nippoRenderer{0, "%s"})
	fmt.Println(fmt.Sprintf("%s", textBytes))
}

func (r *nippoRenderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	switch node := node.(type) {
	case *ast.Text:
		r.text(w, node)
	case *ast.Paragraph:
		if entering == false {
			w.Write([]byte("\n"))
		}
	case *ast.List:
		fmt.Println(entering)
		if entering {
			r.level++
		} else {
			r.level--
		}
		r.format = "・%s"
	case *ast.Heading:
		r.level = 0
		r.format = "【%s】"
		if entering == false {
			w.Write([]byte("\n"))
		}
	case *ast.Link:
		r.format = "%s"
	}
	return ast.GoToNext
}
func (r *nippoRenderer) text(w io.Writer, text *ast.Text) {
	if len(text.Literal) <= 0 {
		return
	}
	if r.level > 0 {
		w.Write([]byte(strings.Repeat("  ", r.level-1)))
	}
	w.Write([]byte(fmt.Sprintf(r.format, text.Literal)))
}

func (r *nippoRenderer) RenderHeader(w io.Writer, ast ast.Node) {
	w.Write(ast.AsContainer().Content)
}
func (r *nippoRenderer) RenderFooter(w io.Writer, ast ast.Node) {
	w.Write(ast.AsContainer().Content)
}
