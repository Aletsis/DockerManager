import type { VolumeInfo, VolumeDiskUsageSummary, VolumePruneResult } from '../../../types';
import { dockerApi } from '../../../shared/services/api';
import { uiStore } from '../../../shared/stores/ui.svelte';

class VolumesStore {
  volumes = $state<VolumeInfo[]>([]);
  diskUsage = $state<VolumeDiskUsageSummary | null>(null);
  loading = $state<boolean>(true);
  actionLoading = $state<string>('');

  async fetchVolumes() {
    try {
      const [volumeList, diskUsageRes] = await Promise.all([
        dockerApi.listVolumes(),
        dockerApi.getVolumeDiskUsage(),
      ]);
      this.volumes = volumeList;
      this.diskUsage = diskUsageRes;
    } catch (err: any) {
      uiStore.showToast(`Error al cargar volúmenes: ${err}`, 'error');
    } finally {
      this.loading = false;
    }
  }

  async removeVolume(name: string, force = false): Promise<boolean> {
    this.actionLoading = name;
    try {
      await dockerApi.removeVolume(name, force);
      uiStore.showToast('Volumen eliminado correctamente');
      await this.fetchVolumes();
      return true;
    } catch (err: any) {
      uiStore.showToast(`Error al eliminar volumen: ${err}`, 'error');
      return false;
    } finally {
      this.actionLoading = '';
    }
  }

  async prune(): Promise<VolumePruneResult | null> {
    try {
      const res = await dockerApi.pruneVolumes();
      await this.fetchVolumes();
      uiStore.showToast('Limpieza de volúmenes completada');
      return res;
    } catch (err: any) {
      uiStore.showToast(`Error al limpiar volúmenes: ${err}`, 'error');
      throw err;
    }
  }

  async inspectVolume(name: string): Promise<string> {
    return dockerApi.inspectVolume(name);
  }
}

export const volumesStore = new VolumesStore();
