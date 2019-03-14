// Package scanner implements a buffered scanner which provides lexical analysis
// (tokenizing) of SPL source code. A scanner takes a bufio.Reader as source
// which can then be tokenized through repeated calls to the Scan() method.
package scanner
