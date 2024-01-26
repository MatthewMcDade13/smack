(println "Ayyye this really something!!!")

(def add (fn (a b) (+ a b)))

(def sub (fn (a b) (- a b)))
(def not (fn (a) if a false true))

(let (x 54 y 46) (add x y))

(def a 100)
(def b 200)

(println (eval (list + a b)))




