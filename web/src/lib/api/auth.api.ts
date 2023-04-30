import type { User } from '../../models/user.model';
import { fetchWithCredential, throwErrorWhenResponseStautsIsNotOk } from './api';

export const getMe = async (): Promise<User> => {
  const resp = await fetchWithCredential('/me', { method: 'GET' });
  throwErrorWhenResponseStautsIsNotOk(resp, 'failed to get me');
  const body: User = await resp.json();
  return body;
};

export const login = async (username: string, password: string) => {
  const resp = await fetchWithCredential('/login', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  });
  throwErrorWhenResponseStautsIsNotOk(resp, 'failed to login');
  const body: User = await resp.json();
  return body;
};

export const logout = async () => {
  const resp = await fetchWithCredential('/logout', { method: 'POST' });
  throwErrorWhenResponseStautsIsNotOk(resp, 'failed to logout');
};
