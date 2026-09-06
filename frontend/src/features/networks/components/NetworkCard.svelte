<script lang="ts">
  import {
    Network,
    Server,
    Trash2,
    Code2,
    Link2,
    Unlink,
    Copy,
    Check,
    Globe,
    Shield,
    ChevronDown,
    ChevronRight,
    Play,
    Square,
    RotateCw,
    Loader2,
    Layers,
    Cpu,
  } from '@lucide/svelte';
  import type { NetworkInfo, NetworkContainerRef } from '../../../types';
  import { networksStore } from '../stores/networks.svelte';
  import { containersStore } from '../../containers/stores/containers.svelte';
  import { uiStore } from '../../../shared/stores/ui.svelte';

  let { network } = $props<{ network: NetworkInfo }>();

  let isExpanded = $state(true);
  let copiedSubnet = $state(false);
  let containerToDisconnect = $state<NetworkContainerRef | null>(null);

  const isActionLoading = $derived(
    networksStore.actionLoading === network.id ||
    containersStore.actionLoading === `network:${network.name}`
  );

  function copyToClipboard(text: string) {
    if (!text) return;
    navigator.clipboard.writeText(text);
    copiedSubnet = true;
    setTimeout(() => {
      copiedSubnet = false;
    }, 2000);
  }

  async function handleDisconnect(c: NetworkContainerRef) {
    await networksStore.disconnectContainer(network.id, c.id, false);
    containerToDisconnect = null;
  }

  // Driver badges styling
  const driverStyles: Record<string, { bg: string; text: string; border: string }> = {
    bridge: { bg: 'bg-blue-500/10', text: 'text-blue-600 dark:text-blue-400', border: 'border-blue-500/20' },
    macvlan: { bg: 'bg-purple-500/10', text: 'text-purple-600 dark:text-purple-400', border: 'border-purple-500/20' },
    ipvlan: { bg: 'bg-indigo-500/10', text: 'text-indigo-600 dark:text-indigo-400', border: 'border-indigo-500/20' },
    overlay: { bg: 'bg-emerald-500/10', text: 'text-emerald-600 dark:text-emerald-400', border: 'border-emerald-500/20' },
    host: { bg: 'bg-amber-500/10', text: 'text-amber-600 dark:text-amber-400', border: 'border-amber-500/20' },
    none: { bg: 'bg-slate-500/10', text: 'text-slate-600 dark:text-slate-400', border: 'border-slate-500/20' },
  };

  const currentDriverStyle = $derived(driverStyles[network.driver] || driverStyles.bridge);
</script>

