docker compose -f deployments/docker-compose.yml up --build запускать в корне проекта 

**Project Overview**

Simple HTTP service for managing key-value pairs with Tarantool as the storage backend.

**Features**
1. CRUD operations for key-value pairs
2. HTTP REST API interface
3. Tarantool storage backend

**Quick Start**
1. Clone the repository
`git clone https://github.com/ItserX/kvStorage.git`
2.  Start the services  
`cd kvStorage`  
`docker compose -f deployments/docker-compose.yml up --build`
