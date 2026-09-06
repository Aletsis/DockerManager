<script lang="ts">
  import { AlertTriangle, X } from '@lucide/svelte';

  let {
    isOpen = false,
    title,
    message,
    confirmLabel = 'Confirmar',
    onConfirm,
    onCancel,
    isDestructive = true,
  } = $props<{
    isOpen: boolean;
    title: string;
    message: string;
    confirmLabel?: string;
    onConfirm: () => void;
    onCancel: () => void;
    isDestructive?: boolean;
  }>();
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 bg-slate-950/60 backdrop-blur-sm flex items-center justify-center p-4">
    <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl w-full max-w-md shadow-2xl p-5 space-y-4 animate-in fade-in zoom-in-95 duration-150">
      <div class="flex items-start justify-between gap-3">
        <div class="flex items-center gap-3">
          <div class="p-2 rounded-xl bg-rose-500/10 text-rose-600 dark:text-rose-400">
            <AlertTriangle class="w-5 h-5" />
          </div>
          <h3 class="font-semibold text-slate-900 dark:text-slate-100 text-sm">
            {title}
          </h3>
        </div>
        <button
          onclick={onCancel}
          class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 p-1"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <p class="text-xs text-slate-600 dark:text-slate-300 leading-relaxed">
        {message}
      </p>

      <div class="flex items-center justify-end gap-2 pt-2">
        <button
          onclick={onCancel}
          class="px-3 py-1.5 rounded-lg border border-slate-200 dark:border-slate-700 text-xs font-medium text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
        >
          Cancelar
        </button>
        <button
          onclick={onConfirm}
          class="px-3 py-1.5 rounded-lg text-xs font-medium text-white transition-colors {isDestructive ? 'bg-rose-600 hover:bg-rose-700' : 'bg-indigo-600 hover:bg-indigo-700'}"
        >
          {confirmLabel}
        </button>
      </div>
    </div>
  </div>
{/if}
