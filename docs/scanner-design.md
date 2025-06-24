# Scanner Design
- Takes raw src as string of chars
- Produces list of tokens
- Tokens posess any fields the parser may need. Such as Lexeme, Literal, Type Const etc

### Main Procs
- peek
- peekNext
- advance
- addToken

### Syntax
- Should the Scanner handle syntax logic?
    - Yes, this makes it easier for the parser
    - No, it has 1 job. Syntax falls on the Parser
- So, should the parser ignore subsequent newlines?
- Or, just only do as its told?

### Matching
- Switch or map?
- Map where each char is assigned a func
    - This makes logic less linear and harder to maintain
- Switch/Logical
    - Linear, straightforward. Most compiled languages optimize switches anyway