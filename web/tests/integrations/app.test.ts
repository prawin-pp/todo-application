import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import App from '../../src/App.svelte';
import { auth } from '../../src/stores/auth.store';

describe('App component', () => {
  it('should redirect to login page when user is not logged in', async () => {
    auth.set(null);
    render(App, {});

    const login = await screen.findByText('Login');

    expect(login).not.toBeNull();
  });

  it('should redirect to main page when user is logged in', async () => {
    auth.set({ id: 'MOCK_ID', username: 'MOCK_USERNAME' });
    render(App, {});

    const main = await screen.findByTestId('main');

    expect(main).not.toBeNull();
  });
});
