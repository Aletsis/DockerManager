package docker

import (
	"testing"

	dockertypes "github.com/docker/docker/api/types"
	imagedockertypes "github.com/docker/docker/api/types/image"
	networktypes "github.com/docker/docker/api/types/network"
)

func TestToDomainContainer(t *testing.T) {
	raw := dockertypes.Container{
		ID:    "1234567890abcdef123456",
		Names: []string{"/web-server"},
		Image: "nginx:alpine",
		State: "running",
		Ports: []dockertypes.Port{
			{IP: "0.0.0.0", PrivatePort: 80, PublicPort: 8080, Type: "tcp"},
		},
		Labels: map[string]string{
			"com.docker.compose.project":             "my_project",
			"com.docker.compose.service":             "web",
			"com.docker.compose.project.working_dir": "/app",
		},
		NetworkSettings: &dockertypes.SummaryNetworkSettings{
			Networks: map[string]*networktypes.EndpointSettings{
				"bridge": {
					NetworkID: "net_bridge_1",
					IPAddress: "172.17.0.2",
					Gateway:   "172.17.0.1",
				},
			},
		},
	}

	domainContainer := toDomainContainer(raw)

	if domainContainer.ID != raw.ID {
		t.Errorf("expected ID %s, got %s", raw.ID, domainContainer.ID)
	}
	if domainContainer.ShortID != "1234567890ab" {
		t.Errorf("expected ShortID '1234567890ab', got %s", domainContainer.ShortID)
	}
	if domainContainer.Name != "web-server" {
		t.Errorf("expected Name 'web-server', got %s", domainContainer.Name)
	}
	if domainContainer.ComposeProject != "my_project" || domainContainer.ComposeService != "web" {
		t.Errorf("unexpected compose project/service: %s/%s", domainContainer.ComposeProject, domainContainer.ComposeService)
	}
	if len(domainContainer.Ports) != 1 || domainContainer.Ports[0].PublicPort != 8080 {
		t.Errorf("unexpected ports: %+v", domainContainer.Ports)
	}
	if len(domainContainer.Networks) != 1 || domainContainer.Networks[0].IPAddress != "172.17.0.2" {
		t.Errorf("unexpected networks: %+v", domainContainer.Networks)
	}
}

func TestToDomainImage(t *testing.T) {
	// 1. Tagged image
	rawTagged := imagedockertypes.Summary{
		ID:         "sha256:abcdef1234567890123456",
		RepoTags:   []string{"postgres:15-alpine"},
		Size:       300000000,
		Containers: 1,
	}
	imgTagged := toDomainImage(rawTagged, true)
	if imgTagged.Repository != "postgres" || imgTagged.Tag != "15-alpine" {
		t.Errorf("unexpected repo/tag: %s:%s", imgTagged.Repository, imgTagged.Tag)
	}
	if imgTagged.ShortID != "abcdef123456" {
		t.Errorf("expected short ID abcdef123456, got %s", imgTagged.ShortID)
	}
	if !imgTagged.InUse || imgTagged.IsDangling {
		t.Errorf("expected inUse=true, isDangling=false, got inUse=%v isDangling=%v", imgTagged.InUse, imgTagged.IsDangling)
	}

	// 2. Dangling image
	rawDangling := imagedockertypes.Summary{
		ID:         "sha256:9999991234567890123456",
		RepoTags:   []string{"<none>:<none>"},
		Size:       150000000,
		Containers: 0,
	}
	imgDangling := toDomainImage(rawDangling, false)
	if !imgDangling.IsDangling {
		t.Errorf("expected dangling image, got isDangling=false")
	}
	if imgDangling.InUse {
		t.Errorf("expected inUse=false, got true")
	}
}
