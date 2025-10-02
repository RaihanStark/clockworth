<script lang="ts">
	import {
		Calendar,
		Settings,
		Clock,
		ChevronRight,
		LogOut,
		ChevronDown,
		User,
		HelpCircle
	} from '@lucide/svelte';
	import { page } from '$app/stores';

	interface NavItem {
		label: string;
		icon: any;
		href: string;
		badge?: string;
	}

	interface Props {
		onNavigate?: () => void;
	}

	let { onNavigate }: Props = $props();
	let dropdownOpen = $state(false);

	const navItems: NavItem[] = [
		{ label: 'Today', icon: Calendar, href: '/' },
		{ label: 'Settings', icon: Settings, href: '/settings' }
	];

	function handleNavClick() {
		if (onNavigate) {
			onNavigate();
		}
	}

	function isActive(href: string): boolean {
		if (href === '/') {
			return $page.url.pathname === '/';
		}
		return $page.url.pathname.startsWith(href);
	}

	function toggleDropdown() {
		dropdownOpen = !dropdownOpen;
	}

	function handleLogout() {
		// Add logout logic here
		console.log('Logout clicked');
		dropdownOpen = false;
	}

	function closeDropdown() {
		dropdownOpen = false;
	}
</script>

<aside class="flex h-screen w-72 flex-col bg-gradient-to-b from-slate-950 to-slate-900 text-white">
	<!-- Logo/Brand -->
	<div class="relative overflow-hidden px-6 py-6">
		<div class="relative z-10">
			<div class="mb-2 flex items-center gap-3">
				<div
					class="flex h-12 w-12 items-center justify-center rounded-2xl bg-gradient-to-br from-emerald-400 to-emerald-600 shadow-lg shadow-emerald-500/50"
				>
					<Clock class="text-white" size={24} />
				</div>
				<div>
					<h1 class="text-2xl font-black tracking-tight">Clockworth</h1>
					<p class="text-xs font-medium text-emerald-400">Time Economics</p>
				</div>
			</div>
		</div>
		<!-- Decorative glow -->
		<div class="absolute -right-10 -top-10 h-40 w-40 rounded-full bg-emerald-500/10 blur-3xl"></div>
	</div>

	<!-- Navigation -->
	<nav class="flex-1 space-y-2 px-4 py-2">
		{#each navItems as item}
			{@const active = isActive(item.href)}
			<a
				href={item.href}
				onclick={handleNavClick}
				class="group relative flex items-center gap-4 rounded-xl px-4 py-3.5 text-base font-semibold transition-all {active
					? 'bg-gradient-to-r from-emerald-500 to-emerald-600 text-white shadow-lg shadow-emerald-500/30'
					: 'text-slate-400 hover:bg-white/5 hover:text-white'}"
			>
				<div
					class="flex h-10 w-10 items-center justify-center rounded-lg {active
						? 'bg-white/20'
						: 'bg-white/5'}"
				>
					<svelte:component this={item.icon} size={22} />
				</div>
				<span class="flex-1">{item.label}</span>
				{#if item.badge}
					<span
						class="rounded-full bg-emerald-400/20 px-2 py-0.5 text-xs font-bold text-emerald-400"
					>
						{item.badge}
					</span>
				{/if}
				{#if active}
					<ChevronRight size={18} class="opacity-70" />
				{/if}
			</a>
		{/each}
	</nav>

	<!-- Footer -->
	<div class="relative border-t border-white/10 px-4 py-4">
		<!-- Dropdown Menu -->
		{#if dropdownOpen}
			<!-- Backdrop overlay -->
			<button onclick={closeDropdown} class="fixed inset-0 z-10" aria-label="Close menu"></button>

			<!-- Menu -->
			<div
				class="to-slate-850 absolute bottom-full left-4 right-4 z-20 overflow-hidden rounded-2xl bg-gradient-to-b from-slate-800 shadow-2xl ring-1 ring-white/10 backdrop-blur-sm"
			>
				<!-- User Info Header -->
				<div class="border-b border-white/10 px-4 py-3">
					<p class="text-sm font-bold text-white">User</p>
					<p class="text-xs text-slate-400">user@email.com</p>
				</div>

				<!-- Menu Items -->
				<div class="py-2">
					<a
						href="/profile"
						onclick={closeDropdown}
						class="flex items-center gap-3 px-4 py-2.5 text-sm font-semibold text-slate-300 transition-colors hover:bg-white/5 hover:text-white"
					>
						<User size={18} />
						<span>Profile</span>
					</a>

					<a
						href="/settings"
						onclick={closeDropdown}
						class="flex items-center gap-3 px-4 py-2.5 text-sm font-semibold text-slate-300 transition-colors hover:bg-white/5 hover:text-white"
					>
						<Settings size={18} />
						<span>Settings</span>
					</a>

					<a
						href="/help"
						onclick={closeDropdown}
						class="flex items-center gap-3 px-4 py-2.5 text-sm font-semibold text-slate-300 transition-colors hover:bg-white/5 hover:text-white"
					>
						<HelpCircle size={18} />
						<span>Help & Support</span>
					</a>
				</div>

				<!-- Logout Section -->
				<div class="border-t border-white/10 py-2">
					<button
						onclick={handleLogout}
						class="flex w-full items-center gap-3 px-4 py-2.5 text-sm font-semibold text-red-400 transition-colors hover:bg-red-500/10 hover:text-red-300"
					>
						<LogOut size={18} />
						<span>Logout</span>
					</button>
				</div>
			</div>
		{/if}

		<!-- Profile Button -->
		<button
			onclick={toggleDropdown}
			class="flex w-full items-center gap-3 rounded-xl px-2 py-2 transition-all hover:bg-white/5 {dropdownOpen
				? 'bg-white/5'
				: ''}"
		>
			<div
				class="flex h-10 w-10 items-center justify-center rounded-full bg-gradient-to-br from-purple-500 to-pink-600 font-bold text-white shadow-lg"
			>
				U
			</div>
			<div class="flex-1 text-left">
				<p class="text-sm font-semibold text-white">User</p>
				<p class="text-xs text-slate-400">Pro Plan</p>
			</div>
			<ChevronDown
				size={18}
				class="text-slate-400 transition-transform duration-200 {dropdownOpen ? 'rotate-180' : ''}"
			/>
		</button>
	</div>
</aside>
