<script lang="ts">
  import {
    X,
    Terminal,
    RefreshCw,
    Copy,
    Check,
    Search,
    ArrowDown,
  } from '@lucide/svelte';
  import { GetContainerLogs } from '../../../../wailsjs/go/main/App';

  let {
    containerId,
    containerName,
    onClose,
  } = $props<{
    containerId: string | null;
    containerName: string;
    onClose: () => void;
  }>();

  let logs = $state<string>('');
  let tail = $state<number>(200);
  let filter = $state<string>('');
  let loading = $state<boolean>(false);
  let autoScroll = $state<boolean>(true);
  let copied = $state<boolean>(false);
  let logsEndEl: HTMLDivElement | null = $state(null);

  async function fetchLogs() {
    if (!containerId) return;
    loading = true;
    try {
      const data = await GetContainerLogs(containerId, tail);
      logs = data || 'Sin registros emitidos.';
    } catch (err) {
      logs = `Error al obtener logs: ${err}`;
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    if (containerId) {
      fetchLogs();
      const timer = setInterval(() => {
        if (autoScroll) {
          fetchLogs();
        }
      }, 3000);
      return () => clearInterval(timer);
    }
  });

  $effect(() => {
    if (autoScroll && logsEndEl && logs) {
      logsEndEl.scrollIntoView({ behavior: 'smooth' });
    }
  });

  function copyLogs() {
    navigator.clipboard.writeText(logs);
    copied = true;
    setTimeout(() => (copied = false), 2000);
  }

  const filteredLogs = $derived(
    filter
      ? logs
          .split('\n')
          .filter((line) => line.toLowerCase().includes(filter.toLowerCase()))
          .join('\n')
      : logs
  );
</script>

{#if containerId}
  <div class="fixed inset-0 z-50 bg-slate-950/60 backdrop-blur-sm flex items-center justify-center p-4">
    <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl w-full max-w-4xl h-[80vh] flex flex-col shadow-2xl overflow-hidden">
      <!-- Modal Header -->
      <div class="flex items-center justify-between px-5 py-3.5 border-b border-slate-200 dark:border-slate-800">
        <div class="flex items-center gap-2.5">
          <div class="p-1.5 rounded-lg bg-indigo-500/10 text-indigo-500">
            <Terminal class="w-4 h-4" />
          </div>
          <div>
            <h2 class="text-sm font-semibold text-slate-900 dark:text-slate-100">
              Logs en tiempo real: <span class="text-indigo-600 dark:text-indigo-400">{containerName}</span>
            </h2>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <!-- Search filter -->
          <div class="relative">
            <Search class="w-3.5 h-3.5 absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" />
            <input
              type="text"
              placeholder="Filtrar..."
              bind:value={filter}
              class="text-xs pl-7 pr-2.5 py-1 rounded-lg border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-1 focus:ring-indigo-500 w-36"
            />
          </div>

          <!-- Lines tail -->
          <select
            aria-label="Líneas de registro"
            bind:value={tail}
            class="text-xs py-1 px-2.5 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-800 dark:text-slate-200"
          >
            <option value={50} class="bg-white dark:bg-slate-800 text-slate-800 dark:text-slate-200">50 líneas</option>
            <option value={100} class="bg-white dark:bg-slate-800 text-slate-800 dark:text-slate-200">100 líneas</option>
            <option value={200} class="bg-white dark:bg-slate-800 text-slate-800 dark:text-slate-200">200 líneas</option>
            <option value={500} class="bg-white dark:bg-slate-800 text-slate-800 dark:text-slate-200">500 líneas</option>
          </select>

          <!-- Auto-scroll toggle -->
          <button
            onclick={() => (autoScroll = !autoScroll)}
            title={autoScroll ? 'Auto-scroll activo' : 'Auto-scroll pausado'}
            class="p-1.5 rounded-lg border text-xs flex items-center gap-1 transition-colors {autoScroll ? 'border-indigo-500/40 bg-indigo-50 dark:bg-indigo-950/40 text-indigo-600 dark:text-indigo-400' : 'border-slate-200 dark:border-slate-700 text-slate-500'}"
          >
            <ArrowDown class="w-3.5 h-3.5" />
          </button>

          <!-- Copy Logs -->
          <button
            onclick={copyLogs}
            title="Copiar todo al portapapeles"
            class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300 transition-colors"
          >
            {#if copied}
              <Check class="w-3.5 h-3.5 text-emerald-500" />
            {:else}
              <Copy class="w-3.5 h-3.5" />
            {/if}
          </button>

          <!-- Refresh -->
          <button
            onclick={fetchLogs}
            title="Refrescar ahora"
            class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300 transition-colors"
          >
            <RefreshCw class="w-3.5 h-3.5 {loading ? 'animate-spin text-indigo-500' : ''}" />
          </button>

          <!-- Close -->
          <button
            onclick={onClose}
            class="p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 transition-colors"
          >
            <X class="w-4 h-4" />
          </button>
        </div>
      </div>

      <!-- Logs Output Area -->
      <div class="flex-1 bg-slate-950 p-4 overflow-auto font-mono text-xs text-slate-200 leading-relaxed selection:bg-indigo-500/40">
        <pre class="whitespace-pre-wrap break-all">{filteredLogs}</pre>
        <div bind:this={logsEndEl}></div>
      </div>
    </div>
  </div>
{/if}
