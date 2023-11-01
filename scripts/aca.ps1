#!/usr/bin/env pwsh

$ErrorActionPreference = 'Stop'

Import-Module Az.App
# Connect to your Azure subscription.
#Connect-AzAccount


# Get App Env Info
Write-Output "Getting App Env"
$env = Get-AzContainerAppManagedEnv | Where-Object { $_.name -eq "consumption" }
$env

# Create App Template
Write-Output "Creating App Template"

$port = 80
# $envVar1 = New-Object -TypeName Microsoft.Azure.PowerShell.Cmdlets.ContainerApp.Models.IEnvironmentVar
# $envVar1.Name = "ARK_WEB_PORT"
# $envVar1.Value = $port

# $envVar2 = New-Object -TypeName Microsoft.Azure.PowerShell.Cmdlets.ContainerApp.Models.IEnvironmentVar
# $envVar2.Name = "MY_ENV_VAR2"
# $envVar2.Value = "Goodbye, world!"

$envVars = @{
    ARK_WEB_PORT = $port
}

$version = git describe --tags --abbrev=0
$containerAppTemplateObjectParams = @{
    Name = "azps-containerapp"
    Image = "ghcr.io/katasec/arkserver:$version"
    ResourceCpu = 0.25
    ResourceMemory = "0.5Gi"
    Command = "/ark web"
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
    IngressTargetPort = $port
}

$app=New-AzContainerApp @containerAppParams
$app | Format-Table


#Remove-AzContainerApp -Name app1 -ResourceGroupName $env.ResourceGroupName