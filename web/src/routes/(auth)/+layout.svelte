<script lang="ts">
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth';

	// Redirect to dashboard if already logged in
	$effect(() => {
		if ($authStore.initialized && $authStore.user) {
			goto('/dashboard');
		}
	});
</script>

<slot />
