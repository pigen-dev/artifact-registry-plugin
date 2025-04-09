package pkg

import (
	"context"
	"fmt"
	"log"
	"strings"

	artifactregistry "cloud.google.com/go/artifactregistry/apiv1"
	artifactregistrypb "cloud.google.com/go/artifactregistry/apiv1/artifactregistrypb"
	"github.com/pigen-dev/artifact-registry-plugin/helpers"
	shared "github.com/pigen-dev/shared"
	"google.golang.org/api/iterator"
)

type ArtifactRegistry struct {
	Config Config `yaml:"config"`
	Output Output `yaml:"output"`
}

type Config struct {
	Region string `yaml:"region"`
	RepoName string `yaml:"repo_name"`
	ProjectId string `yaml:"project_id"`
}

type Output struct {
	RepoUrl string `yaml:"repo_url"`
}

func (ar *ArtifactRegistry) ParseConfig(in map[string] any) error {
	config:=Config{}
	err:= helpers.YamlConfigParser(in, &config)
	if err != nil {
		return err
	}
	ar.Config=config
	log.Println("/////////////////////////////")
	log.Println(ar.Config.RepoName)
	log.Println("/////////////////////////////")
	return nil
}

func (ar *ArtifactRegistry) SetupPlugin() error {
	ctx := context.Background()
	c, err := artifactregistry.NewClient(ctx)
	if err != nil {
		return err
	}

	repo := artifactregistrypb.Repository{
		Name:"projects/"+ar.Config.ProjectId+"/locations/"+ar.Config.Region+"/repositories/"+ar.Config.RepoName,
		Format: artifactregistrypb.Repository_DOCKER,
	}

	repo_request := artifactregistrypb.CreateRepositoryRequest{
		Parent:"projects/"+ar.Config.ProjectId+"/locations/"+ar.Config.Region,
		RepositoryId: ar.Config.RepoName,
		Repository: &repo,
	}
	resp, err := c.CreateRepository(ctx,&repo_request)
	if err != nil {
		return err
	}
	resp.Wait(ctx)
	
	defer c.Close()
	return nil
}

func (ar *ArtifactRegistry) GetOutput() shared.GetOutputResponse {
	ctx := context.Background()
	c, err := artifactregistry.NewClient(ctx)
	if err != nil {
		return shared.GetOutputResponse{Output: nil, Error: err}
	}
	listRepositoriesRequest := &artifactregistrypb.ListRepositoriesRequest{
		Parent:"projects/"+ar.Config.ProjectId+"/locations/"+ar.Config.Region,
	}
	repositories := c.ListRepositories(ctx,listRepositoriesRequest)
	for {
		resp, err := repositories.Next()
		if err == iterator.Done {
			return shared.GetOutputResponse{Output: nil, Error: err}
		}
		if err != nil {
			return shared.GetOutputResponse{Output: nil, Error: err}
		}
		if resp.Name == "projects/"+ar.Config.ProjectId+"/locations/"+ar.Config.Region+"/repositories/"+ar.Config.RepoName {
			ar.Output.RepoUrl = ar.Config.Region + "-" + strings.ToLower(resp.Format.String()) + ".pkg.dev/" + ar.Config.ProjectId + "/" + ar.Config.RepoName
			break
		}
	}
	output, err := helpers.StructToMap(ar.Output)
	if err != nil {
		return shared.GetOutputResponse{Output: nil, Error: err}
	}
	return shared.GetOutputResponse{Output: output, Error: nil}
}


func (ar *ArtifactRegistry) Destroy() error {
	ctx := context.Background()
	c, err := artifactregistry.NewClient(ctx)
	if err != nil {
		return err
	}
	deleteRepositoryRequest := &artifactregistrypb.DeleteRepositoryRequest{
		Name: "projects/"+ar.Config.ProjectId+"/locations/"+ar.Config.Region+"/repositories/"+ar.Config.RepoName,
	}
	resp, err := c.DeleteRepository(ctx, deleteRepositoryRequest)
	if err != nil {
		return fmt.Errorf("error destroying repository")
	}
	err = resp.Wait(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for destroying repository")
	}
	return nil
}