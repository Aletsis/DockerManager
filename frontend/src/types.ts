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
