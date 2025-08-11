# Golang REST API
A REST api writte in Golang, capable of 
- navigating to either a /users or /apps endpoint 
- editing the lists of said users or apps 
- requires authentication for the users endpoint 
- allows you to choose which port to start your server on, throough the use of flags or through CLI 
- has a basic persistant storage configured 
- exposes a prometheus endpoint 
- contanins prometheus and grafana containers that managa visualisation of metrics.
    - to run them: docker-compose up
- contains a script that continuously runs request
    - to run ./background_requests
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
- curl http://localhost:8080/metrics | grep http_requests_total | sort
- curl http://localhost:8080/metrics | grep 'http_requests_total.*apps'
- curl http://localhost:8080/metrics | grep 'http_requests_total.*users'
