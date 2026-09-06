<script lang="ts">
  import {
    X,
    Code2,
    RefreshCw,
    Copy,
    Check,
    Search,
    Filter,
  } from '@lucide/svelte';
  import { dockerApi } from '../../../shared/services/api';

  let {
    containerId,
    containerName,
    onClose,
  } = $props<{
    containerId: string | null;
    containerName: string;
    onClose: () => void;
  }>();

  let rawJson = $state<string>('');
  let parsedData = $state<any>(null);
  let loading = $state<boolean>(false);
  let error = $state<string | null>(null);
  let searchQuery = $state<string>('');
  let selectedSection = $state<'all' | 'Config' | 'NetworkSettings' | 'Mounts' | 'State' | 'HostConfig'>('all');
  let copied = $state<boolean>(false);

  async function loadInspect() {
    if (!containerId) return;
    loading = true;
    error = null;
    try {
      const data = await dockerApi.inspectContainer(containerId);
      rawJson = data;
      try {
        parsedData = JSON.parse(data);
      } catch {
        parsedData = null;
      }
    } catch (err: any) {
      error = err?.message || String(err);
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    if (containerId) {
      loadInspect();
    }
  });

  // Handle ESC key to close
  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      onClose();
    }
  }

  // Determine which JSON string to display based on selected section filter
  const displayedJson = $derived.by(() => {
    if (!parsedData || selectedSection === 'all') {
      return rawJson;
    }

    // Extract section if available
    let target = parsedData;
    // Docker inspect returns either an object or an array with 1 element
    if (Array.isArray(parsedData) && parsedData.length > 0) {
      target = parsedData[0];
    }

    if (target && target[selectedSection] !== undefined) {
      return JSON.stringify(target[selectedSection], null, 2);
    }

    return rawJson;
  });

  // Split lines and calculate search matches
  const lines = $derived(displayedJson ? displayedJson.split('\n') : []);

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
        // count occurrences
        let pos = 0;
        while ((pos = lower.indexOf(q, pos)) !== -1) {
          count++;
          pos += q.length;
        }
      }
    });

    return { matchesCount: count, lineIndices: matched };
  });

  function copyToClipboard() {
    if (!displayedJson) return;
    navigator.clipboard.writeText(displayedJson);
    copied = true;
    setTimeout(() => (copied = false), 2000);
  }

  // Escape HTML entities to prevent XSS
  function escapeHtml(str: string): string {
    return str
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#039;');
  }

  // Syntax highlighter with search highlighting
  function highlightLine(line: string, query: string): string {
    const escaped = escapeHtml(line);

    // Regex to match JSON tokens: key, string, number, boolean, null
    const tokenRegex = /("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g;

    let highlighted = escaped.replace(tokenRegex, (match) => {
      let cls = 'text-slate-800 dark:text-slate-200';
      if (/^"/.test(match)) {
        if (/:$/.test(match)) {
          // Key
          cls = 'text-indigo-400 dark:text-indigo-400 font-semibold';
        } else {
          // String value
          cls = 'text-emerald-400 dark:text-emerald-400';
        }
      } else if (/true|false/.test(match)) {
        // Boolean
        cls = 'text-amber-400 dark:text-amber-400 font-medium';
      } else if (/null/.test(match)) {
        // Null
        cls = 'text-rose-400 dark:text-rose-400 italic';
      } else {
        // Number
        cls = 'text-cyan-400 dark:text-cyan-400 font-mono';
      }
      return `<span class="${cls}">${match}</span>`;
    });

    // If there is an active search query, highlight it with <mark>
    if (query.trim()) {
      const escapedQuery = escapeHtml(query.trim());
      // Case insensitive match outside of html tags
      const searchRegex = new RegExp(`(?![^<]*>)(${escapedQuery.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')})`, 'gi');
      highlighted = highlighted.replace(searchRegex, '<mark class="bg-amber-400/40 text-amber-100 rounded-xs px-0.5 font-bold">$1</mark>');
    }

    return highlighted;
  }
</script>

<svelte:window onkeydown={handleKeyDown} />

{#if containerId}
  <div class="fixed inset-0 z-50 bg-slate-950/60 backdrop-blur-sm flex items-center justify-center p-3 sm:p-6 animate-in fade-in duration-150">
    <div
      class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl w-full max-w-5xl h-[88vh] flex flex-col shadow-2xl overflow-hidden animate-in zoom-in-95 duration-150"
    >
      <!-- Modal Header -->
      <div class="px-5 py-4 border-b border-slate-200 dark:border-slate-800 flex items-center justify-between gap-4 bg-slate-50/50 dark:bg-slate-950/40">
        <div class="flex items-center gap-3 min-w-0">
          <div class="p-2 rounded-xl bg-indigo-50 dark:bg-indigo-950/60 text-indigo-600 dark:text-indigo-400 border border-indigo-200 dark:border-indigo-800/60">
            <Code2 class="w-5 h-5" />
          </div>
          <div class="min-w-0">
            <div class="flex items-center gap-2 flex-wrap">
              <h2 class="text-base font-semibold text-slate-900 dark:text-slate-100 truncate">
                {containerName}
              </h2>
              <span class="text-[11px] font-mono font-medium px-2 py-0.5 rounded-full bg-indigo-50 dark:bg-indigo-950/50 text-indigo-700 dark:text-indigo-300 border border-indigo-200 dark:border-indigo-800">
                docker inspect
              </span>
            </div>
            <p class="text-xs text-slate-500 dark:text-slate-400 truncate font-mono mt-0.5">
              ID: {containerId}
            </p>
          </div>
        </div>

        <!-- Header Actions: Copy, Refresh, Close -->
        <div class="flex items-center gap-1.5 flex-shrink-0">
          <button
            onclick={copyToClipboard}
            disabled={loading || !displayedJson}
            title="Copiar JSON al portapapeles"
            class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-xs font-medium text-slate-700 dark:text-slate-200 transition-colors cursor-pointer disabled:opacity-50"
          >
            {#if copied}
              <Check class="w-3.5 h-3.5 text-emerald-500" />
              <span class="text-emerald-600 dark:text-emerald-400">Copiado</span>
            {:else}
              <Copy class="w-3.5 h-3.5 text-slate-400" />
              <span>Copiar</span>
            {/if}
          </button>

          <button
            onclick={loadInspect}
            disabled={loading}
            title="Refrescar inspección"
            class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300 transition-colors cursor-pointer disabled:opacity-50"
          >
            <RefreshCw class="w-4 h-4 {loading ? 'animate-spin text-indigo-500' : ''}" />
          </button>

          <button
            onclick={onClose}
            title="Cerrar (Esc)"
            class="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-rose-50 dark:hover:bg-rose-950/40 text-slate-400 hover:text-rose-600 dark:hover:text-rose-400 transition-colors cursor-pointer ml-1"
          >
            <X class="w-4 h-4" />
          </button>
        </div>
      </div>

      <!-- Controls Sub-Header: Search Bar & Section Filters -->
      <div class="px-5 py-3 border-b border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
        <!-- Search Input -->
        <div class="relative flex-1 max-w-md">
          <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
          <input
            type="text"
            bind:value={searchQuery}
            placeholder="Buscar clave, variable o valor (ej. ENV, 8080, IP)..."
            class="w-full pl-9 pr-8 py-1.5 text-xs rounded-lg border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/60 text-slate-900 dark:text-slate-100 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500"
          />
          {#if searchQuery}
            <button
              onclick={() => (searchQuery = '')}
              class="absolute right-2.5 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 p-0.5"
            >
              <X class="w-3.5 h-3.5" />
            </button>
          {/if}
        </div>

        <!-- Match Counter -->
        {#if searchQuery}
          <div class="flex items-center text-xs font-medium text-slate-500 dark:text-slate-400">
            {#if searchResults.matchesCount > 0}
              <span class="px-2 py-0.5 rounded-md bg-amber-50 dark:bg-amber-950/40 text-amber-700 dark:text-amber-300 border border-amber-200 dark:border-amber-800/60">
                {searchResults.matchesCount} {searchResults.matchesCount === 1 ? 'coincidencia' : 'coincidencias'}
              </span>
            {:else}
              <span class="text-rose-500 dark:text-rose-400">
                Sin coincidencias
              </span>
            {/if}
          </div>
        {/if}

        <!-- Section Chips -->
        <div class="flex items-center gap-1.5 overflow-x-auto pb-1 sm:pb-0 scrollbar-none text-xs">
          <span class="text-slate-400 text-[11px] flex items-center gap-1 mr-1">
            <Filter class="w-3 h-3" /> Sección:
          </span>

          {#each [
            { id: 'all', label: 'Todo' },
            { id: 'Config', label: 'Config / Env' },
            { id: 'NetworkSettings', label: 'Network' },
            { id: 'Mounts', label: 'Mounts' },
            { id: 'State', label: 'State' },
            { id: 'HostConfig', label: 'HostConfig' },
          ] as sec}
            <button
              onclick={() => (selectedSection = sec.id as any)}
              class="px-2.5 py-1 rounded-lg font-medium transition-colors whitespace-nowrap cursor-pointer {selectedSection === sec.id
                ? 'bg-indigo-600 text-white shadow-xs'
                : 'bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-700'}"
            >
              {sec.label}
            </button>
          {/each}
        </div>
      </div>

      <!-- JSON Viewer Body -->
      <div class="flex-1 overflow-auto bg-slate-950 p-4 font-mono text-xs text-slate-200 selection:bg-indigo-500 selection:text-white">
        {#if loading && !displayedJson}
          <div class="h-full flex flex-col items-center justify-center gap-3 text-slate-400">
            <RefreshCw class="w-6 h-6 animate-spin text-indigo-400" />
            <p>Obteniendo inspección del contenedor...</p>
          </div>
        {:else if error}
          <div class="h-full flex flex-col items-center justify-center gap-3 text-center px-4">
            <div class="p-3 rounded-xl bg-rose-950/50 text-rose-400 border border-rose-800">
              <p class="font-semibold">Error al inspeccionar el contenedor</p>
              <p class="text-xs text-rose-300 mt-1">{error}</p>
            </div>
            <button
              onclick={loadInspect}
              class="px-4 py-2 rounded-lg bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-medium cursor-pointer"
            >
              Reintentar
            </button>
          </div>
        {:else if lines.length > 0}
          <div class="space-y-0.5">
            {#each lines as line, idx}
              {@const isMatch = searchResults.lineIndices.has(idx)}
              <div
                class="flex hover:bg-slate-900/90 rounded-xs px-1 -mx-1 transition-colors {isMatch ? 'bg-amber-950/30 ring-1 ring-amber-500/30' : ''}"
              >
                <!-- Line Number -->
                <span class="w-12 flex-shrink-0 text-slate-600 select-none text-right pr-3 font-mono text-[11px]">
                  {idx + 1}
                </span>
                <!-- Syntax Highlighted Code -->
                <pre class="flex-1 whitespace-pre-wrap break-all leading-relaxed font-mono">{@html highlightLine(line, searchQuery)}</pre>
              </div>
            {/each}
          </div>
        {:else}
          <div class="h-full flex items-center justify-center text-slate-500">
            Sin datos de inspección disponibles.
          </div>
        {/if}
      </div>

      <!-- Footer Info -->
      <div class="px-5 py-2.5 border-t border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-950 text-[11px] text-slate-500 dark:text-slate-400 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <span>{lines.length} líneas</span>
          <span>•</span>
          <span>Sección: <strong class="font-mono text-slate-700 dark:text-slate-300">{selectedSection}</strong></span>
        </div>
        <div>
          Presiona <kbd class="px-1.5 py-0.5 rounded bg-slate-200 dark:bg-slate-800 font-mono text-[10px] text-slate-700 dark:text-slate-300">Esc</kbd> para cerrar
        </div>
      </div>
    </div>
  </div>
{/if}
