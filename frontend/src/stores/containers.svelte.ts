import type {
  ContainerInfo,
  SystemOverview,
  ContainerStats,
  ComposeStackGroup,
  DockerNetworkGroup,
} from '../types';
import {
  ListContainers,
  GetOverview,
  StartContainer,
  StopContainer,
  RestartContainer,
  PauseContainer,
  UnpauseContainer,
  RemoveContainer,
  GetContainerStats,
  StartStack,
  StopStack,
  RestartStack,
  StartNetwork,
  StopNetwork,
  RestartNetwork,
} from '../../wailsjs/go/main/App';
import { uiStore } from './ui.svelte';
import { imagesStore } from './images.svelte';

export type GroupMode = 'flat' | 'stack' | 'network';
export type FilterState = 'all' | 'running' | 'paused' | 'stopped';

class ContainersStore {
  containers = $state<ContainerInfo[]>([]);
  overview = $state<SystemOverview | null>(null);
  statsMap = $state<Record<string, ContainerStats>>({});
  loading = $state<boolean>(true);
  isRefreshing = $state<boolean>(false);
  error = $state<string | null>(null);
  actionLoading = $state<string>('');

  filterState = $state<FilterState>('all');
  refreshInterval = $state<number>(2000);
  groupMode = $state<GroupMode>((() => {
    if (typeof localStorage === 'undefined') return 'stack';
    const saved = localStorage.getItem('containerGroupMode') as GroupMode | null;
    if (saved === 'flat' || saved === 'stack' || saved === 'network') {
      return saved;
    }
    const legacyStack = localStorage.getItem('groupByStack');
    if (legacyStack !== null) {
      return legacyStack === 'true' ? 'stack' : 'flat';
    }
    return 'stack';
  })());

  private pollTimer: any = null;

