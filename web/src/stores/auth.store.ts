import { writable } from 'svelte/store';
import type { User } from '../models/user.model';

const createAuth = () => {
  const { set, subscribe, update } = writable<User | null>(null);
  return {
    set,
    subscribe,
    update,
    clear: () => set(null),
  };
};

export const auth = createAuth();
