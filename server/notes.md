

# Brainstorming


## Concepts

### What's a cloudspace ?

A `cloudspace` is a cloud environment where users can securely deploy cloud infrastructure & services. A `cloudpsace` consists of group of networks that are tagged and managed by ark. The smallest cloudspace contains 1 DMZ network with 1 application network. These networks are built to be conformant with the cloud provider's recommended Well Architected Framework patterns. Users can extend a cloudspace by adding additional cloud networks. 

Users can securely deploy cloud resources in to a `cloudspace` without having to understand the underlying cloud networking infrastructure.


```go
type CloudSpace stuct {
    Name
    HubNetwork
    SpokeNetworks SpokeNetwork
}

type HubNetwork struct {
    Firewall
    LogWorkspaces
    EventGrid
    RelayBridge
}
```

## Things I want to do

### Flow 1
- Create a cloudspace

```
ark create cloudspace coolspace1
```

- Create a vm in the cloudspace

Select the cloudspace:

```
ark cloudspace select coolspace1
```
Launch VM in the cloudspace:

```
ark vm create -projectName "myproject" -name "myvm01" -adminuser "azureuser" -password "mypassword" -imageName "MicrosoftWindowsServer:WindowsServer:2022-Datacenter:latest"
```

This generates a yaml file myvm01.yaml:
```
version: "0.1"
cloudid: azure
projectName: project1
resources:
  vm:
    name: myvm01
    image: MicrosoftWindowsServer:WindowsServer:2022-Datacenter:latest 
    asg: myvm01
    adminuser: "azureuser"
    adminpassword: "mypassword"
```

Apply the generated file:

```
ark apply -f myvm01.yaml
```

- Connect to the VM in the cloudspace
- Create an AKS in the cloudspace
- Run commands against the AKS cluster in the cloudspace



|Path|Verb|Purpose|
|-|-|-|
|/azure/cloudspace| POST | Init Cloud Space
|/azure/cloudspace/vm | POST | Create a VM in the cloud space
|/azure/cloudspace/vm | DELETE | Delte a VM in the cloud space


