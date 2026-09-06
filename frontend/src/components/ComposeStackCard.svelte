<script lang="ts">
  import {
    Layers,
    ChevronDown,
    ChevronRight,
    Play,
    Square,
    RotateCw,
    Folder,
    Loader2,
  } from '@lucide/svelte';
  import type { ComposeStackGroup, ContainerStats } from '../types';
  import ContainerCard from './ContainerCard.svelte';

  let {
    stack,
    statsMap,
    onStart,
    onStop,
    onRestart,
    onPause,
    onUnpause,
    onRemove,
    onOpenTerminal,
    onViewLogs,
    onViewStats,
    onStartStack,
    onStopStack,
    onRestartStack,
    actionLoading = '',
  } = $props<{
    stack: ComposeStackGroup;
    statsMap: Record<string, ContainerStats>;
    onStart: (id: string) => void;
    onStop: (id: string) => void;
    onRestart: (id: string) => void;
    onPause: (id: string) => void;
    onUnpause: (id: string) => void;
    onRemove: (id: string, name: string) => void;
    onOpenTerminal: (id: string, name: string) => void;
    onViewLogs: (id: string, name: string) => void;
    onViewStats: (id: string, name: string) => void;
    onStartStack: (projectName: string) => void;
    onStopStack: (projectName: string) => void;
    onRestartStack: (projectName: string) => void;
    actionLoading?: string;
  }>();

  let isExpanded = $state(true);

  const isStackLoading = $derived(actionLoading === `stack:${stack.name}`);
  const hasRunning = $derived(stack.runningCount > 0);
  const allRunning = $derived(stack.runningCount === stack.totalCount && stack.totalCount > 0);
  const someRunning = $derived(stack.runningCount > 0 && stack.runningCount < stack.totalCount);
  const allStopped = $derived(stack.runningCount === 0);
</script>

