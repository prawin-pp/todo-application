<script lang="ts">
  import { push } from 'svelte-spa-router';
  import { api } from '../lib/api';
  import { auth } from '../stores/auth.store';

  let username: string;
  let password: string;

  $: isFormValid = username && password;

  export const handleSubmitLogin = async () => {
    if (!isFormValid) return;

    try {
      const user = await api.auth.login(username, password);
      auth.set(user);
      push('/');
    } catch (err) {
      console.error(err);
    }
  };
</script>

<div id="login" class="h-full w-full overflow-auto bg-gray-100">
  <div class="mx-auto my-20 w-[550px] rounded-3xl bg-white px-14 py-12 shadow-lg">
    <form action="#" on:submit|preventDefault={handleSubmitLogin}>
      <div class="mb-6 text-3xl font-bold">Log in to your account</div>
      <div class="mb-6">
        <label for="username" class="mb-2 block font-medium text-gray-900">
          Username<span class="text-red-500">*</span>
        </label>
        <input
          data-testid="input_username"
          type="text"
          id="username"
          class="block w-full rounded-lg border border-gray-200 bg-gray-50 p-2.5 text-gray-900 focus:border-primary focus:ring-primary"
          placeholder="Enter your username"
          required
          autocomplete="username"
          bind:value={username}
        />
      </div>
      <div class="mb-10">
        <label for="password" class="mb-2 block font-medium text-gray-900">
          Password<span class="text-red-500">*</span>
        </label>
        <input
          data-testid="input_password"
          type="password"
          id="current-password"
          class="block w-full rounded-lg border border-gray-200 bg-gray-50 p-2.5 text-gray-900 focus:border-primary focus:ring-primary"
          placeholder="Enter your password"
          required
          autocomplete="current-password"
          bind:value={password}
        />
      </div>
      <button
        type="submit"
        class="w-full rounded-lg bg-primary py-2.5 text-white hover:bg-primary-700"
        class:!bg-gray-500={!isFormValid}
        disabled={!isFormValid}
      >
        Login
      </button>
    </form>
  </div>
</div>
