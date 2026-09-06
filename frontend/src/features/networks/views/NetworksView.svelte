<script lang="ts">
  import {
    Network,
    Plus,
    Trash2,
    Search,
    AlertTriangle,
    RefreshCw,
    CheckCircle2,
    Sparkles,
    Link2,
    Layers,
    Server,
  } from '@lucide/svelte';
  import type { NetworkInfo } from '../../../types';
  import ConfirmModal from '../../../shared/components/ConfirmModal.svelte';
  import NetworkCard from '../components/NetworkCard.svelte';
  import { networksStore } from '../stores/networks.svelte';
  import { uiStore } from '../../../shared/stores/ui.svelte';

  // Search & Filter State
  let searchQuery = $state('');
  let filterState = $state<'all' | 'custom' | 'default' | 'inactive'>('all');

  // Prune state
  let isPruneModalOpen = $state(false);
  let isPruning = $state(false);
  let pruneSuccessInfo = $state<{ count: number } | null>(null);

  // Derived calculations
  const inactiveNetworks = $derived(
    networksStore.networks.filter((n) => !n.isDefault && n.containersCount === 0)
  );

  const totalAttachments = $derived(
    networksStore.networks.reduce((acc, n) => acc + n.containersCount, 0)
  );

  const customNetworksCount = $derived(
    networksStore.networks.filter((n) => !n.isDefault).length
  );

  // Driver breakdown summary string
  const driversSummary = $derived.by(() => {
    const counts: Record<string, number> = {};
    for (const n of networksStore.networks) {
      counts[n.driver] = (counts[n.driver] || 0) + 1;
    }
    return Object.entries(counts)
      .map(([d, c]) => `${c} ${d}`)
      .join(', ');
  });

  // Filtered networks
  const filteredNetworks = $derived(
    networksStore.networks.filter((net) => {
      if (filterState === 'custom' && net.isDefault) return false;
      if (filterState === 'default' && !net.isDefault) return false;
      if (filterState === 'inactive' && (net.isDefault || net.containersCount > 0)) return false;

      if (searchQuery.trim()) {
        const q = searchQuery.toLowerCase();
        const matchName = net.name.toLowerCase().includes(q);
        const matchDriver = net.driver.toLowerCase().includes(q);
        const matchSubnet = net.ipam?.some(
          (ip) => ip.subnet?.toLowerCase().includes(q) || ip.gateway?.toLowerCase().includes(q)
        );
        const matchContainer = net.containers.some(
          (c) =>
            c.name.toLowerCase().includes(q) ||
            c.id.toLowerCase().includes(q) ||
            c.ipv4Address?.toLowerCase().includes(q)
        );
        return matchName || matchDriver || matchSubnet || matchContainer;
      }
      return true;
    })
  );

  async function handleConfirmPrune() {
    isPruning = true;
    try {
      const result = await networksStore.prune();
      isPruneModalOpen = false;
      if (result && result.networksDeleted) {
        pruneSuccessInfo = { count: result.networksDeleted.length };
        setTimeout(() => {
          pruneSuccessInfo = null;
        }, 6000);
      }
    } finally {
      isPruning = false;
    }
  }

  function handleRefresh() {
    networksStore.fetchNetworks();
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
            ¡Limpieza de redes completada!
          </h4>
          <p class="text-xs text-emerald-700 dark:text-emerald-300">
            Se eliminaron <strong>{pruneSuccessInfo.count}</strong> redes inactivas del sistema Docker.
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

  <!-- Top Hero / Network Metric Cards -->
  <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
    <!-- Card 1: Total Networks -->
    <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl p-4 shadow-xs relative overflow-hidden">
      <div class="flex items-center justify-between">
        <span class="text-xs font-medium text-slate-500 dark:text-slate-400">Total de Redes</span>
        <div class="p-2 rounded-xl bg-teal-500/10 text-teal-600 dark:text-teal-400">
          <Network class="w-4 h-4" />
        </div>
      </div>
      <div class="mt-2 flex items-baseline gap-2">
        <span class="text-2xl font-bold tracking-tight text-slate-900 dark:text-slate-100">
          {networksStore.networks.length}
        </span>
        <span class="text-xs text-slate-500 dark:text-slate-400">
          ({customNetworksCount} personalizadas)
        </span>
      </div>
      <div class="mt-3 flex items-center gap-1.5 text-xs text-slate-600 dark:text-slate-400 truncate">
        <Layers class="w-3.5 h-3.5 text-teal-500 shrink-0" />
        <span class="truncate">{driversSummary || 'bridge, host, none'}</span>
      </div>
    </div>

    <!-- Card 2: Inactive / Pruneable Networks -->
    <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl p-4 shadow-xs relative overflow-hidden">
      <div class="flex items-center justify-between">
        <span class="text-xs font-medium text-slate-500 dark:text-slate-400">Redes Inactivas (Prune)</span>
        <div class="p-2 rounded-xl bg-amber-500/10 text-amber-600 dark:text-amber-400">
          <AlertTriangle class="w-4 h-4" />
        </div>
      </div>
      <div class="mt-2 flex items-baseline gap-2">
        <span class="text-2xl font-bold tracking-tight {inactiveNetworks.length > 0 ? 'text-amber-600 dark:text-amber-400' : 'text-slate-900 dark:text-slate-100'}">
          {inactiveNetworks.length}
        </span>
        <span class="text-xs text-slate-500 dark:text-slate-400">
          sin contenedores activos
        </span>
      </div>
      <div class="mt-3 flex items-center justify-between gap-2">
        <span class="text-xs text-slate-500 dark:text-slate-400">Redes no predeterminadas</span>
        {#if inactiveNetworks.length > 0}
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

    <!-- Card 3: Connected Container Attachments -->
    <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl p-4 shadow-xs relative overflow-hidden">
      <div class="flex items-center justify-between">
        <span class="text-xs font-medium text-slate-500 dark:text-slate-400">Contenedores Conectados</span>
        <div class="p-2 rounded-xl bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
          <Link2 class="w-4 h-4" />
        </div>
      </div>
      <div class="mt-2 flex items-baseline gap-2">
        <span class="text-2xl font-bold tracking-tight text-emerald-600 dark:text-emerald-400">
          {totalAttachments}
        </span>
        <span class="text-xs text-slate-500 dark:text-slate-400">
          vínculos de red activos
        </span>
      </div>
      <div class="mt-3 flex items-center gap-1.5 text-xs text-emerald-600 dark:text-emerald-400 font-medium">
        <CheckCircle2 class="w-3.5 h-3.5 shrink-0" />
        <span>Enrutamiento y DNS interno activo</span>
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
        placeholder="Buscar por nombre, driver, subred o contenedor..."
        bind:value={searchQuery}
        class="w-full pl-9 pr-3 py-1.5 text-xs rounded-xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/80 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-1 focus:ring-teal-500 transition-colors"
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
          Todas ({networksStore.networks.length})
        </button>
        <button
          onclick={() => (filterState = 'custom')}
          class="px-2.5 py-1 rounded-lg transition-colors cursor-pointer {filterState === 'custom' ? 'bg-white dark:bg-slate-900 text-teal-600 dark:text-teal-400 shadow-2xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-100'}"
        >
          Personalizadas ({customNetworksCount})
        </button>
        <button
          onclick={() => (filterState = 'default')}
          class="px-2.5 py-1 rounded-lg transition-colors cursor-pointer {filterState === 'default' ? 'bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100 shadow-2xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-100'}"
        >
          Predeterminadas ({networksStore.networks.length - customNetworksCount})
        </button>
        <button
          onclick={() => (filterState = 'inactive')}
          class="px-2.5 py-1 rounded-lg transition-colors cursor-pointer {filterState === 'inactive' ? 'bg-white dark:bg-slate-900 text-amber-600 dark:text-amber-400 shadow-2xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-100'}"
        >
          Inactivas ({inactiveNetworks.length})
        </button>
      </div>

      <!-- Action Buttons -->
      <div class="flex items-center gap-2">
        <button
          onclick={handleRefresh}
          class="p-2 rounded-xl border border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-400 transition-colors cursor-pointer"
          title="Refrescar redes"
        >
          <RefreshCw class="w-4 h-4 {networksStore.loading ? 'animate-spin text-teal-500' : ''}" />
        </button>

        {#if inactiveNetworks.length > 0}
          <button
            onclick={() => (isPruneModalOpen = true)}
            class="flex items-center gap-1.5 px-3 py-1.5 rounded-xl border border-amber-300 dark:border-amber-700/60 bg-amber-50 dark:bg-amber-950/30 text-amber-700 dark:text-amber-300 hover:bg-amber-100 dark:hover:bg-amber-900/40 text-xs font-semibold transition-colors cursor-pointer shadow-2xs"
            title="Limpiar redes inactivas"
          >
            <Trash2 class="w-3.5 h-3.5" />
            <span class="hidden sm:inline">Limpiar Inactivas</span>
          </button>
        {/if}

        <button
          onclick={() => uiStore.openCreateNetworkModal()}
          class="flex items-center gap-1.5 px-3.5 py-1.5 rounded-xl bg-teal-600 hover:bg-teal-700 text-white text-xs font-semibold transition-colors cursor-pointer shadow-xs"
        >
          <Plus class="w-4 h-4" />
          <span>Crear Red</span>
        </button>
      </div>
    </div>
  </div>

  <!-- Network Cards List -->
  {#if networksStore.loading && networksStore.networks.length === 0}
    <div class="py-20 text-center space-y-3">
      <RefreshCw class="w-8 h-8 animate-spin text-teal-500 mx-auto" />
      <p class="text-xs text-slate-400">Cargando redes Docker...</p>
    </div>
  {:else if filteredNetworks.length === 0}
    <div class="py-16 text-center border border-dashed border-slate-200 dark:border-slate-800 rounded-2xl bg-white dark:bg-slate-900 p-8 space-y-3">
      <div class="w-12 h-12 rounded-2xl bg-teal-500/10 text-teal-600 dark:text-teal-400 flex items-center justify-center mx-auto">
        <Network class="w-6 h-6" />
      </div>
      <h3 class="text-sm font-semibold text-slate-800 dark:text-slate-200">
        No se encontraron redes
      </h3>
      <p class="text-xs text-slate-500 dark:text-slate-400 max-w-sm mx-auto">
        {searchQuery ? `No hay redes que coincidan con "${searchQuery}".` : 'No hay redes para mostrar con el filtro seleccionado.'}
      </p>
      {#if !searchQuery && filterState === 'inactive'}
        <p class="text-xs text-emerald-600 dark:text-emerald-400 font-medium">
          ¡Excelente! Todas las redes personalizadas están en uso.
        </p>
      {:else}
        <button
          onclick={() => uiStore.openCreateNetworkModal()}
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-teal-600 hover:bg-teal-700 text-white text-xs font-semibold cursor-pointer shadow-xs mt-2"
        >
          <Plus class="w-4 h-4" />
          <span>Crear Primera Red</span>
        </button>
      {/if}
    </div>
  {:else}
    <div class="space-y-4">
      {#each filteredNetworks as network (network.id)}
        <NetworkCard {network} />
      {/each}
    </div>
  {/if}
</div>

<!-- Prune Confirm Modal -->
<ConfirmModal
  isOpen={isPruneModalOpen}
  title="Limpiar Redes Inactivas (Network Prune)"
  message={`¿Estás seguro de que deseas eliminar todas las redes Docker no utilizadas? Esta acción eliminará permanentemente ${inactiveNetworks.length} red(es) que no tienen contenedores conectados. Las redes de sistema predeterminadas no se verán afectadas.`}
  confirmLabel={isPruning ? "Limpiando..." : "Eliminar Redes Inactivas"}
  onConfirm={handleConfirmPrune}
  onCancel={() => (isPruneModalOpen = false)}
  isDestructive={true}
/>
