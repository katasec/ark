# Overview 

Tfrunner will:

- Download an ark resource from an OCI registry, unpack it and then execute a terraform init and terraform apply using tfexec

- The TF for every ark resource needs input data. The struct should therefore inject the configdata before running the TF


