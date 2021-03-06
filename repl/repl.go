package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/SpaceHexagon/ecs/object"

	"github.com/SpaceHexagon/ecs/evaluator"
	"github.com/SpaceHexagon/ecs/lexer"
	"github.com/SpaceHexagon/ecs/parser"
)

const MONKEY_FACE = `
				  	__,__
	   		  .--. .\\-"   "-//. .--.
  			 / .. \  \\\\ ////  / .. \ -/\/\/\--/\/--
  -/\/\/\--/\--              | | '| /    Y    \ |' | |
  			 | \ \ | -O- | -O- |/ /  |
  			 \ '- ,\.____|____-./, -'/
   			  ''-' /_  ^___^ _\ '-''    --/\/\/--/\/\---
     -/\/\/\--/\/--	        |\./|||||\/ |
       				\  \||||| / /   -/\/\/\--/\/--
        			'._' -=-' _.'
          			   '-----'
`

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "You have invoked the wrath of\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		evaluated := evaluator.Eval(program, env, nil)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
