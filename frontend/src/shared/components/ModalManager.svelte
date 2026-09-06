<script lang="ts">
  import TerminalModal from '../../features/terminal/modals/TerminalModal.svelte';
  import LogsModal from '../../features/containers/modals/LogsModal.svelte';
  import StatsModal from '../../features/containers/modals/StatsModal.svelte';
  import InspectModal from '../../features/containers/modals/InspectModal.svelte';
  import ConfirmModal from './ConfirmModal.svelte';
  import CreateContainerModal from '../../features/containers/modals/CreateContainerModal.svelte';
  import { uiStore } from '../stores/ui.svelte';
  import { containersStore } from '../../features/containers/stores/containers.svelte';
  import { imagesStore } from '../../features/images/stores/images.svelte';

  async function handleConfirmDelete() {
    if (!uiStore.confirmDelete) return;
    const { id } = uiStore.confirmDelete;
    uiStore.closeConfirmDelete();
    await containersStore.handleRemove(id);
  }
</script>

<!-- Create Container Modal -->
<CreateContainerModal
  isOpen={uiStore.isCreateModalOpen}
  initialImage={uiStore.createModalInitialImage}
  localImages={imagesStore.images}
  onClose={() => uiStore.closeCreateModal()}
  onSuccess={() => {
    uiStore.showToast('Contenedor creado e iniciado exitosamente');
    containersStore.fetchData(true);
  }}
/>

<!-- Terminal Modal -->
{#if uiStore.activeTerminal}
  <TerminalModal
    containerId={uiStore.activeTerminal.id}
    containerName={uiStore.activeTerminal.name}
    onClose={() => uiStore.closeTerminal()}
  />
{/if}

<!-- Logs Modal -->
<LogsModal
  containerId={uiStore.activeLogs?.id || null}
  containerName={uiStore.activeLogs?.name || ''}
  onClose={() => uiStore.closeLogs()}
/>

<!-- Stats Modal -->
<StatsModal
  containerId={uiStore.activeStats?.id || null}
  containerName={uiStore.activeStats?.name || ''}
  onClose={() => uiStore.closeStats()}
/>

<!-- Inspect Modal -->
<InspectModal
  containerId={uiStore.activeInspect?.id || null}
  containerName={uiStore.activeInspect?.name || ''}
  onClose={() => uiStore.closeInspect()}
/>

<!-- Confirm Delete Container Modal -->
<ConfirmModal
  isOpen={!!uiStore.confirmDelete}
  title="Eliminar Contenedor"
  message={`¿Estás seguro de que deseas eliminar permanentemente el contenedor "${uiStore.confirmDelete?.name}"? Esta acción no se puede deshacer.`}
  confirmLabel="Eliminar Contenedor"
  onConfirm={handleConfirmDelete}
  onCancel={() => uiStore.closeConfirmDelete()}
  isDestructive={true}
/>
