# GoodWords
> Service designed to find specified "bad words" in endpoints responses

## Requirements
#### General:
- Golang v1.24.0 <br/>
![golang](https://badgen.net/static/go/1.24.0/green?icon=github)<br/>
You can install Golang <a href="https://go.dev/doc/install">here</a>

#### With Docker:
- Docker <br/>
 ![docker](https://badgen.net/static/docker/@latest/purple)<br/>
 You can install Docker <a href="https://docs.docker.com/engine/install/">here</a>

 ## Installing
 To install service just clone this repository

## Testing
In root derictory:</br>
To test service run:
```
cd ./service

go test
```

 ## Configuring
Go to ./services/config <br/>
Create and edit the config.toml as in config.example.toml:
```
bad_words = [
  "good python",
  "bad gopher",
]

[[list_of_endpoints]]
url = "http://localhost:8000/one"
max_time = 10
max_retries = 4
return_data = "text"

[[list_of_endpoints]]
url = "http://localhost:8000/two"
max_time = 10
max_retries = 4
return_data = "text"

[[list_of_endpoints]]
url = "http://localhost:8000/three"
max_time = 10
max_retries = 4
return_data = "array"

[[list_of_endpoints]]
url = "http://localhost:8001/three"
max_time = 10
max_retries = 4
return_data = "array"
```

##### "bad_words" is a list of words to be found in endpoint response

##### "list_of_endpoints" is a list of endpoints data where:
- url is a... URL 
- max_time is a maximum time client will wait for endpoint response
- max_retries is a maximum amount of retries to get endpoint data
- return_data is a datatype which endpoint return. Can be either "text" or "array". 

## Running
#### Without Docker
In project's root directory run:
```
cd ./service

go build ./cmd/main.go

//For Windows
Main.exe

//For Linux
./Main
```

#### With Docker
!!!</br> 
If you want to analyze responses from endpoint on your local machine (e.g localhost), you</br>
need to change "localhost" in URL to "host.docker.internal". Or if your running endpoint as</br>
docker container, you need to change "localhost" in URL to </name_of_endpoint_container/></br>
!!!</br> 

In project's root directory run:
```
//For Windows
docker-compose up

//For Linux
docker compose up
```

## Program work example
#### For given input:
List of bad words:

```
bad_words = [
  "good python",
  "bad gopher",
  "bad man"
]
```

Data received from endpoints stated in config.example.toml:

1. From "http://localhost:8000/one"
```
["GooD pYthoN hello what bad GOpher, Bad GOpHer, Good Python"]
```
2. From "http://localhost:8000/two"
```
["hello, golang is good"]
```
3. From "http://localhost:8000/three"
```
[
"bad GOpher hello what bad GOpher",
"hello good python bad GOpher good python",
"good python"
]
```
4. From "http://localhost:8001/three"
```
nothing because nothing runs on localhost:8001
```

####  Program work result:

```
{
    "http://localhost:8000/one": { //output result for "text" return datatype
        "total_count": 4,
        "words": [
            {
                "index": 0,
                "word": "good python"
            },
            {
                "index": 47,
                "word": "good python"
            },
            {
                "index": 23,
                "word": "bad gopher"
            },
            {
                "index": 35,
                "word": "bad gopher"
            }
        ]
    },
    "http://localhost:8000/three": { //output result for "array" return datatype
        "total_count": 6,
        "words": [
            {
                "expr_index": 0,
                "index": 0,
                "word": "bad gopher"
            },
            {
                "expr_index": 0,
                "index": 22,
                "word": "bad gopher"
            },
            {
                "expr_index": 1,
                "index": 6,
                "word": "good python"
            },
            {
                "expr_index": 1,
                "index": 29,
                "word": "good python"
            },
            {
                "expr_index": 1,
                "index": 18,
                "word": "bad gopher"
            },
            {
                "expr_index": 2,
                "index": 0,
                "word": "good python"
            }
        ]
    },
    "http://localhost:8000/two": { // output for data without "bad words"
        "total_count": 0,
        "words": null
    },
    "http://localhost:8000/three": "too many retries" //output result if program can't get data after "max_retries" retries attempts
}
```
