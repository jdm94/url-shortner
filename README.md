# url-shortener
Rest API written in GO that support url shortening

### Steps To Run 
  * Build dockerfile using following command from working dir - docker build -t url-shortener:1.0 .
  * Run Docker-compose using following command from working dir - docker-compose up
  * Test Case #1 : curl --location 'http://localhost:8081/generate-short-url' \
--header 'Content-Type: application/json' \
--data '{
    "url" : "https://www.baeldung.com/curl-rest"
}'
* You can change the put any url in the above json format or alternatively use postman to test

