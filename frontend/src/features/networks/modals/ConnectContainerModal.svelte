<script lang="ts">
  import {
    X,
    Link2,
    Loader2,
    Search,
    AlertCircle,
    CheckCircle2,
    Server,
    Network,
  } from '@lucide/svelte';
  import { containersStore } from '../../containers/stores/containers.svelte';
  import { networksStore } from '../stores/networks.svelte';

  let {
    networkId,
    networkName,
    connectedContainerIds = [],
    onClose,
  } = $props<{
    networkId: string | null;
    networkName: string;
    connectedContainerIds?: string[];
    onClose: () => void;
  }>();

  let selectedContainerId = $state('');
  let staticIp = $state('');
  let searchQuery = $state('');
  let isConnecting = $state(false);
  let error = $state('');

  // Available containers (excluding containers already connected)
  const availableContainers = $derived(
    containersStore.containers.filter((c) => {
      const isAlreadyConnected = connectedContainerIds.includes(c.id) || connectedContainerIds.includes(c.shortId);
      if (isAlreadyConnected) return false;

      if (searchQuery.trim()) {
        const q = searchQuery.toLowerCase();
        return c.name.toLowerCase().includes(q) || c.id.toLowerCase().includes(q) || c.image.toLowerCase().includes(q);
      }
      return true;
    })
  );

  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      onClose();
    }
  }

  async function handleConnect(e: Event) {
    e.preventDefault();
    if (!networkId || !selectedContainerId) {
      error = 'Selecciona un contenedor para conectar.';
      return;
    }

    error = '';
    isConnecting = true;
    try {
      const success = await networksStore.connectContainer(
        networkId,
        selectedContainerId,
        staticIp.trim()
      );
      if (success) {
        selectedContainerId = '';
        staticIp = '';
        onClose();
      }
    } catch (err: any) {
      error = err?.message || String(err);
    } finally {
      isConnecting = false;
    }
  }
</script>

<svelte:window onkeydown={handleKeyDown} />

