# Base Interaction
- What if we let users of Hype interact with Go?
- We let Go do some of the work for us
- This would make Hype a partner for Go
- Use Go src with Hype src. Or Hype to start Go
- This means our types are closely mingling
- But we can also standalone
"""
import go (
    "fmt"
)

import hyp (
    "count"
)

count(fmt.Sprintf(%s\n, "Hello"))
"""

- Basically conditional hoisting
    - Go through all func and var defs and define them
    - If a statement moves a var to a new env then hoist it
    - But if it does not, dont

### Par-For
- Parallel for loop
- Every 'iteration' that would be, gets a goroutine

### Lazy Auto-Par
- Runs things in parallel only if they can be
- Otherwise runs them in single main thread