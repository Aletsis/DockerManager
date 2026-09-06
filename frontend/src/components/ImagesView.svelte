<script lang="ts">
  import {
    Package,
    Download,
    Trash2,
    Sparkles,
    HardDrive,
    Search,
    Copy,
    Check,
    AlertTriangle,
    Layers,
    Clock,
    RefreshCw,
    ShieldCheck,
    CheckCircle2,
    Info,
    Play,
  } from '@lucide/svelte';
  import type { ImageInfo, DiskUsageSummary } from '../types';
  import { formatBytes, formatRelativeTime } from '../utils';
  import ConfirmModal from './ConfirmModal.svelte';
  import PullImageModal from './PullImageModal.svelte';

  let {
    images = [],
    diskUsage = null,
    loading = false,
    onRefresh,
    onRemoveImage,
    onPrune,
    onDeployContainer,
    actionLoading = '',
  } = $props<{
    images: ImageInfo[];
    diskUsage: DiskUsageSummary | null;
    loading: boolean;
    onRefresh: () => void;
    onRemoveImage: (id: string, force: boolean) => Promise<void>;
    onPrune: (danglingOnly: boolean) => Promise<{ imagesDeleted: string[]; spaceReclaimed: number }>;
    onDeployContainer?: (imageTag: string) => void;
    actionLoading?: string;
  }>();

  // Search & Filter State
  let searchQuery = $state('');
  let filterState = $state<'all' | 'inUse' | 'unused' | 'dangling'>('all');

  // Modals state
  let isPullModalOpen = $state(false);
  let imageToDelete = $state<ImageInfo | null>(null);
  let isPruneModalOpen = $state(false);
  let pruneOnlyDangling = $state(true);
  let isPruning = $state(false);
  let pruneSuccessInfo = $state<{ count: number; bytes: number } | null>(null);

  // Copy feedback state
  let copiedId = $state<string | null>(null);

  function copyToClipboard(text: string, id: string, e: MouseEvent) {
    e.stopPropagation();
    navigator.clipboard.writeText(text);
    copiedId = id;
    setTimeout(() => {
      if (copiedId === id) copiedId = null;
    }, 1500);
  }

  // Filter images
  const filteredImages = $derived(
    images.filter((img) => {
      if (filterState === 'inUse' && !img.inUse) return false;
      if (filterState === 'unused' && img.inUse) return false;
      if (filterState === 'dangling' && !img.isDangling) return false;

      if (searchQuery.trim()) {
        const q = searchQuery.toLowerCase();
        const matchRepo = img.repository.toLowerCase().includes(q);
        const matchTag = img.tag.toLowerCase().includes(q);
        const matchId = img.shortId.toLowerCase().includes(q) || img.id.toLowerCase().includes(q);
        return matchRepo || matchTag || matchId;
      }
      return true;
    })
  );

  const danglingImages = $derived(images.filter((img) => img.isDangling));
  const inUseCount = $derived(images.filter((img) => img.inUse).length);
  const unusedCount = $derived(images.filter((img) => !img.inUse).length);

  async function handleConfirmDelete() {
    if (!imageToDelete) return;
    const target = imageToDelete;
    imageToDelete = null;
    await onRemoveImage(target.id, false);
  }

  async function handleConfirmPrune() {
    isPruning = true;
    try {
      const result = await onPrune(pruneOnlyDangling);
      isPruneModalOpen = false;
      pruneSuccessInfo = {
        count: result.imagesDeleted?.length || 0,
        bytes: result.spaceReclaimed || 0,
      };
      setTimeout(() => {
        pruneSuccessInfo = null;
      }, 5000);
    } finally {
      isPruning = false;
    }
  }
</script>

