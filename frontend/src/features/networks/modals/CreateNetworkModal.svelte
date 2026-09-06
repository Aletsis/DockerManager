<script lang="ts">
  import {
    X,
    Network,
    Plus,
    Loader2,
    Shield,
    Globe,
    Layers,
    Server,
    Settings2,
    AlertCircle,
  } from '@lucide/svelte';
  import { networksStore } from '../stores/networks.svelte';

  let {
    isOpen,
    onClose,
  } = $props<{
    isOpen: boolean;
    onClose: () => void;
  }>();

  let name = $state('');
  let driver = $state<'bridge' | 'macvlan' | 'ipvlan' | 'overlay' | 'host'>('bridge');
  let subnet = $state('');
  let gateway = $state('');
  let ipRange = $state('');
  let parentInterface = $state('');
  let isInternal = $state(false);
  let isAttachable = $state(true);
  let enableIPv6 = $state(false);
  let isSubmitting = $state(false);
  let formError = $state('');
  let showAdvanced = $state(false);

  function resetForm() {
    name = '';
    driver = 'bridge';
    subnet = '';
    gateway = '';
    ipRange = '';
    parentInterface = '';
    isInternal = false;
    isAttachable = true;
    enableIPv6 = false;
    formError = '';
    showAdvanced = false;
  }

  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      onClose();
    }
  }

  async function handleSubmit(e: Event) {
    e.preventDefault();
    if (!name.trim()) {
      formError = 'El nombre de la red es requerido.';
      return;
    }

    formError = '';
    isSubmitting = true;

    try {
      const options: Record<string, string> = {};
      if ((driver === 'macvlan' || driver === 'ipvlan') && parentInterface.trim()) {
        options['parent'] = parentInterface.trim();
      }

      const success = await networksStore.createNetwork({
        name: name.trim(),
        driver,
        subnet: subnet.trim() || undefined,
        gateway: gateway.trim() || undefined,
        ipRange: ipRange.trim() || undefined,
        internal: isInternal,
        attachable: isAttachable,
        enableIPv6,
        options: Object.keys(options).length > 0 ? options : undefined,
      });

      if (success) {
        resetForm();
        onClose();
      }
    } finally {
      isSubmitting = false;
    }
  }

  const driverDescriptions = {
    bridge: 'Red de puente local estándar para aislar contenedores en el host.',
    macvlan: 'Asigna una dirección MAC a cada contenedor haciéndolo parecer un dispositivo físico en la red.',
    ipvlan: 'Similar a macvlan pero comparte la misma dirección MAC con diferentes direcciones IP.',
    overlay: 'Red distribuida entre múltiples hosts Docker daemon (Docker Swarm).',
    host: 'Elimina el aislamiento de red entre el contenedor y el host Docker.',
  };
</script>

<svelte:window onkeydown={handleKeyDown} />

