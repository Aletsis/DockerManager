export function formatBytes(bytes: number, decimals = 1): string {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
}

export function formatUptime(status: string): string {
  return status.replace(/Up\s+/, '').replace(/\(healthy\)/, '• Healthy');
}

export function getStateColor(state: string): {
  badgeBg: string;
  badgeText: string;
  dotBg: string;
  glowClass: string;
} {
  switch (state.toLowerCase()) {
    case 'running':
      return {
        badgeBg: 'bg-emerald-500/10 dark:bg-emerald-400/15 border-emerald-500/30 dark:border-emerald-400/30',
        badgeText: 'text-emerald-700 dark:text-emerald-400',
        dotBg: 'bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.7)]',
        glowClass: 'animate-pulse',
      };
    case 'paused':
      return {
        badgeBg: 'bg-amber-500/10 dark:bg-amber-400/15 border-amber-500/30 dark:border-amber-400/30',
        badgeText: 'text-amber-700 dark:text-amber-400',
        dotBg: 'bg-amber-500',
        glowClass: '',
      };
    case 'exited':
    case 'dead':
    default:
      return {
        badgeBg: 'bg-slate-500/10 dark:bg-slate-400/10 border-slate-300 dark:border-slate-700',
        badgeText: 'text-slate-600 dark:text-slate-400',
        dotBg: 'bg-slate-400 dark:bg-slate-600',
        glowClass: '',
      };
  }
}

export function formatRelativeTime(unixSeconds: number): string {
  if (!unixSeconds) return 'Desconocido';
  const now = Math.floor(Date.now() / 1000);
  const diff = Math.max(0, now - unixSeconds);

  if (diff < 60) return 'Hace unos segundos';
  const minutes = Math.floor(diff / 60);
  if (minutes < 60) return `Hace ${minutes} min`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `Hace ${hours} h`;
  const days = Math.floor(hours / 24);
  if (days < 30) return `Hace ${days} d`;
  const months = Math.floor(days / 30);
  if (months < 12) return `Hace ${months} m`;
  const years = Math.floor(months / 12);
  return `Hace ${years} a`;
}

