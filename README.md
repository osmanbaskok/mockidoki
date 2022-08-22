# Mockidoki
A mock API for testing purposes

## Description
Mockidoki is an API providing the following functionalities: 
* Returning pre-defined HTTP responses for specific HTTP requests.
    - It matches the provided request to the corresponding response from the database using regular expressions.
    - Supports all the HTTP methods (GET, POST, DELETE, PUT, PATCH, etc.).
    - Uses request's url, header, and body for matching.
    - Allows customizing HTTP status, body, and header for the response.
* Acting like a message publisher (Currently only for Kafka)
    - Receives a json payload via HTTP request and publishes it onto the specified Kafka topic as a Kafka message 

## Getting Started
Run docker compose and get Postgres DB up
* docker-compose up -d

Run the following command in the root directory
* go run main.go

#### HTTP mock example:
* Run the following Curl script 
````  
curl -v --location --request POST 'http://localhost:8080/http-mocks/warehouses/3/stock-transactions' \
    --header 'x-clientid: mockidoki-user' \
    --header 'Content-Type: application/json' \
    --data-raw '{
    "productId": 12345,
    "amount": 45.37,
    "type": "SALES"
    }'
````

* Response to the Curl script should be as the following (returning 201, and the specified header values)

```` 
* upload completely sent off: 59 out of 59 bytes
  < HTTP/1.1 201 Created
  < Content-Type: application/json
  < Test_header1: test_perfect1
  < Test_header2: test_perfect2
  < Date: Sun, 21 Aug 2022 20:16:04 GMT
  < Content-Length: 0
```` 

PS: There are three pre-defined sample HTTP mock records in the http_mock table. You can check and play with them as you wish.

#### Event mock example:
* Run the following Curl script
````  
curl --location --request POST 'http://localhost:8080/event-mocks/stock-created/process' \
--data-raw '{
	"productId": 12345,
	"price": 49.90,
	"name": "Mockidoki Tshirt"
}'
````  
* Run the following command in the terminal to receive the message sent by the curl script given above (make sure that you first have kafka & zookeper installed on your machine to run the following command)
````  
kafka-console-consumer --bootstrap-server localhost:29092 --topic my-mockidoki-topic --from-beginning
````  

PS: You can also post an array of messages at once to the given API endpoint (http://localhost:8080/event-mocks/{key}}/process-list). They will still be published as a single message on the specified topic.

### Dependencies

* Postgres DB
* Kafka (optional, when event mocks needed)

### License
[GNU General Public License v3.0](LICENSE)