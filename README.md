# Description

A shopping CRUD application that integrates REST-API with gRPC, build with microservice architechture.

# Tech Stacks
- Go 
- Gorm
- Echo 
- gRPC 
- REST
- JWT Auth
- MongoDB
- PostgreSQL
- Swagger
- Docker

# Flow

The application is divided into 2 services.

- Service Shopping runs on http://localhost:8080/
    - /transactions (POST) - Create a new transactions
    - /transactions (GET) - Retrieve all transactions
    - /transactions/{id} (GET) - Retrieve a specific transactions by ID
    - /transactions/{id} (PUT) - Update a transaction informations
    - /transactions/{id} (DELETE) - Delete a transaction informations
    - /products (POST) - Create a new product
    - /products (GET) - Retrieve all products
    - /products/{id} (GET) - Retrieve a specific product by ID
    - /products/{id} (PUT) - Update a product informations
    - /products/{id} (DELETE) - Delete a product informations

- Service Payment runs on http://localhost:50051/
    - /payments (POST) - Create a new payment

# How to Use

1. Clone this repository.

2. Create an `.env` based on the `.env.example`

3. On the root folder of this repository type the following command:

```bash
docker compose up
```

4. Once it has finished, the 2 microservices will run together and you can access the API documentation through http://localhost:8080/swagger/index.html