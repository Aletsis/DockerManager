<script lang="ts">
  import { Plus, Trash2, ArrowRight, ChevronDown } from '@lucide/svelte';
  import { uiStore } from '../../../shared/stores/ui.svelte';

  export interface VolumeRow {
    hostPath: string;
    containerPath: string;
    mode: 'rw' | 'ro';
  }

  let {
    volumes = $bindable<VolumeRow[]>([]),
    disabled = false,
  } = $props<{
    volumes: VolumeRow[];
    disabled?: boolean;
  }>();

  function addVolume(hostPath = '', containerPath = '', mode: 'rw' | 'ro' = 'rw') {
    volumes = [...volumes, { hostPath, containerPath, mode }];
  }

  function removeVolume(index: number) {
    volumes = volumes.filter((_, i) => i !== index);
  }
</script>

<div class="space-y-2.5">
  <div class="flex items-center justify-between">
    <div>
      <h4 class="text-xs font-semibold text-slate-900 dark:text-slate-100 flex items-center gap-1.5">
        <span>Volúmenes y Binds</span>
        {#if volumes.length > 0}
          <span class="px-1.5 py-0.2 rounded-full text-[10px] bg-indigo-100 dark:bg-indigo-900/50 text-indigo-700 dark:text-indigo-300 font-mono font-normal">
            {volumes.length}
          </span>
        {/if}
      </h4>
      <p class="text-[11px] text-slate-500 dark:text-slate-400">
        Monta carpetas locales del host o volúmenes dentro del contenedor
      </p>
    </div>

    <button
      type="button"
      {disabled}
      onclick={() => addVolume()}
      class="px-2.5 py-1 text-xs rounded-lg font-medium text-indigo-600 dark:text-indigo-400 hover:bg-indigo-50 dark:hover:bg-indigo-950/40 border border-indigo-200 dark:border-indigo-900/40 flex items-center gap-1 transition-colors cursor-pointer"
    >
      <Plus class="w-3.5 h-3.5" />
      <span>Agregar Volumen</span>
    </button>
  </div>

  {#if volumes.length === 0}
    <div class="p-3 rounded-xl border border-dashed border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/30 text-center">
      <p class="text-xs text-slate-500">
        Sin volúmenes montados. Los datos almacenados se guardarán temporalmente en la capa del contenedor.
      </p>
    </div>
  {:else}
    <div class="space-y-2">
      {#each volumes as vol, index}
        <div class="flex items-center gap-2 p-2 rounded-xl bg-slate-50 dark:bg-slate-800/40 border border-slate-200 dark:border-slate-800 text-xs">
          <div class="flex-1 space-y-1">
            <span class="text-[10px] text-slate-500 dark:text-slate-400 font-medium">Ruta Host o Volumen</span>
            <input
              type="text"
              placeholder="/home/user/data o pg_vol"
              {disabled}
              bind:value={vol.hostPath}
              class="w-full px-2.5 py-1.5 text-xs rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100 font-mono focus:outline-none focus:ring-1 focus:ring-blue-500"
            />
          </div>

          <div class="pt-4 text-slate-400">
            <ArrowRight class="w-4 h-4" />
          </div>

          <div class="flex-1 space-y-1">
            <span class="text-[10px] text-slate-500 dark:text-slate-400 font-medium">Ruta Contenedor</span>
            <input
              type="text"
              placeholder="/var/lib/data"
              {disabled}
              bind:value={vol.containerPath}
              class="w-full px-2.5 py-1.5 text-xs rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100 font-mono focus:outline-none focus:ring-1 focus:ring-blue-500"
            />
          </div>

          <div class="w-32 sm:w-36 space-y-1">
            <span class="text-[10px] text-slate-500 dark:text-slate-400 font-medium">Modo</span>
            <div class="relative">
              <select
                aria-label="Modo de volumen"
                {disabled}
                bind:value={vol.mode}
                class="w-full appearance-none -webkit-appearance-none text-xs py-1.5 pl-2.5 pr-7 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-800 dark:text-slate-200 font-mono focus:outline-none focus:ring-1 focus:ring-blue-500 cursor-pointer shadow-xs transition-colors truncate"
                style="background-color: {uiStore.isDark ? '#1e293b' : '#ffffff'}; color: {uiStore.isDark ? '#e2e8f0' : '#1e293b'};"
              >
                <option value="rw" style="background-color: {uiStore.isDark ? '#1e293b' : '#ffffff'}; color: {uiStore.isDark ? '#e2e8f0' : '#1e293b'};">Lectura/Escritura (rw)</option>
                <option value="ro" style="background-color: {uiStore.isDark ? '#1e293b' : '#ffffff'}; color: {uiStore.isDark ? '#e2e8f0' : '#1e293b'};">Solo Lectura (ro)</option>
              </select>
              <ChevronDown class="w-3.5 h-3.5 absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" />
            </div>
          </div>

          <div class="pt-4">
            <button
              type="button"
              {disabled}
              onclick={() => removeVolume(index)}
              title="Eliminar volumen"
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
