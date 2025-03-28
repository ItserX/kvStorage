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
Or use a ready-made server:  
http://217.198.5.83/

**API Documentation**  
Create Key-Value Pair  
`POST /kv {"key": "key1", "value": {"v1":1, "v2": true, "v3": [1,2,3,4,5]}`  
Get Value by Key  
`GET /kv/{id}`  
Update Value by key  
`PUT /kv/{id} {"value": {"new_value": 1}}`  
Delete Key  
`DELETE /kv/{id}`  

**Run Tests**  
```bash
go test -cover ./internal/handlers/ ./internal/storage/
```
