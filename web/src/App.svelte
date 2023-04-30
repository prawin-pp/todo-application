<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import Router, { push } from 'svelte-spa-router';
  import { wrap } from 'svelte-spa-router/wrap';
  import { api } from './lib/api';
  import { auth } from './stores/auth.store';

  let unsubscribeAuthChanges: () => void;

  const routes = {
    '/login': wrap({
      asyncComponent: () => import('./routes/Login.svelte'),
      conditions: [
        () => {
          if ($auth) {
            push('/');
          }
          return true;
        },
      ],
    }),
    '/*': wrap({
      asyncComponent: () => import('./routes/Main.svelte'),
      conditions: [
        () => {
          if (!$auth) {
            push('/login');
            return false;
          }
          return true;
        },
      ],
    }),
  };

  onMount(async () => {
    unsubscribeAuthChanges = handleAuthChanges();
    await fetchCurrentUser();
  });

  onDestroy(() => {
    unsubscribeAuthChanges && unsubscribeAuthChanges();
  });

  const fetchCurrentUser = async () => {
    try {
      const user = await api.auth.getMe();
      auth.set(user);
    } catch (err) {
      console.error(err);
    }
  };

  const handleAuthChanges = () => {
    return auth.subscribe((user) => {
      if (!user) return push('/login');
      return push('/');
    });
  };
</script>

<Router {routes} />
