# HomeCloudHTTP

This is a custom HTTP library for my HomeCloud project. The reason why this approach of own custom library was chosen because I
wanted to implement my own HTTP library to better understand how one actually works and how it could be implemented.

## TODO

### POST

Need to implement the POST method workflows. How will this look like and What is needed to be implemented.

1. __INIT__     : Server sets up endpoints for post methods like /api/post/..., listens for connection 
2. ___Client__  : Runs a cli command like `./main.go -p file.txt`, need to implement Command in commands.go, connects to server and sends data
3. __Server__   : Reads the requested action, creating a new file contents.
  - If the file already exists returns HTTP status 409 
  - If file doesn't exist it creates file as requested and responds with 204, (201, 200 are also sometimes used).
  - Sends response to Client 
4. __Client__   : The Client reads the response
