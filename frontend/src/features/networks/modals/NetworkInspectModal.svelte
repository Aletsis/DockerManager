<script lang="ts">
  import {
    X,
    Code2,
    RefreshCw,
    Copy,
    Check,
    Search,
    Network,
  } from '@lucide/svelte';
  import { networksStore } from '../stores/networks.svelte';

  let {
    networkId,
    networkName,
    onClose,
  } = $props<{
    networkId: string | null;
    networkName: string;
    onClose: () => void;
  }>();

  let rawJson = $state<string>('');
  let loading = $state<boolean>(false);
  let error = $state<string | null>(null);
  let searchQuery = $state<string>('');
  let copied = $state<boolean>(false);

  async function loadInspect() {
    if (!networkId) return;
    loading = true;
    error = null;
    try {
      const data = await networksStore.inspectNetwork(networkId);
      rawJson = data;
    } catch (err: any) {
      error = err?.message || String(err);
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    if (networkId) {
      loadInspect();
    }
  });

  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      onClose();
    }
  }

  function handleCopy() {
    if (!rawJson) return;
    navigator.clipboard.writeText(rawJson);
    copied = true;
    setTimeout(() => {
      copied = false;
    }, 2000);
  }

  const lines = $derived(rawJson ? rawJson.split('\n') : []);

  const searchResults = $derived.by(() => {
    const q = searchQuery.trim().toLowerCase();
    if (!q) {
      return { matchesCount: 0, lineIndices: new Set<number>() };
    }

    let count = 0;
    const matched = new Set<number>();

    lines.forEach((line, idx) => {
      const lower = line.toLowerCase();
      if (lower.includes(q)) {
        matched.add(idx);
        let pos = 0;
        while ((pos = lower.indexOf(q, pos)) !== -1) {
          count++;
          pos += q.length;
        }
      }
    });

    return { matchesCount: count, lineIndices: matched };
  });

  function highlightLine(line: string, query: string): string {
    if (!query.trim()) {
      return line
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;');
    }

    const escapedLine = line
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;');

    const escapedQuery = query
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/[.*+?^${}()|[\]\\]/g, '\\$&');

    const regex = new RegExp(`(${escapedQuery})`, 'gi');
    return escapedLine.replace(regex, '<mark class="bg-amber-300 dark:bg-amber-500/50 text-slate-900 dark:text-slate-100 rounded-xs px-0.5">$1</mark>');
  }
</script>

<svelte:window onkeydown={handleKeyDown} />

