import { fireEvent, render, screen, waitFor } from '@testing-library/svelte';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { logout } from '../../src/lib/api/auth.api';
import { User } from '../../src/models/user.model';
import Main from '../../src/routes/Main.svelte';
import { auth } from '../../src/stores/auth.store';

vi.mock('../../src/lib/api/auth.api', () => ({
  logout: vi.fn().mockResolvedValue(undefined),
}));

describe('Main component', () => {
  describe('logout', () => {
    let currentAuth: User | null;
    let unsubscribe: () => void;

    beforeEach(() => {
      unsubscribe = auth.subscribe((value) => (currentAuth = value));
      auth.set({ id: 'TEST_ID', username: 'TEST_USERNAME' });
    });

    afterEach(() => {
      unsubscribe();
    });

    it('should call logout api and set auth store to null', async () => {
      render(Main, { params: {} });

      const btnLogout = screen.getByTestId('btn-logout');
      await fireEvent.click(btnLogout);

      await waitFor(() => expect(logout).toBeCalledTimes(1));
      expect(currentAuth).toEqual(null);
    });
  });
});
