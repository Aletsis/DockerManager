package docker

import (
	"strings"

	containerdomain "dockermanager/internal/domain/container"
	imagedomain "dockermanager/internal/domain/image"

	dockertypes "github.com/docker/docker/api/types"
	imagedockertypes "github.com/docker/docker/api/types/image"
)

func toDomainContainer(c dockertypes.Container) containerdomain.Container {
	name := ""
	if len(c.Names) > 0 {
		name = strings.TrimPrefix(c.Names[0], "/")
	}

	shortID := c.ID
	if len(shortID) > 12 {
		shortID = shortID[:12]
	}

	ports := make([]containerdomain.PortMapping, 0, len(c.Ports))
	for _, p := range c.Ports {
		ports = append(ports, containerdomain.PortMapping{
			IP:          p.IP,
			PrivatePort: p.PrivatePort,
			PublicPort:  p.PublicPort,
			Type:        p.Type,
		})
	}

	composeProject := ""
	composeService := ""
	composeWorkingDir := ""
	composeConfigFile := ""
	if c.Labels != nil {
		composeProject = c.Labels["com.docker.compose.project"]
		composeService = c.Labels["com.docker.compose.service"]
		composeWorkingDir = c.Labels["com.docker.compose.project.working_dir"]
		composeConfigFile = c.Labels["com.docker.compose.project.config_files"]
	}

	networks := make([]containerdomain.NetworkInfo, 0)
	if c.NetworkSettings != nil && c.NetworkSettings.Networks != nil {
		for netName, netEndpoint := range c.NetworkSettings.Networks {
			if netEndpoint == nil {
				continue
			}
			networks = append(networks, containerdomain.NetworkInfo{
				NetworkName: netName,
				NetworkID:   netEndpoint.NetworkID,
				IPAddress:   netEndpoint.IPAddress,
				Gateway:     netEndpoint.Gateway,
				MacAddress:  netEndpoint.MacAddress,
				Aliases:     netEndpoint.Aliases,
			})
		}
	}

	return containerdomain.Container{
		ID:                c.ID,
		ShortID:           shortID,
		Names:             c.Names,
		Name:              name,
		Image:             c.Image,
		ImageID:           c.ImageID,
		Command:           c.Command,
		Created:           c.Created,
		State:             containerdomain.State(c.State),
		Status:            c.Status,
		Ports:             ports,
		Networks:          networks,
		SizeRw:            c.SizeRw,
		SizeRootFs:        c.SizeRootFs,
		Labels:            c.Labels,
		ComposeProject:    composeProject,
		ComposeService:    composeService,
		ComposeWorkingDir: composeWorkingDir,
		ComposeConfigFile: composeConfigFile,
	}
}

func toDomainImage(img imagedockertypes.Summary, inUse bool) imagedomain.Image {
	shortID := img.ID
	if len(shortID) > 19 && strings.HasPrefix(shortID, "sha256:") {
		shortID = shortID[7:19]
	} else if len(shortID) > 12 {
		shortID = shortID[:12]
	}

	repo := "<none>"
	tag := "<none>"
	isDangling := len(img.RepoTags) == 0 || (len(img.RepoTags) == 1 && img.RepoTags[0] == "<none>:<none>")

	if len(img.RepoTags) > 0 && img.RepoTags[0] != "<none>:<none>" {
		repo, tag, isDangling = ParseRepoTag(img.RepoTags[0])
	}

	return imagedomain.Image{
		ID:         img.ID,
		ShortID:    shortID,
		Repository: repo,
		Tag:        tag,
		RepoTags:   img.RepoTags,
		Created:    img.Created,
		Size:       img.Size,
		SharedSize: img.SharedSize,
		Containers: img.Containers,
		InUse:      inUse,
		IsDangling: isDangling,
	}
}
