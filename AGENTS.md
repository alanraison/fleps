You should strive to write idiomatic go code with good test coverage.
Whenever an error is encountered, it should be wrapped with a description of what was happening in 
the function that caught it, before returning it up the stack.