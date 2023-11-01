#!/usr/bin/env pwsh

Import-Module Az.App
# Connect to your Azure subscription.
#Connect-AzAccount


# Get App Env Info
Write-Output "Getting App Env"
$env = Get-AzContainerAppManagedEnv | Where { $_.name -eq "consumption" }
$env

# Create App Template
Write-Output "Creating App Template"

$envVar1 = New-Object -TypeName Microsoft.Azure.PowerShell.Cmdlets.ContainerApp.Models.IEnvironmentVar
$envVar1.Name = "ARK_WEB_PORT"
$envVar1.Value = "80"

$envVar2 = New-Object -TypeName Microsoft.Azure.PowerShell.Cmdlets.ContainerApp.Models.IEnvironmentVar
$envVar2.Name = "MY_ENV_VAR2"
$envVar2.Value = "Goodbye, world!"

$envVars = $envVar1, $envVar2

$containerAppTemplateObjectParams = @{
    Name = "azps-containerapp"
    Image = "ghcr.io/katasec/arkserver:v0.0.11"
    ResourceCpu = 0.25
    ResourceMemory = "0.5Gi"
    Command = "ark web"
    Env = $envVars
}
$image = New-AzContainerAppTemplateObject @containerAppTemplateObjectParams

# Create/Update App
Write-Output "Deploying App..."
$containerAppParams = @{
    Name = "app1"
    ResourceGroupName = $env.ResourceGroupName
    TemplateContainer = $image
    Location = $env.Location
    ManagedEnvironmentId = $env.Id
    IngressExternal = $true
    IngressTransport = "auto"
    IngressTargetPort = 8080
}

$app=New-AzContainerApp @containerAppParams
$app | Format-Table