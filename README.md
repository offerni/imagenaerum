# Imagenaerum 
## Multithreaded image processing using RabbitMQ

This service processes images with various transformations such as blur, crop, resize, grayscale, and invert. It uses RabbitMQ for message queuing and Docker Compose for container orchestration.

**Disclaimer This is still a work in progress and there are some validations and error handling missing.**

## Accepted Parameters

The `ImageProcess` endpoint accepts the following parameters:

- **files**: Multipart file upload.
- **blur**: A string representing the blur intensity.
- **crop_anchor**: A string representing the crop dimensions in the format `width,height`.
- **resize**: A string representing the resize dimensions in the format `width,height`.
- **grayscale**: A string (any non-empty value will apply grayscale).
- **invert**: A string (any non-empty value will apply invert).
- ... More processes soon

## Running RabbitMQ with Docker Compose

To run the service using Docker Compose, follow these steps:

1. Ensure Docker and Docker Compose are installed on your machine.
2. Clone the repository:
    ```sh
    git clone https://github.com/offerni/imagenaerum.git
    cd imagenaerum
    ```
3. Start the RabbitMQ container:
    ```sh
    docker-compose up --build
    ```

## Starting the Worker and Consumer Servers

To start the worker and consumer servers, run the following commands in separate terminals:

1. Start the worker server:
    ```sh
    cd worker
    cp .env.default .env
    go run cmd/main.go
    ```

2. Start the consumer server:
    ```sh
    cd consumer
    cp .env.default .env
    go run cmd/main.go
    ```

## Available Endpoints

### Image Processing Endpoint

- **POST /image_process**
  - **Description**: Processes images with the specified transformations.
  - **Parameters**:
    - `blur`: Optional. Blur intensity.
    - `crop_anchor`: Optional. Crop dimensions in `width,height`.
    - `resize`: Optional. Resize dimensions in `width,height`.
    - `grayscale`: Optional. Apply grayscale.
    - `invert`: Optional. Apply invert.
    - `files`: Required. Multipart file upload.
  - **Example**:
    ```sh
    curl -X POST http://localhost:8080/image_process \
      -F "files=@/path/to/image.jpg" \
      -F "blur=2.5" \
      -F "crop_anchor=100,100" \
      -F "resize=200,200" \
      -F "grayscale=true" \
      -F "invert=true"
    ```

### View Processed Image Endpoint

- **GET /image/{id}**
  - **Description**: Retrieves the processed image by its ID.
  - **Parameters**:
    - `id`: Required. The ID of the processed image.
  - **Example**:
    ```sh
    curl -X GET http://localhost:8080/image/2d3a1b0e-11a7-4bde-8a1b-c83a29a1c653.jpg
    ```