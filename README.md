# buildserver-client

Client library to communicate with different build servers

## Currently supported build server

- Teamcity

## Download library

```
# go get -u github.com/raghuP9/buildserver-client
```

## Trigger builds from command line using CLI

### Generate teamcityctl binary

```
# go get -u github.com/raghuP9/buildserver-client/cmd/teamcityctl
# $GOPATH/bin/teamcityctl --help
```

### Add build to queue

```
# export TEAMCITY_TOKEN=<token>
# $GOPATH/bin/teamcityctl --server http://teamcity.example.com start-build --pipeline <pipeline_id> --branch <branch_name> \
   --param KEY1=VALUE1 --param KEY2=VALUE2 --dependency pipeline_id1:build_id1 --dependency pipeline_id2:build_id2 \
   --comment "<your text comment>"
```

### Get build details by id

```
# export TEAMCITY_TOKEN=<token>
# $GOPATH/bin/teamcityctl --server http://teamcity.example.com get-build --id <build_id>
```

### Stop running build by id

```
# export TEAMCITY_TOKEN=<token>
# $GOPATH/bin/teamcityctl --server http://teamcity.example.com stop-build --id <build_id> --comment "<your text comment>
```

### Cancel queued build

```
# export TEAMCITY_TOKEN=<token>
# $GOPATH/bin/teamcityctl --server http://teamcity.example.com cancel-build --id <build_id> --comment "<your text comment>
```

### Get content from a text file in artifact for a given build id

```
# export TEAMCITY_TOKEN=<token>
# $GOPATH/bin/teamcityctl --server http://teamcity.example.com fetch-artifact --id <build_id> --path <path_relative_to_artifacts_directory>
```
## Make API calls to teamcity build server from your code
### Create client
```
package main

import (
  "time"
  "github.com/raghuP9/buildserver-client/pkg/buildserver/teamcity"
)

func main() {
  client := teamcity.NewTeamcityClient(
    5 * time.Second,               // http request timeout
    5 * time.Second,               // http dial timeout
    5 * time.Second,               // TLS handshake timeout
    "http://myteamcityserver.com", // teamcity server URL with https or http whichever applies
    "<teamcity-token>",            // teamcity token to make authenticated requests to build server
    false,                         // skip certficate validation in case using self signed certificates
  )
}
```

### Trigger builds using client
```
id, err := client.StartBuild(
  "<teamcityBuildTypeID>",                // Provide build configuration ID to trigger that build on
  "<branch-name>",                        // Provide branch name on which you want to trigger build
  "<text-comment-on-build>",              // Provide a comment
  map[string]string{                      // Provide params/env variables that you want to
    "env.MY_VAR1": "MY_VALUE1",           // add/override when triggering build
    "env.MY_VAR2": "MY_VALUE2",
  },
  map[string]int{                         // Provide build IDs of other builds which are specified as 
    "<teamcityBuildTypeID1>": 123456,     // snapshot dependencies and whose already built artifacts you
    "<teamcityBuildTypeID2>": 123457,     // want to use in this build.
  },
)
```
### Get build status by ID(int)
```
statusDetails := teamcity.TCBuildDetails{}
err := client.GetBuild(id, &statusDetails)
```
### Cancel a queued build by ID (int)
```
err := client.CancelQueuedBuild(id, "your-comment-for-cancelling-build")
```
### Stop a running build by ID (int)
```
err := client.StopBuild(id, "your-comment-for-stopping-build")
```
### Get artifact text file (currently only supported for smaller text files)
```
byteArray, contentType, err := client.GetArtifactTextFile("path/to/artifact", id)
```
