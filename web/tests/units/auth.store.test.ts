import { get } from 'svelte/store';
import { describe, expect, it } from 'vitest';
import { auth } from '../../src/stores/auth.store';

describe('Auth Store', () => {
  describe('clear', () => {
    it('should set user to null', () => {
      auth.set({ id: '1', username: 'Test User' });

      auth.clear();

      expect(get(auth)).toBeNull();
    });
  });
});
