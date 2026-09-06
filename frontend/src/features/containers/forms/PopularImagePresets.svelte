<script lang="ts">
  import { Sparkles } from '@lucide/svelte';

  export interface ImagePreset {
    name: string;
    label: string;
    defaultPort?: { host: string; container: string };
    defaultEnv?: { key: string; value: string };
  }

  const popularImages: ImagePreset[] = [
    { name: 'nginx:alpine', label: 'Nginx', defaultPort: { host: '8080', container: '80' } },
    { name: 'postgres:16-alpine', label: 'Postgres', defaultPort: { host: '5432', container: '5432' }, defaultEnv: { key: 'POSTGRES_PASSWORD', value: 'postgres' } },
    { name: 'redis:alpine', label: 'Redis', defaultPort: { host: '6379', container: '6379' } },
    { name: 'node:20-alpine', label: 'Node.js', defaultPort: { host: '3000', container: '3000' } },
    { name: 'python:3.11-slim', label: 'Python' },
  ];

  let {
    disabled = false,
    onSelect,
  } = $props<{
    disabled?: boolean;
    onSelect: (preset: ImagePreset) => void;
  }>();
</script>

<div class="space-y-2">
  <div class="flex items-center gap-1.5 text-xs text-slate-500 dark:text-slate-400 font-medium">
    <Sparkles class="w-3.5 h-3.5 text-amber-500" />
    <span>Imágenes Populares de Inicio Rápido</span>
  </div>
  <div class="flex items-center gap-2 flex-wrap">
    {#each popularImages as preset}
      <button
        type="button"
        {disabled}
        onclick={() => onSelect(preset)}
        class="px-2.5 py-1 text-xs rounded-lg border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/80 hover:bg-blue-50 dark:hover:bg-blue-950/40 hover:border-blue-300 dark:hover:border-blue-700 text-slate-700 dark:text-slate-200 font-medium transition-colors cursor-pointer"
      >
        {preset.label}
      </button>
    {/each}
  </div>
</div>
