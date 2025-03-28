**Project Overview**

Simple HTTP service for managing key-value pairs with Tarantool as the storage backend.

**Features**
1. CRUD operations for key-value pairs
2. HTTP REST API interface
3. Tarantool storage backend

**Quick Start**
1. Clone the repository
```bash
git clone https://github.com/ItserX/kvStorage.git
```  
2.Start the services
```bash  
cd kvStorage  
docker compose -f deployments/docker-compose.yml up --build
```
**Production Server**:  
http://217.198.5.83/

**Run Tests**  
```bash
$ go test -cover ./internal/handlers/ ./internal/storage/  
ok      kvManager/internal/handlers     0.006s          coverage: 64.1% of statements  
ok      kvManager/internal/storage      0.013s          coverage: 92.3% of statements  
```
**API Documentation**  
Create Key-Value Pair  
`POST /kv body: {"key": "key1", "value": {"v1":1, "v2": true, "v3": [1,2,3,4,5]}}`  

Get Value by Key  
`GET /kv/{id}`  

Update Value by key  
`PUT /kv/{id} body: {"value": {"new_value": 1}}`  

Delete Key  
`DELETE /kv/{id}`  

**Configuration**
```ini
APP_PORT=:8080                    #HTTP server port  
TARANTOOL_ADDRESS=tarantool:3301  #DB host:port
TARANTOOL_USER=guest              #Authentication user
```
