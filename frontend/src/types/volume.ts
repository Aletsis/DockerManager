export interface VolumeContainerRef {
  id: string;
  name: string;
  state: string;
  destination: string;
  rw: boolean;
}

export interface VolumeInfo {
  name: string;
  driver: string;
  mountpoint: string;
  createdAt: string;
  labels: Record<string, string>;
  scope: string;
  size: number;
  inUse: boolean;
  containers: VolumeContainerRef[];
}

export interface VolumeDiskUsageSummary {
  totalVolumes: number;
  totalSize: number;
  danglingCount: number;
  danglingSize: number;
  reclaimableSize: number;
}

export interface VolumePruneResult {
  volumesDeleted: string[];
  spaceReclaimed: number;
}
