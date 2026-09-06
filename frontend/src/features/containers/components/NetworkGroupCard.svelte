<script lang="ts">
  import {
    Network,
    ChevronDown,
    ChevronRight,
    Play,
    Square,
    RotateCw,
    Loader2,
  } from '@lucide/svelte';
  import type { DockerNetworkGroup } from '../../../types';
  import ContainerCard from './ContainerCard.svelte';
  import { containersStore } from '../stores/containers.svelte';

  let { network } = $props<{ network: DockerNetworkGroup }>();

  let isExpanded = $state(true);

  const isNetworkLoading = $derived(containersStore.actionLoading === `network:${network.name}`);
  const allRunning = $derived(network.runningCount === network.totalCount && network.totalCount > 0);
  const allStopped = $derived(network.runningCount === 0);
</script>

<div class="rounded-2xl border border-teal-200/70 dark:border-teal-900/40 bg-teal-50/20 dark:bg-teal-950/10 p-4 transition-all duration-200 shadow-xs">
  <!-- Network Header Bar -->
  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 pb-3">
    <!-- Left Section: Expand toggle, Network Name, Info -->
    <div class="flex items-center gap-3 min-w-0">
      <button
        onclick={() => (isExpanded = !isExpanded)}
        class="p-1 rounded-lg hover:bg-teal-100 dark:hover:bg-teal-900/40 text-slate-500 dark:text-slate-400 transition-colors cursor-pointer"
        title={isExpanded ? "Plegar red" : "Desplegar red"}
      >
        {#if isExpanded}
          <ChevronDown class="w-4 h-4" />
        {:else}
          <ChevronRight class="w-4 h-4" />
        {/if}
      </button>

      <div class="w-8 h-8 rounded-xl bg-teal-600/10 dark:bg-teal-400/10 border border-teal-500/20 text-teal-600 dark:text-teal-400 flex items-center justify-center flex-shrink-0">
        <Network class="w-4 h-4" />
      </div>

      <div class="min-w-0">
        <div class="flex items-center gap-2 flex-wrap">
          <h2 class="font-bold text-slate-900 dark:text-slate-100 text-sm truncate tracking-tight">
            {network.name}
          </h2>

          {#if network.isDefault}
            <span class="text-[10px] font-semibold uppercase tracking-wider px-2 py-0.5 rounded-md bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 border border-slate-200 dark:border-slate-700">
              Red Predeterminada
            </span>
          {:else}
            <span class="text-[10px] font-semibold uppercase tracking-wider px-2 py-0.5 rounded-md bg-teal-100 dark:bg-teal-900/50 text-teal-700 dark:text-teal-300 border border-teal-200 dark:border-teal-800">
              Red Personalizada
            </span>
          {/if}

          <span class="text-xs text-slate-500 dark:text-slate-400 font-medium">
            ({network.totalCount} {network.totalCount === 1 ? 'contenedor conectado' : 'contenedores conectados'})
          </span>
        </div>

        {#if network.networkId}
          <div class="text-[11px] text-slate-400 dark:text-slate-500 font-mono truncate mt-0.5" title={network.networkId}>
            ID: {network.networkId.substring(0, 12)}
          </div>
        {/if}
      </div>
    </div>

    <!-- Right Section: Batch Operations & Running count -->
    <div class="flex items-center justify-between sm:justify-end gap-3 flex-shrink-0">
      <!-- Status Badge -->
      <div class="flex items-center gap-1.5 text-xs">
        <span class="font-medium text-slate-600 dark:text-slate-300">
          <span class="font-bold text-emerald-600 dark:text-emerald-400">{network.runningCount}</span>/{network.totalCount} activos
        </span>
      </div>

      <!-- Actions Buttons -->
      <div class="flex items-center gap-1.5">
        {#if isNetworkLoading}
          <div class="flex items-center gap-1.5 px-3 py-1.5 text-xs text-teal-600 dark:text-teal-400 font-medium">
            <Loader2 class="w-3.5 h-3.5 animate-spin" />
            <span>Procesando red...</span>
          </div>
        {:else}
          <!-- Start Network -->
          <button
            onclick={() => containersStore.handleStartNetwork(network.name)}
            disabled={allRunning || !!containersStore.actionLoading}
            title={allRunning ? "Todos los contenedores ya están activos" : "Iniciar todos los contenedores de la red"}
            class="px-2.5 py-1 rounded-lg text-xs font-medium border border-emerald-500/30 bg-emerald-50 dark:bg-emerald-950/30 text-emerald-600 dark:text-emerald-400 hover:bg-emerald-100 dark:hover:bg-emerald-900/50 transition-colors flex items-center gap-1 disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer"
          >
            <Play class="w-3.5 h-3.5 fill-current" />
            <span class="hidden md:inline">Iniciar</span>
          </button>

          <!-- Restart Network -->
          <button
            onclick={() => containersStore.handleRestartNetwork(network.name)}
            disabled={!!containersStore.actionLoading}
            title="Reiniciar todos los contenedores de esta red"
            class="px-2.5 py-1 rounded-lg text-xs font-medium border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-200 transition-colors flex items-center gap-1 disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer"
          >
            <RotateCw class="w-3.5 h-3.5" />
            <span class="hidden md:inline">Reiniciar</span>
          </button>

          <!-- Stop Network -->
          <button
            onclick={() => containersStore.handleStopNetwork(network.name)}
            disabled={allStopped || !!containersStore.actionLoading}
            title={allStopped ? "Todos los contenedores ya están detenidos" : "Detener todos los contenedores de la red"}
            class="px-2.5 py-1 rounded-lg text-xs font-medium border border-rose-200 dark:border-rose-900/60 bg-rose-50 dark:bg-rose-950/30 text-rose-600 dark:text-rose-400 hover:bg-rose-100 dark:hover:bg-rose-900/40 transition-colors flex items-center gap-1 disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer"
          >
            <Square class="w-3.5 h-3.5 fill-current" />
            <span class="hidden md:inline">Detener</span>
          </button>
        {/if}
      </div>
    </div>
  </div>

  <!-- Network Children Containers List -->
  {#if isExpanded}
    <div class="space-y-2.5 pt-2 border-t border-teal-200/50 dark:border-teal-900/30">
      {#each network.containers as container (container.id)}
        <ContainerCard
          {container}
          currentNetworkContext={network.name}
        />
      {/each}
    </div>
  {/if}
</div>
