# HBT Pub Snooker

## tasks

### apply

directory: stack
environment: AWS_PROFILE=snooker

```shell
terraform apply -var-file=tfvars/dev.auto.tfvars -auto-approve
```

### format

requires: format-tf

### format-tf

directory: stack

```shell
terraform fmt --recursive
```

### init

directory: stack
environment: AWS_PROFILE=snooker

```shell
terraform init -backend-config=tfvars/dev.backend.tfvars
```

### plan

directory: stack
environment: AWS_PROFILE=snooker

```shell
terraform plan -var-file=tfvars/dev.auto.tfvars
```

### push-config

environment: AWS_PROFILE=snooker
directory: stack

```shell
aws s3 cp tfvars/dev.auto.tfvars s3://snooker-dev-state-20260609104639313100000001
aws s3 cp tfvars/dev.backend.tfvars s3://snooker-dev-state-20260609104639313100000001
```

### pull-config

environment: AWS_PROFILE=snooker
directory: stack

```shell
mkdir -p tfvars
aws s3 cp s3://snooker-dev-state-20260609104639313100000001/dev.auto.tfvars tfvars/
aws s3 cp s3://snooker-dev-state-20260609104639313100000001/dev.backend.tfvars tfvars/
```
