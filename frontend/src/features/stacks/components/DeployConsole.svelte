<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { Terminal, Trash2, Copy, Check, Loader2, CheckCircle2, AlertCircle } from '@lucide/svelte';
  import { EventsOn, EventsOff } from '../../../../wailsjs/runtime/runtime';

  let {
    isRunning = false,
    status = 'idle', // 'idle' | 'running' | 'success' | 'error'
    statusMessage = '',
  } = $props<{
    isRunning?: boolean;
    status?: 'idle' | 'running' | 'success' | 'error';
    statusMessage?: string;
  }>();

  let logs = $state<string[]>([]);
  let consoleEl: HTMLDivElement | null = $state(null);
  let copied = $state(false);

  function handleLogOutput(line: string) {
    logs.push(line);
    setTimeout(() => {
      if (consoleEl) {
        consoleEl.scrollTop = consoleEl.scrollHeight;
      }
    }, 10);
  }

  export function clear() {
    logs = [];
  }

  export function addLog(line: string) {
    handleLogOutput(line);
  }

  onMount(() => {
    EventsOn('stack:deploy:output', handleLogOutput);
  });

  onDestroy(() => {
    EventsOff('stack:deploy:output');
  });

  async function handleCopy() {
    if (logs.length === 0) return;
    try {
      await navigator.clipboard.writeText(logs.join('\n'));
      copied = true;
      setTimeout(() => (copied = false), 2000);
    } catch {
      // ignore
    }
  }
</script>

<div class="rounded-xl border border-slate-800 bg-slate-950 text-slate-200 overflow-hidden font-mono text-xs flex flex-col shadow-inner">
  <!-- Console Header -->
  <div class="flex items-center justify-between px-3 py-2 bg-slate-900/90 border-b border-slate-800 select-none">
    <div class="flex items-center gap-2 text-[11px]">
      <Terminal class="w-3.5 h-3.5 text-violet-400" />
      <span class="font-semibold text-slate-300">Consola de Despliegue</span>

      {#if status === 'running'}
        <span class="flex items-center gap-1 text-amber-400 bg-amber-950/60 px-2 py-0.5 rounded-full border border-amber-800/40 text-[10px]">
          <Loader2 class="w-2.5 h-2.5 animate-spin" />
          <span>Ejecutando...</span>
        </span>
      {:else if status === 'success'}
        <span class="flex items-center gap-1 text-emerald-400 bg-emerald-950/60 px-2 py-0.5 rounded-full border border-emerald-800/40 text-[10px]">
          <CheckCircle2 class="w-2.5 h-2.5" />
          <span>Completado</span>
        </span>
      {:else if status === 'error'}
        <span class="flex items-center gap-1 text-rose-400 bg-rose-950/60 px-2 py-0.5 rounded-full border border-rose-800/40 text-[10px]">
          <AlertCircle class="w-2.5 h-2.5" />
          <span>Error</span>
        </span>
      {/if}
    </div>

    <div class="flex items-center gap-1 text-[11px]">
      <button
        type="button"
        onclick={handleCopy}
        disabled={logs.length === 0}
        class="flex items-center gap-1 px-2 py-0.5 rounded hover:bg-slate-800 text-slate-400 hover:text-slate-200 disabled:opacity-40 transition-colors cursor-pointer"
        title="Copiar salida"
      >
        {#if copied}
          <Check class="w-3 h-3 text-emerald-400" />
        {:else}
          <Copy class="w-3 h-3" />
        {/if}
        <span>Copiar</span>
      </button>

      <button
        type="button"
        onclick={clear}
        disabled={logs.length === 0 || isRunning}
        class="flex items-center gap-1 px-2 py-0.5 rounded hover:bg-slate-800 text-slate-400 hover:text-rose-400 disabled:opacity-40 transition-colors cursor-pointer"
        title="Limpiar consola"
      >
        <Trash2 class="w-3 h-3" />
        <span>Limpiar</span>
      </button>
    </div>
  </div>

  <!-- Logs Scroll Window -->
  <div
    bind:this={consoleEl}
    class="p-3 h-44 overflow-y-auto space-y-1 font-mono text-[11px] leading-relaxed select-text"
  >
    {#if logs.length === 0}
      <div class="text-slate-500 italic py-2">
        {isRunning ? 'Esperando salida de Docker Compose...' : 'Esperando ejecución del stack...'}
      </div>
    {:else}
      {#each logs as line}
        <div class="whitespace-pre-wrap break-all text-slate-300 font-mono">
          {#if line.toLowerCase().includes('error') || line.toLowerCase().includes('fail')}
            <span class="text-rose-400">{line}</span>
          {:else if line.toLowerCase().includes('running') || line.toLowerCase().includes('started') || line.toLowerCase().includes('healthy')}
            <span class="text-emerald-400">{line}</span>
          {:else if line.toLowerCase().includes('warning') || line.toLowerCase().includes('warn')}
            <span class="text-amber-400">{line}</span>
          {:else}
            <span>{line}</span>
          {/if}
        </div>
      {/each}
    {/if}

    {#if statusMessage}
      <div class="pt-2 text-xs font-semibold {status === 'error' ? 'text-rose-400' : status === 'success' ? 'text-emerald-400' : 'text-slate-400'}">
        {statusMessage}
      </div>
    {/if}
  </div>
</div>
