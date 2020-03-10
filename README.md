# Handler step

## readme

- purpose
- sequence diagram
- how to run
- how to test
- how to deploy

## Todo

- logging
- tracing
- metrics
- graceful shutdown
- liveness handler
- readyness handler
- version handler
- testable
- documents
- circuit breaker

## recommendations

- don't call panic
- log.Fatal only in main.go
- goroutine only in main.go
- defer recovery everytime in goroutine
- global variables beware about race conditions

## main

- read config
- initial framework by config
- initial database by config
- initial http.client by config
- initial logger by config
- initial tracer
- graceful shutdown
- liveness handler
- readyness handler
- version handler

## logging requirement

- X-Request-ID
- singerton and non blocking test

## framework handler

- get request
- call handler then get response & error
- return response to caller

## controller (business logic)

- call dependencies then get response & error
- log all error with tracing-ID, crime-scence, course
- ?rollback

## dependencies

- ?circuit breaker
- ?rollback

