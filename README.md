# Simple Mock Server
This project provides a simple mock server that reads route definitions from a JSON file and serves the corresponding responses when accessed. The server automatically reloads the route definitions whenever the JSON file is modified.

Getting Started
To get started with the Simple Mock Server, you will need to have Go installed on your system.

## Installation
1. Clone the repository:
```
git clone https://github.com/eidenar/backend_mocker.git
```

2. Change into the project directory:
```
cd backend_mocker
```

3. Ensure the structure.json file exists in the project directory with the route definitions you want to use. See the Example JSON File section for an example.

4. Build the project:
```
go build
```

5. Run the server:
```
./main
```

By default, the server will start on 127.0.0.1:8080. You can change the host and port by providing command-line flags:
```
./main -host 0.0.0.0 -port 3000
```


## Example JSON File
Here's an example JSON file containing two routes:

```
[
  {
    "url": "/api/users",
    "response": {
      "users": [
        {
          "id": 1,
          "name": "John Doe",
          "email": "john.doe@example.com"
        },
        {
          "id": 2,
          "name": "Jane Doe",
          "email": "jane.doe@example.com"
        }
      ]
    }
  },
  {
    "url": "/api/products",
    "response": {
      "products": [
        {
          "id": 1,
          "name": "Laptop",
          "price": 1200
        },
        {
          "id": 2,
          "name": "Smartphone",
          "price": 800
        }
      ]
    }
  }
]
```


## Testing
To run the tests for this project, execute the following command in the project directory:
```
go test -v
```
