export type TabType = 'containers' | 'images' | 'volumes';

export interface ToastNotification {
  text: string;
  type: 'success' | 'error';
}

class UiStore {
  isDark = $state<boolean>((() => {
    const saved = localStorage.getItem('theme');
    if (saved) return saved === 'dark';
    return typeof window !== 'undefined' ? window.matchMedia('(prefers-color-scheme: dark)').matches : true;
  })());

  activeTab = $state<TabType>('containers');
  searchQuery = $state<string>('');
  notification = $state<ToastNotification | null>(null);

  // Modals
  activeTerminal = $state<{ id: string; name: string } | null>(null);
  activeLogs = $state<{ id: string; name: string } | null>(null);
  activeStats = $state<{ id: string; name: string } | null>(null);
  activeInspect = $state<{ id: string; name: string } | null>(null);
  activeVolumeInspect = $state<{ name: string } | null>(null);
  confirmDelete = $state<{ id: string; name: string } | null>(null);
  isCreateModalOpen = $state<boolean>(false);
  createModalInitialImage = $state<string>('');

  private toastTimeout: any = null;

  constructor() {
    // Sync theme with DOM initially
    if (typeof document !== 'undefined') {
      if (this.isDark) {
        document.documentElement.classList.add('dark');
      } else {
        document.documentElement.classList.remove('dark');
      }
    }
  }

  toggleTheme() {
    this.isDark = !this.isDark;
    if (typeof document !== 'undefined') {
      if (this.isDark) {
        document.documentElement.classList.add('dark');
        localStorage.setItem('theme', 'dark');
      } else {
        document.documentElement.classList.remove('dark');
        localStorage.setItem('theme', 'light');
      }
    }
  }

  showToast(text: string, type: 'success' | 'error' = 'success') {
    if (this.toastTimeout) {
      clearTimeout(this.toastTimeout);
    }
    this.notification = { text, type };
    this.toastTimeout = setTimeout(() => {
      this.notification = null;
    }, 3000);
  }

  openTerminal(id: string, name: string) {
    this.activeTerminal = { id, name };
  }

  closeTerminal() {
    this.activeTerminal = null;
  }

  openLogs(id: string, name: string) {
    this.activeLogs = { id, name };
  }

  closeLogs() {
    this.activeLogs = null;
  }

  openStats(id: string, name: string) {
    this.activeStats = { id, name };
  }

  closeStats() {
    this.activeStats = null;
  }

  openInspect(id: string, name: string) {
    this.activeInspect = { id, name };
  }

  closeInspect() {
    this.activeInspect = null;
  }

  openVolumeInspect(name: string) {
    this.activeVolumeInspect = { name };
  }

  closeVolumeInspect() {
    this.activeVolumeInspect = null;
  }

  openConfirmDelete(id: string, name: string) {
    this.confirmDelete = { id, name };
  }

  closeConfirmDelete() {
    this.confirmDelete = null;
  }

  openCreateModal(initialImage = '') {
    this.createModalInitialImage = initialImage;
    this.isCreateModalOpen = true;
  }

  closeCreateModal() {
    this.isCreateModalOpen = false;
    this.createModalInitialImage = '';
  }
}

export const uiStore = new UiStore();