<div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl p-4 shadow-xs transition-all duration-200 hover:border-slate-300 dark:hover:border-slate-700 space-y-3.5">
  <!-- Card Header -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
    <!-- Left: Network name, driver & badges -->
    <div class="flex items-center gap-3 min-w-0">
      <div class="w-10 h-10 rounded-xl {currentDriverStyle.bg} border {currentDriverStyle.border} {currentDriverStyle.text} flex items-center justify-center shrink-0">
        <Network class="w-5 h-5" />
      </div>

      <div class="min-w-0">
        <div class="flex items-center gap-2 flex-wrap">
          <h3 class="font-bold text-slate-900 dark:text-slate-100 text-sm tracking-tight truncate">
            {network.name}
          </h3>

          <!-- Driver badge -->
          <span class="text-[10px] font-mono uppercase tracking-wider px-2 py-0.5 rounded-md font-semibold {currentDriverStyle.bg} {currentDriverStyle.text} border {currentDriverStyle.border}">
            {network.driver}
          </span>

          <!-- Scope -->
          <span class="text-[10px] font-mono px-1.5 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-slate-500 dark:text-slate-400 border border-slate-200 dark:border-slate-700">
            {network.scope}
          </span>

          <!-- Default Network Badge -->
          {#if network.isDefault}
            <span class="text-[10px] font-medium px-2 py-0.5 rounded-md bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 border border-slate-200 dark:border-slate-700">
              Predeterminada
            </span>
          {/if}

          <!-- Internal badge -->
          {#if network.internal}
            <span class="text-[10px] font-medium px-1.5 py-0.5 rounded-md bg-rose-50 dark:bg-rose-950/40 text-rose-600 dark:text-rose-400 border border-rose-200 dark:border-rose-800 flex items-center gap-1">
              <Shield class="w-2.5 h-2.5" />
              Aislada Interna
            </span>
          {/if}

          <!-- IPv6 badge -->
          {#if network.enableIPv6}
            <span class="text-[10px] font-medium px-1.5 py-0.5 rounded-md bg-cyan-50 dark:bg-cyan-950/40 text-cyan-600 dark:text-cyan-400 border border-cyan-200 dark:border-cyan-800">
              IPv6
            </span>
          {/if}
        </div>

        <div class="flex items-center gap-2 text-[11px] text-slate-500 dark:text-slate-400 mt-0.5 font-mono">
          <span>ID: {network.shortId || (network.id ? network.id.slice(0, 12) : '')}</span>
          {#if network.created}
            <span>• Creada: {new Date(network.created).toLocaleDateString()}</span>
          {/if}
        </div>
      </div>
    </div>

    <!-- Right: Header Actions -->
    <div class="flex items-center gap-1.5 self-end sm:self-center shrink-0">
      <!-- Start/Stop containers in network (if any) -->
      {#if network.containersCount > 0}
        <button
          onclick={() => containersStore.handleStartNetwork(network.name)}
          disabled={isActionLoading}
          class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-emerald-50 dark:hover:bg-emerald-950/30 text-emerald-600 dark:text-emerald-400 transition-colors cursor-pointer disabled:opacity-50"
          title="Iniciar todos los contenedores de la red"
        >
          <Play class="w-3.5 h-3.5" />
        </button>
        <button
          onclick={() => containersStore.handleStopNetwork(network.name)}
          disabled={isActionLoading}
          class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-rose-50 dark:hover:bg-rose-950/30 text-rose-600 dark:text-rose-400 transition-colors cursor-pointer disabled:opacity-50"
          title="Detener todos los contenedores de la red"
        >
          <Square class="w-3.5 h-3.5" />
        </button>
      {/if}

      <!-- Connect container button -->
      <button
        onclick={() => uiStore.openConnectContainer(network.id, network.name)}
        class="flex items-center gap-1 px-2.5 py-1.5 rounded-xl border border-teal-200 dark:border-teal-800 bg-teal-50 dark:bg-teal-950/30 text-teal-700 dark:text-teal-300 hover:bg-teal-100 dark:hover:bg-teal-900/40 text-xs font-medium transition-colors cursor-pointer"
        title="Conectar un contenedor a esta red en caliente"
      >
        <Link2 class="w-3.5 h-3.5" />
        <span class="hidden md:inline">Conectar</span>
      </button>

      <!-- Inspect Network JSON -->
      <button
        onclick={() => uiStore.openNetworkInspect(network.id, network.name)}
        class="p-1.5 rounded-xl border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-400 transition-colors cursor-pointer"
        title="Inspeccionar configuración JSON"
      >
        <Code2 class="w-4 h-4" />
      </button>

      <!-- Delete Network -->
      {#if !network.isDefault}
        <button
          onclick={() => networksStore.removeNetwork(network.id, network.name)}
          disabled={isActionLoading || network.containersCount > 0}
          class="p-1.5 rounded-xl border border-slate-200 dark:border-slate-700 hover:border-rose-300 dark:hover:border-rose-800 hover:bg-rose-50 dark:hover:bg-rose-950/30 text-slate-400 hover:text-rose-600 dark:hover:text-rose-400 transition-colors cursor-pointer disabled:opacity-30 disabled:cursor-not-allowed"
          title={network.containersCount > 0 ? "No se puede eliminar: tiene contenedores conectados" : "Eliminar red"}
        >
          {#if isActionLoading}
            <Loader2 class="w-4 h-4 animate-spin" />
          {:else}
            <Trash2 class="w-4 h-4" />
          {/if}
        </button>
      {/if}
    </div>
  </div>

  <!-- IPAM Details Bar -->
  {#if network.ipam && network.ipam.length > 0 && (network.ipam[0].subnet || network.ipam[0].gateway)}
    <div class="bg-slate-50 dark:bg-slate-850/50 border border-slate-200/80 dark:border-slate-800/80 rounded-xl px-3.5 py-2.5 flex flex-wrap items-center justify-between gap-3 text-xs">
      <div class="flex items-center gap-4 flex-wrap">
        {#each network.ipam as ipamConfig}
          {#if ipamConfig.subnet}
            <div class="flex items-center gap-1.5">
              <span class="text-slate-500 dark:text-slate-400 font-medium">Subred:</span>
              <button
                onclick={() => copyToClipboard(ipamConfig.subnet || '')}
                class="font-mono text-slate-800 dark:text-slate-200 hover:text-teal-600 dark:hover:text-teal-400 flex items-center gap-1 font-semibold cursor-pointer"
                title="Copiar subred"
              >
                <span>{ipamConfig.subnet}</span>
                {#if copiedSubnet}
                  <Check class="w-3 h-3 text-emerald-500" />
                {:else}
                  <Copy class="w-3 h-3 opacity-50 hover:opacity-100" />
                {/if}
              </button>
            </div>
          {/if}

          {#if ipamConfig.gateway}
            <div class="flex items-center gap-1.5">
              <span class="text-slate-500 dark:text-slate-400 font-medium">Gateway:</span>
              <span class="font-mono text-slate-800 dark:text-slate-200">{ipamConfig.gateway}</span>
            </div>
          {/if}

          {#if ipamConfig.ipRange}
            <div class="flex items-center gap-1.5">
              <span class="text-slate-500 dark:text-slate-400 font-medium">Rango IP:</span>
              <span class="font-mono text-slate-800 dark:text-slate-200">{ipamConfig.ipRange}</span>
            </div>
          {/if}
        {/each}
      </div>

      <div class="text-[11px] text-slate-400 font-mono">
        Driver IPAM: {network.ipam[0]?.driver || 'default'}
      </div>
    </div>
  {/if}

  <!-- Connected Containers Section -->
  <div class="space-y-2">
    <div class="flex items-center justify-between">
      <button
        onclick={() => (isExpanded = !isExpanded)}
        class="flex items-center gap-2 text-xs font-semibold text-slate-700 dark:text-slate-300 hover:text-slate-900 dark:hover:text-slate-100 cursor-pointer"
      >
        {#if isExpanded}
          <ChevronDown class="w-3.5 h-3.5" />
        {:else}
          <ChevronRight class="w-3.5 h-3.5" />
        {/if}
        <span>Contenedores Conectados ({network.containersCount})</span>
      </button>

      {#if network.containersCount === 0}
        <span class="text-[11px] text-amber-600 dark:text-amber-400 font-medium">
          Inactiva (Sin contenedores)
        </span>
      {/if}
    </div>

    {#if isExpanded}
      {#if network.containersCount === 0}
        <div class="py-4 text-center border border-dashed border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50/50 dark:bg-slate-900/30">
          <p class="text-xs text-slate-400">
            No hay contenedores conectados a esta red.
          </p>
          <button
            onclick={() => uiStore.openConnectContainer(network.id, network.name)}
            class="mt-1.5 text-xs font-medium text-teal-600 dark:text-teal-400 hover:underline cursor-pointer"
          >
            + Conectar un contenedor ahora
          </button>
        </div>
      {:else}
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2">
          {#each network.containers as c}
            <div class="flex items-center justify-between p-2.5 rounded-xl border border-slate-200/80 dark:border-slate-800/80 bg-slate-50/70 dark:bg-slate-800/40 text-xs hover:bg-slate-100/70 dark:hover:bg-slate-800/70 transition-colors">
              <div class="min-w-0 pr-2">
                <div class="font-medium text-slate-900 dark:text-slate-100 truncate flex items-center gap-1.5">
                  <span class="w-2 h-2 rounded-full shrink-0 {c.state === 'running' ? 'bg-emerald-500' : 'bg-slate-400'}" title={c.state === 'running' ? 'Activo' : 'Detenido'}></span>
                  <span class="truncate">{c.name}</span>
                  {#if c.state && c.state !== 'running'}
                    <span class="text-[9px] px-1 py-0.2 rounded bg-slate-200 dark:bg-slate-700 text-slate-500 uppercase">detenido</span>
                  {/if}
                </div>
                <div class="font-mono text-[11px] text-slate-500 dark:text-slate-400 flex items-center gap-2 mt-0.5 truncate">
                  {#if c.ipv4Address}
                    <span class="text-teal-600 dark:text-teal-400 font-semibold">{c.ipv4Address}</span>
                  {/if}
                  {#if c.macAddress}
                    <span class="text-slate-400 truncate">{c.macAddress}</span>
                  {/if}
                </div>
              </div>

              <!-- Disconnect button -->
              <button
                onclick={() => handleDisconnect(c)}
                class="p-1 rounded-lg text-slate-400 hover:text-rose-600 dark:hover:text-rose-400 hover:bg-rose-50 dark:hover:bg-rose-950/30 transition-colors cursor-pointer shrink-0"
                title="Desconectar {c.name} de esta red"
              >
                <Unlink class="w-3.5 h-3.5" />
              </button>
            </div>
          {/each}
        </div>
      {/if}
    {/if}
  </div>
</div>
