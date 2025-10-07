<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();

	// Auth initialization is handled by route-specific layouts
	// (dashboard)/+layout.ts for protected routes
	// (auth)/+layout.ts for auth routes
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children?.()}
