import { computed, signal } from '@angular/core';

export type PageActionFeedbackKind = 'idle' | 'loading' | 'success' | 'error';

export interface PageActionFeedbackState {
  kind: PageActionFeedbackKind;
  message: string | null;
}

export function createPageActionFeedback() {
  const state = signal<PageActionFeedbackState>({ kind: 'idle', message: null });

  return {
    state,
    isBusy: computed(() => state().kind === 'loading'),
    clear(): void {
      state.set({ kind: 'idle', message: null });
    },
    loading(message: string): void {
      state.set({ kind: 'loading', message });
    },
    success(message: string): void {
      state.set({ kind: 'success', message });
    },
    error(message: string): void {
      state.set({ kind: 'error', message });
    },
  };
}
