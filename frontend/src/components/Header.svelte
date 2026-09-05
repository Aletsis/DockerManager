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
  } from '@lucide/svelte';
  import type { SystemOverview, DiskUsageSummary } from '../types';

  let {
    overview,
    diskUsage = null,
    activeTab = $bindable<'containers' | 'images'>('containers'),
    searchQuery = $bindable(''),
    isDark,
    onToggleTheme,
    refreshInterval = $bindable(2000),
    onManualRefresh,
    isRefreshing,
  } = $props<{
    overview: SystemOverview | null;
    diskUsage?: DiskUsageSummary | null;
    activeTab: 'containers' | 'images';
    searchQuery: string;
    isDark: boolean;
    onToggleTheme: () => void;
    refreshInterval: number;
    onManualRefresh: () => void;
    isRefreshing: boolean;
  }>();
</script>

<header class="border-b border-slate-200 dark:border-slate-800 bg-white/80 dark:bg-slate-900/80 backdrop-blur sticky top-0 z-30 transition-colors duration-200">
  <div class="max-w-7xl mx-auto px-4 sm:px-6 py-2.5">
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-3">
      <!-- Left: Logo & Navigation Tabs -->
      <div class="flex items-center gap-6">
        <div class="flex items-center gap-3">
          <div class="w-9 h-9 rounded-lg bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white shadow-sm shadow-indigo-500/20">
            <Boxes class="w-5 h-5" />
          </div>
          <div>
            <div class="flex items-center gap-2">
              <h1 class="font-semibold text-lg text-slate-900 dark:text-slate-100 tracking-tight">
                DockerManager
              </h1>
              {#if overview?.serverVersion}
                <span class="text-[11px] font-mono px-1.5 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-slate-500 dark:text-slate-400 border border-slate-200 dark:border-slate-700">
                  v{overview.serverVersion}
                </span>
              {/if}
            </div>
            <p class="text-[11px] text-slate-500 dark:text-slate-400">
              Gestor de Contenedores y Recursos
            </p>
          </div>
        </div>

        <!-- Navigation Tabs Switcher -->
        <nav class="flex items-center p-1 bg-slate-100 dark:bg-slate-800/80 rounded-xl border border-slate-200 dark:border-slate-700/60 text-xs">
          <button
            onclick={() => (activeTab = 'containers')}
            class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg font-medium transition-all {activeTab === 'containers' ? 'bg-white dark:bg-slate-900 text-blue-600 dark:text-blue-400 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-100'}"
          >
            <Boxes class="w-3.5 h-3.5" />
            <span>Contenedores</span>
            {#if overview}
              <span class="text-[10px] font-mono px-1.5 py-0.2 rounded-full {activeTab === 'containers' ? 'bg-blue-100 dark:bg-blue-900/50 text-blue-700 dark:text-blue-300' : 'bg-slate-200 dark:bg-slate-700 text-slate-600 dark:text-slate-400'}">
                {overview.containers}
              </span>
            {/if}
          </button>

          <button
            onclick={() => (activeTab = 'images')}
            class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg font-medium transition-all {activeTab === 'images' ? 'bg-white dark:bg-slate-900 text-indigo-600 dark:text-indigo-400 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-100'}"
          >
            <HardDrive class="w-3.5 h-3.5" />
            <span>Imágenes y Disco</span>
            {#if overview}
              <span class="text-[10px] font-mono px-1.5 py-0.2 rounded-full {activeTab === 'images' ? 'bg-indigo-100 dark:bg-indigo-900/50 text-indigo-700 dark:text-indigo-300' : 'bg-slate-200 dark:bg-slate-700 text-slate-600 dark:text-slate-400'}">
                {overview.images}
              </span>
            {/if}
            {#if diskUsage && diskUsage.danglingCount > 0}
              <span class="w-2 h-2 rounded-full bg-amber-500 animate-pulse" title="{diskUsage.danglingCount} capas huérfanas"></span>
            {/if}
          </button>
        </nav>
      </div>


      <!-- Quick Metrics Bar -->
      {#if overview}
        <div class="flex flex-wrap items-center gap-2 text-xs">
          <div class="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border border-emerald-500/20 font-medium">
            <CheckCircle2 class="w-3.5 h-3.5" />
            <span>{overview.containersRunning} Activos</span>
          </div>
          {#if overview.containersPaused > 0}
            <div class="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-amber-500/10 text-amber-700 dark:text-amber-400 border border-amber-500/20 font-medium">
              <PauseCircle class="w-3.5 h-3.5" />
              <span>{overview.containersPaused} Pausados</span>
            </div>
          {/if}
          <div class="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 border border-slate-200 dark:border-slate-700 font-medium">
            <StopCircle class="w-3.5 h-3.5" />
            <span>{overview.containersStopped} Detenidos</span>
          </div>
          <div class="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 border border-slate-200 dark:border-slate-700 font-medium">
            <Package class="w-3.5 h-3.5" />
            <span>{overview.images} Imágenes</span>
          </div>
        </div>
      {/if}

      <!-- Actions: Search, Refresh, Theme Toggle -->
      <div class="flex items-center gap-2.5">
        <!-- Search for Containers -->
        {#if activeTab === 'containers'}
          <div class="relative flex-1 md:w-56 animate-in fade-in duration-150">
            <Search class="w-4 h-4 absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" />
            <input
              type="text"
              placeholder="Buscar contenedor..."
              bind:value={searchQuery}
              class="w-full pl-8 pr-3 py-1.5 text-xs rounded-lg border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/80 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-1 focus:ring-indigo-500 transition-colors"
            />
          </div>
        {/if}


        <!-- Refresh Interval Selector with Theme-Proof Styling -->
        <div class="relative">
          <select
            aria-label="Frecuencia de actualización"
            bind:value={refreshInterval}
            class="appearance-none -webkit-appearance-none text-xs py-1.5 pl-2.5 pr-7 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-800 dark:text-slate-200 focus:outline-none focus:ring-1 focus:ring-indigo-500 cursor-pointer shadow-xs transition-colors"
            style="background-color: {isDark ? '#1e293b' : '#ffffff'}; color: {isDark ? '#e2e8f0' : '#1e293b'};"
          >
            <option value={1000} style="background-color: {isDark ? '#1e293b' : '#ffffff'}; color: {isDark ? '#e2e8f0' : '#1e293b'};">1s</option>
            <option value={2000} style="background-color: {isDark ? '#1e293b' : '#ffffff'}; color: {isDark ? '#e2e8f0' : '#1e293b'};">2s</option>
            <option value={5000} style="background-color: {isDark ? '#1e293b' : '#ffffff'}; color: {isDark ? '#e2e8f0' : '#1e293b'};">5s</option>
            <option value={10000} style="background-color: {isDark ? '#1e293b' : '#ffffff'}; color: {isDark ? '#e2e8f0' : '#1e293b'};">10s</option>
            <option value={0} style="background-color: {isDark ? '#1e293b' : '#ffffff'}; color: {isDark ? '#e2e8f0' : '#1e293b'};">Pausado</option>
          </select>
          <ChevronDown class="w-3.5 h-3.5 absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" />
        </div>

        <!-- Manual Refresh Button -->
        <button
          onclick={onManualRefresh}
          title="Refrescar datos ahora"
          class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300 transition-colors"
        >
          <RefreshCw class="w-4 h-4 {isRefreshing ? 'animate-spin text-indigo-500' : ''}" />
        </button>

        <!-- Theme Toggle (Dark / Light) -->
        <button
          onclick={onToggleTheme}
          title={isDark ? 'Cambiar a modo claro' : 'Cambiar a modo oscuro'}
          class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300 transition-colors"
        >
          {#if isDark}
            <Sun class="w-4 h-4 text-amber-400" />
          {:else}
            <Moon class="w-4 h-4 text-slate-600" />
          {/if}
        </button>
      </div>
    </div>
  </div>
</header>
