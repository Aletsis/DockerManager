import type { ImageInfo, DiskUsageSummary, PruneResult } from '../types';
import { dockerApi } from '../services/api';
import { uiStore } from './ui.svelte';

class ImagesStore {
  images = $state<ImageInfo[]>([]);
  diskUsage = $state<DiskUsageSummary | null>(null);
  loading = $state<boolean>(true);
  actionLoading = $state<string>('');

  async fetchImages() {
    try {
      const [imageList, diskUsageRes] = await Promise.all([
        dockerApi.listImages(),
        dockerApi.getDiskUsage(),
      ]);
      this.images = imageList;
      this.diskUsage = diskUsageRes;
    } catch (err: any) {
      uiStore.showToast(`Error al cargar imágenes: ${err}`, 'error');
    } finally {
      this.loading = false;
    }
  }

  async removeImage(id: string, force: boolean) {
    this.actionLoading = id;
    try {
      await dockerApi.removeImage(id, force);
      uiStore.showToast('Imagen eliminada correctamente');
      await this.fetchImages();
    } catch (err: any) {
      uiStore.showToast(`Error al eliminar imagen: ${err}`, 'error');
    } finally {
      this.actionLoading = '';
    }
  }

  async prune(danglingOnly: boolean): Promise<PruneResult | null> {
    try {
      const res = await dockerApi.pruneImages(danglingOnly);
      await this.fetchImages();
      uiStore.showToast('Limpieza completada exitosamente');
      return res;
    } catch (err: any) {
      uiStore.showToast(`Error al ejecutar limpieza: ${err}`, 'error');
      throw err;
    }
  }
}

export const imagesStore = new ImagesStore();
