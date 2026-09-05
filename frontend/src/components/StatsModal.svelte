<script lang="ts">
  import {
    X,
    Activity,
    Cpu,
    Database,
    ArrowDownUp,
    HardDrive,
    Users,
  } from '@lucide/svelte';
  import type { ContainerStats } from '../types';
  import { GetContainerStats } from '../../wailsjs/go/main/App';
  import { formatBytes } from '../utils';

  let {
    containerId,
    containerName,
    onClose,
  } = $props<{
    containerId: string | null;
    containerName: string;
    onClose: () => void;
  }>();

  let stats = $state<ContainerStats | null>(null);
  let error = $state<string | null>(null);

  async function fetchStats() {
    if (!containerId) return;
    try {
      const data = await GetContainerStats(containerId);
      stats = data as unknown as ContainerStats;
      error = null;
    } catch (err) {
      error = `Error al obtener estadísticas: ${err}`;
    }
  }

  $effect(() => {
    if (containerId) {
      fetchStats();
      const timer = setInterval(fetchStats, 1500);
      return () => clearInterval(timer);
    }
  });
</script>

{#if containerId}
  <div class="fixed inset-0 z-50 bg-slate-950/60 backdrop-blur-sm flex items-center justify-center p-4">
    <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl w-full max-w-2xl shadow-2xl overflow-hidden animate-in fade-in zoom-in-95 duration-150">
      <!-- Modal Header -->
      <div class="flex items-center justify-between px-5 py-4 border-b border-slate-200 dark:border-slate-800">
        <div class="flex items-center gap-2.5">
          <div class="p-1.5 rounded-lg bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
            <Activity class="w-4 h-4" />
          </div>
          <div>
            <h2 class="text-sm font-semibold text-slate-900 dark:text-slate-100">
              Monitoreo en vivo: <span class="text-emerald-600 dark:text-emerald-400">{containerName}</span>
            </h2>
            <p class="text-[11px] text-slate-500">Métricas actualizadas cada 1.5s</p>
          </div>
        </div>

        <button
          onclick={onClose}
          class="p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 transition-colors"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Modal Body -->
      <div class="p-6 space-y-5">
        {#if error}
          <div class="p-3 rounded-lg bg-rose-500/10 border border-rose-500/20 text-rose-600 dark:text-rose-400 text-xs">
            {error}
          </div>
        {/if}

        {#if stats}
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <!-- CPU Usage Card -->
            <div class="p-4 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/40 space-y-3">
              <div class="flex items-center justify-between text-xs text-slate-500">
                <div class="flex items-center gap-1.5 font-medium">
                  <Cpu class="w-4 h-4 text-indigo-500" />
                  <span>Uso de CPU</span>
                </div>
                <span class="font-mono text-base font-semibold text-slate-900 dark:text-slate-100">
                  {stats.cpuPercentage.toFixed(2)}%
                </span>
              </div>
              <div class="w-full h-2 bg-slate-200 dark:bg-slate-700 rounded-full overflow-hidden">
                <div
                  class="h-full transition-all duration-300 {stats.cpuPercentage > 80 ? 'bg-rose-500' : stats.cpuPercentage > 50 ? 'bg-amber-500' : 'bg-indigo-500'}"
                  style="width: {Math.min(stats.cpuPercentage, 100)}%;"
                ></div>
              </div>
            </div>

            <!-- Memory Usage Card -->
            <div class="p-4 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/40 space-y-3">
              <div class="flex items-center justify-between text-xs text-slate-500">
                <div class="flex items-center gap-1.5 font-medium">
                  <Database class="w-4 h-4 text-emerald-500" />
                  <span>Memoria RAM</span>
                </div>
                <span class="font-mono text-base font-semibold text-slate-900 dark:text-slate-100">
                  {stats.memoryPercentage.toFixed(1)}%
                </span>
              </div>
              <div class="w-full h-2 bg-slate-200 dark:bg-slate-700 rounded-full overflow-hidden">
                <div
                  class="h-full transition-all duration-300 {stats.memoryPercentage > 85 ? 'bg-rose-500' : stats.memoryPercentage > 65 ? 'bg-amber-500' : 'bg-emerald-500'}"
                  style="width: {Math.min(stats.memoryPercentage, 100)}%;"
                ></div>
              </div>
              <div class="flex justify-between text-[11px] text-slate-500 font-mono">
                <span>{formatBytes(stats.memoryUsage)}</span>
                <span>{formatBytes(stats.memoryLimit)}</span>
              </div>
            </div>

            <!-- Network I/O -->
            <div class="p-4 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/40 space-y-2">
              <div class="flex items-center gap-1.5 text-xs font-medium text-slate-500">
                <ArrowDownUp class="w-4 h-4 text-blue-500" />
                <span>Red I/O</span>
              </div>
              <div class="grid grid-cols-2 gap-2 text-xs font-mono pt-1">
                <div>
                  <span class="text-[10px] text-slate-400 block">Descarga (Rx)</span>
                  <span class="font-semibold text-slate-800 dark:text-slate-200">
                    {formatBytes(stats.networkRx)}
                  </span>
                </div>
                <div>
                  <span class="text-[10px] text-slate-400 block">Subida (Tx)</span>
                  <span class="font-semibold text-slate-800 dark:text-slate-200">
                    {formatBytes(stats.networkTx)}
                  </span>
                </div>
              </div>
            </div>

            <!-- Block I/O & PIDs -->
            <div class="p-4 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/40 space-y-2">
              <div class="flex items-center justify-between text-xs font-medium text-slate-500">
                <div class="flex items-center gap-1.5">
                  <HardDrive class="w-4 h-4 text-purple-500" />
                  <span>Disco & Procesos</span>
                </div>
                <div class="flex items-center gap-1 text-slate-600 dark:text-slate-300 font-mono">
                  <Users class="w-3.5 h-3.5 text-slate-400" />
                  <span>{stats.pids} PIDs</span>
                </div>
              </div>
              <div class="grid grid-cols-2 gap-2 text-xs font-mono pt-1">
                <div>
                  <span class="text-[10px] text-slate-400 block">Lectura</span>
                  <span class="font-semibold text-slate-800 dark:text-slate-200">
                    {formatBytes(stats.blockRead)}
                  </span>
                </div>
                <div>
                  <span class="text-[10px] text-slate-400 block">Escritura</span>
                  <span class="font-semibold text-slate-800 dark:text-slate-200">
                    {formatBytes(stats.blockWrite)}
                  </span>
                </div>
              </div>
            </div>
          </div>
        {:else}
          <div class="py-12 text-center text-xs text-slate-400">
            Conectando con el stream de métricas del contenedor...
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}
