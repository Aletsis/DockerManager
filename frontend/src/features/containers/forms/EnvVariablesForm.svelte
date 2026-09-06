<script lang="ts">
  import { Plus, Trash2 } from '@lucide/svelte';

  export interface EnvRow {
    key: string;
    value: string;
  }

  let {
    envVars = $bindable<EnvRow[]>([]),
    disabled = false,
  } = $props<{
    envVars: EnvRow[];
    disabled?: boolean;
  }>();

  function addEnv(key = '', value = '') {
    envVars = [...envVars, { key, value }];
  }

  function removeEnv(index: number) {
    envVars = envVars.filter((_, i) => i !== index);
  }
</script>

<div class="space-y-2.5">
  <div class="flex items-center justify-between">
    <div>
      <h4 class="text-xs font-semibold text-slate-900 dark:text-slate-100 flex items-center gap-1.5">
        <span>Variables de Entorno (ENV)</span>
        {#if envVars.length > 0}
          <span class="px-1.5 py-0.2 rounded-full text-[10px] bg-emerald-100 dark:bg-emerald-900/50 text-emerald-700 dark:text-emerald-300 font-mono font-normal">
            {envVars.length}
          </span>
        {/if}
      </h4>
      <p class="text-[11px] text-slate-500 dark:text-slate-400">
        Configura parámetros de configuración o credenciales para el contenedor
      </p>
    </div>

    <button
      type="button"
      {disabled}
      onclick={() => addEnv()}
      class="px-2.5 py-1 text-xs rounded-lg font-medium text-emerald-600 dark:text-emerald-400 hover:bg-emerald-50 dark:hover:bg-emerald-950/40 border border-emerald-200 dark:border-emerald-900/40 flex items-center gap-1 transition-colors cursor-pointer"
    >
      <Plus class="w-3.5 h-3.5" />
      <span>Agregar Variable</span>
    </button>
  </div>

  {#if envVars.length === 0}
    <div class="p-3 rounded-xl border border-dashed border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/30 text-center">
      <p class="text-xs text-slate-500">
        Sin variables de entorno adicionales.
      </p>
    </div>
  {:else}
    <div class="space-y-2">
      {#each envVars as env, index}
        <div class="flex items-center gap-2 p-2 rounded-xl bg-slate-50 dark:bg-slate-800/40 border border-slate-200 dark:border-slate-800 text-xs">
          <div class="flex-1 space-y-1">
            <span class="text-[10px] text-slate-500 dark:text-slate-400 font-medium">Clave (Key)</span>
            <input
              type="text"
              placeholder="POSTGRES_PASSWORD"
              {disabled}
              bind:value={env.key}
              class="w-full px-2.5 py-1.5 text-xs rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100 font-mono focus:outline-none focus:ring-1 focus:ring-blue-500"
            />
          </div>

          <div class="pt-4 text-slate-400 font-bold">=</div>

          <div class="flex-1 space-y-1">
            <span class="text-[10px] text-slate-500 dark:text-slate-400 font-medium">Valor (Value)</span>
            <input
              type="text"
              placeholder="secreto123"
              {disabled}
              bind:value={env.value}
              class="w-full px-2.5 py-1.5 text-xs rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100 font-mono focus:outline-none focus:ring-1 focus:ring-blue-500"
            />
          </div>

          <div class="pt-4">
            <button
              type="button"
              {disabled}
              onclick={() => removeEnv(index)}
              title="Eliminar variable"
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
