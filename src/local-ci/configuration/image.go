package configuration

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"local-ci/utils"
)

type Image struct {
	name       string
	entrypoint []string
}

func (image Image) ConfigureContainer(containerConfiguration *container.Config) {
	containerConfiguration.Image = image.name

	if len(image.entrypoint) > 0 {
		containerConfiguration.Entrypoint = image.entrypoint
	}
}

func (image Image) Pull(client *client.Client) (io.ReadCloser, error) {
	return client.ImagePull(context.Background(), image.name, types.ImagePullOptions{})
}

func CreateImageFromYaml(yml interface{}) *Image {
	image := Image{}

	if name, ok := yml.(string); ok {
		image.name = name
	}

	if name, ok := yml.(map[interface{}]interface{}); ok {
		image.name = name["name"].(string)

		if nil != name["entrypoint"] {
			entrypoint, _ := utils.InterfaceArrayToStringArray(name["entrypoint"].([]interface{}))

			image.entrypoint = entrypoint
		}
	}

	return &image
}
