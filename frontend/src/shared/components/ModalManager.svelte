<script lang="ts">
  import TerminalModal from '../../features/terminal/modals/TerminalModal.svelte';
  import LogsModal from '../../features/containers/modals/LogsModal.svelte';
  import StatsModal from '../../features/containers/modals/StatsModal.svelte';
  import InspectModal from '../../features/containers/modals/InspectModal.svelte';
  import VolumeInspectModal from '../../features/volumes/modals/VolumeInspectModal.svelte';
  import ConfirmModal from './ConfirmModal.svelte';
  import CreateContainerModal from '../../features/containers/modals/CreateContainerModal.svelte';
  import CreateNetworkModal from '../../features/networks/modals/CreateNetworkModal.svelte';
  import ConnectContainerModal from '../../features/networks/modals/ConnectContainerModal.svelte';
  import NetworkInspectModal from '../../features/networks/modals/NetworkInspectModal.svelte';
  import { uiStore } from '../stores/ui.svelte';
  import { containersStore } from '../../features/containers/stores/containers.svelte';
  import { imagesStore } from '../../features/images/stores/images.svelte';
  import { networksStore } from '../../features/networks/stores/networks.svelte';

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

<!-- Volume Inspect Modal -->
<VolumeInspectModal
  volumeName={uiStore.activeVolumeInspect?.name || null}
  onClose={() => uiStore.closeVolumeInspect()}
/>

<!-- Create Network Modal -->
<CreateNetworkModal
  isOpen={uiStore.isCreateNetworkModalOpen}
  onClose={() => uiStore.closeCreateNetworkModal()}
/>

<!-- Network Inspect Modal -->
<NetworkInspectModal
  networkId={uiStore.activeNetworkInspect?.id || null}
  networkName={uiStore.activeNetworkInspect?.name || ''}
  onClose={() => uiStore.closeNetworkInspect()}
/>

<!-- Connect Container to Network Modal -->
{#if uiStore.activeConnectContainer}
  {@const targetNet = networksStore.networks.find(n => n.id === uiStore.activeConnectContainer?.networkId)}
  <ConnectContainerModal
    networkId={uiStore.activeConnectContainer.networkId}
    networkName={uiStore.activeConnectContainer.networkName}
    connectedContainerIds={targetNet?.containers.map(c => c.id) || []}
    onClose={() => uiStore.closeConnectContainer()}
  />
{/if}

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
