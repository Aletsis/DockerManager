<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { Download, X, AlertCircle, CheckCircle2, Loader2, Sparkles } from '@lucide/svelte';
  import { PullImage } from '../../../../wailsjs/go/main/App';
  import { EventsOn, EventsOff } from '../../../../wailsjs/runtime/runtime';
  import type { PullProgressEvent } from '../../../types';
  import { formatBytes } from '../../../shared/utils/utils';

  let {
    isOpen = false,
    onClose,
    onSuccess,
  } = $props<{
    isOpen: boolean;
    onClose: () => void;
    onSuccess: () => void;
  }>();

  let imageName = $state('');
  let isPulling = $state(false);
  let error = $state<string | null>(null);
  let isDone = $state(false);

  // Layers progress tracking
  let layers = $state<Record<string, {
    status: string;
    current: number;
    total: number;
    progress: string;
  }>>({});
  let generalStatus = $state<string>('');

  const quickSuggestions = [
    'redis:alpine',
    'nginx:alpine',
    'postgres:16-alpine',
    'node:20-alpine',
    'alpine:latest',
    'python:3.11-slim',
  ];

  onMount(() => {
    EventsOn('image:pull:progress', handleProgress);
  });

  onDestroy(() => {
    EventsOff('image:pull:progress');
  });

  function handleProgress(event: PullProgressEvent) {
    if (event.error) {
      error = event.error;
      isPulling = false;
      return;
    }

    if (event.status) {
      generalStatus = event.status;
    }

    if (event.id) {
      const currentLayer = layers[event.id] || { status: '', current: 0, total: 0, progress: '' };
      layers = {
        ...layers,
        [event.id]: {
          status: event.status || currentLayer.status,
          current: event.current || currentLayer.current,
          total: event.total || currentLayer.total,
          progress: event.progress || currentLayer.progress,
        },
      };
    }
  }

  async function handlePull(nameToPull?: string) {
    const target = (nameToPull || imageName).trim();
    if (!target || isPulling) return;

    imageName = target;
    isPulling = true;
    error = null;
    isDone = false;
    layers = {};
    generalStatus = 'Iniciando descarga...';

    try {
      await PullImage(target);
      isDone = true;
      generalStatus = '¡Descarga completada exitosamente!';
      setTimeout(() => {
        onSuccess();
        handleClose();
      }, 1500);
    } catch (err: any) {
      error = err?.toString() || 'Error al descargar la imagen';
    } finally {
      isPulling = false;
    }
  }

  function handleClose() {
    if (isPulling) return;
    imageName = '';
    error = null;
    isDone = false;
    layers = {};
    generalStatus = '';
    onClose();
  }

  const layerList = $derived(Object.entries(layers));
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 bg-slate-950/60 backdrop-blur-sm flex items-center justify-center p-4">
    <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl w-full max-w-lg shadow-2xl overflow-hidden flex flex-col max-h-[85vh] animate-in fade-in zoom-in-95 duration-150">
      
      <!-- Modal Header -->
      <div class="px-5 py-4 border-b border-slate-200 dark:border-slate-800 flex items-center justify-between">
        <div class="flex items-center gap-2.5">
          <div class="p-2 rounded-xl bg-blue-500/10 text-blue-600 dark:text-blue-400">
            <Download class="w-5 h-5" />
          </div>
          <div>
            <h3 class="font-semibold text-slate-900 dark:text-slate-100 text-sm">
              Descargar Imagen de Docker
            </h3>
            <p class="text-xs text-slate-500 dark:text-slate-400">
              Descarga imágenes públicas o privadas desde Docker Hub u otro registro
            </p>
          </div>
        </div>

        {#if !isPulling}
          <button
            onclick={handleClose}
            class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
          >
            <X class="w-4 h-4" />
          </button>
        {/if}
      </div>

      <!-- Modal Body -->
      <div class="p-5 space-y-4 overflow-y-auto flex-1">
        <!-- Input & Pull Form -->
        <div class="space-y-1.5">
          <label for="image-name-input" class="text-xs font-medium text-slate-700 dark:text-slate-300">
            Nombre de la Imagen y Tag
          </label>
          <div class="flex items-center gap-2">
            <input
              id="image-name-input"
              type="text"
              disabled={isPulling || isDone}
              placeholder="ej. redis:alpine, postgres:16, nginx:latest"
              bind:value={imageName}
              onkeydown={(e) => {
                if (e.key === 'Enter') handlePull();
              }}
              class="flex-1 px-3 py-2 text-xs rounded-xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/80 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all disabled:opacity-60 font-mono"
            />
            <button
              onclick={() => handlePull()}
              disabled={!imageName.trim() || isPulling || isDone}
              class="px-4 py-2 rounded-xl text-xs font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:pointer-events-none transition-colors flex items-center gap-1.5 shadow-sm shadow-blue-500/20"
            >
              {#if isPulling}
                <Loader2 class="w-3.5 h-3.5 animate-spin" />
                <span>Descargando...</span>
              {:else}
                <Download class="w-3.5 h-3.5" />
                <span>Descargar</span>
              {/if}
            </button>
          </div>
        </div>

        <!-- Quick Suggestions -->
        {#if !isPulling && !isDone}
          <div class="space-y-1.5">
            <div class="flex items-center gap-1 text-[11px] text-slate-500 dark:text-slate-400">
              <Sparkles class="w-3 h-3 text-amber-500" />
              <span>Sugerencias populares:</span>
            </div>
            <div class="flex flex-wrap gap-1.5">
              {#each quickSuggestions as sug}
                <button
                  type="button"
                  onclick={() => {
                    imageName = sug;
                  }}
                  class="px-2 py-1 text-[11px] font-mono rounded-lg border border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-800 text-slate-600 dark:text-slate-300 hover:border-blue-400 dark:hover:border-blue-600 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
                >
                  {sug}
                </button>
              {/each}
            </div>
          </div>
        {/if}

        <!-- Status or Error Message -->
        {#if error}
          <div class="p-3 rounded-xl bg-rose-50 dark:bg-rose-950/30 border border-rose-200 dark:border-rose-900/50 flex items-start gap-2.5">
            <AlertCircle class="w-4 h-4 text-rose-600 dark:text-rose-400 flex-shrink-0 mt-0.5" />
            <p class="text-xs text-rose-700 dark:text-rose-300 leading-relaxed font-mono break-all">
              {error}
            </p>
          </div>
        {/if}

        <!-- Progress Overview -->
        {#if isPulling || isDone || layerList.length > 0}
          <div class="space-y-3 pt-2">
            <div class="flex items-center justify-between text-xs">
              <span class="font-medium text-slate-700 dark:text-slate-300 flex items-center gap-1.5">
                {#if isDone}
                  <CheckCircle2 class="w-4 h-4 text-emerald-500" />
                  <span class="text-emerald-600 dark:text-emerald-400 font-semibold">Completado</span>
                {:else if isPulling}
                  <Loader2 class="w-3.5 h-3.5 animate-spin text-blue-500" />
                  <span>{generalStatus || 'Descargando capas...'}</span>
                {/if}
              </span>
              {#if layerList.length > 0}
                <span class="text-[11px] font-mono text-slate-400">
                  {layerList.length} {layerList.length === 1 ? 'capa' : 'capas'}
                </span>
              {/if}
            </div>

            <!-- Layer List -->
            {#if layerList.length > 0}
              <div class="space-y-2 max-h-48 overflow-y-auto p-2.5 rounded-xl bg-slate-50 dark:bg-slate-950/60 border border-slate-200 dark:border-slate-800 text-xs font-mono">
                {#each layerList as [id, layer]}
                  <div class="space-y-1">
                    <div class="flex items-center justify-between text-[11px]">
                      <span class="text-slate-600 dark:text-slate-400 font-semibold">{id}</span>
                      <span class="text-slate-500 dark:text-slate-400">
                        {layer.status}
                        {#if layer.total > 0}
                          ({formatBytes(layer.current)} / {formatBytes(layer.total)})
                        {/if}
                      </span>
                    </div>
                    {#if layer.total > 0}
                      <div class="w-full bg-slate-200 dark:bg-slate-800 rounded-full h-1.5 overflow-hidden">
                        <div
                          class="bg-blue-500 h-1.5 rounded-full transition-all duration-300"
                          style="width: {Math.min(100, Math.round((layer.current / layer.total) * 100))}%"
                        ></div>
                      </div>
                    {/if}
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {/if}

      </div>

      <!-- Modal Footer -->
      <div class="px-5 py-3 border-t border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50 flex items-center justify-end">
        <button
          onclick={handleClose}
          disabled={isPulling}
          class="px-3.5 py-1.5 rounded-lg border border-slate-200 dark:border-slate-700 text-xs font-medium text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 disabled:opacity-50 transition-colors"
        >
          {isDone ? 'Cerrar' : 'Cancelar'}
        </button>
      </div>

    </div>
  </div>
{/if}
