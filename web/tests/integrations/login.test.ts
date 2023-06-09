import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import { login } from '../../src/lib/api/auth.api';
import Login from '../../src/routes/Login.svelte';

vi.mock('../../src/lib/api/auth.api', () => ({ login: vi.fn() }));

describe('Login component', () => {
  describe('login', () => {
    it('should not call login api when username or password is empty', async () => {
      render(Login, {});
      const btnLogin = screen.getByText('Login');

      await fireEvent.click(btnLogin);

      expect(login).not.toHaveBeenCalled();
    });
  });
});
