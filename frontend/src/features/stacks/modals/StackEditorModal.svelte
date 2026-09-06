<script lang="ts">
  import {
    Layers,
    X,
    FolderOpen,
    Play,
    Save,
    Trash2,
    Sparkles,
    AlertCircle,
    CheckCircle2,
    FileCode,
    FileText,
    Plus,
    Terminal,
  } from '@lucide/svelte';
  import YamlEditor from '../components/YamlEditor.svelte';
  import DeployConsole from '../components/DeployConsole.svelte';
  import { COMPOSE_PRESETS } from '../constants/presets';
  import { dockerApi } from '../../../shared/services/api';
  import { uiStore } from '../../../shared/stores/ui.svelte';
  import { containersStore } from '../../containers/stores/containers.svelte';

  export interface StackFile {
    id: string;
    name: string;
    path: string;
    type: 'compose' | 'dockerfile';
    serviceName?: string;
    content: string;
    existsOnDisk: boolean;
    isModified?: boolean;
  }

  let {
    isOpen = false,
    mode = 'create',
    initialStackName = '',
    initialWorkingDir = '',
    initialConfigFile = '',
    onClose,
    onSuccess,
  } = $props<{
    isOpen: boolean;
    mode?: 'create' | 'edit';
    initialStackName?: string;
    initialWorkingDir?: string;
    initialConfigFile?: string;
    onClose: () => void;
    onSuccess?: () => void;
  }>();

  let projectName = $state('');
  let workingDir = $state('');
  let configFile = $state('');
  let files = $state<StackFile[]>([]);
  let activeFileId = $state<string>('compose');

  let selectedPreset = $state('minimal');
  let isExecuting = $state(false);
  let executionStatus = $state<'idle' | 'running' | 'success' | 'error'>('idle');
  let executionMessage = $state('');
  let showConsole = $state(false);
  let deployConsoleComponent: any = $state(null);

  const activeFile = $derived(
    files.find((f) => f.id === activeFileId) || files[0]
  );
  const composeFile = $derived(
    files.find((f) => f.type === 'compose')
  );
  const dockerfileCount = $derived(
    files.filter((f) => f.type === 'dockerfile').length
  );

  // Validation
  const isValidName = $derived(/^[a-zA-Z0-9][a-zA-Z0-9_-]*$/.test(projectName.trim()));
  const isValidToDeploy = $derived(
    projectName.trim().length > 0 &&
    isValidName &&
    composeFile &&
    composeFile.content.trim().length > 0 &&
    !isExecuting
  );

  // Initialize or load stack data
  $effect(() => {
    if (isOpen) {
      projectName = initialStackName || '';
      workingDir = initialWorkingDir || '';
      configFile = initialConfigFile || '';
      executionStatus = 'idle';
      executionMessage = '';
      showConsole = false;

      if (mode === 'edit') {
        loadExistingStack();
      } else {
        // Apply default preset
        const defaultPreset = COMPOSE_PRESETS.find((p) => p.id === selectedPreset) || COMPOSE_PRESETS[0];
        if (!projectName) {
          projectName = defaultPreset.defaultName;
        }
        files = [
          {
            id: 'compose',
            name: 'compose.yaml',
            path: '',
            type: 'compose',
            content: defaultPreset.yaml,
            existsOnDisk: false,
            isModified: false,
          },
        ];
        activeFileId = 'compose';
        updateDefaultDirectory(projectName);
      }
    }
  });

  async function updateDefaultDirectory(name: string) {
    if (!workingDir && name.trim()) {
      try {
        const defaultDir = await dockerApi.getDefaultStackDirectory(name.trim());
        workingDir = defaultDir;
      } catch {
        // ignore
      }
    }
  }

  async function loadExistingStack() {
    isExecuting = true;
    try {
      const info = await dockerApi.getStackComposeFile(projectName, workingDir, configFile);
      if (info) {
        workingDir = info.workingDir || workingDir;
        configFile = info.configFile || configFile;

        const resolvedComposeName = configFile
          ? configFile.split(/[/\\]/).pop() || 'compose.yaml'
          : 'compose.yaml';

        const mainCompose: StackFile = {
          id: 'compose',
          name: resolvedComposeName,
          path: info.configFile || '',
          type: 'compose',
          content: info.content || (info.existsOnDisk ? '' : COMPOSE_PRESETS[0].yaml),
          existsOnDisk: info.existsOnDisk,
          isModified: false,
        };

        const discoveredDfs: StackFile[] = (info.dockerfiles || []).map((df, idx) => ({
          id: `dockerfile-${idx}`,
          name: df.name || 'Dockerfile',
          path: df.path,
          type: 'dockerfile',
          serviceName: df.serviceName,
          content: df.content,
          existsOnDisk: df.existsOnDisk,
          isModified: false,
        }));

        files = [mainCompose, ...discoveredDfs];
        activeFileId = 'compose';
      }
    } catch (err: any) {
      uiStore.showToast(`Error al cargar stack: ${err}`, 'error');
    } finally {
      isExecuting = false;
    }
  }

  function handlePresetChange(e: Event) {
    const target = e.target as HTMLSelectElement;
    const preset = COMPOSE_PRESETS.find((p) => p.id === target.value);
    if (preset && composeFile) {
      selectedPreset = preset.id;
      composeFile.content = preset.yaml;
      composeFile.isModified = true;
      if (
        mode === 'create' &&
        (!projectName ||
          projectName === 'nginx-web' ||
          projectName === 'db-stack' ||
          projectName === 'node-mongo-stack' ||
          projectName === 'wordpress-site')
      ) {
        projectName = preset.defaultName;
        updateDefaultDirectory(preset.defaultName);
      }
    }
  }

  async function handleSelectDirectory() {
    try {
      const selected = await dockerApi.selectDirectory();
      if (selected) {
        workingDir = selected;
      }
    } catch (err: any) {
      uiStore.showToast(`Error al seleccionar carpeta: ${err}`, 'error');
    }
  }

  function handleAddDockerfile() {
    const defaultDockerfile = `# Dockerfile
FROM alpine:latest
WORKDIR /app
CMD ["echo", "Hola desde Dockerfile"]
`;
    const newPath = workingDir ? `${workingDir}/Dockerfile` : 'Dockerfile';
    const newDf: StackFile = {
      id: `dockerfile-${Date.now()}`,
      name: 'Dockerfile',
      path: newPath,
      type: 'dockerfile',
      content: defaultDockerfile,
      existsOnDisk: false,
      isModified: true,
    };
    files = [...files, newDf];
    activeFileId = newDf.id;
  }

  async function saveFile(file: StackFile): Promise<void> {
    if (file.type === 'compose') {
      const savedPath = await dockerApi.saveStackComposeFile(
        workingDir,
        file.path || configFile,
        file.content
      );
      file.path = savedPath;
      configFile = savedPath;
      file.existsOnDisk = true;
      file.isModified = false;
    } else {
      const targetPath =
        file.path || (workingDir ? `${workingDir}/${file.name}` : file.name);
      await dockerApi.saveStackFile(targetPath, file.content);
      file.path = targetPath;
      file.existsOnDisk = true;
      file.isModified = false;
    }
  }

  async function handleSaveFile() {
    if (!activeFile) return;
    try {
      isExecuting = true;
      await saveFile(activeFile);
      uiStore.showToast(`"${activeFile.name}" guardado exitosamente`);
    } catch (err: any) {
      uiStore.showToast(`Error al guardar: ${err}`, 'error');
    } finally {
      isExecuting = false;
    }
  }

  async function handleDeployUp() {
    if (!isValidToDeploy || !composeFile) return;

    showConsole = true;
    isExecuting = true;
    executionStatus = 'running';
    executionMessage = '';
    if (deployConsoleComponent) {
      deployConsoleComponent.clear();
    }

    try {
      // Save all files with pending modifications before deploying
      for (const file of files) {
        if (file.isModified || !file.existsOnDisk) {
          await saveFile(file);
        }
      }

      await dockerApi.upStack({
        projectName: projectName.trim(),
        workingDir: workingDir.trim(),
        configFile: composeFile.path || configFile.trim(),
        content: composeFile.content,
        removeOrphans: true,
      });

      executionStatus = 'success';
      executionMessage = 'Stack desplegado exitosamente con Docker Compose.';
      uiStore.showToast(`Stack "${projectName}" desplegado exitosamente`);
      await containersStore.fetchData(true);
      if (onSuccess) onSuccess();
    } catch (err: any) {
      executionStatus = 'error';
      executionMessage = `Error en el despliegue: ${err}`;
      uiStore.showToast(`Fallo al desplegar stack: ${err}`, 'error');
    } finally {
      isExecuting = false;
    }
  }

  async function handleDown() {
    if (!projectName.trim()) return;

    const confirmed = confirm(
      `¿Estás seguro de desmontar y eliminar el stack "${projectName}"? Esto detendrá y eliminará sus contenedores y redes asociadas.`
    );
    if (!confirmed) return;

    showConsole = true;
    isExecuting = true;
    executionStatus = 'running';
    executionMessage = '';
    if (deployConsoleComponent) {
      deployConsoleComponent.clear();
    }

    try {
      await dockerApi.downStack({
        projectName: projectName.trim(),
        workingDir: workingDir.trim(),
        configFile: (composeFile && composeFile.path) || configFile.trim(),
        removeVolumes: false,
      });

      executionStatus = 'success';
      executionMessage = 'Stack desmontado y detenido exitosamente.';
      uiStore.showToast(`Stack "${projectName}" desmontado`);
      await containersStore.fetchData(true);
      if (onSuccess) onSuccess();
    } catch (err: any) {
      executionStatus = 'error';
      executionMessage = `Error al desmontar stack: ${err}`;
      uiStore.showToast(`Error al desmontar: ${err}`, 'error');
    } finally {
      isExecuting = false;
    }
  }
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 sm:p-6 animate-in fade-in duration-200">
    <!-- Backdrop -->
    <div
      class="fixed inset-0 bg-slate-900/60 backdrop-blur-xs transition-opacity"
      onclick={() => {
        if (!isExecuting) onClose();
      }}
      role="presentation"
    ></div>

    <!-- Modal Dialog -->
    <div class="relative w-full max-w-4xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl shadow-2xl overflow-hidden flex flex-col max-h-[92vh] z-10">
      <!-- Modal Header -->
      <div class="flex items-center justify-between px-6 py-4 border-b border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50 select-none flex-shrink-0">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-violet-500/10 text-violet-600 dark:text-violet-400 border border-violet-500/20 flex items-center justify-center">
            <Layers class="w-5 h-5" />
          </div>
          <div>
            <div class="flex items-center gap-2">
              <h3 class="font-bold text-base text-slate-900 dark:text-slate-100">
                {mode === 'create' ? 'Nuevo Stack (Docker Compose)' : `Editar Stack: ${projectName}`}
              </h3>
              <span class="text-[10px] font-semibold uppercase tracking-wider px-2 py-0.5 rounded-md bg-violet-100 dark:bg-violet-900/50 text-violet-700 dark:text-violet-300 border border-violet-200 dark:border-violet-800">
                {mode === 'create' ? 'Creación' : 'Gestión'}
              </span>
            </div>
            <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
              {mode === 'create'
                ? 'Configura y despliega un conjunto de servicios coordinados mediante Docker Compose.'
                : 'Inspecciona, edita y re-despliega compose.yml y Dockerfiles detectados.'}
            </p>
          </div>
        </div>

        <button
          type="button"
          onclick={onClose}
          disabled={isExecuting}
          class="p-2 rounded-xl text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors disabled:opacity-40 cursor-pointer"
          title="Cerrar modal"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Modal Body -->
      <div class="flex-1 overflow-y-auto p-6 space-y-4">
        <!-- Configuration Bar: Project Name, Presets, Working Directory -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- Stack Name Field -->
          <div class="space-y-1.5">
            <label for="stack-name" class="block text-xs font-semibold text-slate-700 dark:text-slate-300">
              Nombre del Stack
            </label>
            <div class="relative">
              <input
                id="stack-name"
                type="text"
                bind:value={projectName}
                disabled={mode === 'edit' || isExecuting}
                oninput={() => updateDefaultDirectory(projectName)}
                placeholder="ej. mi-aplicacion-web"
                class="w-full px-3 py-2 rounded-xl text-xs bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-violet-500/20 focus:border-violet-500 disabled:opacity-60 font-mono"
              />
            </div>
            {#if projectName && !isValidName}
              <p class="text-[11px] text-rose-500">
                El nombre solo puede contener letras, números, guiones y guiones bajos.
              </p>
            {/if}
          </div>

          <!-- Presets Selector (in create mode) or Discovery Status (in edit mode) -->
          {#if mode === 'create'}
            <div class="space-y-1.5">
              <label for="stack-preset" class="block text-xs font-semibold text-slate-700 dark:text-slate-300 flex items-center justify-between">
                <span>Plantilla Inicial (Preset)</span>
                <span class="text-[11px] text-violet-600 dark:text-violet-400 font-normal flex items-center gap-1">
                  <Sparkles class="w-3 h-3" />
                  Rápido
                </span>
              </label>
              <select
                id="stack-preset"
                bind:value={selectedPreset}
                onchange={handlePresetChange}
                disabled={isExecuting}
                class="w-full px-3 py-2 rounded-xl text-xs bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800 text-slate-900 dark:text-slate-100 focus:outline-none focus:ring-2 focus:ring-violet-500/20 focus:border-violet-500 cursor-pointer disabled:opacity-60"
              >
                {#each COMPOSE_PRESETS as preset}
                  <option value={preset.id}>{preset.name} - {preset.description}</option>
                {/each}
              </select>
            </div>
          {:else}
            <!-- In Edit Mode: Discovered Files Status -->
            <div class="space-y-1.5">
              <span class="block text-xs font-semibold text-slate-700 dark:text-slate-300">
                Archivos Detectados
              </span>
              <div class="flex items-center gap-2 p-2 rounded-xl border text-xs {composeFile?.existsOnDisk ? 'bg-emerald-50 dark:bg-emerald-950/20 border-emerald-200 dark:border-emerald-800/40 text-emerald-700 dark:text-emerald-300' : 'bg-amber-50 dark:bg-amber-950/20 border-amber-200 dark:border-amber-800/40 text-amber-700 dark:text-amber-300'}">
                {#if composeFile?.existsOnDisk}
                  <CheckCircle2 class="w-4 h-4 flex-shrink-0" />
                  <span class="truncate">
                    {composeFile.name} detectado
                    {#if dockerfileCount > 0}
                      • {dockerfileCount} {dockerfileCount === 1 ? 'Dockerfile' : 'Dockerfiles'} encontrados
                    {/if}
                  </span>
                {:else}
                  <AlertCircle class="w-4 h-4 flex-shrink-0" />
                  <span class="truncate">No se encontró archivo en la ruta original. Se creará al guardar.</span>
                {/if}
              </div>
            </div>
          {/if}
        </div>

        <!-- Working Directory Selection -->
        <div class="space-y-1.5">
          <label for="working-dir" class="block text-xs font-semibold text-slate-700 dark:text-slate-300 flex items-center justify-between">
            <span>Directorio de Trabajo (Working Directory)</span>
            {#if activeFile?.path}
              <span class="text-[11px] text-slate-400 font-mono truncate max-w-sm" title={activeFile.path}>
                Ruta: {activeFile.path}
              </span>
            {/if}
          </label>
          <div class="flex items-center gap-2">
            <div class="relative flex-1">
              <input
                id="working-dir"
                type="text"
                bind:value={workingDir}
                disabled={isExecuting}
                placeholder="Ruta en tu equipo donde residen los archivos del stack"
                class="w-full px-3 py-2 rounded-xl text-xs bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-violet-500/20 focus:border-violet-500 font-mono"
              />
            </div>
            <button
              type="button"
              onclick={handleSelectDirectory}
              disabled={isExecuting}
              class="px-3 py-2 rounded-xl text-xs font-medium border border-slate-200 dark:border-slate-800 bg-slate-100 dark:bg-slate-800 hover:bg-slate-200 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200 transition-colors flex items-center gap-1.5 flex-shrink-0 cursor-pointer disabled:opacity-50"
              title="Explorar carpetas del equipo"
            >
              <FolderOpen class="w-3.5 h-3.5" />
              <span>Examinar...</span>
            </button>
          </div>
        </div>

        <!-- Multi-file Tabs & Editor Area -->
        <div class="space-y-1.5">
          <div class="flex items-center justify-between">
            <span class="block text-xs font-semibold text-slate-700 dark:text-slate-300">
              Archivos del Stack ({files.length})
            </span>
            <button
              type="button"
              onclick={() => (showConsole = !showConsole)}
              class="text-[11px] font-medium text-violet-600 dark:text-violet-400 hover:underline flex items-center gap-1 cursor-pointer"
            >
              <Terminal class="w-3 h-3" />
              <span>{showConsole ? 'Ocultar Consola' : 'Mostrar Consola de Despliegue'}</span>
            </button>
          </div>

          <!-- File Switcher Tabs Bar -->
          <div class="flex items-center justify-between border-b border-slate-200 dark:border-slate-800 bg-slate-100/70 dark:bg-slate-950 rounded-t-xl px-2 pt-1.5 overflow-x-auto select-none gap-1">
            <div class="flex items-center gap-1 overflow-x-auto">
              {#each files as file (file.id)}
                <button
                  type="button"
                  onclick={() => (activeFileId = file.id)}
                  class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded-t-lg transition-colors border-t border-x cursor-pointer {activeFileId === file.id ? 'bg-white dark:bg-slate-900 text-violet-600 dark:text-violet-400 border-slate-200 dark:border-slate-800 font-bold shadow-xs' : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200 border-transparent hover:bg-slate-200/50 dark:hover:bg-slate-800/50'}"
                >
                  {#if file.type === 'compose'}
                    <FileCode class="w-3.5 h-3.5 text-violet-500" />
                  {:else}
                    <FileText class="w-3.5 h-3.5 text-sky-500" />
                  {/if}
                  <span>{file.name}</span>
                  {#if file.serviceName}
                    <span class="text-[9px] px-1 py-0.2 rounded bg-slate-200 dark:bg-slate-800 text-slate-500 dark:text-slate-400 font-mono font-normal">
                      {file.serviceName}
                    </span>
                  {/if}
                  {#if file.isModified}
                    <span class="w-1.5 h-1.5 rounded-full bg-amber-500" title="Archivo modificado"></span>
                  {/if}
                </button>
              {/each}
            </div>

            <button
              type="button"
              onclick={handleAddDockerfile}
              class="flex items-center gap-1 px-2.5 py-1 text-xs text-slate-500 hover:text-violet-600 dark:hover:text-violet-400 hover:bg-slate-200/60 dark:hover:bg-slate-800/60 rounded-lg transition-colors cursor-pointer flex-shrink-0"
              title="Añadir un Dockerfile al stack"
            >
              <Plus class="w-3 h-3" />
              <span>Dockerfile</span>
            </button>
          </div>

          <!-- Code Editor for the active file -->
          {#if activeFile}
            <div class="h-80">
              <YamlEditor
                bind:value={activeFile.content}
                filename={activeFile.name}
                readOnly={isExecuting}
                oninput={() => (activeFile.isModified = true)}
              />
            </div>
          {/if}
        </div>

        <!-- Real-time Deployment Console -->
        {#if showConsole}
          <div class="pt-2 animate-in fade-in duration-200">
            <DeployConsole
              bind:this={deployConsoleComponent}
              isRunning={isExecuting}
              status={executionStatus}
              statusMessage={executionMessage}
            />
          </div>
        {/if}
      </div>

      <!-- Modal Footer Action Bar -->
      <div class="px-6 py-4 border-t border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50 flex flex-col sm:flex-row items-center justify-between gap-3 flex-shrink-0">
        <div class="flex items-center gap-2">
          {#if mode === 'edit'}
            <button
              type="button"
              onclick={handleDown}
              disabled={isExecuting}
              class="px-3 py-2 rounded-xl text-xs font-medium border border-rose-200 dark:border-rose-900/60 bg-rose-50 dark:bg-rose-950/30 text-rose-600 dark:text-rose-400 hover:bg-rose-100 dark:hover:bg-rose-900/40 transition-colors flex items-center gap-1.5 cursor-pointer disabled:opacity-40"
              title="Detener y desmontar stack con docker compose down"
            >
              <Trash2 class="w-3.5 h-3.5" />
              <span>Destruir Stack (Down)</span>
            </button>
          {/if}
        </div>

        <div class="flex items-center gap-2 w-full sm:w-auto justify-end">
          <button
            type="button"
            onclick={onClose}
            disabled={isExecuting}
            class="px-4 py-2 rounded-xl text-xs font-medium border border-slate-200 dark:border-slate-800 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-300 transition-colors cursor-pointer disabled:opacity-40"
          >
            Cerrar
          </button>

          <button
            type="button"
            onclick={handleSaveFile}
            disabled={isExecuting || !projectName.trim() || !activeFile?.content.trim()}
            class="px-3.5 py-2 rounded-xl text-xs font-medium border border-slate-300 dark:border-slate-700 bg-white dark:bg-slate-800 hover:bg-slate-100 dark:hover:bg-slate-700 text-slate-800 dark:text-slate-200 transition-colors flex items-center gap-1.5 cursor-pointer disabled:opacity-40"
            title="Guardar archivo actual en disco"
          >
            <Save class="w-3.5 h-3.5" />
            <span>Guardar {activeFile ? `(${activeFile.name})` : ''}</span>
          </button>

          <button
            type="button"
            onclick={handleDeployUp}
            disabled={!isValidToDeploy}
            class="px-4 py-2 rounded-xl text-xs font-medium text-white bg-violet-600 hover:bg-violet-700 transition-colors flex items-center gap-1.5 shadow-sm shadow-violet-500/25 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
            title="Guardar archivos y desplegar stack con docker compose up -d"
          >
            <Play class="w-3.5 h-3.5 fill-current" />
            <span>{mode === 'create' ? 'Desplegar Stack (Up -d)' : 'Guardar y Re-desplegar'}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}
