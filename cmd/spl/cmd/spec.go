package cmd

import (
	"github.com/spf13/cobra"
)

const spec = `General

SPL ("Simple Programming Language") is a simple procedural programming language.
It contains a predefined primitive type for integers and a type constructor for
fields. SPL also uses logical values, but variables of this type cannot be
created, nor can there be literals of this type.

SPL knows procedures (but no functions), both with value parameters and with
reference parameters. Values of composite types must be passed as reference
parameters. Procedures can define local variables. There are no global
variables, nor are there nested procedure agreements.

You can use the conditional statement (single- and two-armed), the reject loop,
the assignment, the call of a procedure, and the compound statement for
statement sequences.

Expressions are constructed using integers. Six comparison operators, the four
basic arithmetic operations and the negation are available. Parentheses allow
any combination of partial expressions. Literals can be noted in different
representations for the predefined type "integer".

Fields are indexed by integer expressions, both on the left and right sides of
assignments or in argument expressions. Allowed index expressions return a value
in the range 0...(n-1) if the field has n elements. Indexing a field outside
this range causes a runtime error.

The runtime library provides procedures for entering and outputting whole
numbers and individual characters on the text screen. It also provides a
procedure for terminating a program immediately. Another procedure returns the
time elapsed since the program was started. On the graphics screen, individual
pixels, straight lines or circles can be drawn in any color. It is also possible
to fill the whole graphic screen with any color.

At the end of this language definition there is an example program in which some
of the possibilities of SPL are demonstrated.

Lexical conventions

A distinction is made between upper and lower case.

Spaces and horizontal tokens separate tokens, but have no other meaning for the
program. Line breaks also count as spaces.

A comment is introduced by a double slash ("//") and extends to the end of the
line in which it appears. It also separates tokens and has no meaning for the
program.
Example:

	// That's a comment.

Identifiers consist of any number of letters, digits and underscores. You must
start with a letter, where the underscore is considered as a letter.
Example:

	that_is_an_identifier

Numbers are formed by stringing decimal digits together. Alternatively, numbers
can be specified in hexadecimal. They begin with the prefix "0x" and contain
hexadecimal digits (0-9,a-f). The hexadecimal digits a-f can also be written in
capital letters. A third alternative for the representation of numbers is the
inclusion of exactly one character in apostrophes. The number shown is then the
ASCII code of the character between the apostrophe. The character string \n is
considered a character with the meaning "line break" (0x0a). Numbers must be in
the range 0..2^31-1; the unary operator "-" is available for displaying negative
numbers (see below). The interpretation of numbers that fall outside the
specified range is undefined.
Examples:

	1234
	0x1a2f3F4e
	'a'
	'\n'

SPL contains the following reserved words that cannot be used as identifiers:

	array else if of proc ref type var while

The following characters and character combinations have meaning:

( ) [ ] { } = # < <= > >= := :  , ; + - * /

All non-permitted characters or character strings are recognized and rejected.

The main program

An SPL program consists of a collection of type and procedure agreements without
a specific sequence. Identifiers for types must be agreed before they are used.
This restriction does not apply to procedure identifiers so that mutually
recursive procedures can be formulated.

At least one procedure agreement is required, namely that of the procedure with
the specified name "main". This procedure has no parameters and is activated
automatically when the program is started.

Types and type agreements

There is a predefined primitive type for integers ("int"). The identifier "int"
is not a reserved word. It is considered implicitly agreed before all user
declarations.

The data type constructor "array" constructs a field over a basic type. The
field size is defined statically at the time of translation and is part of the
type. The basic type can be any type.
Examples:

	array[3] of int
	array[3] of array[5] of int

A type agreement combines an identifier with a type. It has the following
structure:

	type <name> = <type> ;

Then <name> can be used as an abbreviation for <type>.
Examples:

	type myInt = int;
	type vector = array[5] of int;
	type matrix = array[3] of vector;
	type mat3 = array[10] of array[20] of vector;

Each type expression constructs a new type. Two types are the same if they were
constructed using the same type expression.
Example of the same types:

	type typ1 = array[5] of int;
	type type type2 = type1;

Example of different types:

	type typ1 = array[5] of int;
	type type type2 = array[5] of int;

Declaratives

A procedure agreement combines an identifier with a procedure. It has the
following form:

	proc <name> ( <parameter list> ) { <declarations> <statement list> }

The optional <parameter list> names the formal parameters of the Prodedur. Each
parameter is specified in one of two forms:

	<name> : <type>

or

	ref <name> : <type>

The first form denotes a value parameter, the second form a reference parameter.
The individual parameters are separated by commas in the list. The names of the
parameters are valid until the end of the procedure. The optional <declarations>
define local variables. Each declaration has the following form:

	var <name> : <type> ;

These names are also valid until the end of the procedure. You must not collide
with a parameter name. It is possible to declare parameters or local variables
with a name that already has a different meaning further outside, i.e. is a type
or procedure name. This external meaning is hidden by the local meaning. The
optional <statement list> consists of any number of statements.
Examples:

	proc nothing() {}
	proc copy(i: int, ref j: int) { j := i; }
	proc swap(ref i: int, ref j: int) {
	    var k: int;
	    k := i; i := j; j := k;
	}

Instructions

Statements are used to achieve an effect (change state, side effect). There are
six different instructions. The empty statement consists of only one semicolon
and does nothing.
The assignment has the form:

	<lhs> := <rhs> ;

where <lhs> must be a variable of type int (indexed field variables are also
allowed) and <rhs> an expression of the same type. At runtime, the right side is
evaluated and assigned to the left side. Any expressions on the left side are
evaluated before any expressions on the right side.
Example:

	x[3] := x[2] * 2;

The conditional statement can come in one of two forms:

	if ( <expression> ) <statement1>
	if ( <expression> ) <statement1> else <statement2>

where <expression> is an expression that returns a logical value and
<statement1> or <statement2> is a statement. The statement behind an "else"
belongs to the innermost "if", to which no "else" has yet been assigned. At
runtime <expression> is evaluated. If it returns "true", the <statement1> is
executed. If it returns "wrong", nothing else is done in the first form; in the
second form, <statement2> is executed in this case. In any case, execution is
continued with the statement following the "if" statement.
Examples:

	if (x < 0) x := 42;
	if (x < 0) x := 42; else x := 43;

The reject loop is formulated as follows:

	while ( <expression> ) <statement>

where <expression> is an expression that returns a logical value and <statement>
is an instruction. At runtime <expression> is evaluated. If it returns "true",
it is executed and the loop is executed again. If it returns "wrong", the
statement following the "while" statement is continued.
Example:

	while (x < 10) x := x + 1;

To make a statement sequence appear syntactically as a single statement, there
is the compound statement. It consists of any number of statements (possibly not
even a single one) that are enclosed in curly brackets.
Example:

  { k := i; i = j; j := k; }

A procedure call is achieved by the following form:
  <name> ( <argument list> ) ;
The optional <argument list> is a comma-separated list of expressions whose
number and types must match the number of parameters and their types in the
procedure agreement. At runtime, the argument expressions are evaluated from
left to right and then the procedure is activated. Only variables (simple
variables or field variables) can appear in the argument list for reference
parameters, whereas any expressions are permitted for value parameters. Fields
must be passed as reference parameters. After the procedure has been executed,
the execution returns after the procedure call.
Example:

	swap(n, m);
	sum(3 * x + 5, 9 - y / 2, z);

Expressions

Expressions are used to calculate values. The values are either integer values
or truth values.

The four basic arithmetic operations with the usual precedents and
associativities are available for calculating integer values. The unary minus is
also available; it binds more strongly than multiplication.

The six relational operators are <, <=, >, >=, = and # (not equal). They compare
the values of two integer expressions and return a logical value. They bind
weaker than the addition. There are no operators available for combining Boolean
values.

At any point of an expression, parentheses can be used to override the built-in
precedents and associativities.

The operands of an expression are other expressions, and ultimately literals or
variables. The latter are either simple variables or indexed field variables. A
field is indexed by an expression with an integer value in square brackets after
the field variable. Operations with one field as a whole are not allowed (except
passing the whole field as a reference).

At runtime, for all expressions with two operands, first the left and then the
right operand are evaluated; the operation is then executed.
Examples:

	1
	x
	3 + x * y
	(3 + x) * y
	5 * -a[n-2]
	i < n
	b - 2 # a + 3

Library procedures

Displays the value of i on the text screen:

	printi(i: int)

Outputs the character with the ASCII code i on the text screen:

	printc(i: int)

Reads an integer from the keyboard and stores it in i. The input is buffered
line by line with echo:

	readi(ref i: int)

Reads a character from the keyboard and saves its ASCII code in i. Input is
unbuffered and without echo:

	readc(ref i: int)

Ends the running program and does not return to the caller:

	exit()

Returns in i the time in seconds since the program was started:

	time(ref i: int)

Clears the graphics screen with the color. Colors are formed by specifying the
R, G and B components according to the pattern 0x00RRRGGBB. The values 0...255
are therefore available for each component:

	clearAll(color: int)

Sets the pixel with coordinates x and y to color. Limits: 0 <= x < 640, 0 <= y <
480:

	setPixel(x: int, y: int, color: int)

Draws a straight line from (x1|y1) to (x2|y2) with the color. Limits like
setPixel:

	drawLine(x1: int, y1: int, x2: int, y2: int, color: int)

Draws a circle around the center (x0|y0) with radius radius and color:

	drawCircle(x0: int, y0: int, radius: int, color: int)
`

// specCmd represents the spec command.
var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "Print the spl language specification",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Print(spec)
	},
}

func init() {
	rootCmd.AddCommand(specCmd)
}
