<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { Terminal } from '@xterm/xterm';
  import { FitAddon } from '@xterm/addon-fit';
  import '@xterm/xterm/css/xterm.css';
  import {
    X,
    Terminal as TerminalIcon,
    RefreshCw,
    Maximize2,
    Minimize2,
    Eraser,
    ChevronDown,
  } from '@lucide/svelte';
  import {
    StartTerminal,
    WriteTerminal,
    ResizeTerminal,
    CloseTerminal,
  } from '../../../wailsjs/go/main/App';
  import { EventsOn, EventsOff } from '../../../wailsjs/runtime/runtime';

  let {
    containerId,
    containerName,
    onClose,
  } = $props<{
    containerId: string | null;
    containerName: string;
    onClose: () => void;
  }>();

  let terminalContainerEl: HTMLDivElement | null = $state(null);
  let isFullscreen = $state<boolean>(false);
  let status = $state<'connecting' | 'connected' | 'disconnected' | 'error'>('connecting');
  let statusMessage = $state<string>('');
  let currentShell = $state<string>('auto');
  let activeShellName = $state<string>('');
  let sessionId = $state<string | null>(null);

  let term: Terminal | null = null;
  let fitAddon: FitAddon | null = null;
  let resizeObserver: ResizeObserver | null = null;
  let resizeDebounceTimer: any = null;

  function base64ToUint8Array(b64: string): Uint8Array {
    const binaryStr = atob(b64);
    const len = binaryStr.length;
    const bytes = new Uint8Array(len);
    for (let i = 0; i < len; i++) {
      bytes[i] = binaryStr.charCodeAt(i);
    }
    return bytes;
  }

  async function initSession() {
    if (!containerId || !terminalContainerEl) return;

    // Teardown previous session if existing
    teardownSession();

    status = 'connecting';
    statusMessage = 'Iniciando terminal PTY en el contenedor...';

    // Initialize xterm
    term = new Terminal({
      cursorBlink: true,
      cursorStyle: 'bar',
      fontSize: 13,
      lineHeight: 1.25,
      fontFamily: "'JetBrains Mono', 'Fira Code', 'Cascadia Code', Menlo, Monaco, Consolas, monospace",
      theme: {
        background: '#090d16',
        foreground: '#e2e8f0',
        cursor: '#6366f1',
        cursorAccent: '#090d16',
        selectionBackground: 'rgba(99, 102, 241, 0.35)',
        black: '#1e293b',
        red: '#f43f5e',
        green: '#10b981',
        yellow: '#f59e0b',
        blue: '#3b82f6',
        magenta: '#d946ef',
        cyan: '#06b6d4',
        white: '#f8fafc',
        brightBlack: '#475569',
        brightRed: '#fb7185',
        brightGreen: '#34d399',
        brightYellow: '#fbbf24',
        brightBlue: '#60a5fa',
        brightMagenta: '#e879f9',
        brightCyan: '#22d3ee',
        brightWhite: '#ffffff',
      },
    });

    fitAddon = new FitAddon();
    term.loadAddon(fitAddon);
    term.open(terminalContainerEl);

    // Initial fit
    try {
      fitAddon.fit();
    } catch {
      // Ignore initial render fit error if DOM is settling
    }

    const rows = term.rows || 24;
    const cols = term.cols || 80;

    try {
      const res = await StartTerminal(containerId, currentShell, rows, cols);
      sessionId = res.sessionId;
      activeShellName = res.shell;
      status = 'connected';
      statusMessage = '';

      // Register event listener for incoming output
      EventsOn(`terminal:data:${res.sessionId}`, (base64Chunk: string) => {
        if (!term) return;
        try {
          const bytes = base64ToUint8Array(base64Chunk);
          term.write(bytes);
        } catch {
          term.write(atob(base64Chunk));
        }
      });

      // Register event listener for exit
      EventsOn(`terminal:exit:${res.sessionId}`, () => {
        status = 'disconnected';
        statusMessage = 'Sesión cerrada';
        if (term) {
          term.write('\r\n\x1b[33m[Proceso finalizado. Puedes cerrar la ventana o reconectar]\x1b[0m\r\n');
        }
      });

      // Handle user keyboard input
      term.onData((data: string) => {
        if (sessionId && status === 'connected') {
          WriteTerminal(sessionId, data).catch((err) => {
            console.error('Error writing to terminal:', err);
          });
        }
      });

      // Re-fit once mounted
      setTimeout(() => {
        handleResize();
      }, 50);

      term.focus();
    } catch (err: any) {
      status = 'error';
      statusMessage = err?.toString() || 'Error al conectar con la terminal';
      if (term) {
        term.write(`\r\n\x1b[31mError al iniciar terminal: ${statusMessage}\x1b[0m\r\n`);
      }
    }
  }

  function handleResize() {
    if (!term || !fitAddon || !sessionId || status !== 'connected') return;
    try {
      fitAddon.fit();
      const rows = term.rows;
      const cols = term.cols;
      if (rows > 0 && cols > 0) {
        ResizeTerminal(sessionId, rows, cols).catch(() => {});
      }
    } catch (e) {
      // Fit calculation can fail if hidden
    }
  }

  function onContainerResize() {
    if (resizeDebounceTimer) clearTimeout(resizeDebounceTimer);
    resizeDebounceTimer = setTimeout(() => {
      handleResize();
    }, 60);
  }

  function teardownSession() {
    if (sessionId) {
      EventsOff(`terminal:data:${sessionId}`);
      EventsOff(`terminal:exit:${sessionId}`);
      CloseTerminal(sessionId).catch(() => {});
      sessionId = null;
    }

    if (term) {
      term.dispose();
      term = null;
    }

    if (terminalContainerEl) {
      terminalContainerEl.innerHTML = '';
    }
  }

  function handleClear() {
    if (term) {
      term.clear();
      term.focus();
    }
  }

  function handleShellChange(newShell: string) {
    currentShell = newShell;
    initSession();
  }

  function handleClose() {
    teardownSession();
    onClose();
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape' && !isFullscreen) {
      handleClose();
    } else if (e.key === 'Enter' && status === 'disconnected') {
      initSession();
    }
  }

  $effect(() => {
    if (containerId && terminalContainerEl) {
      initSession();
    }
  });

  onMount(() => {
    window.addEventListener('keydown', handleKeydown);

    if (terminalContainerEl) {
      resizeObserver = new ResizeObserver(() => {
        onContainerResize();
      });
      resizeObserver.observe(terminalContainerEl);
    }
  });

  onDestroy(() => {
    window.removeEventListener('keydown', handleKeydown);
    if (resizeObserver) {
      resizeObserver.disconnect();
      resizeObserver = null;
    }
    if (resizeDebounceTimer) {
      clearTimeout(resizeDebounceTimer);
    }
    teardownSession();
  });
