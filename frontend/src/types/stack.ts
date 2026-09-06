import type { ContainerInfo } from './container';

export interface ComposeStackGroup {
  name: string;
  workingDir?: string;
  configFile?: string;
  containers: ContainerInfo[];
  runningCount: number;
  totalCount: number;
}
