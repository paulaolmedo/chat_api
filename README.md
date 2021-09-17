# Chat API

Basic API to simulate a chat. 
A static version of the  documentation can be seen [here](https://paulaolmedo.gitlab.io/chat_api/openapi.html). Otherwise, the swagger file can be found [here](https://github.com/paulaolmedo/chat_api/tree/main/docs) along with an entity diagram of the database.


### 💻 Prerequisites
Installed Go version >= 16 or Docker

### ❗ 
Inside the code are notes of functionality, and TODO's of things to improve

### Build/Run instructions
#### Without docker
When running without docker, beware of having the Auth0 properties under **pkg/auth**, since it provides the necessary authentication. OR the properties file under **cmd** can be changed to be a **dev** environment, that allows to run the endpoints w/o token in the header.
```
  go run pkg/server.go
```

#### With 🐳 
```
  docker pull paulabeatrizolmedo/chat_api 
```

```
  docker run -it --rm -p 8080:8080 -t paulabeatrizolmedo/chat_api 
```

### Query examples
#### User creation
```
curl --request POST \
  --url http://localhost:8080/users \
  --header 'Content-Type: application/json' \
  --data '{
	"username": "username",
  "password": "password"
}'
```

#### Message creation
```
curl --request POST \
  --url http://localhost:8080/messages \
  --header 'Authorization: Bearer TOKEN' \
  --header 'Content-Type: application/json' \
  --data '{
	"sender": 1,
	"recipient":1,
	"content":{
		"type": "text",
		"text": {
			"text": "hello"
		}
	}
}'
```

#### Search messages
```
curl --request GET \
  --url 'http://localhost:8080/messages?recipient=1&start=1&limit=1' \
  --header 'Authorization: Bearer TOKEN'
```