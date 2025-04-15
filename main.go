package main

import (
	"encoding/gob"
	"errors"

	"github.com/hashicorp/go-plugin"
	"github.com/pigen-dev/artifact-registry-plugin/pkg"
	shared "github.com/pigen-dev/shared"
)



func main(){
	gob.Register(errors.New(""))
	artifactPlugin := &pkg.ArtifactRegistry{}
	pluginMap := map[string]plugin.Plugin{"pigenPlugin": &shared.PigenPlugin{Impl: artifactPlugin}}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         pluginMap,
	})
}