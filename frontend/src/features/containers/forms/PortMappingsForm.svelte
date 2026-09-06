<script lang="ts">
  import { Plus, Trash2, ArrowRight, ChevronDown } from '@lucide/svelte';
  import { uiStore } from '../../../shared/stores/ui.svelte';

  export interface PortRow {
    host: string;
    container: string;
    protocol: 'tcp' | 'udp';
  }

  let {
    ports = $bindable<PortRow[]>([]),
    disabled = false,
  } = $props<{
    ports: PortRow[];
    disabled?: boolean;
  }>();

  const portPresets = [
    { label: '80:80 (HTTP)', host: '80', container: '80' },
    { label: '8080:80 (Web)', host: '8080', container: '80' },
    { label: '3000:3000 (Dev)', host: '3000', container: '3000' },
    { label: '5432:5432 (Postgres)', host: '5432', container: '5432' },
    { label: '6379:6379 (Redis)', host: '6379', container: '6379' },
  ];

  function addPort(host = '', container = '', protocol: 'tcp' | 'udp' = 'tcp') {
    ports = [...ports, { host, container, protocol }];
  }

  function removePort(index: number) {
    ports = ports.filter((_, i) => i !== index);
  }

  function applyPreset(p: { host: string; container: string }) {
    const exists = ports.some((item) => item.host === p.host && item.container === p.container);
    if (!exists) {
      addPort(p.host, p.container, 'tcp');
    }
  }
</script>

<div class="space-y-2.5">
  <div class="flex items-center justify-between">
    <div>
      <h4 class="text-xs font-semibold text-slate-900 dark:text-slate-100 flex items-center gap-1.5">
        <span>Mapeo de Puertos</span>
        {#if ports.length > 0}
          <span class="px-1.5 py-0.2 rounded-full text-[10px] bg-blue-100 dark:bg-blue-900/50 text-blue-700 dark:text-blue-300 font-mono font-normal">
            {ports.length}
          </span>
        {/if}
      </h4>
      <p class="text-[11px] text-slate-500 dark:text-slate-400">
        Redirige puertos del equipo anfitrión (Host) al contenedor
      </p>
    </div>

    <button
      type="button"
      {disabled}
      onclick={() => addPort()}
      class="px-2.5 py-1 text-xs rounded-lg font-medium text-blue-600 dark:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-950/40 border border-blue-200 dark:border-blue-900/40 flex items-center gap-1 transition-colors cursor-pointer"
    >
      <Plus class="w-3.5 h-3.5" />
      <span>Agregar Puerto</span>
    </button>
  </div>

  {#if ports.length === 0}
    <div class="p-3 rounded-xl border border-dashed border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/30 text-center space-y-2">
      <p class="text-xs text-slate-500">
        Sin puertos configurados. Haz clic en "Agregar Puerto" o selecciona uno común:
      </p>
      <div class="flex items-center justify-center gap-1.5 flex-wrap">
        {#each portPresets as preset}
          <button
            type="button"
            {disabled}
            onclick={() => applyPreset(preset)}
            class="px-2 py-0.5 text-[11px] font-mono rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-300 hover:border-blue-400 hover:text-blue-600 transition-colors cursor-pointer"
          >
            + {preset.label}
          </button>
        {/each}
      </div>
    </div>
  {:else}
    <div class="space-y-2">
      {#each ports as port, index}
        <div class="flex items-center gap-2 p-2 rounded-xl bg-slate-50 dark:bg-slate-800/40 border border-slate-200 dark:border-slate-800 text-xs">
          <div class="flex-1 space-y-1">
            <span class="text-[10px] text-slate-500 dark:text-slate-400 font-medium">Puerto Host</span>
            <input
              type="text"
              placeholder="8080"
              {disabled}
              bind:value={port.host}
              class="w-full px-2.5 py-1.5 text-xs rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100 font-mono focus:outline-none focus:ring-1 focus:ring-blue-500"
            />
          </div>

          <div class="pt-4 text-slate-400">
            <ArrowRight class="w-4 h-4" />
          </div>

          <div class="flex-1 space-y-1">
            <span class="text-[10px] text-slate-500 dark:text-slate-400 font-medium">Puerto Contenedor</span>
            <input
              type="text"
              placeholder="80"
              {disabled}
              bind:value={port.container}
              class="w-full px-2.5 py-1.5 text-xs rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100 font-mono focus:outline-none focus:ring-1 focus:ring-blue-500"
            />
          </div>

          <div class="w-24 space-y-1">
            <span class="text-[10px] text-slate-500 dark:text-slate-400 font-medium">Protocolo</span>
            <div class="relative">
              <select
                aria-label="Protocolo de puerto"
                {disabled}
                bind:value={port.protocol}
                class="w-full appearance-none -webkit-appearance-none text-xs py-1.5 pl-2.5 pr-7 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-800 dark:text-slate-200 font-mono focus:outline-none focus:ring-1 focus:ring-blue-500 cursor-pointer shadow-xs transition-colors"
                style="background-color: {uiStore.isDark ? '#1e293b' : '#ffffff'}; color: {uiStore.isDark ? '#e2e8f0' : '#1e293b'};"
              >
                <option value="tcp" style="background-color: {uiStore.isDark ? '#1e293b' : '#ffffff'}; color: {uiStore.isDark ? '#e2e8f0' : '#1e293b'};">TCP</option>
                <option value="udp" style="background-color: {uiStore.isDark ? '#1e293b' : '#ffffff'}; color: {uiStore.isDark ? '#e2e8f0' : '#1e293b'};">UDP</option>
              </select>
              <ChevronDown class="w-3.5 h-3.5 absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" />
            </div>
          </div>

          <div class="pt-4">
            <button
              type="button"
              {disabled}
              onclick={() => removePort(index)}
              title="Eliminar puerto"
              class="p-1.5 rounded-lg text-slate-400 hover:text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-950/40 transition-colors cursor-pointer"
            >
              <Trash2 class="w-4 h-4" />
            </button>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>
