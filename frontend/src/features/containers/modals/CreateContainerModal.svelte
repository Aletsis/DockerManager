<script lang="ts">
  import {
    X,
    Box,
    Play,
    ChevronDown,
    ChevronRight,
    AlertCircle,
    Loader2,
    Check,
    HardDrive,
  } from '@lucide/svelte';
  import { dockerApi } from '../../../shared/services/api';
  import type { ImageInfo, CreateContainerRequest } from '../../../types';
  import PopularImagePresets, { type ImagePreset } from '../forms/PopularImagePresets.svelte';
  import PortMappingsForm, { type PortRow } from '../forms/PortMappingsForm.svelte';
  import VolumeMappingsForm, { type VolumeRow } from '../forms/VolumeMappingsForm.svelte';
  import EnvVariablesForm, { type EnvRow } from '../forms/EnvVariablesForm.svelte';

  let {
    isOpen = false,
    initialImage = '',
    localImages = [],
    onClose,
    onSuccess,
  } = $props<{
    isOpen: boolean;
    initialImage?: string;
    localImages?: ImageInfo[];
    onClose: () => void;
    onSuccess: (containerId: string) => void;
  }>();

  // Form State
  let imageName = $state('');
  let containerName = $state('');
  let autoStart = $state(true);
  let restartPolicy = $state('unless-stopped');
  let showAdvanced = $state(false);

  // Dynamic rows
  let ports = $state<PortRow[]>([]);
  let volumes = $state<VolumeRow[]>([]);
  let envVars = $state<EnvRow[]>([]);

  // Execution State
  let isSubmitting = $state(false);
  let errorMessage = $state<string | null>(null);

  $effect(() => {
    if (isOpen) {
      if (initialImage) {
        imageName = initialImage;
      }
      errorMessage = null;
    }
  });

  const isImageLocal = $derived.by(() => {
    const trimmed = imageName.trim().toLowerCase();
    if (!trimmed) return false;
    return (localImages || []).some((img) => {
      const fullTag = `${img.repository}:${img.tag}`.toLowerCase();
      return fullTag === trimmed || img.repository.toLowerCase() === trimmed;
    });
  });

  function handleSelectPreset(preset: ImagePreset) {
    imageName = preset.name;
    if (preset.defaultPort) {
      const exists = ports.some(
        (p) => p.host === preset.defaultPort!.host && p.container === preset.defaultPort!.container
      );
      if (!exists) {
        ports = [...ports, { host: preset.defaultPort.host, container: preset.defaultPort.container, protocol: 'tcp' }];
      }
    }
    if (preset.defaultEnv) {
      const exists = envVars.some((e) => e.key === preset.defaultEnv!.key);
      if (!exists) {
        envVars = [...envVars, { key: preset.defaultEnv.key, value: preset.defaultEnv.value }];
      }
    }
  }

  function handleClose() {
    if (isSubmitting) return;
    onClose();
    setTimeout(() => {
      imageName = '';
      containerName = '';
      ports = [];
      volumes = [];
      envVars = [];
      autoStart = true;
      restartPolicy = 'unless-stopped';
      showAdvanced = false;
      errorMessage = null;
    }, 200);
  }

  async function handleSubmit(e: Event) {
    e.preventDefault();
    if (!imageName.trim()) {
      errorMessage = 'El nombre o tag de la imagen es obligatorio';
      return;
    }

    isSubmitting = true;
    errorMessage = null;

    try {
      const portStrings = ports
        .filter((p) => p.host.trim() && p.container.trim())
        .map((p) => `${p.host.trim()}:${p.container.trim()}/${p.protocol}`);

      const volumeStrings = volumes
        .filter((v) => v.hostPath.trim() && v.containerPath.trim())
        .map((v) => `${v.hostPath.trim()}:${v.containerPath.trim()}:${v.mode}`);

      const envStrings = envVars
        .filter((e) => e.key.trim())
        .map((e) => `${e.key.trim()}=${e.value.trim()}`);

      const request: CreateContainerRequest = {
        image: imageName.trim(),
        name: containerName.trim() || undefined,
        ports: portStrings.length > 0 ? portStrings : undefined,
        volumes: volumeStrings.length > 0 ? volumeStrings : undefined,
        env: envStrings.length > 0 ? envStrings : undefined,
        restartPolicy,
        autoStart,
      };

      const result = await dockerApi.createContainer(request);
      onSuccess(result.id);
      handleClose();
    } catch (err: any) {
      errorMessage = err?.toString() || 'Ocurrió un error inesperado al crear el contenedor';
    } finally {
      isSubmitting = false;
    }
  }
