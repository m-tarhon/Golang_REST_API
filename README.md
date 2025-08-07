# Golang REST API
A REST api writte in Golang, capable of 
- navigating to either a /users or /apps endpoint 
- editing the lists of said users or apps 
- requires authentication for the users endpoint 
- allows you to choose which port to start your server on, throough the use of flags or through CLI 
- has a basic persistant storage configured 

## To run, do:
Open two terminals
- first run: go run ./server
    - can use -port=8080 to start server on port 8080
    - without it, the CLI will ask you to choose
- then, in the second terminal window: go run ./client 
    - can use the same flag, or the CLI will ask you to specify
    - if you wish to access the /users endpoint, the credentials are boss:man
- type help to navigate the client side 
- hve fun! please try breaking this so i can improve it 

## Testing out the prometheus endpoint
Useful commands
- curl http://localhost:8080/metrics | grep http_requests_total
- curl -X POST -u boss:man http://localhost:8080/users -H "Content-Type: application/json" -d '{"name":"alyce","age":25}'
- curl -u boss:man http://localhost:8080/users/alyce
- curl -X DELETE -u boss:man http://localhost:8080/users -H "Content-Type: application/json" -d '{"name":"alyce","age":25}'
- curl -u boss:man http://localhost:8080/apps
- curl -u boss:man http://localhost:8080/apps/Quran
- curl -X POST -u boss:man http://localhost:8080/users -H "Content-Type: application/json" -d '{"name":"bible", "born":1666,"price":76.2 }'
- curl -X POST -u boss:man http://localhost:8080/apps -H "Content-Type: application/json" -d '{"name":"Hazel", "born":1666,"price":76.2 }'
- curl http://localhost:8080/metrics | grep http_requests_total | sort
- curl http://localhost:8080/metrics | grep 'http_requests_total.*apps'
- curl http://localhost:8080/metrics | grep 'http_requests_total.*users'
