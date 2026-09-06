<script lang="ts">
  import {
    X,
    Box,
    Play,
    Plus,
    Trash2,
    Sparkles,
    ChevronDown,
    ChevronRight,
    AlertCircle,
    Loader2,
    Check,
    HardDrive,
    ArrowRight,
  } from '@lucide/svelte';
  import { CreateContainer } from '../../wailsjs/go/main/App';
  import type { ImageInfo, CreateContainerRequest } from '../types';

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

  // Dynamic lists
  interface PortRow {
    host: string;
    container: string;
    protocol: 'tcp' | 'udp';
  }
  let ports = $state<PortRow[]>([]);

  interface VolumeRow {
    hostPath: string;
    containerPath: string;
    mode: 'rw' | 'ro';
  }
  let volumes = $state<VolumeRow[]>([]);

  interface EnvRow {
    key: string;
    value: string;
  }
  let envVars = $state<EnvRow[]>([]);

  // Execution & Error State
  let isSubmitting = $state(false);
  let errorMessage = $state<string | null>(null);

  // Popular image presets
  const popularImages = [
    { name: 'nginx:alpine', label: 'Nginx', defaultPort: { host: '8080', container: '80' } },
    { name: 'postgres:16-alpine', label: 'Postgres', defaultPort: { host: '5432', container: '5432' }, defaultEnv: { key: 'POSTGRES_PASSWORD', value: 'postgres' } },
    { name: 'redis:alpine', label: 'Redis', defaultPort: { host: '6379', container: '6379' } },
    { name: 'node:20-alpine', label: 'Node.js', defaultPort: { host: '3000', container: '3000' } },
    { name: 'python:3.11-slim', label: 'Python' },
  ];

  // Quick port presets
  const portPresets = [
    { label: '80:80 (HTTP)', host: '80', container: '80' },
    { label: '8080:80 (Web)', host: '8080', container: '80' },
    { label: '3000:3000 (Dev)', host: '3000', container: '3000' },
    { label: '5432:5432 (Postgres)', host: '5432', container: '5432' },
    { label: '6379:6379 (Redis)', host: '6379', container: '6379' },
  ];

  // When modal opens or initialImage changes
  $effect(() => {
    if (isOpen) {
      if (initialImage) {
        imageName = initialImage;
      }
      errorMessage = null;
    }
  });

  // Verify if current image exists locally
  const isImageLocal = $derived.by(() => {
    const trimmed = imageName.trim().toLowerCase();
    if (!trimmed) return false;
    return (localImages || []).some((img) => {
      const fullTag = `${img.repository}:${img.tag}`.toLowerCase();
      return fullTag === trimmed || img.repository.toLowerCase() === trimmed;
    });
  });

  // Port helpers
  function addPort(host = '', container = '', protocol: 'tcp' | 'udp' = 'tcp') {
    ports = [...ports, { host, container, protocol }];
  }

  function removePort(index: number) {
    ports = ports.filter((_, i) => i !== index);
  }

  function applyPortPreset(p: { host: string; container: string }) {
    // Avoid exact duplicate
    const exists = ports.some((item) => item.host === p.host && item.container === p.container);
    if (!exists) {
      addPort(p.host, p.container, 'tcp');
    }
  }

  // Volume helpers
  function addVolume(hostPath = '', containerPath = '', mode: 'rw' | 'ro' = 'rw') {
    volumes = [...volumes, { hostPath, containerPath, mode }];
  }

  function removeVolume(index: number) {
    volumes = volumes.filter((_, i) => i !== index);
  }

  // Environment variable helpers
  function addEnv(key = '', value = '') {
    envVars = [...envVars, { key, value }];
  }

  function removeEnv(index: number) {
    envVars = envVars.filter((_, i) => i !== index);
  }

  function applyPreset(preset: typeof popularImages[0]) {
    imageName = preset.name;
    if (preset.defaultPort && ports.length === 0) {
      addPort(preset.defaultPort.host, preset.defaultPort.container, 'tcp');
    }
    if (preset.defaultEnv && envVars.length === 0) {
      addEnv(preset.defaultEnv.key, preset.defaultEnv.value);
    }
  }

  function resetForm() {
    imageName = '';
    containerName = '';
    ports = [];
    volumes = [];
    envVars = [];
    restartPolicy = 'unless-stopped';
    autoStart = true;
    showAdvanced = false;
    errorMessage = null;
    isSubmitting = false;
  }

  function handleClose() {
    if (isSubmitting) return;
    resetForm();
    onClose();
  }

  async function handleSubmit(e: Event) {
    e.preventDefault();
    const targetImage = imageName.trim();
    if (!targetImage) {
      errorMessage = 'Por favor especifica una imagen de Docker.';
      return;
    }

    // Format ports
    const formattedPorts: string[] = [];
    for (const p of ports) {
      const h = p.host.trim();
      const c = p.container.trim();
      if (c) {
        if (h) {
          formattedPorts.push(`${h}:${c}/${p.protocol}`);
        } else {
          formattedPorts.push(`${c}/${p.protocol}`);
        }
      }
    }

    // Format volumes
    const formattedVolumes: string[] = [];
    for (const v of volumes) {
      const h = v.hostPath.trim();
      const c = v.containerPath.trim();
      if (h && c) {
        formattedVolumes.push(`${h}:${c}:${v.mode}`);
      }
    }

    // Format env vars
    const formattedEnv: string[] = [];
    for (const env of envVars) {
      const k = env.key.trim();
      if (k) {
        formattedEnv.push(`${k}=${env.value}`);
      }
    }

    const payload: CreateContainerRequest = {
      image: targetImage,
      name: containerName.trim() || undefined,
      ports: formattedPorts.length > 0 ? formattedPorts : undefined,
      volumes: formattedVolumes.length > 0 ? formattedVolumes : undefined,
      env: formattedEnv.length > 0 ? formattedEnv : undefined,
      restartPolicy: restartPolicy || 'no',
      autoStart,
    };

    isSubmitting = true;
    errorMessage = null;

    try {
      const res = await CreateContainer(payload);
      onSuccess(res.id);
      handleClose();
    } catch (err: any) {
      errorMessage = err?.toString() || 'Error al crear el contenedor';
    } finally {
      isSubmitting = false;
    }
  }
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 bg-slate-950/60 backdrop-blur-sm flex items-center justify-center p-4">
    <div
      class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl w-full max-w-2xl shadow-2xl overflow-hidden flex flex-col max-h-[90vh] animate-in fade-in zoom-in-95 duration-150"
    >
      <!-- Modal Header -->
      <div class="px-6 py-4 border-b border-slate-200 dark:border-slate-800 flex items-center justify-between bg-slate-50/50 dark:bg-slate-900/50">
        <div class="flex items-center gap-3">
          <div class="w-9 h-9 rounded-xl bg-blue-500/10 text-blue-600 dark:text-blue-400 flex items-center justify-center font-bold">
            <Box class="w-5 h-5" />
          </div>
          <div>
            <h3 class="font-semibold text-slate-900 dark:text-slate-100 text-sm">
              Crear Nuevo Contenedor
            </h3>
            <p class="text-xs text-slate-500 dark:text-slate-400">
              Configura imagen, puertos, volúmenes y variables de entorno para desplegar
            </p>
          </div>
        </div>

        {#if !isSubmitting}
          <button
            onclick={handleClose}
            aria-label="Cerrar modal"
            class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer"
          >
            <X class="w-4 h-4" />
          </button>
        {/if}
      </div>

      <!-- Modal Body (Scrollable Form) -->
      <form onsubmit={handleSubmit} class="p-6 space-y-5 overflow-y-auto flex-1">
        <!-- Error Alert -->
        {#if errorMessage}
          <div class="p-3.5 rounded-xl bg-rose-50 dark:bg-rose-950/30 border border-rose-200 dark:border-rose-900/50 flex items-start gap-2.5">
            <AlertCircle class="w-4 h-4 text-rose-600 dark:text-rose-400 flex-shrink-0 mt-0.5" />
            <div class="space-y-1">
              <h4 class="text-xs font-semibold text-rose-800 dark:text-rose-300">
                Fallo al crear contenedor
              </h4>
              <p class="text-xs text-rose-700 dark:text-rose-300 leading-relaxed font-mono break-all">
                {errorMessage}
              </p>
            </div>
          </div>
        {/if}

        <!-- Section 1: Image & Name -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- Image Input -->
          <div class="space-y-1.5">
            <div class="flex items-center justify-between">
              <label for="create-image-input" class="text-xs font-medium text-slate-700 dark:text-slate-300 flex items-center gap-1">
                <span>Imagen Docker</span>
                <span class="text-rose-500">*</span>
              </label>
              {#if imageName.trim()}
                {#if isImageLocal}
                  <span class="text-[10px] font-medium text-emerald-600 dark:text-emerald-400 flex items-center gap-0.5">
                    <Check class="w-3 h-3" /> Local
                  </span>
                {:else}
                  <span class="text-[10px] font-medium text-blue-600 dark:text-blue-400 flex items-center gap-0.5">
                    <HardDrive class="w-3 h-3" /> Se descargará
                  </span>
                {/if}
              {/if}
            </div>

            <input
              id="create-image-input"
              type="text"
              required
              disabled={isSubmitting}
              placeholder="ej. nginx:alpine, postgres:16"
              bind:value={imageName}
              class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/80 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all disabled:opacity-60 font-mono"
            />

            <!-- Suggestions chips -->
            <div class="flex items-center gap-1 pt-1 flex-wrap">
              <span class="text-[10px] text-slate-400 flex items-center gap-0.5">
                <Sparkles class="w-2.5 h-2.5 text-amber-500" /> Populares:
              </span>
              {#each popularImages as pop}
                <button
                  type="button"
                  disabled={isSubmitting}
                  onclick={() => applyPreset(pop)}
                  class="px-1.5 py-0.5 text-[10px] font-mono rounded-md border border-slate-200 dark:border-slate-800 bg-slate-100/60 dark:bg-slate-800/60 text-slate-600 dark:text-slate-300 hover:border-blue-400 hover:text-blue-600 transition-colors cursor-pointer"
                >
                  {pop.label}
                </button>
              {/each}
            </div>
          </div>

          <!-- Container Name -->
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
              class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/80 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all disabled:opacity-60 font-mono"
            />
            <p class="text-[10px] text-slate-400">
              Debe contener solo letras, números, guiones y guiones bajos.
            </p>
          </div>
        </div>

        <hr class="border-slate-200 dark:border-slate-800" />

        <!-- Section 2: Port Mappings -->
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
              disabled={isSubmitting}
              onclick={() => addPort()}
              class="px-2.5 py-1 text-xs rounded-lg font-medium text-blue-600 dark:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-950/40 border border-blue-200 dark:border-blue-900/40 flex items-center gap-1 transition-colors cursor-pointer"
            >
              <Plus class="w-3.5 h-3.5" />
              <span>Agregar Puerto</span>
            </button>
          </div>

          <!-- Port Presets -->
          {#if ports.length === 0}
            <div class="p-3 rounded-xl border border-dashed border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/30 text-center space-y-2">
              <p class="text-xs text-slate-500">
                Sin puertos configurados. Haz clic en "Agregar Puerto" o selecciona uno común:
              </p>
              <div class="flex items-center justify-center gap-1.5 flex-wrap">
                {#each portPresets as preset}
                  <button
                    type="button"
                    disabled={isSubmitting}
                    onclick={() => applyPortPreset(preset)}
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
                      disabled={isSubmitting}
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
                      disabled={isSubmitting}
                      bind:value={port.container}
                      class="w-full px-2.5 py-1.5 text-xs rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100 font-mono focus:outline-none focus:ring-1 focus:ring-blue-500"
                    />
                  </div>

                  <div class="w-24 space-y-1">
                    <span class="text-[10px] text-slate-500 dark:text-slate-400 font-medium">Protocolo</span>
                    <select
                      aria-label="Protocolo de puerto"
                      disabled={isSubmitting}
                      bind:value={port.protocol}
                      class="w-full px-2 py-1.5 text-xs rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100 font-mono focus:outline-none focus:ring-1 focus:ring-blue-500 cursor-pointer"
                    >
                      <option value="tcp">TCP</option>
                      <option value="udp">UDP</option>
                    </select>
                  </div>

                  <div class="pt-4">
                    <button
                      type="button"
                      disabled={isSubmitting}
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

        <hr class="border-slate-200 dark:border-slate-800" />

        <!-- Section 3: Volumes / Binds -->
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
              disabled={isSubmitting}
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
                      disabled={isSubmitting}
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
                      disabled={isSubmitting}
                      bind:value={vol.containerPath}
                      class="w-full px-2.5 py-1.5 text-xs rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100 font-mono focus:outline-none focus:ring-1 focus:ring-blue-500"
                    />
                  </div>

                  <div class="w-28 space-y-1">
                    <span class="text-[10px] text-slate-500 dark:text-slate-400 font-medium">Modo</span>
                    <select
                      aria-label="Modo de volumen"
                      disabled={isSubmitting}
                      bind:value={vol.mode}
                      class="w-full px-2 py-1.5 text-xs rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100 font-mono focus:outline-none focus:ring-1 focus:ring-blue-500 cursor-pointer"
                    >
                      <option value="rw">Lectura/Escritura (rw)</option>
                      <option value="ro">Solo Lectura (ro)</option>
                    </select>
                  </div>

                  <div class="pt-4">
                    <button
                      type="button"
                      disabled={isSubmitting}
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

        <hr class="border-slate-200 dark:border-slate-800" />

        <!-- Section 4: Environment Variables (ENV) -->
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
              disabled={isSubmitting}
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
                      disabled={isSubmitting}
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
                      disabled={isSubmitting}
                      bind:value={env.value}
                      class="w-full px-2.5 py-1.5 text-xs rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100 font-mono focus:outline-none focus:ring-1 focus:ring-blue-500"
                    />
                  </div>

                  <div class="pt-4">
                    <button
                      type="button"
                      disabled={isSubmitting}
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

        <!-- Section 5: Advanced (Collapsible) -->
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
              <!-- Restart Policy -->
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

              <!-- Auto Start Checkbox -->
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

      <!-- Modal Footer Actions -->
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
