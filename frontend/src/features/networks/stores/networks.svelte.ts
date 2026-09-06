import type { NetworkInfo, CreateNetworkRequest, NetworkPruneResult } from '../../../types';
import { dockerApi } from '../../../shared/services/api';
import { uiStore } from '../../../shared/stores/ui.svelte';

class NetworksStore {
  networks = $state<NetworkInfo[]>([]);
  loading = $state<boolean>(true);
  actionLoading = $state<string>('');

  async fetchNetworks() {
    try {
      const list = await dockerApi.listNetworks();
      this.networks = list;
    } catch (err: any) {
      uiStore.showToast(`Error al cargar redes: ${err}`, 'error');
    } finally {
      this.loading = false;
    }
  }

  async createNetwork(req: CreateNetworkRequest): Promise<boolean> {
    this.actionLoading = 'create';
    try {
      await dockerApi.createNetwork(req);
      uiStore.showToast(`Red "${req.name}" creada exitosamente`);
      await this.fetchNetworks();
      return true;
    } catch (err: any) {
      uiStore.showToast(`Error al crear red: ${err}`, 'error');
      return false;
    } finally {
      this.actionLoading = '';
    }
  }

  async removeNetwork(id: string, name: string): Promise<boolean> {
    this.actionLoading = id;
    try {
      await dockerApi.removeNetwork(id);
      uiStore.showToast(`Red "${name}" eliminada correctamente`);
      await this.fetchNetworks();
      return true;
    } catch (err: any) {
      uiStore.showToast(`Error al eliminar red: ${err}`, 'error');
      return false;
    } finally {
      this.actionLoading = '';
    }
  }

  async prune(): Promise<NetworkPruneResult | null> {
    try {
      const res = await dockerApi.pruneNetworks();
      await this.fetchNetworks();
      uiStore.showToast('Limpieza de redes inactivas completada');
      return res;
    } catch (err: any) {
      uiStore.showToast(`Error al limpiar redes: ${err}`, 'error');
      throw err;
    }
  }

  async connectContainer(networkId: string, containerId: string, ipAddress = ''): Promise<boolean> {
    this.actionLoading = `${networkId}:${containerId}`;
    try {
      await dockerApi.connectContainerToNetwork(networkId, containerId, ipAddress);
      uiStore.showToast('Contenedor conectado a la red en caliente');
      await this.fetchNetworks();
      return true;
    } catch (err: any) {
      uiStore.showToast(`Error al conectar contenedor: ${err}`, 'error');
      return false;
    } finally {
      this.actionLoading = '';
    }
  }

  async disconnectContainer(networkId: string, containerId: string, force = false): Promise<boolean> {
    this.actionLoading = `${networkId}:${containerId}`;
    try {
      await dockerApi.disconnectContainerFromNetwork(networkId, containerId, force);
      uiStore.showToast('Contenedor desconectado de la red');
      await this.fetchNetworks();
      return true;
    } catch (err: any) {
      uiStore.showToast(`Error al desconectar contenedor: ${err}`, 'error');
      return false;
    } finally {
      this.actionLoading = '';
    }
  }

  async inspectNetwork(id: string): Promise<string> {
    return dockerApi.inspectNetwork(id);
  }
}

export const networksStore = new NetworksStore();