  setGroupMode(mode: GroupMode) {
    this.groupMode = mode;
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem('containerGroupMode', mode);
      localStorage.setItem('groupByStack', String(mode === 'stack'));
    }
  }

  setFilterState(filter: FilterState) {
    this.filterState = filter;
  }

  setRefreshInterval(ms: number) {
    this.refreshInterval = ms;
    this.setupPolling();
  }

  async fetchData(manual = false) {
    if (manual) this.isRefreshing = true;
    try {
      const [containerList, hostOverview] = await Promise.all([
        ListContainers(true),
        GetOverview(),
      ]);

      this.containers = (containerList || []) as unknown as ContainerInfo[];
      this.overview = hostOverview as unknown as SystemOverview;
      this.error = null;

      // Also refresh images in background if manual
      if (manual) {
        imagesStore.fetchImages();
      }

      // Query live stats for running containers
      if (containerList && containerList.length > 0) {
        const running = containerList.filter((c: any) => c.state === 'running');
        const statsPromises = running.slice(0, 8).map(async (c: any) => {
          try {
            const s = await GetContainerStats(c.id);
            return { id: c.id, stats: s as unknown as ContainerStats };
          } catch {
            return null;
          }
        });

        const results = await Promise.all(statsPromises);
        const newMap = { ...this.statsMap };
        results.forEach((r) => {
          if (r) newMap[r.id] = r.stats;
        });
        this.statsMap = newMap;
      }
    } catch (err: any) {
      this.error = err?.toString() || 'Error al conectar con el daemon de Docker';
    } finally {
      this.loading = false;
      if (manual) this.isRefreshing = false;
    }
  }

  setupPolling() {
    if (this.pollTimer) {
      clearInterval(this.pollTimer);
      this.pollTimer = null;
    }
    if (this.refreshInterval > 0) {
      this.pollTimer = setInterval(() => {
        this.fetchData();
      }, this.refreshInterval);
    }
  }

  stopPolling() {
    if (this.pollTimer) {
      clearInterval(this.pollTimer);
      this.pollTimer = null;
    }
  }

  // Filtered containers based on search query and state filter
  get filteredContainers(): ContainerInfo[] {
    const q = uiStore.searchQuery.trim().toLowerCase();
    const filter = this.filterState;

    return this.containers.filter((c) => {
      if (filter === 'running' && c.state !== 'running') return false;
      if (filter === 'paused' && c.state !== 'paused') return false;
      if (filter === 'stopped' && c.state === 'running') return false;

      if (q) {
        const matchName = c.name?.toLowerCase().includes(q);
        const matchImage = c.image?.toLowerCase().includes(q);
        const matchId = c.shortId?.toLowerCase().includes(q);
        const matchProject = c.composeProject?.toLowerCase().includes(q);
        const matchService = c.composeService?.toLowerCase().includes(q);
        const matchNetwork = c.networks?.some(
          (n) => n.networkName?.toLowerCase().includes(q) || n.ipAddress?.toLowerCase().includes(q)
        );
        return matchName || matchImage || matchId || !!matchProject || !!matchService || !!matchNetwork;
      }
      return true;
    });
  }

  // Grouped by Docker Compose stacks
  get stackGroupsData(): { stackGroups: ComposeStackGroup[]; standaloneContainers: ContainerInfo[] } {
    const groupsMap = new Map<string, ComposeStackGroup>();
    const standalone: ContainerInfo[] = [];

    for (const c of this.filteredContainers) {
      if (c.composeProject) {
        let group = groupsMap.get(c.composeProject);
        if (!group) {
          group = {
            name: c.composeProject,
            workingDir: c.composeWorkingDir,
            configFile: c.composeConfigFile,
            containers: [],
            runningCount: 0,
            totalCount: 0,
          };
          groupsMap.set(c.composeProject, group);
        }
        group.containers.push(c);
        group.totalCount++;
        if (c.state === 'running') {
          group.runningCount++;
        }
      } else {
        standalone.push(c);
      }
    }

    return {
      stackGroups: Array.from(groupsMap.values()),
      standaloneContainers: standalone,
    };
  }

  // Grouped by Docker networks
  get networkGroupsData(): { networkGroups: DockerNetworkGroup[]; isolatedContainers: ContainerInfo[] } {
    const groupsMap = new Map<string, DockerNetworkGroup>();
    const isolated: ContainerInfo[] = [];

    for (const c of this.filteredContainers) {
      if (c.networks && c.networks.length > 0) {
        for (const net of c.networks) {
          let group = groupsMap.get(net.networkName);
          if (!group) {
            const isDefault = ['bridge', 'host', 'none'].includes(net.networkName);
            group = {
              name: net.networkName,
              networkId: net.networkId,
              isDefault,
              containers: [],
              runningCount: 0,
              totalCount: 0,
            };
            groupsMap.set(net.networkName, group);
          }
          group.containers.push(c);
          group.totalCount++;
          if (c.state === 'running') {
            group.runningCount++;
          }
        }
      } else {
        isolated.push(c);
      }
    }

    const sortedGroups = Array.from(groupsMap.values()).sort((a, b) => {
      if (a.isDefault && !b.isDefault) return 1;
      if (!a.isDefault && b.isDefault) return -1;
      return a.name.localeCompare(b.name);
    });

    return {
      networkGroups: sortedGroups,
      isolatedContainers: isolated,
    };
  }

  // Container Actions
  async handleStart(id: string) {
    this.actionLoading = id;
    try {
      await StartContainer(id);
      uiStore.showToast('Contenedor iniciado exitosamente');
      await this.fetchData();
    } catch (err: any) {
      uiStore.showToast(`Error al iniciar: ${err}`, 'error');
    } finally {
      this.actionLoading = '';
    }
  }

  async handleStop(id: string) {
    this.actionLoading = id;
    try {
      await StopContainer(id);
      uiStore.showToast('Contenedor detenido');
      await this.fetchData();
    } catch (err: any) {
      uiStore.showToast(`Error al detener: ${err}`, 'error');
    } finally {
      this.actionLoading = '';
    }
  }

  async handleRestart(id: string) {
    this.actionLoading = id;
    try {
      await RestartContainer(id);
      uiStore.showToast('Contenedor reiniciado');
      await this.fetchData();
    } catch (err: any) {
      uiStore.showToast(`Error al reiniciar: ${err}`, 'error');
    } finally {
      this.actionLoading = '';
    }
  }

  async handlePause(id: string) {
    this.actionLoading = id;
    try {
      await PauseContainer(id);
      uiStore.showToast('Contenedor pausado');
      await this.fetchData();
    } catch (err: any) {
      uiStore.showToast(`Error al pausar: ${err}`, 'error');
    } finally {
      this.actionLoading = '';
    }
  }

  async handleUnpause(id: string) {
    this.actionLoading = id;
    try {
      await UnpauseContainer(id);
      uiStore.showToast('Contenedor reanudado');
      await this.fetchData();
    } catch (err: any) {
      uiStore.showToast(`Error al reanudar: ${err}`, 'error');
    } finally {
      this.actionLoading = '';
    }
  }

  async handleRemove(id: string) {
    this.actionLoading = id;
    try {
      await RemoveContainer(id, true);
      uiStore.showToast('Contenedor eliminado');
      await this.fetchData();
    } catch (err: any) {
      uiStore.showToast(`Error al eliminar: ${err}`, 'error');
    } finally {
      this.actionLoading = '';
    }
  }

  // Stack Actions
  async handleStartStack(projectName: string) {
    this.actionLoading = `stack:${projectName}`;
    try {
      await StartStack(projectName);
      uiStore.showToast(`Stack "${projectName}" iniciado exitosamente`);
      await this.fetchData();
    } catch (err: any) {
      uiStore.showToast(`Error al iniciar stack: ${err}`, 'error');
    } finally {
      this.actionLoading = '';
    }
  }

  async handleStopStack(projectName: string) {
    this.actionLoading = `stack:${projectName}`;
    try {
      await StopStack(projectName);
      uiStore.showToast(`Stack "${projectName}" detenido`);
      await this.fetchData();
    } catch (err: any) {
      uiStore.showToast(`Error al detener stack: ${err}`, 'error');
    } finally {
      this.actionLoading = '';
    }
  }

  async handleRestartStack(projectName: string) {
    this.actionLoading = `stack:${projectName}`;
    try {
      await RestartStack(projectName);
      uiStore.showToast(`Stack "${projectName}" reiniciado`);
      await this.fetchData();
    } catch (err: any) {
      uiStore.showToast(`Error al reiniciar stack: ${err}`, 'error');
    } finally {
      this.actionLoading = '';
    }
  }

  // Network Actions
  async handleStartNetwork(networkName: string) {
    this.actionLoading = `network:${networkName}`;
    try {
      await StartNetwork(networkName);
      uiStore.showToast(`Red "${networkName}" iniciada exitosamente`);
      await this.fetchData();
    } catch (err: any) {
      uiStore.showToast(`Error al iniciar red: ${err}`, 'error');
    } finally {
      this.actionLoading = '';
    }
  }

  async handleStopNetwork(networkName: string) {
    this.actionLoading = `network:${networkName}`;
    try {
      await StopNetwork(networkName);
      uiStore.showToast(`Red "${networkName}" detenida`);
      await this.fetchData();
    } catch (err: any) {
      uiStore.showToast(`Error al detener red: ${err}`, 'error');
    } finally {
      this.actionLoading = '';
    }
  }

  async handleRestartNetwork(networkName: string) {
    this.actionLoading = `network:${networkName}`;
    try {
      await RestartNetwork(networkName);
      uiStore.showToast(`Red "${networkName}" reiniciada`);
      await this.fetchData();
    } catch (err: any) {
      uiStore.showToast(`Error al reiniciar red: ${err}`, 'error');
    } finally {
      this.actionLoading = '';
    }
  }
}

export const containersStore = new ContainersStore();
