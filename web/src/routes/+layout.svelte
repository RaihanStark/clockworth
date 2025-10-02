<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import NavSidebar from '$lib/components/NavSidebar.svelte';
	import { Menu, Clock } from '@lucide/svelte';

	let { children } = $props();
	let mobileMenuOpen = $state(false);

	function toggleMobileMenu() {
		mobileMenuOpen = !mobileMenuOpen;
	}

	function closeMobileMenu() {
		mobileMenuOpen = false;
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="flex h-screen flex-col overflow-hidden lg:flex-row">
	<!-- Mobile Top Bar -->
	<header class="flex items-center justify-between border-b border-slate-200 bg-white px-4 py-3 lg:hidden">
		<button
			onclick={toggleMobileMenu}
			class="rounded-lg p-2 text-slate-900 transition-colors hover:bg-slate-100"
			aria-label="Toggle menu"
		>
			<Menu size={24} />
		</button>

		<div class="flex items-center gap-2">
			<div class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-emerald-400 to-emerald-600">
				<Clock size={16} class="text-white" />
			</div>
			<h1 class="text-lg font-black text-slate-900">Clockworth</h1>
		</div>

		<div class="w-10"></div> <!-- Spacer for centering -->
	</header>

	<!-- Overlay -->
	{#if mobileMenuOpen}
		<div
			class="fixed inset-0 z-30 bg-black/50 backdrop-blur-sm lg:hidden"
			onclick={closeMobileMenu}
		></div>
	{/if}

	<!-- Sidebar -->
	<div
		class="fixed inset-y-0 left-0 z-40 transform transition-transform duration-300 lg:relative lg:translate-x-0 {mobileMenuOpen
			? 'translate-x-0'
			: '-translate-x-full'}"
	>
		<NavSidebar onNavigate={closeMobileMenu} />
	</div>

	<!-- Main Content -->
	<main class="flex-1 overflow-auto">
		{@render children?.()}
	</main>
</div>
