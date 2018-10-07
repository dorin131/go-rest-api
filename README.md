# This is a demo Go RESTful API for #guild-golang

Usage:

`curl -X "GET" http://localhost:8080/campaigns`

`curl -X "GET" http://localhost:8080/campaign/1`

`curl -X "DELETE" http://localhost:8080/campaign/1`

`curl -X "POST" --data '{
    "name": "Campaign 777",
    "comapany": "SuperAwesome",
    "io": "1009",
    "house": false
  }' http://localhost:8080/campaign`
