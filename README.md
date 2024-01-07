# smack
- Lisp dialect with minimal parenthesis. 
- Curley braces denoting the begin and end of an expression block.

# Potential Syntax

```awk

let x 20;
let y 50;

fn add_or_increase_by_70 {
	(a) => { 
	    let scalar { x + y };
	    let result { scalar + a };
	    result
	}
	(a, b) => {
	    a + b
	} 
};

fn add {
    (a, b) => {
	a + b
    }
};

println { 
	add_or-increase_by_70 { x + y } { y * y } 
};

println { add 20 50 };


 
```