{#if networkId}
  <!-- Backdrop -->
  <div
    class="fixed inset-0 z-50 bg-slate-950/60 backdrop-blur-xs flex items-center justify-center p-3 sm:p-4 animate-in fade-in duration-150"
    role="dialog"
    aria-modal="true"
  >
    <!-- Modal Dialog -->
    <div
      class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl w-full max-w-lg max-h-[90vh] flex flex-col shadow-2xl overflow-hidden animate-in zoom-in-95 duration-150"
    >
      <!-- Header -->
      <div class="px-5 py-4 border-b border-slate-200 dark:border-slate-800 flex items-center justify-between shrink-0 bg-slate-50/50 dark:bg-slate-800/30">
        <div class="flex items-center gap-3">
          <div class="p-2 rounded-xl bg-teal-500/10 text-teal-600 dark:text-teal-400">
            <Link2 class="w-5 h-5" />
          </div>
          <div>
            <h3 class="font-semibold text-base text-slate-900 dark:text-slate-100 tracking-tight leading-tight">
              Conectar Contenedor en Caliente
            </h3>
            <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
              Red de destino: <strong class="font-mono text-teal-600 dark:text-teal-400">{networkName}</strong>
            </p>
          </div>
        </div>

        <button
          onclick={onClose}
          class="p-2 rounded-xl text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer"
          title="Cerrar (Esc)"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Form Body -->
      <form onsubmit={handleConnect} class="flex-1 overflow-y-auto p-5 space-y-4 text-xs">
        {#if error}
          <div class="p-3 rounded-xl bg-rose-50 dark:bg-rose-950/30 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 flex items-center gap-2">
            <AlertCircle class="w-4 h-4 shrink-0" />
            <span>{error}</span>
          </div>
        {/if}

        <!-- Container Search & List -->
        <div class="space-y-2">
          <label for="container-search" class="font-medium text-slate-700 dark:text-slate-300 flex items-center justify-between">
            <span>Seleccionar Contenedor *</span>
            <span class="text-[10px] text-slate-400 font-normal">{availableContainers.length} disponibles</span>
          </label>

          <!-- Search input -->
          <div class="relative">
            <Search class="w-3.5 h-3.5 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" />
            <input
              id="container-search"
              type="text"
              placeholder="Filtrar por nombre o imagen..."
              bind:value={searchQuery}
              class="w-full pl-8.5 pr-3 py-1.5 rounded-lg border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/80 text-slate-900 dark:text-slate-100 placeholder-slate-400 text-xs focus:outline-none focus:ring-1 focus:ring-teal-500"
            />
          </div>

          <!-- Container Selection Scrollable Box -->
          <div class="border border-slate-200 dark:border-slate-700 rounded-xl max-h-48 overflow-y-auto divide-y divide-slate-100 dark:divide-slate-800 bg-white dark:bg-slate-900/60">
            {#if availableContainers.length === 0}
              <div class="p-4 text-center text-slate-400 text-xs">
                No hay contenedores disponibles para conectar a esta red.
              </div>
            {:else}
              {#each availableContainers as c}
                <button
                  type="button"
                  onclick={() => (selectedContainerId = c.id)}
                  class="w-full px-3 py-2.5 flex items-center justify-between text-left transition-colors cursor-pointer {selectedContainerId === c.id ? 'bg-teal-50 dark:bg-teal-950/40 text-teal-900 dark:text-teal-200' : 'hover:bg-slate-50 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-300'}"
                >
                  <div class="flex items-center gap-2.5 min-w-0">
                    <span class="w-2 h-2 rounded-full shrink-0 {c.state === 'running' ? 'bg-emerald-500' : 'bg-slate-400'}"></span>
                    <div class="min-w-0">
                      <div class="font-medium text-xs truncate">{c.name}</div>
                      <div class="text-[10px] text-slate-400 font-mono truncate">{c.image}</div>
                    </div>
                  </div>

                  {#if selectedContainerId === c.id}
                    <CheckCircle2 class="w-4 h-4 text-teal-600 dark:text-teal-400 shrink-0" />
                  {/if}
                </button>
              {/each}
            {/if}
          </div>
        </div>

        <!-- Optional Static IPv4 Address -->
        <div class="space-y-1.5 pt-1">
          <label for="static-ip" class="font-medium text-slate-700 dark:text-slate-300 flex items-center justify-between">
            <span>Dirección IPv4 Estática</span>
            <span class="text-[10px] text-slate-400 font-normal">Opcional (DHCP automático si se omite)</span>
          </label>
          <input
            id="static-ip"
            type="text"
            placeholder="ej. 172.28.0.50"
            bind:value={staticIp}
            class="w-full px-3 py-2 rounded-xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/80 text-slate-900 dark:text-slate-100 placeholder-slate-400 font-mono text-xs focus:outline-none focus:ring-1 focus:ring-teal-500"
          />
        </div>

        <!-- Footer -->
        <div class="pt-3 border-t border-slate-200 dark:border-slate-800 flex items-center justify-end gap-2 shrink-0">
          <button
            type="button"
            onclick={onClose}
            class="px-4 py-2 rounded-xl border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-300 font-medium transition-colors cursor-pointer"
          >
            Cancelar
          </button>
          <button
            type="submit"
            disabled={isConnecting || !selectedContainerId}
            class="flex items-center gap-2 px-4 py-2 rounded-xl bg-teal-600 hover:bg-teal-700 text-white font-medium transition-colors cursor-pointer shadow-xs disabled:opacity-50"
          >
            {#if isConnecting}
              <Loader2 class="w-4 h-4 animate-spin" />
              <span>Conectando...</span>
            {:else}
              <Link2 class="w-4 h-4" />
              <span>Conectar en Caliente</span>
            {/if}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
