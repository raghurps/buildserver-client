# buildserver-client

Client library to communicate with different build servers

## Currently supported build server

- Teamcity

## Download library

```
# go get -u github.com/raghuP9/buildserver-client
```

## Test with example

### Generate teamcityctl binary

```
# go get -u github.com/raghuP9/buildserver-client/cmd/teamcityctl
# $GOPATH/bin/teamcityctl --help
```

### Add build to queue

```
# export TEAMCITY_TOKEN=<token>
# $GOPATH/bin/teamcityctl --server http://teamcity.example.com start-build --pipeline <pipeline_id> --branch <branch_name> --param KEY1=VALUE1 --param KEY2=VALUE2 --dependency pipeline_id1:build_id1 --dependency pipeline_id2:build_id2
```

### Get build details by id

```
# export TEAMCITY_TOKEN=<token>
# $GOPATH/bin/teamcityctl --server http://teamcity.example.com get-build --id <build_id>
```

### Stop running build by id

```
# export TEAMCITY_TOKEN=<token>
# $GOPATH/bin/teamcityctl --server http://teamcity.example.com stop-build --id <build_id>
```

### Cancel queued build

```
# export TEAMCITY_TOKEN=<token>
# $GOPATH/bin/teamcityctl --server http://teamcity.example.com cancel-build --id <build_id>
```

### Get content from a text file in artifact for a given build id

```
# export TEAMCITY_TOKEN=<token>
# $GOPATH/bin/teamcityctl --server http://teamcity.example.com fetch-artifact --id <build_id> --path <path_relative_to_artifacts_directory>
```
