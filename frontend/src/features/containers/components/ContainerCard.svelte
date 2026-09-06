<script lang="ts">
  import {
    Play,
    Square,
    RotateCw,
    Pause,
    Trash2,
    Terminal,
    ScrollText,
    Activity,
    Copy,
    Check,
    Globe,
    Clock,
    Network,
    ExternalLink,
    Code2,
  } from '@lucide/svelte';
  import { BrowserOpenURL } from '../../../../wailsjs/runtime';
  import type { ContainerInfo } from '../../../types';
  import { formatBytes, getStateColor, formatUptime } from '../../../shared/utils/utils';
  import { containersStore } from '../stores/containers.svelte';
  import { uiStore } from '../../../shared/stores/ui.svelte';

  let {
    container,
    currentNetworkContext,
  } = $props<{
    container: ContainerInfo;
    currentNetworkContext?: string;
  }>();

  let copied = $state(false);

  const stats = $derived(containersStore.statsMap[container.id]);
  const isRunning = $derived(container.state === 'running');
  const isPaused = $derived(container.state === 'paused');
  const isExited = $derived(container.state === 'exited' || container.state === 'dead');
  const colors = $derived(getStateColor(container.state));
  const isLoadingThis = $derived(containersStore.actionLoading === container.id);

  function copyId(e: MouseEvent) {
    e.stopPropagation();
    navigator.clipboard.writeText(container.id);
    copied = true;
    setTimeout(() => (copied = false), 1500);
  }

  function openPortUrl(e: MouseEvent, publicPort: number) {
    e.stopPropagation();
    const url = `http://localhost:${publicPort}`;
    try {
      BrowserOpenURL(url);
    } catch {
      window.open(url, '_blank');
    }
  }
</script>

