(+)            ;; -> 0
(+ 1)          ;; -> 1
(+ 1 2)        ;; -> 3
(+ 1 2 3)      ;; -> 6
(+ 1 0.2)      ;; -> 1.2
(+ "a")        ;; -> error

(apply + (list 1 2))

(apropos "apply")

(display (if #f 1))

(car 5)

labels

(/ 0.0 0.0)

(+ undefined)

(eval (cons cons (cons 1 (cons 2 '()))) (interaction-environment))
(eval (cons cons (cons 1 2)) (interaction-environment))

(define a 1)
(eval (define a 2) (interaction-environment))


(define b 1)
(define x (lambda ()
	    (+ b 1)))
(define b #f)
(x)


(define c 1)
c
(begin 
  (define c 2))
((lambda ()
   (define c 2)))
(begin 
  (define c 2)
  (set! c 3))
c
((lambda ()
   (define c 1)
   (set! c 3)))
((lambda ()
   (set! c 3)))
c

(set! asdfasdf 5)

;; Internal definition => error
(let ((a 1))
  (define (f x)
    (define b (+ a x))
    (define a 5)
    (+ a b))
  (f 10))

(eq? (list) '())

(begin)

display
lambda
format

(define (a b)
  (display b))
a