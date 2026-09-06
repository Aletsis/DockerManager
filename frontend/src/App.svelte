<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import Header from './components/Header.svelte';
  import ContainerCard from './components/ContainerCard.svelte';
  import ComposeStackCard from './components/ComposeStackCard.svelte';
  import NetworkGroupCard from './components/NetworkGroupCard.svelte';
  import LogsModal from './components/LogsModal.svelte';
  import StatsModal from './components/StatsModal.svelte';
  import ConfirmModal from './components/ConfirmModal.svelte';
  import TerminalModal from './components/TerminalModal.svelte';
  import ImagesView from './components/ImagesView.svelte';
  import CreateContainerModal from './components/CreateContainerModal.svelte';
  import type {
    ContainerInfo,
    SystemOverview,
    ContainerStats,
    ImageInfo,
    DiskUsageSummary,
    ComposeStackGroup,
    DockerNetworkGroup,
  } from './types';
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
    ListImages,
    GetDiskUsage,
    RemoveImage,
    PruneImages,
    StartStack,
    StopStack,
    RestartStack,
    StartNetwork,
    StopNetwork,
    RestartNetwork,
  } from '../wailsjs/go/main/App';
  import { AlertCircle, Box, Plus, Layers, List, Network } from '@lucide/svelte';

  // Theme State
  let isDark = $state<boolean>(() => {
    const saved = localStorage.getItem('theme');
    if (saved) return saved === 'dark';
    return window.matchMedia('(prefers-color-scheme: dark)').matches;
  });

  // Navigation State
  let activeTab = $state<'containers' | 'images'>('containers');

  // Data States
  let containers = $state<ContainerInfo[]>([]);
  let overview = $state<SystemOverview | null>(null);
  let statsMap = $state<Record<string, ContainerStats>>({});
  let images = $state<ImageInfo[]>([]);
  let diskUsage = $state<DiskUsageSummary | null>(null);
  let loading = $state<boolean>(true);
  let imagesLoading = $state<boolean>(true);
  let isRefreshing = $state<boolean>(false);
  let error = $state<string | null>(null);

  // Group Mode
  type GroupMode = 'flat' | 'stack' | 'network';

  // Filters & Search
  let searchQuery = $state<string>('');
  let filterState = $state<'all' | 'running' | 'paused' | 'stopped'>('all');
  let refreshInterval = $state<number>(2000);
  let groupMode = $state<GroupMode>(() => {
    const saved = localStorage.getItem('containerGroupMode') as GroupMode | null;
    if (saved === 'flat' || saved === 'stack' || saved === 'network') {
      return saved;
    }
    const legacyStack = localStorage.getItem('groupByStack');
    if (legacyStack !== null) {
      return legacyStack === 'true' ? 'stack' : 'flat';
    }
    return 'stack';
  });

  // Modals & Actions
  let activeTerminal = $state<{ id: string; name: string } | null>(null);
  let activeLogs = $state<{ id: string; name: string } | null>(null);
  let activeStats = $state<{ id: string; name: string } | null>(null);
  let confirmDelete = $state<{ id: string; name: string } | null>(null);
  let isCreateModalOpen = $state<boolean>(false);
  let createModalInitialImage = $state<string>('');
  let actionLoading = $state<string>('');
  let notification = $state<{ text: string; type: 'success' | 'error' } | null>(null);

  let pollTimer: any = null;

  // React to theme changes
  $effect(() => {
    if (isDark) {
      document.documentElement.classList.add('dark');
      localStorage.setItem('theme', 'dark');
    } else {
      document.documentElement.classList.remove('dark');
      localStorage.setItem('theme', 'light');
    }
  });

  function showToast(text: string, type: 'success' | 'error' = 'success') {
    notification = { text, type };
    setTimeout(() => (notification = null), 3000);
  }

  async function fetchData(manual = false) {
    if (manual) isRefreshing = true;
    try {
      const [containerList, hostOverview, imageList, diskUsageRes] = await Promise.all([
        ListContainers(true),
        GetOverview(),
        ListImages(),
        GetDiskUsage(),
      ]);

      containers = (containerList || []) as unknown as ContainerInfo[];
      overview = hostOverview as unknown as SystemOverview;
      images = (imageList || []) as unknown as ImageInfo[];
      diskUsage = diskUsageRes as unknown as DiskUsageSummary;
      error = null;

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
        const newMap = { ...statsMap };
        results.forEach((r) => {
          if (r) newMap[r.id] = r.stats;
        });
        statsMap = newMap;
      }
    } catch (err: any) {
      error = err?.toString() || 'Error al conectar con el daemon de Docker';
    } finally {
      loading = false;
      imagesLoading = false;
      if (manual) isRefreshing = false;
    }
  }


  function setupPolling() {
    if (pollTimer) clearInterval(pollTimer);
    if (refreshInterval > 0) {
      pollTimer = setInterval(() => {
        fetchData();
      }, refreshInterval);
    }
  }

  $effect(() => {
    refreshInterval;
    setupPolling();
  });

  onMount(() => {
    fetchData();
    setupPolling();
  });

  onDestroy(() => {
    if (pollTimer) clearInterval(pollTimer);
  });

  // Action handlers
  async function handleStart(id: string) {
    actionLoading = id;
    try {
      await StartContainer(id);
      showToast('Contenedor iniciado exitosamente');
      fetchData();
    } catch (err: any) {
      showToast(`Error al iniciar: ${err}`, 'error');
    } finally {
      actionLoading = '';
    }
  }

  async function handleStop(id: string) {
    actionLoading = id;
    try {
      await StopContainer(id);
      showToast('Contenedor detenido');
      fetchData();
    } catch (err: any) {
      showToast(`Error al detener: ${err}`, 'error');
    } finally {
      actionLoading = '';
    }
  }

  async function handleRestart(id: string) {
    actionLoading = id;
    try {
      await RestartContainer(id);
      showToast('Contenedor reiniciado');
      fetchData();
    } catch (err: any) {
      showToast(`Error al reiniciar: ${err}`, 'error');
    } finally {
      actionLoading = '';
    }
  }

  async function handlePause(id: string) {
    actionLoading = id;
    try {
      await PauseContainer(id);
      showToast('Contenedor pausado');
      fetchData();
    } catch (err: any) {
      showToast(`Error al pausar: ${err}`, 'error');
    } finally {
      actionLoading = '';
    }
  }

  async function handleUnpause(id: string) {
    actionLoading = id;
    try {
      await UnpauseContainer(id);
      showToast('Contenedor reanudado');
      fetchData();
    } catch (err: any) {
      showToast(`Error al reanudar: ${err}`, 'error');
    } finally {
      actionLoading = '';
    }
  }

  async function handleRemove() {
    if (!confirmDelete) return;
    const { id } = confirmDelete;
    actionLoading = id;
    try {
      await RemoveContainer(id, true);
      showToast('Contenedor eliminado');
      confirmDelete = null;
      fetchData();
    } catch (err: any) {
      showToast(`Error al eliminar: ${err}`, 'error');
    } finally {
      actionLoading = '';
    }
  }

  // Image action handlers
  async function handleRemoveImage(id: string, force: boolean) {
    actionLoading = id;
    try {
      await RemoveImage(id, force);
      showToast('Imagen eliminada correctamente');
      await fetchData();
    } catch (err: any) {
      showToast(`Error al eliminar imagen: ${err}`, 'error');
    } finally {
      actionLoading = '';
    }
  }

  async function handlePrune(danglingOnly: boolean) {
    try {
      const res = await PruneImages(danglingOnly);
      await fetchData();
      return res as unknown as { imagesDeleted: string[]; spaceReclaimed: number };
    } catch (err: any) {
      showToast(`Error al ejecutar limpieza: ${err}`, 'error');
      throw err;
    }
  }

  // Stack action handlers
  async function handleStartStack(projectName: string) {
    actionLoading = `stack:${projectName}`;
    try {
      await StartStack(projectName);
      showToast(`Stack "${projectName}" iniciado exitosamente`);
      await fetchData();
    } catch (err: any) {
      showToast(`Error al iniciar stack: ${err}`, 'error');
    } finally {
      actionLoading = '';
    }
  }

  async function handleStopStack(projectName: string) {
    actionLoading = `stack:${projectName}`;
    try {
      await StopStack(projectName);
      showToast(`Stack "${projectName}" detenido`);
      await fetchData();
    } catch (err: any) {
      showToast(`Error al detener stack: ${err}`, 'error');
    } finally {
      actionLoading = '';
    }
  }

  async function handleRestartStack(projectName: string) {
    actionLoading = `stack:${projectName}`;
    try {
      await RestartStack(projectName);
      showToast(`Stack "${projectName}" reiniciado`);
      await fetchData();
    } catch (err: any) {
      showToast(`Error al reiniciar stack: ${err}`, 'error');
    } finally {
      actionLoading = '';
    }
  }

  // Network action handlers
  async function handleStartNetwork(networkName: string) {
    actionLoading = `network:${networkName}`;
    try {
      await StartNetwork(networkName);
      showToast(`Red "${networkName}" iniciada exitosamente`);
      await fetchData();
    } catch (err: any) {
      showToast(`Error al iniciar red: ${err}`, 'error');
    } finally {
      actionLoading = '';
    }
  }

  async function handleStopNetwork(networkName: string) {
    actionLoading = `network:${networkName}`;
    try {
      await StopNetwork(networkName);
      showToast(`Red "${networkName}" detenida`);
      await fetchData();
    } catch (err: any) {
      showToast(`Error al detener red: ${err}`, 'error');
    } finally {
      actionLoading = '';
    }
  }

  async function handleRestartNetwork(networkName: string) {
    actionLoading = `network:${networkName}`;
    try {
      await RestartNetwork(networkName);
      showToast(`Red "${networkName}" reiniciada`);
      await fetchData();
    } catch (err: any) {
      showToast(`Error al reiniciar red: ${err}`, 'error');
    } finally {
      actionLoading = '';
    }
  }

  // Persist group mode preference
  $effect(() => {
    localStorage.setItem('containerGroupMode', groupMode);
    localStorage.setItem('groupByStack', String(groupMode === 'stack'));
  });

  // Filtered containers
  const filteredContainers = $derived(
    containers.filter((c) => {
      if (filterState === 'running' && c.state !== 'running') return false;
      if (filterState === 'paused' && c.state !== 'paused') return false;
      if (filterState === 'stopped' && c.state === 'running') return false;

      if (searchQuery.trim()) {
        const q = searchQuery.toLowerCase();
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
    })
  );

  // Grouped Compose Stacks and Standalone Containers
  const { stackGroups, standaloneContainers } = $derived.by(() => {
    const groupsMap = new Map<string, ComposeStackGroup>();
    const standalone: ContainerInfo[] = [];

    for (const c of filteredContainers) {
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
  });

  // Grouped Networks and Isolated Containers
  const { networkGroups, isolatedContainers } = $derived.by(() => {
    const groupsMap = new Map<string, DockerNetworkGroup>();
    const isolated: ContainerInfo[] = [];

    for (const c of filteredContainers) {
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

    // Sort: User/custom networks first alphabetically, default networks (bridge, host, none) at bottom
    const sortedGroups = Array.from(groupsMap.values()).sort((a, b) => {
      if (a.isDefault && !b.isDefault) return 1;
      if (!a.isDefault && b.isDefault) return -1;
      return a.name.localeCompare(b.name);
    });

    return {
      networkGroups: sortedGroups,
      isolatedContainers: isolated,
    };
  });
</script>

<div class="h-screen bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 flex flex-col transition-colors duration-200 overflow-hidden">
  <!-- Top Header -->
  <Header
    {overview}
    {diskUsage}
    bind:activeTab
    bind:searchQuery
    {isDark}
    onToggleTheme={() => (isDark = !isDark)}
    bind:refreshInterval
    onManualRefresh={() => fetchData(true)}
    {isRefreshing}
  />

  <!-- Main Content Area -->
  <main class="flex-1 overflow-y-auto px-4 sm:px-6 py-6">
    {#if activeTab === 'containers'}
      <div class="max-w-7xl w-full mx-auto space-y-4 pb-16">

      <!-- Filter Chips Bar -->
      <div class="flex items-center justify-between gap-2 flex-wrap">
        <div class="flex items-center gap-2 flex-wrap">
          <div class="flex items-center gap-1.5 p-1 bg-slate-200/60 dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 text-xs">
            <button
              onclick={() => (filterState = 'all')}
              class="px-3 py-1 rounded-lg font-medium transition-colors {filterState === 'all' ? 'bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200'}"
            >
              Todos ({containers.length})
            </button>
            <button
              onclick={() => (filterState = 'running')}
              class="px-3 py-1 rounded-lg font-medium transition-colors {filterState === 'running' ? 'bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200'}"
            >
              Activos ({overview?.containersRunning || 0})
            </button>
            <button
              onclick={() => (filterState = 'stopped')}
              class="px-3 py-1 rounded-lg font-medium transition-colors {filterState === 'stopped' ? 'bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200'}"
            >
              Detenidos ({overview?.containersStopped || 0})
            </button>
            <button
              onclick={() => (filterState = 'paused')}
              class="px-3 py-1 rounded-lg font-medium transition-colors {filterState === 'paused' ? 'bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200'}"
            >
              Pausados ({overview?.containersPaused || 0})
            </button>
          </div>

          <!-- Grouping Toggle -->
          <div class="flex items-center gap-1 p-1 bg-slate-200/60 dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 text-xs">
            <button
              onclick={() => (groupMode = 'flat')}
              title="Ver lista plana de contenedores"
              class="px-2.5 py-1 rounded-lg font-medium transition-colors flex items-center gap-1.5 cursor-pointer {groupMode === 'flat' ? 'bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200'}"
            >
              <List class="w-3.5 h-3.5" />
              <span class="hidden sm:inline">Lista Plana</span>
            </button>
            <button
              onclick={() => (groupMode = 'stack')}
              title="Agrupar contenedores por Docker Compose (Stacks)"
              class="px-2.5 py-1 rounded-lg font-medium transition-colors flex items-center gap-1.5 cursor-pointer {groupMode === 'stack' ? 'bg-white dark:bg-slate-800 text-violet-600 dark:text-violet-400 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200'}"
            >
              <Layers class="w-3.5 h-3.5" />
              <span class="hidden sm:inline">Por Stacks</span>
            </button>
            <button
              onclick={() => (groupMode = 'network')}
              title="Agrupar contenedores por Red Docker compartida"
              class="px-2.5 py-1 rounded-lg font-medium transition-colors flex items-center gap-1.5 cursor-pointer {groupMode === 'network' ? 'bg-white dark:bg-slate-800 text-teal-600 dark:text-teal-400 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200'}"
            >
              <Network class="w-3.5 h-3.5" />
              <span class="hidden sm:inline">Por Redes</span>
            </button>
          </div>
        </div>

        <div class="flex items-center gap-3">
          <div class="text-xs text-slate-500 dark:text-slate-400 font-mono">
            {filteredContainers.length} de {containers.length} contenedores
          </div>
          <button
            onclick={() => {
              createModalInitialImage = '';
              isCreateModalOpen = true;
            }}
            class="px-3 py-1.5 rounded-xl text-xs font-medium text-white bg-blue-600 hover:bg-blue-700 transition-colors flex items-center gap-1.5 shadow-sm shadow-blue-500/20 cursor-pointer"
          >
            <Plus class="w-3.5 h-3.5" />
            <span>Nuevo Contenedor</span>
          </button>
        </div>
      </div>

      <!-- Error banner -->
      {#if error}
        <div class="p-4 rounded-xl bg-rose-50 dark:bg-rose-950/30 border border-rose-200 dark:border-rose-900/50 flex items-start gap-3">
          <AlertCircle class="w-5 h-5 text-rose-600 dark:text-rose-400 flex-shrink-0 mt-0.5" />
          <div class="space-y-1">
            <h4 class="text-sm font-semibold text-rose-800 dark:text-rose-300">
              Fallo de conexión con Docker
            </h4>
            <p class="text-xs text-rose-600 dark:text-rose-400 leading-relaxed">
              {error}. Verifica que el servicio de Docker esté iniciado (`systemctl status docker`).
            </p>
          </div>
        </div>
      {/if}

      <!-- Loading Skeleton -->
      {#if loading}
        <div class="space-y-3">
          {#each [1, 2, 3] as n}
            <div class="h-24 rounded-xl border border-slate-200 dark:border-slate-800 bg-white/60 dark:bg-slate-900/60 animate-pulse"></div>
          {/each}
        </div>
      {/if}

      <!-- Empty State -->
      {#if !loading && filteredContainers.length === 0}
        <div class="p-12 text-center rounded-2xl border border-dashed border-slate-300 dark:border-slate-800 bg-white/40 dark:bg-slate-900/40 space-y-3">
          <div class="w-12 h-12 rounded-2xl bg-slate-100 dark:bg-slate-800 flex items-center justify-center mx-auto text-slate-400">
            <Box class="w-6 h-6" />
          </div>
          <h3 class="font-semibold text-sm text-slate-800 dark:text-slate-200">
            No se encontraron contenedores
          </h3>
          <p class="text-xs text-slate-500 max-w-sm mx-auto">
            {searchQuery
              ? 'Ningún contenedor coincide con el filtro de búsqueda actual.'
              : 'No hay contenedores registrados en el daemon de Docker.'}
          </p>
          {#if !searchQuery}
            <div class="pt-2">
              <button
                onclick={() => {
                  createModalInitialImage = '';
                  isCreateModalOpen = true;
                }}
                class="px-4 py-2 rounded-xl text-xs font-medium text-white bg-blue-600 hover:bg-blue-700 transition-colors inline-flex items-center gap-1.5 shadow-sm shadow-blue-500/20 cursor-pointer"
              >
                <Plus class="w-3.5 h-3.5" />
                <span>Crear Primer Contenedor</span>
              </button>
            </div>
          {/if}
        </div>
      {/if}

      <!-- Containers List -->
      {#if groupMode === 'network' && networkGroups.length > 0}
        <div class="space-y-4">
          <!-- Networks Groups -->
          {#each networkGroups as network (network.name)}
            <NetworkGroupCard
              {network}
              {statsMap}
              onStart={handleStart}
              onStop={handleStop}
              onRestart={handleRestart}
              onPause={handlePause}
              onUnpause={handleUnpause}
              onRemove={(id, name) => (confirmDelete = { id, name })}
              onOpenTerminal={(id, name) => (activeTerminal = { id, name })}
              onViewLogs={(id, name) => (activeLogs = { id, name })}
              onViewStats={(id, name) => (activeStats = { id, name })}
              onStartNetwork={handleStartNetwork}
              onStopNetwork={handleStopNetwork}
              onRestartNetwork={handleRestartNetwork}
              {actionLoading}
            />
          {/each}

          <!-- Isolated Containers (Containers without any network) -->
          {#if isolatedContainers.length > 0}
            <div class="pt-2">
              <div class="flex items-center gap-3 mb-3">
                <span class="text-xs font-bold uppercase tracking-wider text-slate-500 dark:text-slate-400">
                  Contenedores Sin Red / Aislados ({isolatedContainers.length})
                </span>
                <div class="h-px bg-slate-200 dark:bg-slate-800 flex-1"></div>
              </div>
              <div class="space-y-3">
                {#each isolatedContainers as container (container.id)}
                  <ContainerCard
                    {container}
                    stats={statsMap[container.id]}
                    onStart={handleStart}
                    onStop={handleStop}
                    onRestart={handleRestart}
                    onPause={handlePause}
                    onUnpause={handleUnpause}
                    onRemove={(id, name) => (confirmDelete = { id, name })}
                    onOpenTerminal={(id, name) => (activeTerminal = { id, name })}
                    onViewLogs={(id, name) => (activeLogs = { id, name })}
                    onViewStats={(id, name) => (activeStats = { id, name })}
                    {actionLoading}
                  />
                {/each}
              </div>
            </div>
          {/if}
        </div>
      {:else if groupMode === 'stack' && stackGroups.length > 0}
        <div class="space-y-4">
          <!-- Stacks Groups -->
          {#each stackGroups as stack (stack.name)}
            <ComposeStackCard
              {stack}
              {statsMap}
              onStart={handleStart}
              onStop={handleStop}
              onRestart={handleRestart}
              onPause={handlePause}
              onUnpause={handleUnpause}
              onRemove={(id, name) => (confirmDelete = { id, name })}
              onOpenTerminal={(id, name) => (activeTerminal = { id, name })}
              onViewLogs={(id, name) => (activeLogs = { id, name })}
              onViewStats={(id, name) => (activeStats = { id, name })}
              onStartStack={handleStartStack}
              onStopStack={handleStopStack}
              onRestartStack={handleRestartStack}
              {actionLoading}
            />
          {/each}

          <!-- Standalone Containers (Containers without Compose Stack) -->
          {#if standaloneContainers.length > 0}
            <div class="pt-2">
              <div class="flex items-center gap-3 mb-3">
                <span class="text-xs font-bold uppercase tracking-wider text-slate-500 dark:text-slate-400">
                  Contenedores Individuales ({standaloneContainers.length})
                </span>
                <div class="h-px bg-slate-200 dark:bg-slate-800 flex-1"></div>
              </div>
              <div class="space-y-3">
                {#each standaloneContainers as container (container.id)}
                  <ContainerCard
                    {container}
                    stats={statsMap[container.id]}
                    onStart={handleStart}
                    onStop={handleStop}
                    onRestart={handleRestart}
                    onPause={handlePause}
                    onUnpause={handleUnpause}
                    onRemove={(id, name) => (confirmDelete = { id, name })}
                    onOpenTerminal={(id, name) => (activeTerminal = { id, name })}
                    onViewLogs={(id, name) => (activeLogs = { id, name })}
                    onViewStats={(id, name) => (activeStats = { id, name })}
                    {actionLoading}
                  />
                {/each}
              </div>
            </div>
          {/if}
        </div>
      {:else}
        <!-- Flat Containers List -->
        <div class="space-y-3">
          {#each filteredContainers as container (container.id)}
            <ContainerCard
              {container}
              stats={statsMap[container.id]}
              onStart={handleStart}
              onStop={handleStop}
              onRestart={handleRestart}
              onPause={handlePause}
              onUnpause={handleUnpause}
              onRemove={(id, name) => (confirmDelete = { id, name })}
              onOpenTerminal={(id, name) => (activeTerminal = { id, name })}
              onViewLogs={(id, name) => (activeLogs = { id, name })}
              onViewStats={(id, name) => (activeStats = { id, name })}
              {actionLoading}
            />
          {/each}
        </div>
      {/if}
    </div>
    {:else}
      <ImagesView
        {images}
        {diskUsage}
        loading={imagesLoading}
        onRefresh={() => fetchData(true)}
        onRemoveImage={handleRemoveImage}
        onPrune={handlePrune}
        onDeployContainer={(tag) => {
          createModalInitialImage = tag;
          isCreateModalOpen = true;
          activeTab = 'containers';
        }}
        {actionLoading}
      />
    {/if}
  </main>


  <!-- Floating Notification Toast -->
  {#if notification}
    <div
      class="fixed bottom-5 right-5 z-50 px-4 py-2.5 rounded-xl shadow-lg border text-xs font-medium animate-in fade-in slide-in-from-bottom-2 duration-200 {notification.type === 'error' ? 'bg-rose-600 text-white border-rose-700' : 'bg-slate-900 text-white dark:bg-slate-100 dark:text-slate-900 border-slate-800 dark:border-slate-200'}"
    >
      {notification.text}
    </div>
  {/if}

  <!-- Modals -->
  <CreateContainerModal
    isOpen={isCreateModalOpen}
    initialImage={createModalInitialImage}
    localImages={images}
    onClose={() => {
      isCreateModalOpen = false;
      createModalInitialImage = '';
    }}
    onSuccess={() => {
      showToast('Contenedor creado e iniciado exitosamente');
      fetchData(true);
    }}
  />

  {#if activeTerminal}
    <TerminalModal
      containerId={activeTerminal.id}
      containerName={activeTerminal.name}
      onClose={() => (activeTerminal = null)}
    />
  {/if}

  <LogsModal
    containerId={activeLogs?.id || null}
    containerName={activeLogs?.name || ''}
    onClose={() => (activeLogs = null)}
  />

  <StatsModal
    containerId={activeStats?.id || null}
    containerName={activeStats?.name || ''}
    onClose={() => (activeStats = null)}
  />

  <ConfirmModal
    isOpen={!!confirmDelete}
    title="Eliminar Contenedor"
    message={`¿Estás seguro de que deseas eliminar permanentemente el contenedor "${confirmDelete?.name}"? Esta acción no se puede deshacer.`}
    confirmLabel="Eliminar Contenedor"
    onConfirm={handleRemove}
    onCancel={() => (confirmDelete = null)}
    isDestructive={true}
  />
</div>
