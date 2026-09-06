import type { ContainerInfo } from './container';

export interface DockerNetworkGroup {
  name: string;
  networkId?: string;
  isDefault?: boolean;
  containers: ContainerInfo[];
  runningCount: number;
  totalCount: number;
}
