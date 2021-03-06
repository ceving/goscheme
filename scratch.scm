(+)            ;; -> 0
(+ 1)          ;; -> 1
(+ 1 2)        ;; -> 3
(+ 1 2 3)      ;; -> 6
(+ 1 0.2)      ;; -> 1.2
(+ "a")        ;; -> error

(apply + (list 1 2))

(apropos "apply")

(display (if #f 1))
(define unspecified (if #f 1))

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
(display (begin))

display
lambda
format
define

(define (a b)
  (display b))
a

(length 5)
(length (cons 1 2))
(length (list 1 2))
(length)
(length (list))
(length (list 1) (list 1 2))

(list->vector (list 1 2 3))

(lambda () 1)
(pair? '())
(null? '())
(list? '())

(map (lambda (x)
       1)
     (list))

(define hello "hello")
((lambda (str)
   (string-set! str 0 #\H)) hello)
hello
((lambda (str)
   (string-set! str 0 #\H)) ((lambda () "hello")))
(string #\h #\e #\l #\l #\o)

(list #\null(list 1))

(char->integer #\x7F)

(with-input-from-string ""
  (lambda ()
    (read)))

'(1 . 2 3)

(number? 1/2)
(number? +1)
(number? -1)


(define the-continuation #f)

(define call/cc call-with-current-continuation)

(define (test)
  (let ((i 0))
                                        ; call/cc calls its first function argument, passing
                                        ; a continuation variable representing this point in
                                        ; the program as the argument to that function.
                                        ;
                                        ; In this case, the function argument assigns that
                                        ; continuation to the variable the-continuation.
                                        ;
    (call/cc (lambda (k) (set! test k)))
                                        ;
                                        ; The next time the-continuation is called, we start here.
    (set! i (+ i 1))
    i))

(test)
(the-continuation)

(eq? 2 (+ 1 1))


(let* ((yin
        ((lambda (cc) (display #\@) cc) (call/cc (lambda (c) c))))
       (yang
        ((lambda (cc) (display #\*) cc) (call/cc (lambda (c) c)))))
  (yin yang))