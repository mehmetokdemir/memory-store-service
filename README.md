# Memory Store Service
## _A REST-API that works as an in memory key-value store service._

Written in Golang using standard libraries

##Features
 - Post values with a given key
 - Get the value of that key with a specific key
 - Flush all values in the json file and state
 - Every minute should save the data on the memory to the file
 - When the application stops and starts the program again, if there is a file(tmp/TIMESTAMP-data.json) in place, it should load the updated data into the memory
 - Server log file(server.log) that shows http requests with timing
 - API documentation made with swagger

##Installation And Run 
git clone https://github.com/mehmetokdemir/memory-store-service.git

GO run 
   ```sh
cd memory-store-service
PORT=8080 go run .
```

DOCKER run
   ```sh
cd memory-store-service
docker build -t memory-app .
docker run -p 8080:8080 --name memory-store-srv memory-app 
```


##Endpoints
Base path = /memory
- Read value with a key
    > curl -X 'GET' \
  'http://localhost:8080/memory?key=foo' \
  -H 'accept: application/json'
- Post value with a given key
    >curl -X 'POST' \
  'http://localhost:8080/memory' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "key": "foo",
  "value": "bar"
  }'
- Flush data
    > curl -X 'DELETE' \
  'http://localhost:8080/memory' \
  -H 'accept: application/json'
  

##Api Documentation
    http://localhost:8080/docs