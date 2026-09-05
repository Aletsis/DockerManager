<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import Header from './components/Header.svelte';
  import ContainerCard from './components/ContainerCard.svelte';
  import LogsModal from './components/LogsModal.svelte';
  import StatsModal from './components/StatsModal.svelte';
  import ConfirmModal from './components/ConfirmModal.svelte';
  import TerminalModal from './components/TerminalModal.svelte';
  import ImagesView from './components/ImagesView.svelte';
  import type {
    ContainerInfo,
    SystemOverview,
    ContainerStats,
    ImageInfo,
    DiskUsageSummary,
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
  } from '../wailsjs/go/main/App';
  import { AlertCircle, Box } from '@lucide/svelte';

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


  // Filters & Search
  let searchQuery = $state<string>('');
  let filterState = $state<'all' | 'running' | 'paused' | 'stopped'>('all');
  let refreshInterval = $state<number>(2000);

  // Modals & Actions
  let activeTerminal = $state<{ id: string; name: string } | null>(null);
  let activeLogs = $state<{ id: string; name: string } | null>(null);
  let activeStats = $state<{ id: string; name: string } | null>(null);
  let confirmDelete = $state<{ id: string; name: string } | null>(null);
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
        return matchName || matchImage || matchId;
      }
      return true;
    })
  );
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

        <div class="text-xs text-slate-500 dark:text-slate-400 font-mono">
          {filteredContainers.length} de {containers.length} contenedores
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
        </div>
      {/if}

      <!-- Containers List -->
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
    </div>
    {:else}
      <ImagesView
        {images}
        {diskUsage}
        loading={imagesLoading}
        onRefresh={() => fetchData(true)}
        onRemoveImage={handleRemoveImage}
        onPrune={handlePrune}
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
