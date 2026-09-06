<script lang="ts">
  import { Copy, Check, FileCode2 } from '@lucide/svelte';

  let {
    value = $bindable(''),
    filename = 'compose.yaml',
    readOnly = false,
    placeholder = 'Escribe el contenido del archivo aquí...',
    oninput,
  } = $props<{
    value: string;
    filename?: string;
    readOnly?: boolean;
    placeholder?: string;
    oninput?: () => void;
  }>();

  let textareaEl: HTMLTextAreaElement | null = $state(null);
  let lineNumbersEl: HTMLDivElement | null = $state(null);
  let copied = $state(false);

  const lines = $derived(value ? value.split('\n') : ['']);
  const lineCount = $derived(lines.length);

  function handleScroll() {
    if (textareaEl && lineNumbersEl) {
      lineNumbersEl.scrollTop = textareaEl.scrollTop;
    }
  }

  function handleKeyDown(e: KeyboardEvent) {
    if (readOnly) return;
    if (e.key === 'Tab') {
      e.preventDefault();
      if (!textareaEl) return;

      const start = textareaEl.selectionStart;
      const end = textareaEl.selectionEnd;

      if (e.shiftKey) {
        // Shift+Tab: Unindent
        const before = value.substring(0, start);
        const selected = value.substring(start, end);
        const after = value.substring(end);

        // Simple unindent if 2 spaces precede cursor
        if (before.endsWith('  ')) {
          value = before.slice(0, -2) + selected + after;
          setTimeout(() => {
            if (textareaEl) {
              textareaEl.selectionStart = textareaEl.selectionEnd = start - 2;
            }
          }, 0);
        }
      } else {
        // Tab: Insert 2 spaces
        const tabSpaces = '  ';
        value = value.substring(0, start) + tabSpaces + value.substring(end);
        setTimeout(() => {
          if (textareaEl) {
            textareaEl.selectionStart = textareaEl.selectionEnd = start + tabSpaces.length;
          }
        }, 0);
      }
    }
  }

  async function handleCopy() {
    if (!value) return;
    try {
      await navigator.clipboard.writeText(value);
      copied = true;
      setTimeout(() => (copied = false), 2000);
    } catch {
      // ignore
    }
  }
</script>

<div class="flex flex-col h-full rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-950 overflow-hidden font-mono text-xs shadow-xs">
  <!-- Top Editor Toolbar -->
  <div class="flex items-center justify-between px-3 py-1.5 bg-slate-100/70 dark:bg-slate-900 border-b border-slate-200 dark:border-slate-800 text-[11px] text-slate-500 dark:text-slate-400 select-none">
    <div class="flex items-center gap-1.5">
      <FileCode2 class="w-3.5 h-3.5 text-violet-500" />
      <span class="font-semibold text-slate-700 dark:text-slate-300">{filename}</span>
      <span class="text-slate-400">•</span>
      <span>{lineCount} {lineCount === 1 ? 'línea' : 'líneas'}</span>
      <span class="text-slate-400">•</span>
      <span>{value.length} carácteres</span>
    </div>

    <div class="flex items-center gap-1">
      <button
        type="button"
        onclick={handleCopy}
        class="flex items-center gap-1 px-2 py-0.5 rounded-md hover:bg-slate-200 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300 transition-colors cursor-pointer"
        title="Copiar YAML"
      >
        {#if copied}
          <Check class="w-3 h-3 text-emerald-500" />
          <span class="text-emerald-500 font-medium">Copiado</span>
        {:else}
          <Copy class="w-3 h-3" />
          <span>Copiar</span>
        {/if}
      </button>
    </div>
  </div>

  <!-- Editor Body: Line Numbers + Textarea -->
  <div class="relative flex-1 flex min-h-[300px] overflow-hidden bg-slate-50/30 dark:bg-slate-950">
    <!-- Line numbers gutter -->
    <div
      bind:this={lineNumbersEl}
      class="w-10 flex-shrink-0 select-none overflow-hidden py-3 bg-slate-100/50 dark:bg-slate-900/50 text-right pr-2 text-slate-400 dark:text-slate-600 border-r border-slate-200/60 dark:border-slate-800/80 leading-5"
    >
      {#each lines as _, i}
        <div class="h-5">{i + 1}</div>
      {/each}
    </div>

    <!-- Textarea Input Area -->
    <textarea
      bind:this={textareaEl}
      bind:value
      onscroll={handleScroll}
      onkeydown={handleKeyDown}
      {oninput}
      disabled={readOnly}
      {placeholder}
      spellcheck="false"
      class="flex-1 w-full h-full p-3 bg-transparent text-slate-900 dark:text-slate-100 resize-none outline-none leading-5 font-mono text-xs overflow-auto whitespace-pre tab-2 border-0 focus:ring-0 disabled:opacity-60"
    ></textarea>
  </div>
</div>

<style>
  .tab-2 {
    tab-size: 2;
  }
</style>
