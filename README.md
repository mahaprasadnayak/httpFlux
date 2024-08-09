# httpFlux
Streamline your web traffic with httpFlux, a dynamic and efficient HTTP load balancing solution.

# Description
-> Balances web traffic dynamically across multiple servers, improving response times and handling capacity.•
-> Monitors backend servers in real-time, automatically routing traffic away from unhealthy servers to maintain high availability.•
-> Offers customizable load-balancing strategies, such as round-robin and weighted round-robin, to match specific traffic and server
requirements.

# run cmds
go  run .\proxy\proxy.go 
go run server.go -p 8081
go run server.go -p 8082
go run server.go -p 8083
cUrl http://localhost:8080 



