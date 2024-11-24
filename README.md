**To run servers**
- in a new terminal run: go run server/server.go -nodeId \<nodeId\> -maxNode \<maxNoOfNodes\>
- nodeId: Current node's ID, most likely 0 for the first node
- maxNode: The maximum number of nodes expected to run

**To run client**
- in a new terminal run: go run client/client.go -nodeId \<nodeId\>
- nodeId: The ID of the node of which the client tries to connect to

**To crash a node**
- Press ctrl + c inside the terminal of the node
