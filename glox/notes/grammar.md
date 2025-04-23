# Grammar

### Non Terminals
- Symbols that can be further expanded into other symbols

### Terminals
- Actual tokens and symbols used in the langauge 
- Cannot be broken down further

### Production Rules
- Describe how a non terminal can be futher expanded upon

### Notation Operators
- *, Indicated 0 or more repitions of the preceding element
    When a * is used at the end of an element or grouping, 0 or more of that element can occur to the right "x or y or z or p"
    for logic_or -> logic_and ( "or" logic_and )*
    logic_and can exist as a single expression "A" or a string "A or B or C"
- () for grouping expressions

### Right now
statement -> exprStmt | printStmt | ifStmt | whileStmt | whileStmt | block ; 

whileStmt -> "while" "(" expression ")" statement ;

unary -> ("!" | "-" ) unary | call ;

call -> primary ( "(" arguments? ")" )* ;
    - An infinent number of () can be prepended because functions can return callbacks which can be called

arguments -> expression ( "," expression )* ;
    - Because the * precedes a grouping, the grouping can occur 0+ times
    - Meaning an arguments can be just one expression

forStmt -> "for" "(" ( varDecl | exprStmt | ";" ) expression? ";" expression? ")" statement ;
    - The last 2 expressions are optional as shown by the preceding "?", a statement is then run each iter
    - 3 clauses, initializer, executed once
        - Next is condition, when to exit the loop, evaluated at each iter, decided to exit or continue
        - increment, at the end of each loop this expression is done. The value is discarded, so it must have a side effect
    - A for loop can be done in the syntax we already have, but this is *syntactic sugar*, making it easier to use

### Precedence
- A peice of syntax with higher precedence is evaluated before those of lower precedence
- Operators of the same precedence are evaluated from left to right
- Parentheses obviously change the scope of the precedence

### Arity
- The number of args a function or operation expects

### Native Funcs
- Functions impl in the host lang (Go) but exposed to the user, kinda like pre-std
- They count as the imple runtime

### Func Decl 
- Functions are literally a named statement, like a var is a named expr

declaration -> funcDecl | varDecl | statement ;

funDecl -> "fun" function ;
function -> IDENTIFIER "(" parameters? ")" block ;
parameters -> IDENTIFIER ( "," IDENTIFIER )* ; 
    // Can happen 0+ times