<div class="rounded-2xl border border-violet-200/70 dark:border-violet-900/40 bg-violet-50/20 dark:bg-violet-950/10 p-4 transition-all duration-200 shadow-xs">
  <!-- Stack Header Bar -->
  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 pb-3">
    <!-- Left Section: Expand toggle, Stack Name, Info -->
    <div class="flex items-center gap-3 min-w-0">
      <button
        onclick={() => (isExpanded = !isExpanded)}
        class="p-1 rounded-lg hover:bg-violet-100 dark:hover:bg-violet-900/40 text-slate-500 dark:text-slate-400 transition-colors"
        title={isExpanded ? "Plegar stack" : "Desplegar stack"}
      >
        {#if isExpanded}
          <ChevronDown class="w-4 h-4" />
        {:else}
          <ChevronRight class="w-4 h-4" />
        {/if}
      </button>

      <div class="w-8 h-8 rounded-xl bg-violet-600/10 dark:bg-violet-400/10 border border-violet-500/20 text-violet-600 dark:text-violet-400 flex items-center justify-center flex-shrink-0">
        <Layers class="w-4 h-4" />
      </div>

      <div class="min-w-0">
        <div class="flex items-center gap-2 flex-wrap">
          <h2 class="font-bold text-slate-900 dark:text-slate-100 text-sm truncate tracking-tight">
            {stack.name}
          </h2>
          <span class="text-[10px] font-semibold uppercase tracking-wider px-2 py-0.5 rounded-md bg-violet-100 dark:bg-violet-900/50 text-violet-700 dark:text-violet-300 border border-violet-200 dark:border-violet-800">
            Compose Stack
          </span>
          <span class="text-xs text-slate-500 dark:text-slate-400 font-medium">
            ({stack.totalCount} {stack.totalCount === 1 ? 'servicio' : 'servicios'})
          </span>
        </div>

        {#if stack.workingDir || stack.configFile}
          <div class="flex items-center gap-1.5 text-[11px] text-slate-400 dark:text-slate-500 font-mono truncate mt-0.5" title={stack.configFile || stack.workingDir}>
            <Folder class="w-3 h-3 flex-shrink-0 text-slate-400" />
            <span class="truncate">{stack.configFile || stack.workingDir}</span>
          </div>
        {/if}
      </div>
    </div>

    <!-- Right Section: Status Indicator & Unified Action Buttons -->
    <div class="flex items-center gap-3 justify-between sm:justify-end flex-wrap pl-7 sm:pl-0">
      <!-- Status Badge -->
      <div class="flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium border {allRunning ? 'bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-400 border-emerald-200 dark:border-emerald-800/60' : someRunning ? 'bg-amber-50 dark:bg-amber-950/40 text-amber-700 dark:text-amber-400 border-amber-200 dark:border-amber-800/60' : 'bg-slate-100 dark:bg-slate-800/60 text-slate-600 dark:text-slate-400 border-slate-200 dark:border-slate-700'}">
        <span class="w-2 h-2 rounded-full {allRunning ? 'bg-emerald-500 animate-pulse' : someRunning ? 'bg-amber-500' : 'bg-slate-400'}"></span>
        <span>{stack.runningCount} de {stack.totalCount} activos</span>
      </div>

      <!-- Unified Stack Controls -->
      <div class="flex items-center gap-1.5">
        {#if isStackLoading}
          <div class="flex items-center gap-1.5 px-3 py-1.5 text-xs text-violet-600 dark:text-violet-400 font-medium">
            <Loader2 class="w-3.5 h-3.5 animate-spin" />
            <span>Procesando stack...</span>
          </div>
        {:else}
          <!-- Start Stack (if some or all stopped) -->
          <button
            onclick={() => onStartStack(stack.name)}
            disabled={allRunning || !!actionLoading}
            title={allRunning ? "Todos los servicios ya están activos" : "Iniciar todos los servicios del stack"}
            class="px-2.5 py-1 rounded-lg text-xs font-medium border border-emerald-500/30 bg-emerald-50 dark:bg-emerald-950/30 text-emerald-600 dark:text-emerald-400 hover:bg-emerald-100 dark:hover:bg-emerald-900/50 transition-colors flex items-center gap-1 disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer"
          >
            <Play class="w-3.5 h-3.5 fill-current" />
            <span class="hidden md:inline">Iniciar</span>
          </button>

          <!-- Restart Stack -->
          <button
            onclick={() => onRestartStack(stack.name)}
            disabled={!!actionLoading}
            title="Reiniciar todos los contenedores de este stack"
            class="px-2.5 py-1 rounded-lg text-xs font-medium border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-200 transition-colors flex items-center gap-1 disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer"
          >
            <RotateCw class="w-3.5 h-3.5" />
            <span class="hidden md:inline">Reiniciar</span>
          </button>

          <!-- Stop Stack -->
          <button
            onclick={() => onStopStack(stack.name)}
            disabled={allStopped || !!actionLoading}
            title={allStopped ? "Todos los servicios ya están detenidos" : "Detener todos los servicios del stack"}
            class="px-2.5 py-1 rounded-lg text-xs font-medium border border-rose-200 dark:border-rose-900/60 bg-rose-50 dark:bg-rose-950/30 text-rose-600 dark:text-rose-400 hover:bg-rose-100 dark:hover:bg-rose-900/40 transition-colors flex items-center gap-1 disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer"
          >
            <Square class="w-3.5 h-3.5 fill-current" />
            <span class="hidden md:inline">Detener</span>
          </button>
        {/if}
      </div>
    </div>
  </div>

  <!-- Stack Children Containers List -->
  {#if isExpanded}
    <div class="space-y-2.5 pt-2 border-t border-violet-200/50 dark:border-violet-900/30">
      {#each stack.containers as container (container.id)}
        <ContainerCard
          {container}
          stats={statsMap[container.id]}
          {onStart}
          {onStop}
          {onRestart}
          {onPause}
          {onUnpause}
          {onRemove}
          {onOpenTerminal}
          {onViewLogs}
          {onViewStats}
          {actionLoading}
        />
      {/each}
    </div>
  {/if}
</div>
