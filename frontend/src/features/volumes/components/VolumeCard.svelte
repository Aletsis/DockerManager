<script lang="ts">
  import {
    Database,
    Folder,
    HardDrive,
    Trash2,
    Code2,
    Copy,
    Check,
    AlertTriangle,
    Link2,
    Unlink,
    Box,
    ExternalLink,
  } from '@lucide/svelte';
  import type { VolumeInfo } from '../../../types';
  import { formatBytes } from '../../../shared/utils/utils';

  let {
    volume,
    onInspect,
    onDelete,
    isDeleting = false,
  } = $props<{
    volume: VolumeInfo;
    onInspect: (vol: VolumeInfo) => void;
    onDelete: (vol: VolumeInfo) => void;
    isDeleting?: boolean;
  }>();

  let copiedName = $state(false);
  let copiedPath = $state(false);

  function copyText(text: string, type: 'name' | 'path', e: MouseEvent) {
    e.stopPropagation();
    navigator.clipboard.writeText(text);
    if (type === 'name') {
      copiedName = true;
      setTimeout(() => (copiedName = false), 1500);
    } else {
      copiedPath = true;
      setTimeout(() => (copiedPath = false), 1500);
    }
  }
</script>

<div
  class="bg-white dark:bg-slate-900 border rounded-2xl p-4 sm:p-5 shadow-xs transition-all duration-200 flex flex-col justify-between hover:shadow-md hover:border-slate-300 dark:hover:border-slate-700 {volume.inUse ? 'border-slate-200 dark:border-slate-800' : 'border-amber-200/80 dark:border-amber-900/40 bg-gradient-to-b from-white to-amber-50/20 dark:from-slate-900 dark:to-amber-950/10'}"
