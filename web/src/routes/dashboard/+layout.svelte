<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import NavSidebar from '$lib/components/NavSidebar.svelte';
	import { Menu, Clock } from '@lucide/svelte';
	import { authStore } from '$lib/stores/auth';

	let mobileMenuOpen = $state(false);

	function toggleMobileMenu() {
		mobileMenuOpen = !mobileMenuOpen;
	}

	function closeMobileMenu() {
		mobileMenuOpen = false;
	}

	// Protect dashboard routes
	$effect(() => {
		if ($authStore.initialized && !$authStore.user) {
			goto('/login');
		}
	});
</script>

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
