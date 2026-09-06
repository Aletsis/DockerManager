<script lang="ts">
  import {
    Boxes,
    Sun,
    Moon,
    RefreshCw,
    Search,
    CheckCircle2,
    PauseCircle,
    StopCircle,
    Package,
    ChevronDown,
    HardDrive,
    Database,
    Network,
  } from '@lucide/svelte';
  import type { SystemOverview, DiskUsageSummary, VolumeDiskUsageSummary } from '../../../types';
  import { APP_VERSION } from '../../../shared/version';

  let {
    overview,
    diskUsage = null,
    volumeDiskUsage = null,
    volumesCount = 0,
    networksCount = 0,
    inactiveNetworksCount = 0,
    activeTab = $bindable<'containers' | 'images' | 'volumes' | 'networks'>('containers'),
    searchQuery = $bindable(''),
    isDark,
    onToggleTheme,
    refreshInterval = $bindable(2000),
    onManualRefresh,
    isRefreshing,
  } = $props<{
    overview: SystemOverview | null;
    diskUsage?: DiskUsageSummary | null;
    volumeDiskUsage?: VolumeDiskUsageSummary | null;
    volumesCount?: number;
    networksCount?: number;
    inactiveNetworksCount?: number;
    activeTab: 'containers' | 'images' | 'volumes' | 'networks';
    searchQuery: string;
    isDark: boolean;
    onToggleTheme: () => void;
    refreshInterval: number;
    onManualRefresh: () => void;
    isRefreshing: boolean;
  }>();
</script>

