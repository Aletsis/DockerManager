<script lang="ts">
  import {
    Database,
    Trash2,
    HardDrive,
    Search,
    AlertTriangle,
    RefreshCw,
    CheckCircle2,
    Sparkles,
    Link2,
    Unlink,
  } from '@lucide/svelte';
  import type { VolumeInfo } from '../../../types';
  import { formatBytes } from '../../../shared/utils/utils';
  import ConfirmModal from '../../../shared/components/ConfirmModal.svelte';
  import VolumeCard from '../components/VolumeCard.svelte';
  import { volumesStore } from '../stores/volumes.svelte';
  import { uiStore } from '../../../shared/stores/ui.svelte';

  // Search & Filter State
  let searchQuery = $state('');
  let filterState = $state<'all' | 'inUse' | 'dangling'>('all');

  // Modals state
  let volumeToDelete = $state<VolumeInfo | null>(null);
  let isPruneModalOpen = $state(false);
  let isPruning = $state(false);
  let pruneSuccessInfo = $state<{ count: number; bytes: number } | null>(null);

  // Derived filtered volumes
  const filteredVolumes = $derived(
    volumesStore.volumes.filter((vol) => {
      if (filterState === 'inUse' && !vol.inUse) return false;
      if (filterState === 'dangling' && vol.inUse) return false;

      if (searchQuery.trim()) {
        const q = searchQuery.toLowerCase();
        const matchName = vol.name.toLowerCase().includes(q);
        const matchDriver = vol.driver.toLowerCase().includes(q);
        const matchContainer = vol.containers.some((c) =>
          c.name.toLowerCase().includes(q) || c.id.toLowerCase().includes(q) || c.destination.toLowerCase().includes(q)
        );
        return matchName || matchDriver || matchContainer;
      }
      return true;
    })
  );

  const inUseCount = $derived(volumesStore.volumes.filter((v) => v.inUse).length);
  const danglingCount = $derived(volumesStore.volumes.filter((v) => !v.inUse).length);

  async function handleConfirmDelete() {
    if (!volumeToDelete) return;
    const target = volumeToDelete;
    volumeToDelete = null;
    await volumesStore.removeVolume(target.name, target.inUse);
  }

  async function handleConfirmPrune() {
    isPruning = true;
    try {
      const result = await volumesStore.prune();
      isPruneModalOpen = false;
      if (result) {
        pruneSuccessInfo = {
          count: result.volumesDeleted?.length || 0,
          bytes: result.spaceReclaimed || 0,
        };
        setTimeout(() => {
          pruneSuccessInfo = null;
        }, 6000);
      }
    } finally {
      isPruning = false;
    }
  }

  function handleInspect(vol: VolumeInfo) {
    uiStore.openVolumeInspect(vol.name);
  }

  function handleRefresh() {
    volumesStore.fetchVolumes();
  }
</script>

