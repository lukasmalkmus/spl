package repl

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/lukasmalkmus/spl/internal/app/spl/parser"
)

const prompt = ">> "

// Start the Read Evaluate Print Loop.
func Start(in io.Reader, out io.Writer) error {
	printPrompt(out)
	s := bufio.NewScanner(in)
	for s.Scan() {
		text := string(bytes.TrimSpace(s.Bytes()))
		if text == "" {
			printPrompt(out)
			continue
		}

		stmt, err := parser.ParseStatement(text)
		if err != nil {
			parser.PrintError(out, err)
		}
		if b, err := json.MarshalIndent(stmt, "", "    "); err == nil {
			_, _ = fmt.Fprintf(out, "%s\n", b)
		} else {
			_, _ = fmt.Fprintf(out, "%+#v\n", stmt)
		}
		printPrompt(out)
	}
	return s.Err()
}

func printPrompt(w io.Writer) {
	_, _ = fmt.Fprint(w, prompt)
}