<!-- Snippet: Brand (Logo, App Name, Version & Subtitle) -->
{#snippet brand()}
  <div class="flex items-center gap-2.5 sm:gap-3 shrink-0">
    <div class="w-8 h-8 sm:w-9 sm:h-9 rounded-lg bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white shadow-sm shadow-indigo-500/20 shrink-0">
      <Boxes class="w-4 h-4 sm:w-5 sm:h-5 shrink-0" />
    </div>
    <div>
      <div class="flex items-center gap-1.5 sm:gap-2">
        <h1 class="font-semibold text-base sm:text-lg text-slate-900 dark:text-slate-100 tracking-tight leading-none">
          DockerManager
        </h1>
        <span class="text-[10px] sm:text-[11px] font-mono px-1.5 py-0.5 rounded bg-blue-50 dark:bg-blue-950/60 text-blue-600 dark:text-blue-400 border border-blue-200 dark:border-blue-800/60 font-medium leading-none">
          v{APP_VERSION}
        </span>
      </div>
      <p class="text-[11px] text-slate-500 dark:text-slate-400 hidden sm:block mt-0.5">
        Gestor de Contenedores y Recursos
      </p>
    </div>
  </div>
{/snippet}

<!-- Snippet: Navigation Tabs -->
{#snippet navTabs()}
  <nav class="flex items-center p-1 bg-slate-100 dark:bg-slate-800/80 rounded-xl border border-slate-200 dark:border-slate-700/60 text-xs shrink-0">
    <button
      onclick={() => (activeTab = 'containers')}
      class="flex items-center gap-1.5 px-2.5 sm:px-3 py-1.5 rounded-lg font-medium transition-all cursor-pointer {activeTab === 'containers' ? 'bg-white dark:bg-slate-900 text-blue-600 dark:text-blue-400 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-100'}"
    >
      <Boxes class="w-3.5 h-3.5 shrink-0" />
      <span>Contenedores</span>
      {#if overview}
        <span class="text-[10px] font-mono px-1.5 py-0.2 rounded-full {activeTab === 'containers' ? 'bg-blue-100 dark:bg-blue-900/50 text-blue-700 dark:text-blue-300' : 'bg-slate-200 dark:bg-slate-700 text-slate-600 dark:text-slate-400'}">
          {overview.containers}
        </span>
      {/if}
    </button>

    <button
      onclick={() => (activeTab = 'images')}
      class="flex items-center gap-1.5 px-2.5 sm:px-3 py-1.5 rounded-lg font-medium transition-all cursor-pointer {activeTab === 'images' ? 'bg-white dark:bg-slate-900 text-indigo-600 dark:text-indigo-400 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-100'}"
    >
      <HardDrive class="w-3.5 h-3.5 shrink-0" />
      <span>Imágenes<span class="hidden md:inline"> y Disco</span></span>
      {#if overview}
        <span class="text-[10px] font-mono px-1.5 py-0.2 rounded-full {activeTab === 'images' ? 'bg-indigo-100 dark:bg-indigo-900/50 text-indigo-700 dark:text-indigo-300' : 'bg-slate-200 dark:bg-slate-700 text-slate-600 dark:text-slate-400'}">
          {overview.images}
        </span>
      {/if}
      {#if diskUsage && diskUsage.danglingCount > 0}
        <span class="w-2 h-2 rounded-full bg-amber-500 animate-pulse shrink-0" title="{diskUsage.danglingCount} capas huérfanas"></span>
      {/if}
    </button>

    <button
      onclick={() => (activeTab = 'volumes')}
      class="flex items-center gap-1.5 px-2.5 sm:px-3 py-1.5 rounded-lg font-medium transition-all cursor-pointer {activeTab === 'volumes' ? 'bg-white dark:bg-slate-900 text-cyan-600 dark:text-cyan-400 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-100'}"
    >
      <Database class="w-3.5 h-3.5 shrink-0" />
      <span>Volúmenes</span>
      {#if volumesCount !== undefined && volumesCount > 0}
        <span class="text-[10px] font-mono px-1.5 py-0.2 rounded-full {activeTab === 'volumes' ? 'bg-cyan-100 dark:bg-cyan-900/50 text-cyan-700 dark:text-cyan-300' : 'bg-slate-200 dark:bg-slate-700 text-slate-600 dark:text-slate-400'}">
          {volumesCount}
        </span>
      {/if}
      {#if volumeDiskUsage && volumeDiskUsage.danglingCount > 0}
        <span class="w-2 h-2 rounded-full bg-amber-500 animate-pulse shrink-0" title="{volumeDiskUsage.danglingCount} volúmenes huérfanos"></span>
      {/if}
    </button>

    <button
      onclick={() => (activeTab = 'networks')}
      class="flex items-center gap-1.5 px-2.5 sm:px-3 py-1.5 rounded-lg font-medium transition-all cursor-pointer {activeTab === 'networks' ? 'bg-white dark:bg-slate-900 text-teal-600 dark:text-teal-400 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-100'}"
    >
      <Network class="w-3.5 h-3.5 shrink-0" />
      <span>Redes</span>
      {#if networksCount !== undefined && networksCount > 0}
        <span class="text-[10px] font-mono px-1.5 py-0.2 rounded-full {activeTab === 'networks' ? 'bg-teal-100 dark:bg-teal-900/50 text-teal-700 dark:text-teal-300' : 'bg-slate-200 dark:bg-slate-700 text-slate-600 dark:text-slate-400'}">
          {networksCount}
        </span>
      {/if}
      {#if inactiveNetworksCount !== undefined && inactiveNetworksCount > 0}
        <span class="w-2 h-2 rounded-full bg-amber-500 animate-pulse shrink-0" title="{inactiveNetworksCount} redes inactivas"></span>
      {/if}
    </button>
  </nav>
{/snippet}

<!-- Snippet: Quick Metrics Bar -->
{#snippet metrics()}
  {#if overview}
    <div class="flex items-center gap-1.5 sm:gap-2 text-xs overflow-x-auto no-scrollbar py-0.5 shrink-0 max-w-full">
      <div class="flex items-center gap-1.5 px-2 sm:px-2.5 py-1 rounded-full bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border border-emerald-500/20 font-medium whitespace-nowrap shrink-0">
        <CheckCircle2 class="w-3.5 h-3.5 shrink-0" />
        <span>{overview.containersRunning} <span class="hidden sm:inline">Activos</span></span>
      </div>
      {#if overview.containersPaused > 0}
        <div class="flex items-center gap-1.5 px-2 sm:px-2.5 py-1 rounded-full bg-amber-500/10 text-amber-700 dark:text-amber-400 border border-amber-500/20 font-medium whitespace-nowrap shrink-0">
          <PauseCircle class="w-3.5 h-3.5 shrink-0" />
          <span>{overview.containersPaused} <span class="hidden sm:inline">Pausados</span></span>
        </div>
      {/if}
      <div class="flex items-center gap-1.5 px-2 sm:px-2.5 py-1 rounded-full bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 border border-slate-200 dark:border-slate-700 font-medium whitespace-nowrap shrink-0">
        <StopCircle class="w-3.5 h-3.5 shrink-0" />
        <span>{overview.containersStopped} <span class="hidden sm:inline">Detenidos</span></span>
      </div>
      <div class="flex items-center gap-1.5 px-2 sm:px-2.5 py-1 rounded-full bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 border border-slate-200 dark:border-slate-700 font-medium whitespace-nowrap shrink-0">
        <Package class="w-3.5 h-3.5 shrink-0" />
        <span>{overview.images} <span class="hidden sm:inline">Imágenes</span></span>
      </div>
      {#if volumesCount > 0}
        <div class="flex items-center gap-1.5 px-2 sm:px-2.5 py-1 rounded-full bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 border border-slate-200 dark:border-slate-700 font-medium whitespace-nowrap shrink-0">
          <Database class="w-3.5 h-3.5 shrink-0" />
          <span>{volumesCount} <span class="hidden sm:inline">Volúmenes</span></span>
        </div>
      {/if}
      {#if networksCount > 0}
        <div class="flex items-center gap-1.5 px-2 sm:px-2.5 py-1 rounded-full bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 border border-slate-200 dark:border-slate-700 font-medium whitespace-nowrap shrink-0">
          <Network class="w-3.5 h-3.5 shrink-0" />
          <span>{networksCount} <span class="hidden sm:inline">Redes</span></span>
        </div>
      {/if}
    </div>
  {/if}
{/snippet}

<!-- Snippet: Action Controls (Search, Refresh Interval, Refresh Button, Theme Toggle) -->
{#snippet actions()}
  <div class="flex items-center gap-1.5 sm:gap-2 shrink-0">
    <!-- Desktop / Tablet Search for Containers -->
    {#if activeTab === 'containers'}
      <div class="relative hidden sm:block w-36 md:w-48 lg:w-56 animate-in fade-in duration-150 shrink-0">
        <Search class="w-4 h-4 absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" />
        <input
          type="text"
          placeholder="Buscar contenedor..."
          bind:value={searchQuery}
          class="w-full pl-8 pr-3 py-1.5 text-xs rounded-lg border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/80 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-1 focus:ring-indigo-500 transition-colors"
        />
      </div>
    {/if}

    <!-- Refresh Interval Selector -->
    <div class="relative shrink-0">
      <select
        aria-label="Frecuencia de actualización"
        bind:value={refreshInterval}
        class="appearance-none -webkit-appearance-none text-xs py-1.5 pl-2 sm:pl-2.5 pr-6 sm:pr-7 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-800 dark:text-slate-200 focus:outline-none focus:ring-1 focus:ring-indigo-500 cursor-pointer shadow-xs transition-colors"
        style="background-color: {isDark ? '#1e293b' : '#ffffff'}; color: {isDark ? '#e2e8f0' : '#1e293b'};"
      >
        <option value={1000} style="background-color: {isDark ? '#1e293b' : '#ffffff'}; color: {isDark ? '#e2e8f0' : '#1e293b'};">1s</option>
        <option value={2000} style="background-color: {isDark ? '#1e293b' : '#ffffff'}; color: {isDark ? '#e2e8f0' : '#1e293b'};">2s</option>
        <option value={5000} style="background-color: {isDark ? '#1e293b' : '#ffffff'}; color: {isDark ? '#e2e8f0' : '#1e293b'};">5s</option>
        <option value={10000} style="background-color: {isDark ? '#1e293b' : '#ffffff'}; color: {isDark ? '#e2e8f0' : '#1e293b'};">10s</option>
        <option value={0} style="background-color: {isDark ? '#1e293b' : '#ffffff'}; color: {isDark ? '#e2e8f0' : '#1e293b'};">Pausado</option>
      </select>
      <ChevronDown class="w-3.5 h-3.5 absolute right-1.5 sm:right-2 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" />
    </div>

    <!-- Manual Refresh Button -->
    <button
      onclick={onManualRefresh}
      title="Refrescar datos ahora"
      class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300 transition-colors shrink-0 cursor-pointer"
    >
      <RefreshCw class="w-4 h-4 {isRefreshing ? 'animate-spin text-indigo-500' : ''}" />
    </button>

    <!-- Theme Toggle (Dark / Light) -->
    <button
      onclick={onToggleTheme}
      title={isDark ? 'Cambiar a modo claro' : 'Cambiar a modo oscuro'}
      class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300 transition-colors shrink-0 cursor-pointer"
    >
      {#if isDark}
        <Sun class="w-4 h-4 text-amber-400" />
      {:else}
        <Moon class="w-4 h-4 text-slate-600" />
      {/if}
    </button>
  </div>
{/snippet}

<header class="border-b border-slate-200 dark:border-slate-800 bg-white/80 dark:bg-slate-900/80 backdrop-blur sticky top-0 z-30 transition-colors duration-200 shrink-0">
  <div class="max-w-7xl mx-auto px-3 sm:px-6 py-2 sm:py-2.5">
    <!-- Desktop Layout (>= 1280px / xl): Single Row -->
    <div class="hidden xl:flex items-center justify-between gap-4">
      <div class="flex items-center gap-6 shrink-0">
        {@render brand()}
        {@render navTabs()}
      </div>

      <div class="flex items-center justify-center shrink-0">
        {@render metrics()}
      </div>

      <div class="flex items-center justify-end shrink-0">
        {@render actions()}
      </div>
    </div>

    <!-- Responsive Layout (< 1280px / xl): Multi-tier -->
    <div class="flex xl:hidden flex-col gap-2">
      <!-- Tier 1: Brand on Left, Actions on Right -->
      <div class="flex items-center justify-between gap-2.5">
        {@render brand()}
        {@render actions()}
      </div>

      <!-- Tier 2: Tabs on Left, Metrics on Right (or stacked on mobile) -->
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 pt-1.5 border-t border-slate-100 dark:border-slate-800/60">
        <div class="flex items-center justify-between sm:justify-start gap-2">
          {@render navTabs()}
        </div>
        <div class="flex items-center sm:justify-end min-w-0">
          {@render metrics()}
        </div>
      </div>

      <!-- Mobile Search Bar (< sm only, when activeTab is containers) -->
      {#if activeTab === 'containers'}
        <div class="sm:hidden relative w-full pt-0.5">
          <Search class="w-4 h-4 absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" />
          <input
            type="text"
            placeholder="Buscar contenedor..."
            bind:value={searchQuery}
            class="w-full pl-8 pr-3 py-1.5 text-xs rounded-lg border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/80 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-1 focus:ring-indigo-500 transition-colors"
          />
        </div>
      {/if}
    </div>
  </div>
</header>

