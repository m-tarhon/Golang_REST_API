#!/bin/bash

go run ./server -http=8080 -https=8081 &
sleep 5

while true; do
    # Send a request to the /users endpoint
    curl -X GET -u boss:man http://localhost:8080/users
    curl -X POST -u boss:man http://localhost:8080/users -H "Content-Type: application/json" -d '{"name":"user1","age":25}'
    curl -X POST -u boss:man http://localhost:8080/users -H "Content-Type: application/json" -d '{"name":"user2","age":30}'
    curl -X GET -u boss:man http://localhost:8080/users/user1
    curl -X GET -u boss:man http://localhost:8080/users
    curl -X DELETE -u boss:man http://localhost:8080/users -H "Content-Type: application/json" -d '{"name":"user1","age":25}'
    curl -X GET -u boss:man http://localhost:8080/users
    
    # Send request to the /apps endpoint
    curl -X GET http://localhost:8080/apps
    curl -X POST http://localhost:8080/apps -H "Content-Type: application/json" -d '{"name":"app1","born":2000,"price":10.99}'
    curl -X POST http://localhost:8080/apps -H "Content-Type: application/json" -d '{"name":"app2","born":2010,"price":15.49}'
    curl -X GET http://localhost:8080/apps/app1
    curl -X GET http://localhost:8080/apps
    curl -X DELETE http://localhost:8080/apps -H "Content-Type: application/json" -d '{"name":"app1","born":2000,"price":10.99}'
    curl -X GET http://localhost:8080/apps

    sleep 15

done
