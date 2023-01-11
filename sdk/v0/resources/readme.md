# Resources

## What are resources ?

### Step 1:
- User requests the creation/deletion of resources via the cli to the API server.  For e.g. AKS deployment, Azure VM etc.

### Step 2:
- The API server generates a json model of the required resources and passes it to workers in the form of a `message` via MQ.

### Step 3:
- The worker Unmarshalls the request in to the appropriate resource type
- Downloads the remote pulumi program corresponding to the message type
- Injects the `message` as input in to the `Pulumi.<stack>.yaml` of the remote pulumi program
- Runs the remote pulumi program
- Pulumi program updates API server with status upon completion.