<div class="max-w-7xl w-full mx-auto space-y-5 pb-16">
  <!-- Feedback Banner for Prune Success -->
  {#if pruneSuccessInfo}
    <div class="p-4 rounded-2xl bg-emerald-50 dark:bg-emerald-950/30 border border-emerald-200 dark:border-emerald-800 flex items-center justify-between gap-3 animate-in fade-in slide-in-from-top-2 duration-200 shadow-sm">
      <div class="flex items-center gap-3">
        <div class="p-2 rounded-xl bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
          <CheckCircle2 class="w-5 h-5" />
        </div>
        <div>
          <h4 class="text-xs font-semibold text-emerald-900 dark:text-emerald-200">
            ¡Limpieza de volúmenes completada!
          </h4>
          <p class="text-xs text-emerald-700 dark:text-emerald-300">
            Se liberaron <strong class="font-bold">{formatBytes(pruneSuccessInfo.bytes)}</strong> de almacenamiento y se eliminaron {pruneSuccessInfo.count} volúmenes huérfanos.
          </p>
        </div>
      </div>
      <button
        onclick={() => (pruneSuccessInfo = null)}
        class="text-xs font-medium px-3 py-1.5 rounded-lg border border-emerald-300 dark:border-emerald-700 text-emerald-800 dark:text-emerald-200 hover:bg-emerald-100 dark:hover:bg-emerald-900/40 transition-colors cursor-pointer"
      >
        Entendido
      </button>
    </div>
  {/if}

  <!-- Top Hero / Storage Metric Cards -->
  <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
    <!-- Card 1: Total Volume Storage -->
    <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl p-4 shadow-xs relative overflow-hidden">
      <div class="flex items-center justify-between">
        <span class="text-xs font-medium text-slate-500 dark:text-slate-400">Almacenamiento Total</span>
        <div class="p-2 rounded-xl bg-cyan-500/10 text-cyan-600 dark:text-cyan-400">
          <Database class="w-4 h-4" />
        </div>
      </div>
      <div class="mt-2 flex items-baseline gap-2">
        <span class="text-2xl font-bold tracking-tight text-slate-900 dark:text-slate-100">
          {formatBytes(volumesStore.diskUsage?.totalSize || 0)}
        </span>
        <span class="text-xs text-slate-500 dark:text-slate-400">
          en {volumesStore.volumes.length} volumen{volumesStore.volumes.length === 1 ? '' : 'es'}
        </span>
      </div>
      <div class="mt-3 flex items-center gap-1.5 text-xs text-slate-600 dark:text-slate-400">
        <HardDrive class="w-3.5 h-3.5 text-cyan-500 shrink-0" />
        <span>Gestionados por Docker Engine</span>
      </div>
    </div>

    <!-- Card 2: Reclaimable Space (Dangling Volumes) -->
    <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl p-4 shadow-xs relative overflow-hidden">
      <div class="flex items-center justify-between">
        <span class="text-xs font-medium text-slate-500 dark:text-slate-400">Espacio Huérfano Recuperable</span>
        <div class="p-2 rounded-xl bg-amber-500/10 text-amber-600 dark:text-amber-400">
          <AlertTriangle class="w-4 h-4" />
        </div>
      </div>
      <div class="mt-2 flex items-baseline gap-2">
        <span class="text-2xl font-bold tracking-tight {danglingCount > 0 ? 'text-amber-600 dark:text-amber-400' : 'text-slate-900 dark:text-slate-100'}">
          {formatBytes(volumesStore.diskUsage?.reclaimableSize || 0)}
        </span>
        <span class="text-xs text-slate-500 dark:text-slate-400">
          en {danglingCount} volumen{danglingCount === 1 ? '' : 'es'}
        </span>
      </div>
      <div class="mt-3 flex items-center justify-between gap-2">
        <span class="text-xs text-slate-500 dark:text-slate-400">Sin contenedores asociados</span>
        {#if danglingCount > 0}
          <button
            onclick={() => (isPruneModalOpen = true)}
            class="text-[11px] font-semibold text-amber-600 dark:text-amber-400 hover:underline flex items-center gap-1 cursor-pointer"
          >
            <Sparkles class="w-3 h-3" />
            Limpiar ahora
          </button>
        {/if}
      </div>
    </div>

    <!-- Card 3: Active In-Use Volumes -->
    <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl p-4 shadow-xs relative overflow-hidden">
      <div class="flex items-center justify-between">
        <span class="text-xs font-medium text-slate-500 dark:text-slate-400">Volúmenes en Uso</span>
        <div class="p-2 rounded-xl bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
          <Link2 class="w-4 h-4" />
        </div>
      </div>
      <div class="mt-2 flex items-baseline gap-2">
        <span class="text-2xl font-bold tracking-tight text-emerald-600 dark:text-emerald-400">
          {inUseCount}
        </span>
        <span class="text-xs text-slate-500 dark:text-slate-400">
          activos vinculados
        </span>
      </div>
      <div class="mt-3 flex items-center gap-1.5 text-xs text-emerald-600 dark:text-emerald-400 font-medium">
        <CheckCircle2 class="w-3.5 h-3.5 shrink-0" />
        <span>Almacenamiento persistente en uso</span>
      </div>
    </div>
  </div>

  <!-- Search, Filter & Action Toolbar -->
  <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl p-3 sm:p-4 shadow-xs flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
    <!-- Left: Search Box -->
    <div class="relative flex-1 max-w-md">
      <Search class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" />
      <input
        type="text"
        placeholder="Buscar por nombre, driver o contenedor..."
        bind:value={searchQuery}
        class="w-full pl-9 pr-3 py-1.5 text-xs rounded-xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/80 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-1 focus:ring-cyan-500 transition-colors"
      />
    </div>

    <!-- Center/Right: Filter Pills & Action Buttons -->
    <div class="flex items-center flex-wrap gap-2 justify-between sm:justify-end">
      <!-- Filter tabs -->
      <div class="flex items-center p-1 bg-slate-100 dark:bg-slate-800 rounded-xl text-xs font-medium">
        <button
          onclick={() => (filterState = 'all')}
          class="px-2.5 py-1 rounded-lg transition-colors cursor-pointer {filterState === 'all' ? 'bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100 shadow-2xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-100'}"
        >
          Todos ({volumesStore.volumes.length})
        </button>
        <button
          onclick={() => (filterState = 'inUse')}
          class="px-2.5 py-1 rounded-lg transition-colors cursor-pointer {filterState === 'inUse' ? 'bg-white dark:bg-slate-900 text-emerald-600 dark:text-emerald-400 shadow-2xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-100'}"
        >
          En uso ({inUseCount})
        </button>
        <button
          onclick={() => (filterState = 'dangling')}
          class="px-2.5 py-1 rounded-lg transition-colors cursor-pointer {filterState === 'dangling' ? 'bg-white dark:bg-slate-900 text-amber-600 dark:text-amber-400 shadow-2xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-100'}"
        >
          Huérfanos ({danglingCount})
        </button>
      </div>

      <!-- Refresh Button -->
      <button
        onclick={handleRefresh}
        title="Refrescar volúmenes"
        class="p-2 rounded-xl border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300 transition-colors cursor-pointer shrink-0"
      >
        <RefreshCw class="w-3.5 h-3.5 {volumesStore.loading ? 'animate-spin text-cyan-500' : ''}" />
      </button>

      <!-- Prune / Clean Dangling Volumes Button -->
      <button
        onclick={() => (isPruneModalOpen = true)}
        disabled={danglingCount === 0}
        class="flex items-center gap-1.5 px-3 py-1.5 rounded-xl text-xs font-medium transition-all shadow-xs cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed bg-amber-500 hover:bg-amber-600 text-white dark:bg-amber-600 dark:hover:bg-amber-500"
        title="Eliminar todos los volúmenes huérfanos sin contenedor"
      >
        <Trash2 class="w-3.5 h-3.5" />
        <span>Limpiar Huérfanos</span>
        {#if danglingCount > 0}
          <span class="px-1.5 py-0.2 rounded-full bg-white/20 text-[10px] font-mono">
            {danglingCount}
          </span>
        {/if}
      </button>
    </div>
  </div>

  <!-- Main Content: Volumes Grid or Empty State -->
  {#if volumesStore.loading && volumesStore.volumes.length === 0}
    <div class="h-64 flex flex-col items-center justify-center gap-3 text-slate-400">
      <RefreshCw class="w-7 h-7 animate-spin text-cyan-500" />
      <span class="text-xs">Cargando volúmenes de Docker...</span>
    </div>
  {:else if filteredVolumes.length === 0}
    <div class="p-12 text-center bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl shadow-xs">
      <div class="w-12 h-12 rounded-2xl bg-cyan-500/10 text-cyan-600 dark:text-cyan-400 flex items-center justify-center mx-auto mb-3">
        <Database class="w-6 h-6" />
      </div>
      <h3 class="text-sm font-semibold text-slate-900 dark:text-slate-100">
        No se encontraron volúmenes
      </h3>
      <p class="text-xs text-slate-500 dark:text-slate-400 mt-1 max-w-sm mx-auto">
        {#if searchQuery}
          No hay volúmenes que coincidan con los términos de búsqueda "{searchQuery}".
        {:else if filterState === 'dangling'}
          ¡Excelente! No tienes volúmenes huérfanos consumiendo espacio innecesario.
        {:else}
          No hay volúmenes locales creados en este host Docker.
        {/if}
      </p>
    </div>
  {:else}
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      {#each filteredVolumes as volume (volume.name)}
        <VolumeCard
          {volume}
          onInspect={handleInspect}
          onDelete={(vol) => (volumeToDelete = vol)}
          isDeleting={volumesStore.actionLoading === volume.name}
        />
      {/each}
    </div>
  {/if}
</div>

<!-- Modal: Confirm Single Volume Delete -->
<ConfirmModal
  isOpen={!!volumeToDelete}
  title="Eliminar Volumen"
  message={`¿Estás seguro de que deseas eliminar permanentemente el volumen "${volumeToDelete?.name}"? ${volumeToDelete?.inUse ? 'ADVERTENCIA: Este volumen está actualmente enlazado a uno o más contenedores.' : 'Esta acción liberará espacio en disco de forma irreversible.'}`}
  confirmLabel={volumeToDelete?.inUse ? 'Forzar Eliminación' : 'Eliminar Volumen'}
  onConfirm={handleConfirmDelete}
  onCancel={() => (volumeToDelete = null)}
  isDestructive={true}
/>

<!-- Modal: Confirm Prune All Dangling Volumes -->
<ConfirmModal
  isOpen={isPruneModalOpen}
  title="Limpiar Volúmenes Huérfanos"
  message={`Esta operación eliminará permanentemente todos los volúmenes (${danglingCount}) que no estén asociados a ningún contenedor. Se estima recuperar aproximadamente ${formatBytes(volumesStore.diskUsage?.reclaimableSize || 0)} de espacio en disco.`}
  confirmLabel={isPruning ? 'Limpiando...' : 'Confirmar Limpieza'}
  onConfirm={handleConfirmPrune}
  onCancel={() => (isPruneModalOpen = false)}
  isDestructive={true}
/>
