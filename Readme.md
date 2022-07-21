## Golang Exercise v21

Service provides API for the company with following attributes:
 - Name
 - Code
 - Country
 - Website
 - Phone

The service has environment parameters:

 - PORT Default: 8888
 - HOST Default: 0.0.0.0

DB parameters:
 - DSN  Default: postgres://app:pass@localhost:5433/goex21   

Basic Auth parameters:
 - USER  Default: user
 - PASS  Default: pass
 - COUNTRY Default: Cyprus

Message Broker parameters:
 - AMQP_URL: Default: amqp://quest:quest@localhost:5672/

Default port: 8888

Base path: /api/v1

### Http Request AddCompany

Add company to DB

POST http://url:PORT/api/v1/company/add

Auth: BasicAuth

GeoFilter: Yes

Format: JSON 

Example Request:

```azure
POST http://localhost:8888/api/v1/company/add
Content-Type: application/json

{
  "name": "atlassian",
  "code": "team",
  "country": "usa",
 "website": "www.atlassian.com",
  "phone": "33353333"
}
```

Response format: JSON
Example response:
```azure
{
"id": 1
"name": "atlassian",
"code": "team",
"country": "usa",
"website": "www.atlassian.com",
"phone": "33353333",
"isactive": true
}
```

### Http Request EditCompany

Edit company(-ies) in DB according queue params. 

POST http://url:PORT/api/v1/company/edit?param1=..

Auth: No

GeoFilter: No

Format: JSON

Queue params(1 or more params is mandatory):
 - name
 - code
 - country
 - website
 - phone

Example Request:

```azure
POST http://localhost:8888/api/v1/company/edit?name=atlassian
Content-Type: application/json

{
"website": "www.salattian.com"
}
```

Response format: JSON
Example response:
```azure
{
"id": 1
"name": "atlassian",
"code": "team",
"country": "usa",
"website": "www.salattian.com",
"phone": "33353333",
"isactive": true
}
```

### Http Request DeleteCompany

Delete company(-ies) in DB according queue params.

GET http://url:PORT/api/v1/company/delete?param1=..

Auth: Yes

GeoFilter: Yes

Queue params(1 or more params is mandatory):
- name
- code
- country
- website
- phone

Example Request:

```azure
GET http://localhost:8888/api/v1/company/delete?country=usa&phone=123
```

Response format: JSON
Example response:
```azure
{
"id": 1
"name": "atlassian",
"code": "team",
"country": "usa",
"website": "www.salattian.com",
"phone": "33353333",
"isactive": false
}
```

### Http Request SearchCompany

Search company(-ies) in DB according queue params.

GET http://url:PORT/api/v1/company/search?param1=..

Auth: No

GeoFilter: No

Queue params(or without parameters for full list):
- name
- code
- country
- website
- phone

Example Request:

```azure
GET http://localhost:8888/api/v1/company/search?country=usa
```

Response format: JSON
Example response:
```azure
[
  {
    "id": 2,
    "name": "netflix",
    "code": "nflx",
    "country": "usa",
    "website": "www.netflix.com",
    "phone": "1234567"
    "isactive": true
  },
  {
    "id": 7,
    "name": "atlassian",
    "code": "team",
    "country": "usa",
    "website": "www.atlassian.com",
    "phone": "33353333"
    "isactive": true
  }
]
```