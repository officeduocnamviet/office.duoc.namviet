import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface GlobalState {
  theme: 'light' | 'dark';
  sidebarCollapsed: boolean;
  activeBranchId: string | null;
  setTheme: (theme: 'light' | 'dark') => void;
  toggleSidebar: () => void;
  setActiveBranch: (branchId: string) => void;
}

export const useGlobalStore = create<GlobalState>()(
  persist(
    (set) => ({
      theme: 'light',
      sidebarCollapsed: false,
      activeBranchId: null,
      setTheme: (theme) => set({ theme }),
      toggleSidebar: () => set((state) => ({ sidebarCollapsed: !state.sidebarCollapsed })),
      setActiveBranch: (branchId) => set({ activeBranchId: branchId }),
    }),
    {
      name: 'namviet-erp-global-storage',
    }
  )
);
