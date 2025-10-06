<script lang="ts">
	import { onMount } from 'svelte';
	import { authStore } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { Clock } from '@lucide/svelte';

	let error = '';
	let loading = true;

	onMount(() => {
		// Cookie is set by backend, verify session by fetching user
		(async () => {
			await authStore.init();
		})();

		// Check if we successfully authenticated
		const unsubscribe = authStore.subscribe((state) => {
			if (state.initialized) {
				if (state.user) {
					goto('/dashboard');
				} else {
					error = 'Failed to authenticate. Redirecting to login...';
					loading = false;
					setTimeout(() => goto('/login'), 2000);
				}
			}
		});

		return unsubscribe;
	});
</script>

<div
	class="flex min-h-screen items-center justify-center bg-gradient-to-br from-slate-50 to-slate-100 p-4"
>
	<div class="w-full max-w-md text-center">
		<!-- Logo -->
		<div class="mb-8">
			<div
				class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-gradient-to-br from-emerald-400 to-emerald-600 shadow-lg shadow-emerald-500/50"
			>
				<Clock size={32} class="text-white" />
			</div>
			<h1 class="text-3xl font-black text-slate-900">Clockworth</h1>
			<p class="mt-1 text-sm font-semibold text-emerald-600">Time Economics</p>
		</div>

		<!-- Content Card -->
		<div class="rounded-2xl bg-white p-8 shadow-xl">
			{#if loading && !error}
				<div class="flex flex-col items-center gap-4">
					<div
						class="h-12 w-12 animate-spin rounded-full border-4 border-slate-200 border-t-emerald-600"
					></div>
					<p class="text-lg font-semibold text-slate-900">Authenticating...</p>
					<p class="text-sm text-slate-600">Please wait while we sign you in</p>
				</div>
			{:else if error}
				<div class="flex flex-col items-center gap-4">
					<div class="flex h-12 w-12 items-center justify-center rounded-full bg-red-100">
						<svg class="h-6 w-6 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M6 18L18 6M6 6l12 12"
							/>
						</svg>
					</div>
					<p class="text-lg font-semibold text-red-900">{error}</p>
				</div>
			{/if}
		</div>
	</div>
</div>
