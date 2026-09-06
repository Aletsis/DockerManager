<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import Header from './features/system/components/Header.svelte';
  import ContainersView from './views/ContainersView.svelte';
  import ImagesView from './features/images/views/ImagesView.svelte';
  import VolumesView from './features/volumes/views/VolumesView.svelte';
  import ModalManager from './shared/components/ModalManager.svelte';
  import { uiStore } from './shared/stores/ui.svelte';
  import { containersStore } from './features/containers/stores/containers.svelte';
  import { imagesStore } from './features/images/stores/images.svelte';
  import { volumesStore } from './features/volumes/stores/volumes.svelte';

  onMount(() => {
    containersStore.fetchData();
    containersStore.setupPolling();
    imagesStore.fetchImages();
    volumesStore.fetchVolumes();
  });

  onDestroy(() => {
    containersStore.stopPolling();
  });

  function handleManualRefresh() {
    containersStore.fetchData(true);
    imagesStore.fetchImages();
    volumesStore.fetchVolumes();
  }
</script>

<div class="h-screen bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 flex flex-col transition-colors duration-200 overflow-hidden">
  <!-- Top Application Header -->
  <Header
    overview={containersStore.overview}
    diskUsage={imagesStore.diskUsage}
    volumeDiskUsage={volumesStore.diskUsage}
    volumesCount={volumesStore.volumes.length}
    bind:activeTab={uiStore.activeTab}
    bind:searchQuery={uiStore.searchQuery}
    isDark={uiStore.isDark}
    onToggleTheme={() => uiStore.toggleTheme()}
    bind:refreshInterval={containersStore.refreshInterval}
    onManualRefresh={handleManualRefresh}
    isRefreshing={containersStore.isRefreshing}
  />

  <!-- Main View Area -->
  <main class="flex-1 overflow-y-auto px-4 sm:px-6 py-6">
    {#if uiStore.activeTab === 'containers'}
      <ContainersView />
    {:else if uiStore.activeTab === 'images'}
      <ImagesView />
    {:else if uiStore.activeTab === 'volumes'}
      <VolumesView />
    {/if}
  </main>

  <!-- Centralized Modals Container -->
  <ModalManager />

  <!-- Floating Toast Notification -->
  {#if uiStore.notification}
    <div
      class="fixed bottom-5 right-5 z-50 px-4 py-2.5 rounded-xl shadow-lg border text-xs font-medium animate-in fade-in slide-in-from-bottom-2 duration-200 {uiStore.notification.type === 'error' ? 'bg-rose-600 text-white border-rose-700' : 'bg-slate-900 text-white dark:bg-slate-100 dark:text-slate-900 border-slate-800 dark:border-slate-200'}"
    >
      {uiStore.notification.text}
    </div>
  {/if}
</div>
