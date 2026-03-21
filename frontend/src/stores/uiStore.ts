import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface UIState {
  theme: 'light' | 'dark';
  sidebarOpen: boolean;
  notificationPreferences: {
    whatsapp: boolean;
    sms: boolean;
    email: boolean;
  };
  toggleTheme: () => void;
  toggleSidebar: () => void;
  updateNotificationPreferences: (preferences: Partial<UIState['notificationPreferences']>) => void;
}

export const useUIStore = create<UIState>()(
  persist(
    (set) => ({
      theme: 'dark',
      sidebarOpen: true,
      notificationPreferences: {
        whatsapp: true,
        sms: false,
        email: true,
      },
      toggleTheme: () => set((state) => ({ theme: state.theme === 'light' ? 'dark' : 'light' })),
      toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
      updateNotificationPreferences: (preferences) =>
        set((state) => ({
          notificationPreferences: { ...state.notificationPreferences, ...preferences },
        })),
    }),
    {
      name: 'ui-storage',
    }
  )
);
