To support loading access logs and evaluating traffic, you have to modify `src/mirage/mirage.go`. The `mirage.go` generates content requests from each client connected to the network in the lines 71-76 and 82-87. The lines send content requests by `client.RandomRequests()` that generates a content request at the specific access pattern (zipf or gamma) based on your configuration file. The `client` variable also has a method `RequestByID(int)` that generates a request for the content with the specified ID (see `src/graph/client.go` for detailed implementation). You can check the cache server's ID by `client.Upstream().ID()` before `client.RequestByID(int)`.

I also recommend you to generate a client map like the following code:
```
// make a map for easy access to the client variables by their upstream node IDs.
clients := make(map[string]Client)
for _, client := range network.Clients() {
    clients[client.Upstream().ID()] = client
}
// request for a content with ID=152 from a client below "node01"
clients["node05"].RequestByID(152)
```

