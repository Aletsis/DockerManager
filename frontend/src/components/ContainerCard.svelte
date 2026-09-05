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
  } from '@lucide/svelte';
  import type { ContainerInfo, ContainerStats } from '../types';
  import { formatBytes, getStateColor, formatUptime } from '../utils';

  let {
    container,
    stats,
    onStart,
    onStop,
    onRestart,
    onPause,
    onUnpause,
    onRemove,
    onOpenTerminal,
    onViewLogs,
    onViewStats,
    actionLoading,
  } = $props<{
    container: ContainerInfo;
    stats?: ContainerStats;
    onStart: (id: string) => void;
    onStop: (id: string) => void;
    onRestart: (id: string) => void;
    onPause: (id: string) => void;
    onUnpause: (id: string) => void;
    onRemove: (id: string, name: string) => void;
    onOpenTerminal: (id: string, name: string) => void;
    onViewLogs: (id: string, name: string) => void;
    onViewStats: (id: string, name: string) => void;
    actionLoading?: string;
  }>();

  let copied = $state(false);

  const isRunning = $derived(container.state === 'running');
  const isPaused = $derived(container.state === 'paused');
  const isExited = $derived(container.state === 'exited' || container.state === 'dead');
  const colors = $derived(getStateColor(container.state));
  const isLoadingThis = $derived(actionLoading === container.id);

  function copyId(e: MouseEvent) {
    e.stopPropagation();
    navigator.clipboard.writeText(container.id);
    copied = true;
    setTimeout(() => (copied = false), 1500);
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
              <span
                class="text-[10px] font-mono px-1.5 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 border border-slate-200 dark:border-slate-700"
              >
                {p.publicPort ? `${p.publicPort}:${p.privatePort}` : `${p.privatePort}/${p.type}`}
              </span>
            {/each}
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
          onclick={() => onStop(container.id)}
          disabled={isLoadingThis}
          title="Detener contenedor"
          class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-rose-50 dark:hover:bg-rose-950/40 text-slate-600 hover:text-rose-600 dark:text-slate-300 dark:hover:text-rose-400 transition-colors"
        >
          <Square class="w-4 h-4 fill-current" />
        </button>

        <button
          onclick={() => onRestart(container.id)}
          disabled={isLoadingThis}
          title="Reiniciar contenedor"
          class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300 transition-colors"
        >
          <RotateCw class="w-4 h-4" />
        </button>

        <button
          onclick={() => onPause(container.id)}
          disabled={isLoadingThis}
          title="Pausar contenedor"
          class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-amber-50 dark:hover:bg-amber-950/40 text-slate-600 hover:text-amber-600 dark:text-slate-300 dark:hover:text-amber-400 transition-colors"
        >
          <Pause class="w-4 h-4" />
        </button>
      {:else if isPaused}
        <button
          onclick={() => onUnpause(container.id)}
          disabled={isLoadingThis}
          title="Reanudar contenedor"
          class="p-1.5 rounded-lg border border-emerald-500/30 bg-emerald-50 dark:bg-emerald-950/30 text-emerald-600 dark:text-emerald-400 hover:bg-emerald-100 dark:hover:bg-emerald-900/40 transition-colors"
        >
          <Play class="w-4 h-4 fill-current" />
        </button>
      {:else}
        <button
          onclick={() => onStart(container.id)}
          disabled={isLoadingThis}
          title="Iniciar contenedor"
          class="p-1.5 rounded-lg border border-emerald-500/30 bg-emerald-50 dark:bg-emerald-950/30 text-emerald-600 dark:text-emerald-400 hover:bg-emerald-100 dark:hover:bg-emerald-900/40 transition-colors"
        >
          <Play class="w-4 h-4 fill-current" />
        </button>
      {/if}

      <!-- Terminal (Interactive Shell CLI) -->
      <button
        onclick={() => isRunning && onOpenTerminal(container.id, container.name)}
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
        onclick={() => onViewLogs(container.id, container.name)}
        title="Ver registros (Logs)"
        class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300 transition-colors"
      >
        <ScrollText class="w-4 h-4" />
      </button>

      <!-- Stats Monitor -->
      {#if isRunning}
        <button
          onclick={() => onViewStats(container.id, container.name)}
          title="Monitoreo de recursos detallado"
          class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300 transition-colors"
        >
          <Activity class="w-4 h-4" />
        </button>
      {/if}

      <!-- Remove -->
      {#if isExited || isPaused}
        <button
          onclick={() => onRemove(container.id, container.name)}
          disabled={isLoadingThis}
          title="Eliminar contenedor"
          class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-rose-50 dark:hover:bg-rose-950/40 text-slate-400 hover:text-rose-600 dark:hover:text-rose-400 transition-colors"
        >
          <Trash2 class="w-4 h-4" />
        </button>
      {/if}
    </div>
  </div>
</div>
