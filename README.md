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

 ## Configuring
Go to ./services/config <br/>
Create and edit the config.toml as in following example:
```
bad_words = [
  "first bad word",
  "second bad word".
]

[[list_of_endpoints]]
url = "http://your_url_1"
max_time = 10
max_retries = 4
return_data = "text"

[[list_of_endpoints]]
url = "http://your_url_2"
max_time = 10
max_retries = 4
return_data = "text"
```

##### bad_words is a list of words to be found in endpoint response

##### list_of_endpoints is a list of endpoints data where:
- url is a... URL 
- max_time is a maximum time client will wait for endpoint response
- max_retries is a... maximum amount of retries to get endpoint data
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

In project's root directory run:
```
//For Windows
docker-compose up

//For Linux
docker compose up
```

# Testing

To test service run:
```
cd ./service

go test
```

## Response examples
```
{
    "http://host.docker.internal:8000/three": {
        "total_count": 3,
        "words": [
            {
                "expr_index": "0",
                "index": "11",
                "word": "bad gopher"
            },
            {
                "expr_index": "1",
                "index": "6",
                "word": "good python"
            },
            {
                "expr_index": "1",
                "index": "18",
                "word": "bad gopher"
            }
        ]
    },
    "http://host.docker.internal:8000/two": {
        "total_count": 1,
        "words": [
            {
                "index": "21",
                "word": "bad gopher"
            }
        ]
    }
}
```
