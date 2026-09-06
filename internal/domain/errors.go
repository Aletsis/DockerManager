package domain

import "errors"

var (
	ErrContainerNotFound       = errors.New("container not found")
	ErrContainerAlreadyRunning = errors.New("container is already running")
	ErrContainerNotRunning     = errors.New("container is not running")
	ErrContainerPaused         = errors.New("container is paused")
	ErrContainerNotPaused      = errors.New("container is not paused")
	ErrInvalidContainerSpec    = errors.New("invalid container specification")
	ErrCannotRemoveRunning     = errors.New("cannot remove a running container without force")
	ErrImageNotFound           = errors.New("image not found")
	ErrImageInUse              = errors.New("image is in use by a container")
	ErrNetworkNotFound         = errors.New("network not found")
	ErrStackNotFound           = errors.New("compose stack not found")
	ErrSessionNotFound         = errors.New("terminal session not found")
	ErrVolumeNotFound          = errors.New("volume not found")
	ErrVolumeInUse             = errors.New("volume is in use by a container")
)
