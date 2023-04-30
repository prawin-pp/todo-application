import { auth } from '../../stores/auth.store';
import { config } from '../config';

export const fetchWithCredential = async (path: string, init?: RequestInit): Promise<Response> => {
  const url = `${config.api.baseUrl}${path}`;
  const options: RequestInit = { ...(init || {}), credentials: 'include' };
  return fetch(url, options);
};

export const throwErrorWhenResponseStautsIsNotOk = (response: Response, error: string): void => {
  if (response.status === 401) {
    auth.clear();
    throw new Error('session timeout please login again');
  }
  if (response.status < 200 || response.status > 299) {
    throw new Error(error);
  }
};
