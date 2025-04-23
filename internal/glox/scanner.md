## Scanner Notes

#### Components
- Keywords of reserved control words mapped from str to const
- Error capture
- Keep track of position with start, current, and per line
- List of tokens after things have been broken up

#### Scanning tokens
- Append EOF token to tell later the parser where to stop
- Look at each character until we have complete token, then add it to the tokens array
- We look at chars until we have reached the end of the file
- Iterate line upon finding a /n token, also set end token
    - Wab inline stments? {doSome()}
    - Have expected chars })