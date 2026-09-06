import * as WailsApp from '../../wailsjs/go/main/App';
import type {
  ContainerInfo,
  SystemOverview,
  ContainerStats,
  ImageInfo,
  DiskUsageSummary,
  PruneResult,
  CreateContainerRequest,
  CreateContainerResult,
} from '../types';

/**
 * Strongly-typed Docker API client service.
 * Eliminates repeated type casting across stores and components.
 */
export const dockerApi = {
  // System
  getOverview: async (): Promise<SystemOverview> => {
    const res = await WailsApp.GetOverview();
    return res as unknown as SystemOverview;
  },

  // Containers
  listContainers: async (all = true): Promise<ContainerInfo[]> => {
    const res = await WailsApp.ListContainers(all);
    return (res || []) as unknown as ContainerInfo[];
  },
  getContainerStats: async (id: string): Promise<ContainerStats> => {
    const res = await WailsApp.GetContainerStats(id);
    return res as unknown as ContainerStats;
  },
  getContainerLogs: async (id: string, tail = 200): Promise<string> => {
    return WailsApp.GetContainerLogs(id, tail);
  },
  startContainer: (id: string): Promise<void> => WailsApp.StartContainer(id),
  stopContainer: (id: string): Promise<void> => WailsApp.StopContainer(id),
  restartContainer: (id: string): Promise<void> => WailsApp.RestartContainer(id),
  pauseContainer: (id: string): Promise<void> => WailsApp.PauseContainer(id),
  unpauseContainer: (id: string): Promise<void> => WailsApp.UnpauseContainer(id),
  removeContainer: (id: string, force = true): Promise<void> => WailsApp.RemoveContainer(id, force),
  createContainer: async (req: CreateContainerRequest): Promise<CreateContainerResult> => {
    const res = await WailsApp.CreateContainer(req as any);
    return res as unknown as CreateContainerResult;
  },

  // Stacks
  startStack: (name: string): Promise<void> => WailsApp.StartStack(name),
  stopStack: (name: string): Promise<void> => WailsApp.StopStack(name),
  restartStack: (name: string): Promise<void> => WailsApp.RestartStack(name),

  // Networks
  startNetwork: (name: string): Promise<void> => WailsApp.StartNetwork(name),
  stopNetwork: (name: string): Promise<void> => WailsApp.StopNetwork(name),
  restartNetwork: (name: string): Promise<void> => WailsApp.RestartNetwork(name),

  // Images
  listImages: async (): Promise<ImageInfo[]> => {
    const res = await WailsApp.ListImages();
    return (res || []) as unknown as ImageInfo[];
  },
  getDiskUsage: async (): Promise<DiskUsageSummary> => {
    const res = await WailsApp.GetDiskUsage();
    return res as unknown as DiskUsageSummary;
  },
  removeImage: (id: string, force = false): Promise<void> => WailsApp.RemoveImage(id, force),
  pruneImages: async (danglingOnly = true): Promise<PruneResult> => {
    const res = await WailsApp.PruneImages(danglingOnly);
    return res as unknown as PruneResult;
  },
  pullImage: (name: string): Promise<void> => WailsApp.PullImage(name),
};
