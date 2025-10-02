<script lang="ts">
	import { Calendar, BarChart3, Settings, Clock, Sparkles, TrendingUp, ChevronRight } from '@lucide/svelte';

	interface NavItem {
		label: string;
		icon: any;
		href: string;
		active: boolean;
		badge?: string;
	}

	const navItems: NavItem[] = [
		{ label: 'Today', icon: Calendar, href: '/', active: true },
		{ label: 'Analytics', icon: BarChart3, href: '/analytics', active: false, badge: 'Soon' },
		{ label: 'Settings', icon: Settings, href: '/settings', active: false }
	];
</script>

<aside class="flex h-screen w-72 flex-col bg-gradient-to-b from-slate-950 to-slate-900 text-white">
	<!-- Logo/Brand -->
	<div class="relative overflow-hidden px-6 py-6">
		<div class="relative z-10">
			<div class="mb-2 flex items-center gap-3">
				<div class="flex h-12 w-12 items-center justify-center rounded-2xl bg-gradient-to-br from-emerald-400 to-emerald-600 shadow-lg shadow-emerald-500/50">
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
			<a
				href={item.href}
				class="group relative flex items-center gap-4 rounded-xl px-4 py-3.5 text-base font-semibold transition-all {item.active
					? 'bg-gradient-to-r from-emerald-500 to-emerald-600 text-white shadow-lg shadow-emerald-500/30'
					: 'text-slate-400 hover:bg-white/5 hover:text-white'}"
			>
				<div class="flex h-10 w-10 items-center justify-center rounded-lg {item.active ? 'bg-white/20' : 'bg-white/5'}">
					<svelte:component this={item.icon} size={22} />
				</div>
				<span class="flex-1">{item.label}</span>
				{#if item.badge}
					<span class="rounded-full bg-emerald-400/20 px-2 py-0.5 text-xs font-bold text-emerald-400">
						{item.badge}
					</span>
				{/if}
				{#if item.active}
					<ChevronRight size={18} class="opacity-70" />
				{/if}
			</a>
		{/each}
	</nav>

	<!-- Stats Card -->
	<div class="mx-4 mb-4 space-y-3">
		<!-- Daily Streak -->
		<div class="relative overflow-hidden rounded-2xl bg-gradient-to-br from-amber-500 to-orange-600 p-4 shadow-lg">
			<div class="relative z-10">
				<div class="mb-1 flex items-center gap-2">
					<Sparkles size={16} />
					<span class="text-xs font-bold uppercase tracking-wide">Daily Streak</span>
				</div>
				<p class="text-3xl font-black">7 Days</p>
				<p class="mt-1 text-xs opacity-90">Keep it up! 🔥</p>
			</div>
			<div class="absolute -bottom-4 -right-4 h-24 w-24 rounded-full bg-white/10"></div>
		</div>

		<!-- Goal Progress -->
		<div class="rounded-2xl bg-white/5 p-4 backdrop-blur-sm">
			<div class="mb-2 flex items-center justify-between">
				<div class="flex items-center gap-2">
					<TrendingUp size={16} class="text-emerald-400" />
					<span class="text-xs font-bold text-slate-300">Today's Goal</span>
				</div>
				<span class="text-xs font-bold text-emerald-400">62%</span>
			</div>
			<div class="h-2.5 overflow-hidden rounded-full bg-white/10">
				<div class="h-full w-[62%] rounded-full bg-gradient-to-r from-emerald-400 to-emerald-500 shadow-lg shadow-emerald-500/50"></div>
			</div>
			<p class="mt-2 text-xs text-slate-400">5.0 / 8.0 hours tracked</p>
		</div>
	</div>

	<!-- Footer -->
	<div class="border-t border-white/10 px-6 py-4">
		<div class="flex items-center gap-3">
			<div class="flex h-10 w-10 items-center justify-center rounded-full bg-gradient-to-br from-purple-500 to-pink-600 font-bold text-white">
				U
			</div>
			<div class="flex-1">
				<p class="text-sm font-semibold">User</p>
				<p class="text-xs text-slate-400">Pro Plan</p>
			</div>
		</div>
	</div>
</aside>