{#if networkId}
  <!-- Backdrop -->
  <div
    class="fixed inset-0 z-50 bg-slate-950/60 backdrop-blur-xs flex items-center justify-center p-3 sm:p-4 animate-in fade-in duration-150"
    role="dialog"
    aria-modal="true"
  >
    <!-- Modal Container -->
    <div
      class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl w-full max-w-4xl max-h-[90vh] flex flex-col shadow-2xl overflow-hidden animate-in zoom-in-95 duration-150"
    >
      <!-- Modal Header -->
      <div class="px-5 py-4 border-b border-slate-200 dark:border-slate-800 flex items-center justify-between shrink-0 bg-slate-50/50 dark:bg-slate-800/30">
        <div class="flex items-center gap-2.5 min-w-0">
          <div class="p-2 rounded-xl bg-teal-500/10 text-teal-600 dark:text-teal-400 shrink-0">
            <Network class="w-5 h-5" />
          </div>
          <div class="min-w-0">
            <div class="flex items-center gap-2">
              <h3 class="font-semibold text-base text-slate-900 dark:text-slate-100 tracking-tight leading-tight">
                Inspección de Red
              </h3>
              <span class="text-xs font-mono px-2 py-0.5 rounded bg-teal-100 dark:bg-teal-900/40 text-teal-700 dark:text-teal-300">
                JSON
              </span>
            </div>
            <p class="text-xs text-slate-500 dark:text-slate-400 font-mono truncate max-w-md sm:max-w-xl">
              {networkName} ({networkId.length > 12 ? networkId.slice(0, 12) : networkId})
            </p>
          </div>
        </div>

        <button
          onclick={onClose}
          class="p-2 rounded-xl text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer"
          title="Cerrar (Esc)"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Controls Toolbar -->
      <div class="px-5 py-3 border-b border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 flex flex-wrap items-center justify-between gap-3 shrink-0">
        <!-- Search bar -->
        <div class="relative flex-1 min-w-[200px]">
          <Search class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" />
          <input
            type="text"
            placeholder="Buscar en el JSON..."
            bind:value={searchQuery}
            class="w-full pl-9 pr-8 py-1.5 text-xs rounded-xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/60 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-1 focus:ring-teal-500"
          />
          {#if searchQuery}
            <button
              onclick={() => (searchQuery = '')}
              class="absolute right-2.5 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 cursor-pointer"
            >
              <X class="w-3.5 h-3.5" />
            </button>
          {/if}
        </div>

        <!-- Matches feedback -->
        {#if searchQuery.trim()}
          <span class="text-xs font-medium text-slate-500 dark:text-slate-400">
            {searchResults.matchesCount} coincidencia{searchResults.matchesCount === 1 ? '' : 's'}
          </span>
        {/if}

        <!-- Actions -->
        <div class="flex items-center gap-2">
          <button
            onclick={loadInspect}
            disabled={loading}
            class="flex items-center gap-1.5 px-3 py-1.5 rounded-xl border border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-800 text-xs font-medium text-slate-700 dark:text-slate-300 transition-colors cursor-pointer disabled:opacity-50"
            title="Recargar inspección"
          >
            <RefreshCw class="w-3.5 h-3.5 {loading ? 'animate-spin text-teal-500' : ''}" />
            <span>Recargar</span>
          </button>

          <button
            onclick={handleCopy}
            disabled={!rawJson || loading}
            class="flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-slate-900 text-white dark:bg-slate-100 dark:text-slate-900 hover:bg-slate-800 dark:hover:bg-slate-200 text-xs font-medium transition-colors cursor-pointer shadow-xs disabled:opacity-50"
          >
            {#if copied}
              <Check class="w-3.5 h-3.5 text-emerald-400 dark:text-emerald-600" />
              <span>Copiado</span>
            {:else}
              <Copy class="w-3.5 h-3.5" />
              <span>Copiar JSON</span>
            {/if}
          </button>
        </div>
      </div>

      <!-- JSON Viewer Body -->
      <div class="flex-1 overflow-auto bg-slate-950 text-slate-200 p-4 font-mono text-xs leading-relaxed select-text">
        {#if loading}
          <div class="h-64 flex flex-col items-center justify-center gap-3 text-slate-400">
            <RefreshCw class="w-6 h-6 animate-spin text-teal-400" />
            <span class="text-xs">Inspeccionando red en Docker daemon...</span>
          </div>
        {:else if error}
          <div class="h-64 flex flex-col items-center justify-center gap-2 text-rose-400 p-6 text-center">
            <span class="font-semibold text-sm">Error al inspeccionar red</span>
            <p class="text-xs max-w-md text-rose-300/80">{error}</p>
          </div>
        {:else if lines.length > 0}
          <pre class="whitespace-pre"><code>{#each lines as line, index}
<div class="table-row {searchResults.lineIndices.has(index) ? 'bg-teal-950/50' : ''}"><span class="table-cell text-right pr-4 select-none text-slate-600 opacity-60 w-10">{index + 1}</span><span class="table-cell">{@html highlightLine(line, searchQuery)}</span></div>{/each}</code></pre>
        {:else}
          <div class="h-64 flex items-center justify-center text-slate-500 text-xs">
            Sin datos de inspección disponibles.
          </div>
        {/if}
      </div>

      <!-- Modal Footer -->
      <div class="px-5 py-3 border-t border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/90 flex items-center justify-between text-xs text-slate-500 dark:text-slate-400">
        <div>
          Líneas: <span class="font-mono text-slate-700 dark:text-slate-300">{lines.length}</span>
        </div>
        <button
          onclick={onClose}
          class="px-4 py-1.5 rounded-xl border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-300 font-medium transition-colors cursor-pointer"
        >
          Cerrar
        </button>
      </div>
    </div>
  </div>
{/if}
