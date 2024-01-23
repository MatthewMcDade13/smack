# smack
- Lisp dialect with optional parenthesis. 
- Curley braces denoting the begin and end of an expression block to disambiguate when no parenthesis.

# Potential Syntax

```awk

let x 20
let y 50

fn add-or-increase-by-70 {
	(a) => { 
	    let scalar { + x y }
	    let result { + scalar a }
	    result
	}
	(a, b) => { + a b } 
}

fn add {
    (a, b) => { + a b }
};

println { add-or-increase-by-70 { + x y } { * y y }  }

println { add-or-increase-by-70 30 }

println { add 20 50 }


 
```


# Sources
[Make a Lisp (MaL)](https://github.com/kanaka/mal/blob/master/process/guide.md)
