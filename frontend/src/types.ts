export interface PortMapping {
  ip: string;
  privatePort: number;
  publicPort: number;
  type: string;
}

export interface ContainerInfo {
  id: string;
  shortId: string;
  names: string[];
  name: string;
  image: string;
  imageId: string;
  command: string;
  created: number;
  state: 'running' | 'paused' | 'exited' | 'restarting' | 'dead' | string;
  status: string;
  ports: PortMapping[];
  sizeRw: number;
  sizeRootFs: number;
}

export interface SystemOverview {
  containers: number;
  containersRunning: number;
  containersPaused: number;
  containersStopped: number;
  images: number;
  serverVersion: string;
  operatingSystem: string;
  ncpu: number;
  memTotal: number;
}

export interface ContainerStats {
  id: string;
  name: string;
  cpuPercentage: number;
  memoryUsage: number;
  memoryLimit: number;
  memoryPercentage: number;
  networkRx: number;
  networkTx: number;
  blockRead: number;
  blockWrite: number;
  pids: number;
}

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

export interface CreateContainerRequest {
  image: string;
  name?: string;
  ports?: string[];
  volumes?: string[];
  env?: string[];
  restartPolicy?: string;
  autoStart?: boolean;
}

export interface CreateContainerResult {
  id: string;
  warnings?: string[];
}