</script>

{#if containerId}
  <div class="fixed inset-0 z-50 bg-slate-950/70 backdrop-blur-sm flex items-center justify-center p-3 sm:p-5 transition-all">
    <div
      class="bg-[#090d16] border border-slate-800 rounded-2xl w-full flex flex-col shadow-2xl overflow-hidden transition-all duration-200 {isFullscreen
        ? 'fixed inset-2 w-[calc(100%-1rem)] h-[calc(100%-1rem)] rounded-xl z-50'
        : 'max-w-5xl h-[85vh]'}"
    >
      <!-- Terminal Header Bar -->
      <div class="flex items-center justify-between px-4 py-3 bg-slate-900/90 border-b border-slate-800 flex-wrap gap-2 select-none">
        <!-- Left title & container metadata -->
        <div class="flex items-center gap-2.5 min-w-0">
          <div class="p-1.5 rounded-lg bg-indigo-500/10 text-indigo-400 flex-shrink-0">
            <TerminalIcon class="w-4 h-4" />
          </div>

          <div class="flex items-center gap-2 min-w-0">
            <h2 class="text-sm font-semibold text-slate-100 truncate">
              {containerName}
            </h2>

            <!-- Status Pill -->
            <div class="flex items-center gap-1.5 px-2 py-0.5 rounded-full text-[11px] font-medium border {status === 'connected' ? 'bg-emerald-950/40 text-emerald-400 border-emerald-800/60' : status === 'connecting' ? 'bg-amber-950/40 text-amber-400 border-amber-800/60' : 'bg-slate-800 text-slate-400 border-slate-700'}">
              <span class="w-1.5 h-1.5 rounded-full {status === 'connected' ? 'bg-emerald-400 animate-pulse' : status === 'connecting' ? 'bg-amber-400 animate-ping' : 'bg-slate-400'}"></span>
              <span>
                {status === 'connected' ? 'Conectado' : status === 'connecting' ? 'Conectando...' : 'Desconectado'}
              </span>
            </div>

            <!-- Active Shell Badge -->
            {#if activeShellName && status === 'connected'}
              <span class="hidden sm:inline-block px-2 py-0.5 rounded text-[11px] font-mono text-indigo-300 bg-indigo-950/50 border border-indigo-800/50">
                {activeShellName}
              </span>
            {/if}
          </div>
        </div>

        <!-- Right Toolbar Actions -->
        <div class="flex items-center gap-1.5 flex-wrap">
          <!-- Shell Selector -->
          <div class="relative">
            <select
              aria-label="Seleccionar Shell"
              value={currentShell}
              onchange={(e) => handleShellChange((e.target as HTMLSelectElement).value)}
              class="text-xs py-1 pl-2.5 pr-7 rounded-lg border border-slate-700 bg-slate-800 text-slate-200 focus:outline-none focus:ring-1 focus:ring-indigo-500 appearance-none cursor-pointer"
            >
              <option value="auto">Shell: Auto (bash/sh)</option>
              <option value="/bin/bash">/bin/bash</option>
              <option value="/bin/sh">/bin/sh</option>
              <option value="/bin/zsh">/bin/zsh</option>
              <option value="sh">sh</option>
            </select>
            <ChevronDown class="w-3.5 h-3.5 text-slate-400 absolute right-2 top-1/2 -translate-y-1/2 pointer-events-none" />
          </div>

          <!-- Clear Terminal -->
          <button
            onclick={handleClear}
            title="Limpiar pantalla"
            class="p-1.5 rounded-lg border border-slate-700 hover:bg-slate-800 text-slate-300 hover:text-slate-100 transition-colors"
          >
            <Eraser class="w-3.5 h-3.5" />
          </button>

          <!-- Reconnect -->
          <button
            onclick={initSession}
            title="Reconectar sesión"
            class="px-2 py-1 rounded-lg border flex items-center gap-1.5 text-xs transition-colors cursor-pointer {status === 'disconnected' ? 'border-indigo-500 bg-indigo-600 text-white hover:bg-indigo-500 shadow-sm' : 'border-slate-700 hover:bg-slate-800 text-slate-300 hover:text-slate-100'}"
          >
            <RefreshCw class="w-3.5 h-3.5 {status === 'connecting' ? 'animate-spin text-indigo-400' : ''}" />
            {#if status === 'disconnected'}
              <span class="font-sans font-medium text-xs">Reconectar</span>
            {/if}
          </button>

          <!-- Toggle Fullscreen -->
          <button
            onclick={() => {
              isFullscreen = !isFullscreen;
              setTimeout(handleResize, 120);
            }}
            title={isFullscreen ? 'Restaurar tamaño' : 'Pantalla completa'}
            class="p-1.5 rounded-lg border border-slate-700 hover:bg-slate-800 text-slate-300 hover:text-slate-100 transition-colors"
          >
            {#if isFullscreen}
              <Minimize2 class="w-3.5 h-3.5" />
            {:else}
              <Maximize2 class="w-3.5 h-3.5" />
            {/if}
          </button>

          <!-- Close Modal -->
          <button
            onclick={handleClose}
            title="Cerrar terminal"
            class="p-1.5 rounded-lg hover:bg-rose-950/40 border border-transparent hover:border-rose-800/60 text-slate-400 hover:text-rose-400 transition-colors ml-1"
          >
            <X class="w-4 h-4" />
          </button>
        </div>
      </div>

      <!-- Terminal Body / Canvas Mount Point -->
      <div
        bind:this={terminalContainerEl}
        class="flex-1 w-full h-full p-2 bg-[#090d16] overflow-hidden focus:outline-none"
      ></div>

      <!-- Footer Info Bar -->
      <div class="px-4 py-2 bg-slate-900/90 border-t border-slate-800 text-[11px] text-slate-400 flex items-center justify-between font-mono flex-wrap gap-2">
        <div class="flex items-center gap-3">
          <span>Terminal interactiva xterm.js</span>
          {#if status === 'disconnected'}
            <span class="text-amber-400 font-sans font-medium flex items-center gap-1.5">
              <span class="w-1.5 h-1.5 rounded-full bg-amber-400"></span>
              Sesión finalizada
            </span>
            <button
              onclick={initSession}
              class="px-2.5 py-0.5 rounded-md bg-indigo-600 hover:bg-indigo-500 text-white font-sans text-xs font-medium flex items-center gap-1 transition-colors shadow-xs cursor-pointer"
            >
              <RefreshCw class="w-3 h-3" />
              <span>Reconectar terminal</span>
            </button>
          {:else if statusMessage}
            <span class="text-amber-400 font-sans">{statusMessage}</span>
          {/if}
        </div>
        <div class="flex items-center gap-2 text-slate-500">
          {#if status === 'disconnected'}
            <span class="text-indigo-400 font-sans">Presiona Enter para reconectar</span>
            <span>•</span>
          {/if}
          <span>Esc para salir</span>
          <span>•</span>
          <span>Ctrl+C para interrumpir</span>
        </div>
      </div>
    </div>
  </div>
{/if}

<style>
  :global(.xterm) {
    height: 100%;
    padding: 2px 4px;
  }
  :global(.xterm .xterm-viewport) {
    overflow-y: auto !important;
  }
</style>
