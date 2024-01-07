# smack
- Lisp dialect with minimal parenthesis. 
- Curley braces denoting the begin and end of an expression block.

# Potential Syntax

```awk

def x 20;
def y 50;

defn add_or_increase_by_70 {
	(a) => { 
	    def result { x + y };
	    result
	},
	(a, b) => {
	    a + b
	}, 
};

println { 
	add_or-increase_by_70 { x + y } { y * y } 
};

 
```