>
  <div>
    <!-- Card Header: Title, Driver, Status Badge -->
    <div class="flex items-start justify-between gap-3 mb-3">
      <div class="min-w-0 flex-1">
        <div class="flex items-center gap-2 flex-wrap">
          <div class="p-1.5 rounded-lg {volume.inUse ? 'bg-cyan-500/10 text-cyan-600 dark:text-cyan-400' : 'bg-amber-500/10 text-amber-600 dark:text-amber-400'} shrink-0">
            <Database class="w-4 h-4" />
          </div>
          <h3
            class="font-mono text-sm font-semibold text-slate-900 dark:text-slate-100 truncate tracking-tight"
            title={volume.name}
          >
            {volume.name}
          </h3>
          <button
            onclick={(e) => copyText(volume.name, 'name', e)}
            class="p-1 text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 transition-colors rounded cursor-pointer"
            title="Copiar nombre del volumen"
          >
            {#if copiedName}
              <Check class="w-3.5 h-3.5 text-emerald-500" />
            {:else}
              <Copy class="w-3.5 h-3.5" />
            {/if}
          </button>
          <span class="text-[11px] font-mono px-2 py-0.5 rounded-md bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 border border-slate-200 dark:border-slate-700">
            {volume.driver || 'local'}
          </span>
        </div>
      </div>

      <!-- In-Use vs Dangling Status Badge -->
      <div class="shrink-0">
        {#if volume.inUse}
          <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border border-emerald-500/20 shadow-2xs">
            <Link2 class="w-3 h-3" />
            <span>En uso</span>
          </span>
        {:else}
          <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium bg-amber-500/15 text-amber-800 dark:text-amber-300 border border-amber-500/30 shadow-2xs">
            <AlertTriangle class="w-3 h-3 text-amber-500 animate-pulse" />
            <span>Huérfano</span>
          </span>
        {/if}
      </div>
    </div>

    <!-- Metadata Section: Size & Host Mountpoint -->
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-2.5 py-2.5 border-y border-slate-100 dark:border-slate-800/80 text-xs">
      <!-- Estimated Disk Size -->
      <div class="flex items-center gap-2">
        <div class="p-1 rounded bg-slate-100 dark:bg-slate-800 text-slate-500 dark:text-slate-400">
          <HardDrive class="w-3.5 h-3.5" />
        </div>
        <div>
          <span class="text-slate-400 dark:text-slate-500 block text-[11px]">Tamaño en disco</span>
          <span class="font-semibold text-slate-800 dark:text-slate-200">
            {volume.size >= 0 ? formatBytes(volume.size) : 'No disponible'}
          </span>
        </div>
      </div>

      <!-- Scope / Created -->
      {#if volume.createdAt}
        <div class="flex items-center gap-2">
          <div class="p-1 rounded bg-slate-100 dark:bg-slate-800 text-slate-500 dark:text-slate-400">
            <Folder class="w-3.5 h-3.5" />
          </div>
          <div class="truncate">
            <span class="text-slate-400 dark:text-slate-500 block text-[11px]">Fecha de creación</span>
            <span class="text-slate-700 dark:text-slate-300 font-mono text-[11px] truncate block" title={volume.createdAt}>
              {volume.createdAt.length > 19 ? volume.createdAt.substring(0, 19).replace('T', ' ') : volume.createdAt}
            </span>
          </div>
        </div>
      {/if}
    </div>

    <!-- Host Mountpoint Row -->
    <div class="mt-2.5 bg-slate-50 dark:bg-slate-800/50 rounded-xl p-2.5 border border-slate-100 dark:border-slate-800 text-xs">
      <div class="flex items-center justify-between gap-2 mb-1">
        <span class="text-[11px] font-medium text-slate-500 dark:text-slate-400 flex items-center gap-1.5">
          <Folder class="w-3 h-3" />
          Ruta en el host:
        </span>
        <button
          onclick={(e) => copyText(volume.mountpoint, 'path', e)}
          class="p-0.5 text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 transition-colors rounded cursor-pointer"
          title="Copiar ruta de montaje en el host"
        >
          {#if copiedPath}
            <span class="text-[10px] text-emerald-500 font-sans flex items-center gap-0.5">
              <Check class="w-3 h-3" /> Copiado
            </span>
          {:else}
            <Copy class="w-3 h-3" />
          {/if}
        </button>
      </div>
      <div
        class="font-mono text-[11px] text-slate-700 dark:text-slate-300 truncate select-all"
        title={volume.mountpoint}
      >
        {volume.mountpoint}
      </div>
    </div>

    <!-- Associated Containers Section -->
    <div class="mt-3 text-xs">
      <span class="text-[11px] font-medium text-slate-500 dark:text-slate-400 block mb-1.5">
        Contenedores asociados:
      </span>

      {#if volume.containers && volume.containers.length > 0}
        <div class="space-y-1.5 max-h-36 overflow-y-auto pr-1">
          {#each volume.containers as c}
            <div class="flex items-center justify-between gap-2 px-2.5 py-1.5 rounded-lg bg-slate-100/80 dark:bg-slate-800/80 border border-slate-200/60 dark:border-slate-700/60">
              <div class="flex items-center gap-2 min-w-0">
                <span
                  class="w-2 h-2 rounded-full shrink-0 {c.state === 'running' ? 'bg-emerald-500' : 'bg-slate-400'}"
                  title="Estado: {c.state}"
                ></span>
                <Box class="w-3 h-3 text-slate-400 shrink-0" />
                <span class="font-medium text-slate-800 dark:text-slate-200 truncate" title={c.name}>
                  {c.name}
                </span>
                <span class="text-[10px] font-mono text-slate-400 dark:text-slate-500 shrink-0">
                  ({c.id})
                </span>
              </div>
              <div class="text-[11px] font-mono text-slate-500 dark:text-slate-400 truncate text-right max-w-[150px]" title="Punto de montaje en contenedor: {c.destination}">
                {c.destination}
              </div>
            </div>
          {/each}
        </div>
      {:else}
        <div class="p-2.5 rounded-xl bg-amber-500/10 border border-amber-500/20 text-amber-800 dark:text-amber-300 flex items-center gap-2">
          <Unlink class="w-3.5 h-3.5 shrink-0 text-amber-600 dark:text-amber-400" />
          <span class="text-[11px]">
            Sin contenedores vinculados. Es seguro eliminarlo para liberar espacio.
          </span>
        </div>
      {/if}
    </div>
  </div>

  <!-- Card Actions Footer -->
  <div class="mt-4 pt-3 border-t border-slate-100 dark:border-slate-800 flex items-center justify-between gap-2">
    <!-- Inspect Button -->
    <button
      onclick={() => onInspect(volume)}
      class="flex items-center gap-1.5 px-3 py-1.5 rounded-xl text-xs font-medium text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer border border-slate-200 dark:border-slate-700"
      title="Inspeccionar configuración JSON del volumen"
    >
      <Code2 class="w-3.5 h-3.5" />
      <span>Inspeccionar</span>
    </button>

    <!-- Delete Button -->
    <button
      onclick={() => onDelete(volume)}
      disabled={isDeleting}
      class="flex items-center gap-1.5 px-3 py-1.5 rounded-xl text-xs font-medium text-rose-600 dark:text-rose-400 hover:bg-rose-50 dark:hover:bg-rose-950/40 transition-colors cursor-pointer border border-rose-200 dark:border-rose-900/40 disabled:opacity-50"
      title={volume.inUse ? 'El volumen está en uso (requiere forzar)' : 'Eliminar volumen'}
    >
      <Trash2 class="w-3.5 h-3.5" />
      <span>{isDeleting ? 'Eliminando...' : 'Eliminar'}</span>
    </button>
  </div>
</div>
