# swaggerui example
This is a little webapp with mockup database to demonstrate how you might integrate swaggerui into your projects and interact with them through the embedded UI.

To start it, simply `go run *.go` or `go build` in this directory and then interact with your local server at `localhost:8080/swagger/` (note the slash)

Both routes should work, though note that the first route requires "authorization", wich means you have to click the little lock icon and "sign in", which will start attaching the auth header to your requests.
