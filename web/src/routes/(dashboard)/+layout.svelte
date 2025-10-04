<script lang="ts">
	import { page } from '$app/stores';
	import { navigating } from '$app/stores';
	import NavSidebar from '$lib/components/NavSidebar.svelte';
	import { Menu, Clock } from '@lucide/svelte';

	let mobileMenuOpen = $state(false);

	function toggleMobileMenu() {
		mobileMenuOpen = !mobileMenuOpen;
	}

	function closeMobileMenu() {
		mobileMenuOpen = false;
	}
</script>

{#if $navigating || !$page.data.user}
	<!-- Loading State -->
	<div class="flex h-screen items-center justify-center bg-gradient-to-br from-slate-50 to-slate-100">
		<div class="flex flex-col items-center gap-4">
			<div class="h-12 w-12 animate-spin rounded-full border-4 border-slate-200 border-t-emerald-600"></div>
			<p class="text-sm font-semibold text-slate-600">Loading...</p>
		</div>
	</div>
{:else}
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
		<div class="w-10"></div>
	</header>

	<!-- Desktop Sidebar -->
	<div class="hidden lg:block">
		<NavSidebar />
	</div>

	<!-- Mobile Sidebar -->
	{#if mobileMenuOpen}
		<button
			onclick={closeMobileMenu}
			class="fixed inset-0 z-40 bg-black/50 lg:hidden"
			aria-label="Close menu"
		></button>
		<div class="fixed inset-y-0 left-0 z-50 w-72 lg:hidden">
			<NavSidebar onNavigate={closeMobileMenu} />
		</div>
	{/if}

	<!-- Main Content -->
	<main class="flex-1 overflow-hidden">
		<slot />
	</main>
</div>
{/if}
