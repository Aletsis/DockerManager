import type { ContainerInfo } from './container';

export interface ComposeStackGroup {
  name: string;
  workingDir?: string;
  configFile?: string;
  containers: ContainerInfo[];
  runningCount: number;
  totalCount: number;
}

export interface ComposeDeployRequest {
  projectName: string;
  workingDir?: string;
  configFile?: string;
  content?: string;
  removeOrphans?: boolean;
}

export interface ComposeDownRequest {
  projectName: string;
  workingDir?: string;
  configFile?: string;
  removeVolumes?: boolean;
}

export interface DockerfileInfo {
  name: string;
  path: string;
  serviceName?: string;
  content: string;
  existsOnDisk: boolean;
}

export interface ComposeFileInfo {
  projectName: string;
  workingDir: string;
  configFile: string;
  content: string;
  existsOnDisk: boolean;
  dockerfiles: DockerfileInfo[];
}


export interface ComposePreset {
  id: string;
  name: string;
  description: string;
  defaultName: string;
  yaml: string;
}