</script>

{#if isOpen}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/70 backdrop-blur-xs animate-in fade-in duration-200"
    role="dialog"
    aria-modal="true"
  >
    <div
      class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl w-full max-w-2xl max-h-[90vh] shadow-2xl flex flex-col overflow-hidden"
    >
      <!-- Header -->
      <div class="px-6 py-4 border-b border-slate-200 dark:border-slate-800 flex items-center justify-between bg-slate-50/50 dark:bg-slate-900/50">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 rounded-xl bg-blue-600/10 dark:bg-blue-400/10 text-blue-600 dark:text-blue-400 flex items-center justify-center">
            <Box class="w-4 h-4" />
          </div>
          <div>
            <h3 class="font-bold text-sm text-slate-900 dark:text-slate-100">
              Crear Nuevo Contenedor
            </h3>
            <p class="text-[11px] text-slate-500">
              Despliega un contenedor a partir de una imagen Docker
            </p>
          </div>
        </div>
        <button
          onclick={handleClose}
          disabled={isSubmitting}
          class="p-1.5 rounded-lg text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Scrollable Form Body -->
      <form onsubmit={handleSubmit} class="p-6 overflow-y-auto space-y-5 flex-1 text-slate-900 dark:text-slate-100">
        {#if errorMessage}
          <div class="p-3 rounded-xl bg-rose-50 dark:bg-rose-950/30 border border-rose-200 dark:border-rose-900/50 flex items-start gap-2.5 text-xs text-rose-600 dark:text-rose-400">
            <AlertCircle class="w-4 h-4 flex-shrink-0 mt-0.5" />
            <div class="leading-relaxed flex-1">{errorMessage}</div>
          </div>
        {/if}

        <PopularImagePresets disabled={isSubmitting} onSelect={handleSelectPreset} />

        <!-- Image & Name Inputs -->
        <div class="space-y-3 pt-1">
          <div class="space-y-1.5">
            <label for="create-image-input" class="text-xs font-semibold text-slate-800 dark:text-slate-200 flex items-center justify-between">
              <span>Imagen Docker *</span>
              {#if imageName.trim()}
                {#if isImageLocal}
                  <span class="text-[10px] text-emerald-600 dark:text-emerald-400 font-normal flex items-center gap-1">
                    <Check class="w-3 h-3" /> Disponible localmente
                  </span>
                {:else}
                  <span class="text-[10px] text-blue-600 dark:text-blue-400 font-normal flex items-center gap-1">
                    <HardDrive class="w-3 h-3" /> Se descargará automáticamente
                  </span>
                {/if}
              {/if}
            </label>
            <input
              id="create-image-input"
              type="text"
              required
              disabled={isSubmitting}
              placeholder="ej. redis:latest, postgres:16-alpine"
              bind:value={imageName}
              class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/80 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono"
            />
          </div>

          <div class="space-y-1.5">
            <label for="create-name-input" class="text-xs font-medium text-slate-700 dark:text-slate-300">
              Nombre del Contenedor <span class="text-slate-400 font-normal">(Opcional)</span>
            </label>
            <input
              id="create-name-input"
              type="text"
              disabled={isSubmitting}
              placeholder="ej. mi-servidor-web (automático si vacío)"
              bind:value={containerName}
              class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/80 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono"
            />
          </div>
        </div>

        <hr class="border-slate-200 dark:border-slate-800" />

        <!-- Subforms -->
        <PortMappingsForm bind:ports disabled={isSubmitting} />

        <hr class="border-slate-200 dark:border-slate-800" />

        <VolumeMappingsForm bind:volumes disabled={isSubmitting} />

        <hr class="border-slate-200 dark:border-slate-800" />

        <EnvVariablesForm bind:envVars disabled={isSubmitting} />

        <!-- Advanced Options -->
        <div class="pt-1">
          <button
            type="button"
            onclick={() => (showAdvanced = !showAdvanced)}
            class="flex items-center gap-1.5 text-xs font-semibold text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200 transition-colors cursor-pointer"
          >
            {#if showAdvanced}
              <ChevronDown class="w-3.5 h-3.5" />
            {:else}
              <ChevronRight class="w-3.5 h-3.5" />
            {/if}
            <span>Opciones Avanzadas</span>
          </button>

          {#if showAdvanced}
            <div class="mt-3 p-4 rounded-xl bg-slate-50 dark:bg-slate-800/40 border border-slate-200 dark:border-slate-800 space-y-3.5 animate-in fade-in duration-150">
              <div class="space-y-1.5">
                <label for="create-restart-policy" class="text-xs font-medium text-slate-700 dark:text-slate-300">
                  Política de Reinicio
                </label>
                <select
                  id="create-restart-policy"
                  disabled={isSubmitting}
                  bind:value={restartPolicy}
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500 cursor-pointer"
                >
                  <option value="no">No reiniciar automáticamente (no)</option>
                  <option value="unless-stopped">Reiniciar a menos que se detenga manualmente (unless-stopped)</option>
                  <option value="always">Reiniciar siempre (always)</option>
                  <option value="on-failure">Reiniciar solo si falla (on-failure)</option>
                </select>
              </div>

              <label class="flex items-center gap-2 cursor-pointer select-none">
                <input
                  type="checkbox"
                  disabled={isSubmitting}
                  bind:checked={autoStart}
                  class="rounded border-slate-300 dark:border-slate-700 text-blue-600 focus:ring-blue-500"
                />
                <span class="text-xs text-slate-700 dark:text-slate-300">
                  Iniciar contenedor automáticamente tras crearlo
                </span>
              </label>
            </div>
          {/if}
        </div>
      </form>

      <!-- Footer Actions -->
      <div class="px-6 py-4 border-t border-slate-200 dark:border-slate-800 bg-slate-50/80 dark:bg-slate-900/80 flex items-center justify-between">
        <div class="text-[11px] text-slate-500">
          {#if isSubmitting}
            <span class="flex items-center gap-1.5 text-blue-600 dark:text-blue-400">
              <Loader2 class="w-3.5 h-3.5 animate-spin" />
              <span>Preparando imagen y desplegando contenedor...</span>
            </span>
          {/if}
        </div>

        <div class="flex items-center gap-2">
          <button
            type="button"
            onclick={handleClose}
            disabled={isSubmitting}
            class="px-4 py-2 rounded-xl border border-slate-200 dark:border-slate-700 text-xs font-medium text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 disabled:opacity-50 transition-colors cursor-pointer"
          >
            Cancelar
          </button>

          <button
            type="button"
            onclick={handleSubmit}
            disabled={!imageName.trim() || isSubmitting}
            class="px-4 py-2 rounded-xl text-xs font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:pointer-events-none transition-colors flex items-center gap-1.5 shadow-sm shadow-blue-500/20 cursor-pointer"
          >
            {#if isSubmitting}
              <Loader2 class="w-3.5 h-3.5 animate-spin" />
              <span>Creando...</span>
            {:else}
              <Play class="w-3.5 h-3.5 fill-current" />
              <span>{autoStart ? 'Crear e Iniciar' : 'Crear Contenedor'}</span>
            {/if}
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}
