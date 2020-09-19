package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/raghuP9/buildserver-client/pkg/buildserver/teamcity"
	"github.com/urfave/cli/v2"
)

func startBuild(c *cli.Context) error {
	client := teamcity.NewTeamcityClient(
		5*time.Second,
		5*time.Second,
		5*time.Second,
		c.String("server"),
		fmt.Sprintf("Bearer %s", c.String("token")), c.Bool("secure"),
	)
	paramsMap := map[string]string{}
	for _, v := range c.StringSlice("param") {
		param := strings.Split(v, "=")
		if len(param) != 2 {
			err := errors.New("Params not provided in the form of KEY=VALUE")
			log.Println(err.Error())
			return err
		}
		paramsMap[param[0]] = param[1]
	}

	snapDependencyMap := map[string]int{}
	for _, v := range c.StringSlice("snapshot-dependency") {
		dependency := strings.Split(v, "=")
		if len(dependency) != 2 {
			err := errors.New("Snapshot dependency not provided in the form of buildPipelineID=buildID")
			log.Println(err.Error())
			return err
		}
		snapDependencyMap[dependency[0]], _ = strconv.Atoi(dependency[1])
	}

	artfDependencyMap := map[string]int{}
	for _, v := range c.StringSlice("artifact-dependency") {
		dependency := strings.Split(v, "=")
		if len(dependency) != 2 {
			err := errors.New("Artifact dependency not provided in the form of buildPipelineID=buildID")
			log.Println(err.Error())
			return err
		}
		artfDependencyMap[dependency[0]], _ = strconv.Atoi(dependency[1])
	}

	id, err := client.StartBuild(
		c.String("pipeline"),
		c.String("branch"),
		c.String("comment"),
		paramsMap,
		snapDependencyMap,
		artfDependencyMap,
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	log.Printf("Started build with ID: %d\n", id)
	return nil
}

func cancelBuild(c *cli.Context) error {
	client := teamcity.NewTeamcityClient(
		5*time.Second,
		5*time.Second,
		5*time.Second,
		c.String("server"),
		fmt.Sprintf("Bearer %s", c.String("token")), c.Bool("secure"),
	)
	id := c.Int("id")
	err := client.CancelQueuedBuild(id, c.String("comment"))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Printf("Successfully cancelled queued build with id: %d\n", id)
	return nil
}

func stopBuild(c *cli.Context) error {
	client := teamcity.NewTeamcityClient(
		5*time.Second,
		5*time.Second,
		5*time.Second,
		c.String("server"),
		fmt.Sprintf("Bearer %s", c.String("token")), c.Bool("secure"),
	)
	id := c.Int("id")
	err := client.StopBuild(id, c.String("comment"))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Printf("Successfully cancelled queued build with id: %d\n", id)
	return nil
}

func statusBuild(c *cli.Context) error {
	client := teamcity.NewTeamcityClient(
		5*time.Second,
		5*time.Second,
		5*time.Second,
		c.String("server"),
		fmt.Sprintf("Bearer %s", c.String("token")), c.Bool("secure"),
	)
	id := c.Int("id")
	details := &teamcity.TCBuildDetails{}
	err := client.GetBuild(id, details)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Printf("Successfully fetched status for build with id: %d\n", id)
	log.Println(details)
	return nil
}

func fetchArtifact(c *cli.Context) error {
	client := teamcity.NewTeamcityClient(
		5*time.Second,
		5*time.Second,
		5*time.Second,
		c.String("server"),
		fmt.Sprintf("Bearer %s", c.String("token")), c.Bool("secure"),
	)
	id := c.Int("id")
	content, contentType, err := client.GetArtifactTextFile(c.String("path"), c.Int("id"))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Printf("Successfully fetched artifact file from build with id: %d\n", id)
	log.Printf("Content: \n%s\n\n", string(content))
	log.Printf("Content-Type: %s\n", contentType)
	return nil
}

func main() {
	app := &cli.App{
		Name:    "teamcityctl",
		Version: "1.0.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "token",
				Usage:    "Provide auth token to talk to teamcity build server",
				EnvVars:  []string{"TEAMCITY_TOKEN"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "server",
				Aliases:  []string{"serverURL"},
				Usage:    "Provide teamcity server URL",
				Required: true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "start-build",
				Usage: "Start teamcity build by adding build into queue",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "pipeline",
						Usage:    "Provide build pipeline ID",
						Aliases:  []string{"buildPipeline"},
						Required: true,
					},
					&cli.StringFlag{
						Name:     "branch",
						Usage:    "Provide branch name to perform build upon",
						Required: true,
					},
					&cli.StringSliceFlag{
						Name:  "param",
						Usage: "Provide multiple params as key:value, e.g. --param key1=value1 --param key2=value2",
					},
					&cli.StringSliceFlag{
						Name: "snapshot-dependency",
						Usage: "Provide multiple build snapshot dependencies as buildPipelineID:buildID," +
							" e.g --snapshot-dependency myBuildConfigID1:uniqeBuildID1 --snapshot-dependency  myBuildConfigID2:uniqeBuildID2",
					},
					&cli.StringSliceFlag{
						Name: "artifact-dependency",
						Usage: "Provide multiple build artifact dependencies as buildPipelineID:buildID," +
							" e.g --artifact-dependency myBuildConfigID1:uniqeBuildID1 --artifact-dependency  myBuildConfigID2:uniqeBuildID2",
					},
					&cli.StringFlag{
						Name:  "comment",
						Usage: "Provide text comment",
						Value: "Build started by teamcityctl CLI",
					},
				},
				Action: startBuild,
			},
			{
				Name:  "get-build",
				Usage: "get details of a build",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "id",
						Usage:    "Provide unique build ID whose details is required",
						Required: true,
					},
				},
				Action: statusBuild,
			},
			{
				Name:  "cancel-build",
				Usage: "Cancel a queued build that is not yet running or completed",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "id",
						Usage:    "Provide unique build ID that needs to be cancelled",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "comment",
						Usage: "Provide text comment",
						Value: "Build started by teamcityctl CLI",
					},
				},
				Action: cancelBuild,
			},
			{
				Name:  "stop-build",
				Usage: "stop a running build",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "id",
						Usage:    "Provide unique build ID that needs to be stopped",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "comment",
						Usage: "Provide text comment",
						Value: "Build started by teamcityctl CLI",
					},
				},
				Action: stopBuild,
			},
			{
				Name:  "fetch-artifact",
				Usage: "Cancel a queued build that is not yet running or completed",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "id",
						Usage:    "Provide unique build ID that needs to be cancelled",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "path",
						Usage:    "Provide artifact file path relative to artifacts directory",
						Required: true,
					},
				},
				Action: fetchArtifact,
			},
		},
	}
	app.Run(os.Args)
}
