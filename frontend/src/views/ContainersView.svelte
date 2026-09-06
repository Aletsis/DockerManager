<script lang="ts">
  import ContainerCard from '../features/containers/components/ContainerCard.svelte';
  import ComposeStackCard from '../features/containers/components/ComposeStackCard.svelte';
  import NetworkGroupCard from '../features/containers/components/NetworkGroupCard.svelte';
  import { containersStore } from '../features/containers/stores/containers.svelte';
  import { uiStore } from '../shared/stores/ui.svelte';
  import { AlertCircle, Box, Plus, Layers, List, Network } from '@lucide/svelte';

  const { stackGroups, standaloneContainers } = $derived(containersStore.stackGroupsData);
  const { networkGroups, isolatedContainers } = $derived(containersStore.networkGroupsData);
  const filteredContainers = $derived(containersStore.filteredContainers);
</script>

<div class="max-w-7xl w-full mx-auto space-y-4 pb-16">
  <!-- Filter Chips Bar & Grouping Switcher -->
  <div class="flex items-center justify-between gap-2 flex-wrap">
    <div class="flex items-center gap-2 flex-wrap">
      <div class="flex items-center gap-1.5 p-1 bg-slate-200/60 dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 text-xs">
        <button
          onclick={() => containersStore.setFilterState('all')}
          class="px-3 py-1 rounded-lg font-medium transition-colors cursor-pointer {containersStore.filterState === 'all' ? 'bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200'}"
        >
          Todos ({containersStore.containers.length})
        </button>
        <button
          onclick={() => containersStore.setFilterState('running')}
          class="px-3 py-1 rounded-lg font-medium transition-colors cursor-pointer {containersStore.filterState === 'running' ? 'bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200'}"
        >
          Activos ({containersStore.overview?.containersRunning || 0})
        </button>
        <button
          onclick={() => containersStore.setFilterState('stopped')}
          class="px-3 py-1 rounded-lg font-medium transition-colors cursor-pointer {containersStore.filterState === 'stopped' ? 'bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200'}"
        >
          Detenidos ({containersStore.overview?.containersStopped || 0})
        </button>
        <button
          onclick={() => containersStore.setFilterState('paused')}
          class="px-3 py-1 rounded-lg font-medium transition-colors cursor-pointer {containersStore.filterState === 'paused' ? 'bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200'}"
        >
          Pausados ({containersStore.overview?.containersPaused || 0})
        </button>
      </div>

      <!-- Grouping Mode Switcher -->
      <div class="flex items-center gap-1 p-1 bg-slate-200/60 dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 text-xs">
        <button
          onclick={() => containersStore.setGroupMode('flat')}
          title="Ver lista plana de contenedores"
          class="px-2.5 py-1 rounded-lg font-medium transition-colors flex items-center gap-1.5 cursor-pointer {containersStore.groupMode === 'flat' ? 'bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200'}"
        >
          <List class="w-3.5 h-3.5" />
          <span class="hidden sm:inline">Lista Plana</span>
        </button>
        <button
          onclick={() => containersStore.setGroupMode('stack')}
          title="Agrupar contenedores por Docker Compose (Stacks)"
          class="px-2.5 py-1 rounded-lg font-medium transition-colors flex items-center gap-1.5 cursor-pointer {containersStore.groupMode === 'stack' ? 'bg-white dark:bg-slate-800 text-violet-600 dark:text-violet-400 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200'}"
        >
          <Layers class="w-3.5 h-3.5" />
          <span class="hidden sm:inline">Por Stacks</span>
        </button>
        <button
          onclick={() => containersStore.setGroupMode('network')}
          title="Agrupar contenedores por Red Docker compartida"
          class="px-2.5 py-1 rounded-lg font-medium transition-colors flex items-center gap-1.5 cursor-pointer {containersStore.groupMode === 'network' ? 'bg-white dark:bg-slate-800 text-teal-600 dark:text-teal-400 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200'}"
        >
          <Network class="w-3.5 h-3.5" />
          <span class="hidden sm:inline">Por Redes</span>
        </button>
      </div>
    </div>

    <div class="flex items-center gap-3">
      <div class="text-xs text-slate-500 dark:text-slate-400 font-mono">
        {filteredContainers.length} de {containersStore.containers.length} contenedores
      </div>
      <button
        onclick={() => uiStore.openCreateModal()}
        class="px-3 py-1.5 rounded-xl text-xs font-medium text-white bg-blue-600 hover:bg-blue-700 transition-colors flex items-center gap-1.5 shadow-sm shadow-blue-500/20 cursor-pointer"
      >
        <Plus class="w-3.5 h-3.5" />
        <span>Nuevo Contenedor</span>
      </button>
    </div>
  </div>

  <!-- Error banner -->
  {#if containersStore.error}
    <div class="p-4 rounded-xl bg-rose-50 dark:bg-rose-950/30 border border-rose-200 dark:border-rose-900/50 flex items-start gap-3">
      <AlertCircle class="w-5 h-5 text-rose-600 dark:text-rose-400 flex-shrink-0 mt-0.5" />
      <div class="space-y-1">
        <h4 class="text-sm font-semibold text-rose-800 dark:text-rose-300">
          Fallo de conexión con Docker
        </h4>
        <p class="text-xs text-rose-600 dark:text-rose-400 leading-relaxed">
          {containersStore.error}. Verifica que el servicio de Docker esté iniciado (`systemctl status docker`).
        </p>
      </div>
    </div>
  {/if}

  <!-- Loading Skeleton -->
  {#if containersStore.loading}
    <div class="space-y-3">
      {#each [1, 2, 3] as _}
        <div class="h-24 rounded-xl border border-slate-200 dark:border-slate-800 bg-white/60 dark:bg-slate-900/60 animate-pulse"></div>
      {/each}
    </div>
  {/if}

  <!-- Empty State -->
  {#if !containersStore.loading && filteredContainers.length === 0}
    <div class="p-12 text-center rounded-2xl border border-dashed border-slate-300 dark:border-slate-800 bg-white/40 dark:bg-slate-900/40 space-y-3">
      <div class="w-12 h-12 rounded-2xl bg-slate-100 dark:bg-slate-800 flex items-center justify-center mx-auto text-slate-400">
        <Box class="w-6 h-6" />
      </div>
      <h3 class="font-semibold text-sm text-slate-800 dark:text-slate-200">
        No se encontraron contenedores
      </h3>
      <p class="text-xs text-slate-500 max-w-sm mx-auto">
        {uiStore.searchQuery
          ? 'Ningún contenedor coincide con el filtro de búsqueda actual.'
          : 'No hay contenedores registrados en el daemon de Docker.'}
      </p>
      {#if !uiStore.searchQuery}
        <div class="pt-2">
          <button
            onclick={() => uiStore.openCreateModal()}
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
  {#if containersStore.groupMode === 'network' && networkGroups.length > 0}
    <div class="space-y-4">
      <!-- Network Groups -->
      {#each networkGroups as network (network.name)}
        <NetworkGroupCard {network} />
      {/each}

      <!-- Isolated Containers -->
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
              <ContainerCard {container} />
            {/each}
          </div>
        </div>
      {/if}
    </div>
  {:else if containersStore.groupMode === 'stack' && stackGroups.length > 0}
    <div class="space-y-4">
      <!-- Stacks Groups -->
      {#each stackGroups as stack (stack.name)}
        <ComposeStackCard {stack} />
      {/each}

      <!-- Standalone Containers -->
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
              <ContainerCard {container} />
            {/each}
          </div>
        </div>
      {/if}
    </div>
  {:else}
    <!-- Flat Containers List -->
    <div class="space-y-3">
      {#each filteredContainers as container (container.id)}
        <ContainerCard {container} />
      {/each}
    </div>
  {/if}
</div>
