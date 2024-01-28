# smack

A lisp dialect, inspired by Clojure, hosted in Go.

- Interpreted / REPL development/scripting flow
- Transpiles to Go
- Native Go concurrency primitives
- easy Go interop


# TODO
- Macros
- Transpile to Go
- Golang standard lib interop
- Golang channel type in Smack
- File/IO
- Option/Result types
- Compound Data Structures (struct)
- Traits/Behaviors/Protocols
- Namespaces/Modules
- AST Optimization
- Ref (mutable) types
- Basic Standard Lib (Currently not fully implemented)



# Tentative Syntax Example

```lisp

(def add (fn (a b) (+ a b)))
(def sub (fn (a b) (- a b)))
(def not (fn (a) if a false true))

(let (x 54 y 46) (add x y))

(def a 100)
(def b 200)

(println (do (
    (let 
        (res (sub a b)) 
            >= res 50)
)))

(def table {:a 12 :b "bee"})

(def array [1 2 3])

(println "Ayyye this really something!!!")

```