{#if isOpen}
  <!-- Backdrop -->
  <div
    class="fixed inset-0 z-50 bg-slate-950/60 backdrop-blur-xs flex items-center justify-center p-3 sm:p-4 animate-in fade-in duration-150"
    role="dialog"
    aria-modal="true"
  >
    <!-- Modal Dialog -->
    <div
      class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl w-full max-w-xl max-h-[90vh] flex flex-col shadow-2xl overflow-hidden animate-in zoom-in-95 duration-150"
    >
      <!-- Header -->
      <div class="px-5 py-4 border-b border-slate-200 dark:border-slate-800 flex items-center justify-between shrink-0 bg-slate-50/50 dark:bg-slate-800/30">
        <div class="flex items-center gap-3">
          <div class="p-2 rounded-xl bg-teal-500/10 text-teal-600 dark:text-teal-400">
            <Network class="w-5 h-5" />
          </div>
          <div>
            <h3 class="font-semibold text-base text-slate-900 dark:text-slate-100 tracking-tight leading-tight">
              Crear Red Personalizada
            </h3>
            <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
              Configura los parámetros de topología e IPAM en Docker Engine
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
      <form onsubmit={handleSubmit} class="flex-1 overflow-y-auto p-5 space-y-4 text-xs">
        {#if formError}
          <div class="p-3 rounded-xl bg-rose-50 dark:bg-rose-950/30 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 flex items-center gap-2">
            <AlertCircle class="w-4 h-4 shrink-0" />
            <span>{formError}</span>
          </div>
        {/if}

        <!-- Network Name -->
        <div class="space-y-1.5">
          <label for="net-name" class="font-medium text-slate-700 dark:text-slate-300 flex items-center justify-between">
            <span>Nombre de la Red *</span>
            <span class="text-[10px] text-slate-400 font-normal">Obligatorio</span>
          </label>
          <input
            id="net-name"
            type="text"
            placeholder="ej. backend-net, db-cluster, dev-bridge"
            bind:value={name}
            required
            class="w-full px-3 py-2 rounded-xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/80 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-1 focus:ring-teal-500 font-mono text-xs"
          />
        </div>

        <!-- Driver Selector -->
        <div class="space-y-1.5">
          <span class="block font-medium text-slate-700 dark:text-slate-300">
            Controlador de Red (Driver)
          </span>
          <div class="grid grid-cols-2 sm:grid-cols-3 gap-2">
            {#each (['bridge', 'macvlan', 'ipvlan', 'overlay', 'host'] as const) as d}
              <button
                type="button"
                onclick={() => (driver = d)}
                class="px-3 py-2 rounded-xl border text-left flex flex-col gap-0.5 cursor-pointer transition-all {driver === d ? 'border-teal-500 bg-teal-50/50 dark:bg-teal-950/30 text-teal-700 dark:text-teal-300 font-semibold' : 'border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800'}"
              >
                <span class="uppercase tracking-wider text-[11px] font-mono">{d}</span>
              </button>
            {/each}
          </div>
          <p class="text-[11px] text-slate-500 dark:text-slate-400 mt-1 italic">
            {driverDescriptions[driver]}
          </p>
        </div>

        <!-- Macvlan / Ipvlan parent interface -->
        {#if driver === 'macvlan' || driver === 'ipvlan'}
          <div class="space-y-1.5 p-3 rounded-xl bg-amber-500/10 border border-amber-500/20 animate-in fade-in duration-150">
            <label for="parent-iface" class="font-medium text-amber-900 dark:text-amber-200">
              Interfaz Física Padre (Parent Interface)
            </label>
            <input
              id="parent-iface"
              type="text"
              placeholder="ej. eth0, enp3s0, wlan0"
              bind:value={parentInterface}
              class="w-full px-3 py-1.5 rounded-lg border border-amber-300 dark:border-amber-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-1 focus:ring-amber-500 font-mono text-xs"
            />
            <p class="text-[10px] text-amber-700 dark:text-amber-300">
              Interfaz física de red del host a la que se vincularán los contenedores.
            </p>
          </div>
        {/if}

        <!-- IPAM Configuration Box -->
        <div class="border border-slate-200 dark:border-slate-800 rounded-xl p-3.5 space-y-3 bg-slate-50/30 dark:bg-slate-850/30">
          <div class="flex items-center justify-between">
            <h4 class="font-semibold text-slate-800 dark:text-slate-200 text-xs flex items-center gap-1.5">
              <Globe class="w-3.5 h-3.5 text-teal-500" />
              <span>Configuración IPAM (Subred y Enrutamiento)</span>
            </h4>
            <span class="text-[10px] text-slate-400">Opcional</span>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div class="space-y-1">
              <label for="net-subnet" class="text-[11px] font-medium text-slate-600 dark:text-slate-400">
                Subred (CIDR)
              </label>
              <input
                id="net-subnet"
                type="text"
                placeholder="172.28.0.0/16"
                bind:value={subnet}
                class="w-full px-2.5 py-1.5 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100 placeholder-slate-400 font-mono text-xs focus:ring-1 focus:ring-teal-500 focus:outline-none"
              />
            </div>

            <div class="space-y-1">
              <label for="net-gateway" class="text-[11px] font-medium text-slate-600 dark:text-slate-400">
                Puerta de Enlace (Gateway)
              </label>
              <input
                id="net-gateway"
                type="text"
                placeholder="172.28.0.1"
                bind:value={gateway}
                class="w-full px-2.5 py-1.5 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100 placeholder-slate-400 font-mono text-xs focus:ring-1 focus:ring-teal-500 focus:outline-none"
              />
            </div>
          </div>

          <div class="space-y-1">
            <label for="net-iprange" class="text-[11px] font-medium text-slate-600 dark:text-slate-400">
              Rango de Asignación IP (IP-Range)
            </label>
            <input
              id="net-iprange"
              type="text"
              placeholder="172.28.5.0/24 (opcional)"
              bind:value={ipRange}
              class="w-full px-2.5 py-1.5 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100 placeholder-slate-400 font-mono text-xs focus:ring-1 focus:ring-teal-500 focus:outline-none"
            />
          </div>
        </div>

        <!-- Advanced Options Toggle -->
        <button
          type="button"
          onclick={() => (showAdvanced = !showAdvanced)}
          class="flex items-center gap-1.5 text-xs text-slate-500 dark:text-slate-400 hover:text-slate-800 dark:hover:text-slate-200 cursor-pointer pt-1"
        >
          <Settings2 class="w-3.5 h-3.5" />
          <span>{showAdvanced ? 'Ocultar opciones avanzadas' : 'Mostrar opciones avanzadas'}</span>
        </button>

        {#if showAdvanced}
          <div class="space-y-2.5 p-3 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/40 animate-in fade-in duration-150">
            <!-- Internal Network -->
            <label class="flex items-start gap-2.5 cursor-pointer">
              <input
                type="checkbox"
                bind:checked={isInternal}
                class="mt-0.5 rounded text-teal-600 focus:ring-teal-500"
              />
              <div>
                <span class="font-medium text-slate-800 dark:text-slate-200">Red Aislada Interna (--internal)</span>
                <p class="text-[11px] text-slate-500 dark:text-slate-400">
                  Bloquea el tráfico saliente a Internet. Solo comunicación entre contenedores de esta red.
                </p>
              </div>
            </label>

            <!-- Attachable -->
            <label class="flex items-start gap-2.5 cursor-pointer">
              <input
                type="checkbox"
                bind:checked={isAttachable}
                class="mt-0.5 rounded text-teal-600 focus:ring-teal-500"
              />
              <div>
                <span class="font-medium text-slate-800 dark:text-slate-200">Permitir Conexión Manual (--attachable)</span>
                <p class="text-[11px] text-slate-500 dark:text-slate-400">
                  Permite conectar contenedores independientes en caliente a esta red.
                </p>
              </div>
            </label>

            <!-- IPv6 -->
            <label class="flex items-start gap-2.5 cursor-pointer">
              <input
                type="checkbox"
                bind:checked={enableIPv6}
                class="mt-0.5 rounded text-teal-600 focus:ring-teal-500"
              />
              <div>
                <span class="font-medium text-slate-800 dark:text-slate-200">Habilitar IPv6 (--ipv6)</span>
                <p class="text-[11px] text-slate-500 dark:text-slate-400">
                  Habilita el direccionamiento y enrutamiento IPv6 para los contenedores asociados.
                </p>
              </div>
            </label>
          </div>
        {/if}

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
            disabled={isSubmitting || !name.trim()}
            class="flex items-center gap-2 px-4 py-2 rounded-xl bg-teal-600 hover:bg-teal-700 text-white font-medium transition-colors cursor-pointer shadow-xs disabled:opacity-50"
          >
            {#if isSubmitting}
              <Loader2 class="w-4 h-4 animate-spin" />
              <span>Creando Red...</span>
            {:else}
              <Plus class="w-4 h-4" />
              <span>Crear Red</span>
            {/if}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
