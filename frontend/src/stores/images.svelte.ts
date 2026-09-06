import type { ImageInfo, DiskUsageSummary, PruneResult } from '../types';
import {
  ListImages,
  GetDiskUsage,
  RemoveImage,
  PruneImages,
} from '../../wailsjs/go/main/App';
import { uiStore } from './ui.svelte';

class ImagesStore {
  images = $state<ImageInfo[]>([]);
  diskUsage = $state<DiskUsageSummary | null>(null);
  loading = $state<boolean>(true);
  actionLoading = $state<string>('');

  async fetchImages() {
    try {
      const [imageList, diskUsageRes] = await Promise.all([
        ListImages(),
        GetDiskUsage(),
      ]);
      this.images = (imageList || []) as unknown as ImageInfo[];
      this.diskUsage = diskUsageRes as unknown as DiskUsageSummary;
    } catch (err: any) {
      uiStore.showToast(`Error al cargar imágenes: ${err}`, 'error');
    } finally {
      this.loading = false;
    }
  }

  async removeImage(id: string, force: boolean) {
    this.actionLoading = id;
    try {
      await RemoveImage(id, force);
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
      const res = await PruneImages(danglingOnly);
      await this.fetchImages();
      uiStore.showToast('Limpieza completada exitosamente');
      return res as unknown as PruneResult;
    } catch (err: any) {
      uiStore.showToast(`Error al ejecutar limpieza: ${err}`, 'error');
      throw err;
    }
  }
}

export const imagesStore = new ImagesStore();
