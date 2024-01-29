# Deployment

MongoDB  is deployed using Mongo Atlas. The URI can be found on `.env` inside the Payment and Shopping folder.

## Google Cloud Run Step by Step

1. For each service (`payment` and `shopping`), run: 

```bash
docker build --platform=linux/amd64 -t SERVICE-NAME .
```

2. Tag the image:

```bash
docker tag APP-NAME gcr.io/GCP-PROJECT-ID/APP-NAME
```

3. Push the image to Google Cloud Run:

```bash
docker push gcr.io/GCP-PROJECT-ID/APP-NAME
```

4. Create a new service on https://console.cloud.google.com/run

5. Select the container that you've uploaded on step 3.

6. Configuration setup (Port, Access, etc.).

7. Service has been deployed successfully!

## Using the Deployed App

### Shopping Service

https://graded1shopping-7cijv37eka-et.a.run.app

> The base URL will automatically redirect you to the API documentation (Swagger)

### Payment Service

https://graded1payment-7cijv37eka-et.a.run.app

> The base URL won't show anything since it's a gRPC service for the Shopping app

## Using the app locally with Docker

On `./shopping/.env`, change the variable for GRPC_URI into:

```bash
GRPC_URI=host.docker.internal:50051
```

On `./shopping/config/grpc.go` replace the existing InitGrpc function with this instead:

```go
func InitGrpc() (*grpc.ClientConn, pb.PaymentClient) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(os.Getenv("GRPC_URI"), opts...)
	if err != nil {
		log.Fatal(err)
	}
	return conn, pb.NewPaymentClient(conn)
}
```

> This function removes TLS credentials from the gRPC client when running __locally__. Otherwise, the app will occur an invalid TLS handshake.

On this repository root path, run:

```bash
docker-compose up
```

When image building has finished, go to http://localhost:8080/swagger/index.html to see the API documentation.

## Screenshots

### Docker Compose Logs (For local running)

![docker compose](docker%20compose.png)

### Google Cloud Run Dashboard

![GCR Dashboard](cloud%20run%20dashboard.png)

![Shopping Dashboard](shopping%20dashboard.png)

![Payment Dashboard](payment%20dashboard.png)