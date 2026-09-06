<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import Header from './components/layout/Header.svelte';
  import ContainersView from './views/ContainersView.svelte';
  import ImagesView from './views/ImagesView.svelte';
  import ModalManager from './components/modals/ModalManager.svelte';
  import { uiStore } from './stores/ui.svelte';
  import { containersStore } from './stores/containers.svelte';
  import { imagesStore } from './stores/images.svelte';

  onMount(() => {
    containersStore.fetchData();
    containersStore.setupPolling();
    imagesStore.fetchImages();
  });

  onDestroy(() => {
    containersStore.stopPolling();
  });
</script>

<div class="h-screen bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 flex flex-col transition-colors duration-200 overflow-hidden">
  <!-- Top Application Header -->
  <Header
    overview={containersStore.overview}
    diskUsage={imagesStore.diskUsage}
    bind:activeTab={uiStore.activeTab}
    bind:searchQuery={uiStore.searchQuery}
    isDark={uiStore.isDark}
    onToggleTheme={() => uiStore.toggleTheme()}
    bind:refreshInterval={containersStore.refreshInterval}
    onManualRefresh={() => containersStore.fetchData(true)}
    isRefreshing={containersStore.isRefreshing}
  />

  <!-- Main View Area -->
  <main class="flex-1 overflow-y-auto px-4 sm:px-6 py-6">
    {#if uiStore.activeTab === 'containers'}
      <ContainersView />
    {:else}
      <ImagesView />
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
