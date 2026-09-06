import type { ContainerInfo } from './container';

export interface NetworkContainerRef {
  id: string;
  name: string;
  state?: string;
  ipv4Address: string;
  ipv6Address: string;
  macAddress: string;
  endpointId: string;
}

export interface NetworkIPAMConfig {
  driver?: string;
  subnet?: string;
  gateway?: string;
  ipRange?: string;
}

export interface NetworkInfo {
  id: string;
  shortId: string;
  name: string;
  driver: string;
  scope: string;
  internal: boolean;
  attachable: boolean;
  enableIPv6: boolean;
  ipam: NetworkIPAMConfig[];
  containers: NetworkContainerRef[];
  containersCount: number;
  labels: Record<string, string>;
  options: Record<string, string>;
  created: string;
  isDefault: boolean;
}

export interface CreateNetworkRequest {
  name: string;
  driver: string;
  subnet?: string;
  gateway?: string;
  ipRange?: string;
  internal?: boolean;
  attachable?: boolean;
  enableIPv6?: boolean;
  labels?: Record<string, string>;
  options?: Record<string, string>;
}

export interface NetworkPruneResult {
  networksDeleted: string[];
}

export interface DockerNetworkGroup {
  name: string;
  networkId?: string;
  isDefault?: boolean;
  containers: ContainerInfo[];
  runningCount: number;
  totalCount: number;
}
