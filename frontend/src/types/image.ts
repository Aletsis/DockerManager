export interface ImageInfo {
  id: string;
  shortId: string;
  repository: string;
  tag: string;
  repoTags: string[];
  created: number;
  size: number;
  sharedSize: number;
  containers: number;
  inUse: boolean;
  isDangling: boolean;
}

export interface DiskUsageSummary {
  totalImages: number;
  totalSize: number;
  danglingCount: number;
  danglingSize: number;
  reclaimableSize: number;
}

export interface PruneResult {
  imagesDeleted: string[];
  spaceReclaimed: number;
}

export interface PullProgressEvent {
  id: string;
  status: string;
  progress: string;
  current: number;
  total: number;
  error?: string;
}