<div class="max-w-7xl w-full mx-auto space-y-5 pb-16">

  <!-- Feedback Banner for Prune Success -->
  {#if pruneSuccessInfo}
    <div class="p-4 rounded-2xl bg-emerald-50 dark:bg-emerald-950/30 border border-emerald-200 dark:border-emerald-800 flex items-center justify-between gap-3 animate-in fade-in slide-in-from-top-2 duration-200 shadow-sm">
      <div class="flex items-center gap-3">
        <div class="p-2 rounded-xl bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
          <CheckCircle2 class="w-5 h-5" />
        </div>
        <div>
          <h4 class="text-xs font-semibold text-emerald-900 dark:text-emerald-200">
            ¡Limpieza de disco completada con éxito!
          </h4>
          <p class="text-xs text-emerald-700 dark:text-emerald-300">
            Se liberaron <strong class="font-bold">{formatBytes(pruneSuccessInfo.bytes)}</strong> de espacio en disco y se eliminaron {pruneSuccessInfo.count} capas huérfanas.
          </p>
        </div>
      </div>
      <button
        onclick={() => (pruneSuccessInfo = null)}
        class="text-xs font-medium px-3 py-1.5 rounded-lg border border-emerald-300 dark:border-emerald-700 text-emerald-800 dark:text-emerald-200 hover:bg-emerald-100 dark:hover:bg-emerald-900/40 transition-colors"
      >
        Entendido
      </button>
    </div>
  {/if}

  <!-- Top Hero / Disk Recovery Summary Cards -->
  <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
    <!-- Card 1: Total Storage Used -->
    <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl p-4 shadow-xs relative overflow-hidden">
      <div class="flex items-center justify-between">
        <span class="text-xs font-medium text-slate-500 dark:text-slate-400">Almacenamiento Total</span>
        <div class="p-2 rounded-xl bg-blue-500/10 text-blue-600 dark:text-blue-400">
          <HardDrive class="w-4 h-4" />
        </div>
      </div>
      <div class="mt-2 flex items-baseline gap-2">
        <span class="text-2xl font-bold tracking-tight text-slate-900 dark:text-slate-100">
          {formatBytes(diskUsage?.totalSize || 0)}
        </span>
        <span class="text-xs text-slate-500 dark:text-slate-400">
          en {images.length} imágenes
        </span>
      </div>
      <div class="mt-3 text-[11px] text-slate-500 dark:text-slate-400 flex items-center gap-1.5">
        <Layers class="w-3.5 h-3.5 text-slate-400" />
        <span>{inUseCount} en uso activo • {unusedCount} inactivas</span>
      </div>
    </div>

    <!-- Card 2: Dangling / Reclaimable Storage -->
    <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl p-4 shadow-xs relative overflow-hidden">
      <div class="flex items-center justify-between">
        <span class="text-xs font-medium text-slate-500 dark:text-slate-400">Imágenes Huérfanas (Dangling)</span>
        <div class="p-2 rounded-xl bg-amber-500/10 text-amber-600 dark:text-amber-400">
          <AlertTriangle class="w-4 h-4" />
        </div>
      </div>
      <div class="mt-2 flex items-baseline gap-2">
        <span class="text-2xl font-bold tracking-tight text-amber-600 dark:text-amber-400">
          {danglingImages.length}
        </span>
        <span class="text-xs text-slate-500 dark:text-slate-400">
          capas sin etiqueta ({formatBytes(diskUsage?.danglingSize || 0)})
        </span>
      </div>
      <div class="mt-3 text-[11px] text-slate-500 dark:text-slate-400 flex items-center gap-1.5">
        <ShieldCheck class="w-3.5 h-3.5 text-emerald-500" />
        <span>100% seguras de eliminar sin afectar contenedores</span>
      </div>
    </div>

    <!-- Card 3: Quick Prune Action Hero -->
    <div class="bg-gradient-to-br from-indigo-500/10 via-blue-500/5 to-purple-500/10 dark:from-indigo-950/40 dark:via-blue-950/20 dark:to-purple-950/30 border border-indigo-200/80 dark:border-indigo-800/60 rounded-2xl p-4 shadow-xs flex flex-col justify-between">
      <div>
        <div class="flex items-center justify-between">
          <span class="text-xs font-semibold text-indigo-900 dark:text-indigo-300 flex items-center gap-1.5">
            <Sparkles class="w-4 h-4 text-indigo-500" />
            Liberación de Disco
          </span>
          {#if (diskUsage?.danglingSize || 0) > 0}
            <span class="text-[10px] uppercase tracking-wider font-bold px-2 py-0.5 rounded-full bg-emerald-500/15 text-emerald-700 dark:text-emerald-300 border border-emerald-500/30">
              {formatBytes(diskUsage?.danglingSize || 0)} recuperables
            </span>
          {/if}
        </div>
        <p class="text-xs text-slate-600 dark:text-slate-400 mt-1 leading-relaxed">
          Ejecuta una poda para liberar espacio ocupado por capas intermedias obsoletas.
        </p>
      </div>

      <div class="mt-3">
        <button
          onclick={() => (isPruneModalOpen = true)}
          class="w-full py-2 px-3.5 rounded-xl bg-gradient-to-r from-indigo-600 to-blue-600 hover:from-indigo-700 hover:to-blue-700 text-white font-medium text-xs shadow-sm shadow-indigo-500/25 transition-all flex items-center justify-center gap-2 active:scale-[0.99]"
        >
          <Trash2 class="w-4 h-4" />
          <span>Limpieza Rápida (Prune)</span>
        </button>
      </div>
    </div>
  </div>

  <!-- Control & Filter Bar -->
  <div class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
    <!-- Filter Chips -->
    <div class="flex items-center gap-1.5 p-1 bg-slate-200/60 dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 text-xs overflow-x-auto">
      <button
        onclick={() => (filterState = 'all')}
        class="px-3 py-1 rounded-lg font-medium transition-colors whitespace-nowrap {filterState === 'all' ? 'bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200'}"
      >
        Todas ({images.length})
      </button>
      <button
        onclick={() => (filterState = 'inUse')}
        class="px-3 py-1 rounded-lg font-medium transition-colors whitespace-nowrap {filterState === 'inUse' ? 'bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200'}"
      >
        En Uso ({inUseCount})
      </button>
      <button
        onclick={() => (filterState = 'unused')}
        class="px-3 py-1 rounded-lg font-medium transition-colors whitespace-nowrap {filterState === 'unused' ? 'bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200'}"
      >
        Sin Usar ({unusedCount})
      </button>
      <button
        onclick={() => (filterState = 'dangling')}
        class="px-3 py-1 rounded-lg font-medium transition-colors whitespace-nowrap {filterState === 'dangling' ? 'bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100 shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200'}"
      >
        Huérfanas ({danglingImages.length})
      </button>
    </div>

    <!-- Right Actions: Search + Download Button -->
    <div class="flex items-center gap-2">
      <!-- Image Search -->
      <div class="relative flex-1 sm:w-64">
        <Search class="w-3.5 h-3.5 absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" />
        <input
          type="text"
          placeholder="Buscar imagen o tag..."
          bind:value={searchQuery}
          class="w-full pl-8 pr-3 py-1.5 text-xs rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800/80 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-colors"
        />
      </div>

      <!-- Pull Button -->
      <button
        onclick={() => (isPullModalOpen = true)}
        class="px-3.5 py-1.5 rounded-xl bg-blue-600 hover:bg-blue-700 text-white font-medium text-xs transition-colors flex items-center gap-1.5 shadow-sm shadow-blue-500/20 whitespace-nowrap cursor-pointer"
      >
        <Download class="w-3.5 h-3.5" />
        <span>Descargar Imagen</span>
      </button>
    </div>
  </div>

  <!-- Loading State Skeleton -->
  {#if loading}
    <div class="space-y-3">
      {#each [1, 2, 3, 4] as n}
        <div class="h-20 rounded-xl border border-slate-200 dark:border-slate-800 bg-white/60 dark:bg-slate-900/60 animate-pulse"></div>
      {/each}
    </div>
  {/if}

  <!-- Empty State -->
  {#if !loading && filteredImages.length === 0}
    <div class="p-12 text-center rounded-2xl border border-dashed border-slate-300 dark:border-slate-800 bg-white/40 dark:bg-slate-900/40 space-y-3">
      <div class="w-12 h-12 rounded-2xl bg-slate-100 dark:bg-slate-800 flex items-center justify-center mx-auto text-slate-400">
        <Package class="w-6 h-6" />
      </div>
      <h3 class="font-semibold text-sm text-slate-800 dark:text-slate-200">
        No se encontraron imágenes
      </h3>
      <p class="text-xs text-slate-500 max-w-sm mx-auto">
        {searchQuery
          ? 'Ninguna imagen coincide con tu búsqueda actual.'
          : 'No hay imágenes descargadas en tu máquina actualmente.'}
      </p>
      {#if !searchQuery}
        <button
          onclick={() => (isPullModalOpen = true)}
          class="mt-2 inline-flex items-center gap-1.5 px-3.5 py-1.5 text-xs font-medium rounded-xl text-white bg-blue-600 hover:bg-blue-700 transition-colors"
        >
          <Download class="w-3.5 h-3.5" />
          <span>Descargar tu primera imagen</span>
        </button>
      {/if}
    </div>
  {/if}

  <!-- Images List Cards -->
  {#if !loading && filteredImages.length > 0}
    <div class="space-y-2.5">
      {#each filteredImages as img (img.id)}
        <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-3.5 shadow-xs hover:shadow transition-all duration-200 hover:border-slate-300 dark:hover:border-slate-700 flex flex-col md:flex-row md:items-center md:justify-between gap-3">
          
          <!-- Image Title & Tags -->
          <div class="flex items-start gap-3 min-w-0 flex-1">
            <!-- Icon Avatar -->
            <div class="w-9 h-9 rounded-xl flex items-center justify-center flex-shrink-0 mt-0.5 {img.isDangling ? 'bg-amber-500/10 text-amber-600 dark:text-amber-400' : 'bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300'}">
              <Package class="w-4 h-4" />
            </div>

            <div class="min-w-0 space-y-1">
              <div class="flex items-center gap-2 flex-wrap">
                {#if img.isDangling}
                  <span class="font-mono text-xs font-bold text-amber-600 dark:text-amber-400">
                    &lt;none&gt;:&lt;none&gt;
                  </span>
                  <span class="text-[10px] font-semibold px-2 py-0.5 rounded-full bg-amber-500/15 text-amber-700 dark:text-amber-300 border border-amber-500/30">
                    Huérfana
                  </span>
                {:else}
                  <span class="font-semibold text-slate-900 dark:text-slate-100 text-xs truncate max-w-xs md:max-w-md">
                    {img.repository}
                  </span>
                  <span class="text-[11px] font-mono px-2 py-0.5 rounded-md bg-blue-500/10 text-blue-700 dark:text-blue-300 border border-blue-500/20 font-medium">
                    {img.tag}
                  </span>
                {/if}

                <!-- Status inUse Badge -->
                {#if img.inUse}
                  <span class="text-[10px] font-medium px-2 py-0.5 rounded-full bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border border-emerald-500/20 flex items-center gap-1">
                    <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
                    <span>En uso ({img.containers > 0 ? `${img.containers} cont.` : 'asignada'})</span>
                  </span>
                {:else}
                  <span class="text-[10px] font-medium px-2 py-0.5 rounded-full bg-slate-100 dark:bg-slate-800 text-slate-500 dark:text-slate-400 border border-slate-200 dark:border-slate-700">
                    Sin usar
                  </span>
                {/if}
              </div>

              <!-- Metadata row: ID, Size, Created date -->
              <div class="flex items-center gap-3 text-[11px] text-slate-500 dark:text-slate-400 flex-wrap">
                <!-- Short ID with Copy -->
                <button
                  onclick={(e) => copyToClipboard(img.id, img.id, e)}
                  title="Copiar ID completo de la imagen"
                  class="flex items-center gap-1 font-mono text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 px-1.5 py-0.5 rounded bg-slate-50 dark:bg-slate-800/80 border border-slate-200 dark:border-slate-800 transition-colors"
                >
                  <span>{img.shortId}</span>
                  {#if copiedId === img.id}
                    <Check class="w-3 h-3 text-emerald-500" />
                  {:else}
                    <Copy class="w-3 h-3" />
                  {/if}
                </button>

                <span class="text-slate-300 dark:text-slate-700">•</span>

                <!-- Real Disk Size -->
                <span class="font-medium text-slate-700 dark:text-slate-300 flex items-center gap-1">
                  <HardDrive class="w-3 h-3 text-slate-400" />
                  {formatBytes(img.size)}
                </span>

                <span class="text-slate-300 dark:text-slate-700">•</span>

                <!-- Created Date -->
                <span class="flex items-center gap-1">
                  <Clock class="w-3 h-3 text-slate-400" />
                  {formatRelativeTime(img.created)}
                </span>
              </div>
            </div>
          </div>

          <!-- Actions -->
          <div class="flex items-center gap-2 self-end md:self-center">
            <!-- Deploy Container with this Image -->
            {#if onDeployContainer && !img.isDangling}
              {@const imgRef = img.tag && img.tag !== '<none>' ? `${img.repository}:${img.tag}` : img.repository}
              <button
                onclick={() => onDeployContainer(imgRef)}
                title="Crear y levantar un contenedor con esta imagen"
                class="px-3 py-1.5 rounded-xl bg-blue-600 hover:bg-blue-700 text-white font-medium text-xs transition-all flex items-center gap-1.5 shadow-sm shadow-blue-500/20 cursor-pointer whitespace-nowrap active:scale-[0.98]"
              >
                <Play class="w-3.5 h-3.5 fill-current" />
                <span>Desplegar</span>
              </button>
            {/if}

            <!-- Copy pull command -->
            {#if !img.isDangling}
              <button
                onclick={(e) => copyToClipboard(`docker pull ${img.repository}:${img.tag}`, `cmd-${img.id}`, e)}
                title="Copiar comando 'docker pull'"
                class="p-2 rounded-lg border border-slate-200 dark:border-slate-800 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-500 dark:text-slate-400 hover:text-slate-800 dark:hover:text-slate-200 transition-colors text-xs flex items-center gap-1 cursor-pointer"
              >
                {#if copiedId === `cmd-${img.id}`}
                  <Check class="w-3.5 h-3.5 text-emerald-500" />
                  <span class="text-[11px] text-emerald-600 dark:text-emerald-400">Copiado</span>
                {:else}
                  <Copy class="w-3.5 h-3.5" />
                  <span class="text-[11px] hidden sm:inline">Comando</span>
                {/if}
              </button>
            {/if}

            <!-- Delete Image Button -->
            <button
              onclick={() => (imageToDelete = img)}
              disabled={actionLoading === img.id}
              title={img.inUse ? 'Esta imagen está en uso por uno o más contenedores' : 'Eliminar esta imagen del disco'}
              class="p-2 rounded-lg border border-slate-200 dark:border-slate-800 hover:bg-rose-50 dark:hover:bg-rose-950/30 text-slate-400 hover:text-rose-600 dark:hover:text-rose-400 hover:border-rose-200 dark:hover:border-rose-900 transition-colors disabled:opacity-50"
            >
              <Trash2 class="w-3.5 h-3.5" />
            </button>
          </div>

        </div>
      {/each}
    </div>
  {/if}

</div>

<!-- Pull Image Modal -->
<PullImageModal
  isOpen={isPullModalOpen}
  onClose={() => (isPullModalOpen = false)}
  onSuccess={() => {
    onRefresh();
  }}
/>

<!-- Delete Image Confirmation Modal -->
<ConfirmModal
  isOpen={!!imageToDelete}
  title={imageToDelete?.inUse ? '¿Eliminar imagen en uso?' : 'Eliminar Imagen'}
  message={imageToDelete?.inUse
    ? `Atención: La imagen "${imageToDelete?.repository}:${imageToDelete?.tag}" está actualmente en uso por contenedores. Eliminarla forzadamente puede afectar los contenedores asociados. ¿Deseas continuar?`
    : `¿Estás seguro de que deseas eliminar permanentemente la imagen "${imageToDelete?.isDangling ? imageToDelete?.shortId : `${imageToDelete?.repository}:${imageToDelete?.tag}`}" (${formatBytes(imageToDelete?.size || 0)})?`}
  confirmLabel={imageToDelete?.inUse ? 'Forzar Eliminación' : 'Eliminar Imagen'}
  onConfirm={handleConfirmDelete}
  onCancel={() => (imageToDelete = null)}
  isDestructive={true}
/>

<!-- Prune / Clean Confirmation Modal -->
{#if isPruneModalOpen}
  <div class="fixed inset-0 z-50 bg-slate-950/60 backdrop-blur-sm flex items-center justify-center p-4">
    <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl w-full max-w-md shadow-2xl p-5 space-y-4 animate-in fade-in zoom-in-95 duration-150">
      <div class="flex items-start justify-between gap-3">
        <div class="flex items-center gap-3">
          <div class="p-2 rounded-xl bg-indigo-500/10 text-indigo-600 dark:text-indigo-400">
            <Sparkles class="w-5 h-5" />
          </div>
          <div>
            <h3 class="font-semibold text-slate-900 dark:text-slate-100 text-sm">
              Limpieza Rápida de Disco (Prune)
            </h3>
            <p class="text-xs text-slate-500">
              Recupera espacio eliminando capas no utilizadas
            </p>
          </div>
        </div>
      </div>

      <div class="p-3 rounded-xl bg-slate-50 dark:bg-slate-800/60 border border-slate-200 dark:border-slate-700/60 space-y-2">
        <div class="flex items-center gap-2 text-xs font-medium text-slate-800 dark:text-slate-200">
          <ShieldCheck class="w-4 h-4 text-emerald-500 flex-shrink-0" />
          <span>Acción segura garantizada</span>
        </div>
        <p class="text-[11px] text-slate-500 dark:text-slate-400 leading-relaxed">
          Esta operación ejecutará <code class="px-1 py-0.5 rounded bg-slate-200 dark:bg-slate-700 text-slate-800 dark:text-slate-200 font-mono">docker image prune</code> para eliminar capas huérfanas (<span class="font-mono">&lt;none&gt;:&lt;none&gt;</span>). <strong>Tus contenedores activos y detenidos permanecerán completamente a salvo.</strong>
        </p>
      </div>

      <div class="space-y-2 pt-1">
        <label class="flex items-start gap-2 cursor-pointer select-none">
          <input
            type="checkbox"
            bind:checked={pruneOnlyDangling}
            class="mt-0.5 rounded border-slate-300 text-indigo-600 focus:ring-indigo-500"
          />
          <span class="text-xs text-slate-600 dark:text-slate-300">
            <strong>Modo Seguro:</strong> Eliminar únicamente imágenes huérfanas (dangling).
            <span class="block text-[11px] text-slate-400">Desmarca esto si deseas eliminar todas las imágenes que no estén asociadas a ningún contenedor actual.</span>
          </span>
        </label>
      </div>

      <div class="flex items-center justify-end gap-2 pt-2">
        <button
          onclick={() => (isPruneModalOpen = false)}
          disabled={isPruning}
          class="px-3 py-1.5 rounded-lg border border-slate-200 dark:border-slate-700 text-xs font-medium text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
        >
          Cancelar
        </button>
        <button
          onclick={handleConfirmPrune}
          disabled={isPruning}
          class="px-3.5 py-1.5 rounded-lg text-xs font-medium text-white bg-indigo-600 hover:bg-indigo-700 transition-colors flex items-center gap-1.5 shadow-sm shadow-indigo-500/20"
        >
          {#if isPruning}
            <RefreshCw class="w-3.5 h-3.5 animate-spin" />
            <span>Limpiando disco...</span>
          {:else}
            <Sparkles class="w-3.5 h-3.5" />
            <span>Iniciar Limpieza</span>
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}
