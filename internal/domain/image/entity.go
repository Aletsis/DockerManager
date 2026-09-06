package image

import "dockermanager/internal/domain"

// Image represents a local Docker/OCI container image
type Image struct {
	ID         string   `json:"id"`
	ShortID    string   `json:"shortId"`
	Repository string   `json:"repository"`
	Tag        string   `json:"tag"`
	RepoTags   []string `json:"repoTags"`
	Created    int64    `json:"created"`
	Size       int64    `json:"size"`
	SharedSize int64    `json:"sharedSize"`
	Containers int64    `json:"containers"`
	InUse      bool     `json:"inUse"`
	IsDangling bool     `json:"isDangling"`
}

// CanRemove checks whether an image can be removed based on usage invariants
func (img *Image) CanRemove(force bool) error {
	if img.InUse && !force {
		return domain.ErrImageInUse
	}
	return nil
}

// DiskUsageSummary provides aggregated disk metrics
type DiskUsageSummary struct {
	TotalImages     int   `json:"totalImages"`
	TotalSize       int64 `json:"totalSize"`
	DanglingCount   int   `json:"danglingCount"`
	DanglingSize    int64 `json:"danglingSize"`
	ReclaimableSize int64 `json:"reclaimableSize"`
}

// PruneResult contains results of image pruning
type PruneResult struct {
	ImagesDeleted  []string `json:"imagesDeleted"`
	SpaceReclaimed uint64   `json:"spaceReclaimed"`
}

// PullProgressEvent represents layer-level download progress
type PullProgressEvent struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Progress string `json:"progress"`
	Current  int64  `json:"current"`
	Total    int64  `json:"total"`
	Error    string `json:"error,omitempty"`
}
