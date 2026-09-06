export interface PortMapping {
  ip: string;
  privatePort: number;
  publicPort: number;
  type: string;
}

export interface ContainerNetworkInfo {
  networkName: string;
  networkId: string;
  ipAddress: string;
  gateway: string;
  macAddress: string;
  aliases?: string[];
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
  networks?: ContainerNetworkInfo[];
  sizeRw: number;
  sizeRootFs: number;
  labels?: Record<string, string>;
  composeProject?: string;
  composeService?: string;
  composeWorkingDir?: string;
  composeConfigFile?: string;
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