<div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-4 shadow-sm hover:shadow transition-all duration-200 hover:border-slate-300 dark:hover:border-slate-700">
  <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
    <!-- Container Main Info -->
    <div class="flex items-start gap-3.5 min-w-0">
      <!-- Status Dot -->
      <div class="mt-1 flex-shrink-0">
        <div class="w-3 h-3 rounded-full {colors.dotBg} {colors.glowClass}"></div>
      </div>

      <div class="min-w-0 space-y-1">
        <div class="flex items-center gap-2 flex-wrap">
          <h3 class="font-semibold text-slate-900 dark:text-slate-100 text-sm truncate">
            {container.name}
          </h3>

          <!-- State Badge -->
          <span
            class="text-[11px] font-medium px-2 py-0.5 rounded-full border capitalize {colors.badgeBg} {colors.badgeText}"
          >
            {container.state}
          </span>

          {#if container.composeService}
            <span
              title="Servicio Docker Compose"
              class="text-[11px] font-medium px-2 py-0.5 rounded-full border bg-violet-50 dark:bg-violet-950/40 text-violet-700 dark:text-violet-300 border-violet-200 dark:border-violet-800/60 font-mono"
            >
              svc: {container.composeService}
            </span>
          {/if}

          <!-- Short ID -->
          <button
            onclick={copyId}
            title="Copiar ID completo"
            class="flex items-center gap-1 font-mono text-[11px] text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 px-1.5 py-0.5 rounded bg-slate-100 dark:bg-slate-800 transition-colors"
          >
            <span>{container.shortId}</span>
            {#if copied}
              <Check class="w-3 h-3 text-emerald-500" />
            {:else}
              <Copy class="w-3 h-3" />
            {/if}
          </button>
        </div>

        <!-- Image name & Uptime -->
        <div class="flex items-center gap-3 text-xs text-slate-500 dark:text-slate-400 flex-wrap">
          <span class="font-mono text-slate-600 dark:text-slate-300 truncate max-w-xs">
            {container.image}
          </span>
          <span class="flex items-center gap-1">
            <Clock class="w-3.5 h-3.5 text-slate-400" />
            <span>{formatUptime(container.status)}</span>
          </span>
        </div>

        <!-- Ports -->
        {#if container.ports && container.ports.length > 0}
          <div class="flex items-center gap-1.5 flex-wrap pt-0.5">
            <Globe class="w-3.5 h-3.5 text-slate-400 flex-shrink-0" />
            {#each container.ports as p, idx}
              {#if p.publicPort}
                <button
                  onclick={(e) => openPortUrl(e, p.publicPort)}
                  title={`Abrir http://localhost:${p.publicPort} en el navegador`}
                  class="flex items-center gap-1 text-[10px] font-mono px-1.5 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-indigo-600 dark:text-indigo-400 border border-slate-200 dark:border-slate-700 hover:border-indigo-400 dark:hover:border-indigo-500 hover:bg-indigo-50/80 dark:hover:bg-indigo-950/40 transition-colors cursor-pointer group"
                >
                  <span>{p.publicPort}:{p.privatePort}</span>
                  <ExternalLink class="w-2.5 h-2.5 opacity-60 group-hover:opacity-100 transition-opacity" />
                </button>
              {:else}
                <span
                  class="text-[10px] font-mono px-1.5 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 border border-slate-200 dark:border-slate-700"
                >
                  {p.privatePort}/{p.type}
                </span>
              {/if}
            {/each}
          </div>
        {/if}

        <!-- Networks & IP Details -->
        {#if container.networks && container.networks.length > 0}
          <div class="flex items-center gap-1.5 flex-wrap pt-0.5">
            <Network class="w-3.5 h-3.5 text-teal-500/80 flex-shrink-0" />
            {#if currentNetworkContext}
              {@const currentNet = container.networks.find(n => n.networkName === currentNetworkContext)}
              {#if currentNet && currentNet.ipAddress}
                <span
                  title={currentNet.gateway ? `Gateway: ${currentNet.gateway}` : undefined}
                  class="text-[10px] font-mono px-1.5 py-0.5 rounded bg-teal-50 dark:bg-teal-950/40 text-teal-700 dark:text-teal-300 border border-teal-200 dark:border-teal-800/60 font-medium"
                >
                  IP: {currentNet.ipAddress}
                </span>
              {:else if currentNet}
                <span class="text-[10px] font-mono px-1.5 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-slate-500">
                  conectado
                </span>
              {/if}

              {#if container.networks.length > 1}
                <span
                  title={`Otras redes: ${container.networks.filter(n => n.networkName !== currentNetworkContext).map(n => n.networkName).join(', ')}`}
                  class="text-[10px] font-medium px-1.5 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-slate-500 dark:text-slate-400 border border-slate-200 dark:border-slate-700"
                >
                  +{container.networks.length - 1} {container.networks.length - 1 === 1 ? 'red más' : 'redes más'}
                </span>
              {/if}
            {:else}
              {#each container.networks as net}
                <span
                  title={net.ipAddress ? `IP: ${net.ipAddress}${net.gateway ? ` | Gateway: ${net.gateway}` : ''}` : undefined}
                  class="text-[10px] font-mono px-1.5 py-0.5 rounded bg-teal-50 dark:bg-teal-950/40 text-teal-700 dark:text-teal-300 border border-teal-200 dark:border-teal-800/60"
                >
                  {net.networkName}{net.ipAddress ? ` (${net.ipAddress})` : ''}
                </span>
              {/each}
            {/if}
          </div>
        {/if}
      </div>
    </div>

    <!-- Live Resource Meters (If running) -->
    {#if isRunning && stats}
      <div class="flex items-center gap-4 py-1.5 px-3 rounded-lg bg-slate-50 dark:bg-slate-800/60 border border-slate-200/70 dark:border-slate-800 text-xs">
        <!-- CPU -->
        <div class="space-y-1 w-24">
          <div class="flex justify-between text-[10px] text-slate-500 dark:text-slate-400">
            <span>CPU</span>
            <span class="font-mono font-medium text-slate-700 dark:text-slate-200">
              {stats.cpuPercentage.toFixed(1)}%
            </span>
          </div>
          <div class="w-full h-1.5 bg-slate-200 dark:bg-slate-700 rounded-full overflow-hidden">
            <div
              class="h-full transition-all duration-300 {stats.cpuPercentage > 80 ? 'bg-rose-500' : stats.cpuPercentage > 50 ? 'bg-amber-500' : 'bg-indigo-500'}"
              style="width: {Math.min(stats.cpuPercentage, 100)}%;"
            ></div>
          </div>
        </div>

        <!-- RAM -->
        <div class="space-y-1 w-28">
          <div class="flex justify-between text-[10px] text-slate-500 dark:text-slate-400">
            <span>RAM</span>
            <span class="font-mono font-medium text-slate-700 dark:text-slate-200 truncate">
              {formatBytes(stats.memoryUsage)}
            </span>
          </div>
          <div class="w-full h-1.5 bg-slate-200 dark:bg-slate-700 rounded-full overflow-hidden">
            <div
              class="h-full transition-all duration-300 {stats.memoryPercentage > 85 ? 'bg-rose-500' : stats.memoryPercentage > 65 ? 'bg-amber-500' : 'bg-emerald-500'}"
              style="width: {Math.min(stats.memoryPercentage, 100)}%;"
            ></div>
          </div>
        </div>
      </div>
    {/if}

    <!-- Action Buttons Toolbar -->
    <div class="flex items-center gap-1.5 flex-wrap justify-end">
      <!-- Start / Stop / Restart / Pause -->
      {#if isRunning}
        <button
          onclick={() => containersStore.handleStop(container.id)}
          disabled={isLoadingThis}
          title="Detener contenedor"
          class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-rose-50 dark:hover:bg-rose-950/40 text-slate-600 hover:text-rose-600 dark:text-slate-300 dark:hover:text-rose-400 transition-colors cursor-pointer"
        >
          <Square class="w-4 h-4 fill-current" />
        </button>

        <button
          onclick={() => containersStore.handleRestart(container.id)}
          disabled={isLoadingThis}
          title="Reiniciar contenedor"
          class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300 transition-colors cursor-pointer"
        >
          <RotateCw class="w-4 h-4" />
        </button>

        <button
          onclick={() => containersStore.handlePause(container.id)}
          disabled={isLoadingThis}
          title="Pausar contenedor"
          class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-amber-50 dark:hover:bg-amber-950/40 text-slate-600 hover:text-amber-600 dark:text-slate-300 dark:hover:text-amber-400 transition-colors cursor-pointer"
        >
          <Pause class="w-4 h-4" />
        </button>
      {:else if isPaused}
        <button
          onclick={() => containersStore.handleUnpause(container.id)}
          disabled={isLoadingThis}
          title="Reanudar contenedor"
          class="p-1.5 rounded-lg border border-emerald-500/30 bg-emerald-50 dark:bg-emerald-950/30 text-emerald-600 dark:text-emerald-400 hover:bg-emerald-100 dark:hover:bg-emerald-900/40 transition-colors cursor-pointer"
        >
          <Play class="w-4 h-4 fill-current" />
        </button>
      {:else}
        <button
          onclick={() => containersStore.handleStart(container.id)}
          disabled={isLoadingThis}
          title="Iniciar contenedor"
          class="p-1.5 rounded-lg border border-emerald-500/30 bg-emerald-50 dark:bg-emerald-950/30 text-emerald-600 dark:text-emerald-400 hover:bg-emerald-100 dark:hover:bg-emerald-900/40 transition-colors cursor-pointer"
        >
          <Play class="w-4 h-4 fill-current" />
        </button>
      {/if}

      <!-- Terminal (Interactive Shell CLI) -->
      <button
        onclick={() => isRunning && uiStore.openTerminal(container.id, container.name)}
        disabled={!isRunning || !!isLoadingThis}
        title={isRunning ? "Terminal interactiva (Shell CLI)" : "Terminal interactiva (Inicia el contenedor para acceder)"}
        class="p-1.5 rounded-lg border transition-colors {isRunning
          ? 'border-indigo-500/30 bg-indigo-50/60 dark:bg-indigo-950/40 hover:bg-indigo-100 dark:hover:bg-indigo-900/60 text-indigo-600 dark:text-indigo-400 cursor-pointer'
          : 'border-slate-200 dark:border-slate-800 text-slate-300 dark:text-slate-600 opacity-60 cursor-not-allowed'}"
      >
        <Terminal class="w-4 h-4" />
      </button>

      <!-- Logs -->
      <button
        onclick={() => uiStore.openLogs(container.id, container.name)}
        title="Ver registros (Logs)"
        class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300 transition-colors cursor-pointer"
      >
        <ScrollText class="w-4 h-4" />
      </button>

      <!-- Inspect JSON -->
      <button
        onclick={() => uiStore.openInspect(container.id, container.name)}
        title="Inspeccionar configuración (Docker Inspect JSON)"
        class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-indigo-50 dark:hover:bg-indigo-950/40 text-slate-600 hover:text-indigo-600 dark:text-slate-300 dark:hover:text-indigo-400 transition-colors cursor-pointer"
      >
        <Code2 class="w-4 h-4" />
      </button>

      <!-- Stats Monitor -->
      {#if isRunning}
        <button
          onclick={() => uiStore.openStats(container.id, container.name)}
          title="Monitoreo de recursos detallado"
          class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300 transition-colors cursor-pointer"
        >
          <Activity class="w-4 h-4" />
        </button>
      {/if}

      <!-- Remove -->
      {#if isExited || isPaused}
        <button
          onclick={() => uiStore.openConfirmDelete(container.id, container.name)}
          disabled={isLoadingThis}
          title="Eliminar contenedor"
          class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-rose-50 dark:hover:bg-rose-950/40 text-slate-400 hover:text-rose-600 dark:hover:text-rose-400 transition-colors cursor-pointer"
        >
          <Trash2 class="w-4 h-4" />
        </button>
      {/if}
    </div>
  </div>
</div>
