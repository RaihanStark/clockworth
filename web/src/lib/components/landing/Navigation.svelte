<script lang="ts">
	import { ArrowRight, Clock } from '@lucide/svelte';
</script>

<nav class="border-b border-slate-200 bg-white/80 backdrop-blur-sm">
	<div class="mx-auto max-w-7xl px-6 py-4">
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-2">
				<Clock class="text-emerald-600" size={32} />
				<span class="text-2xl font-black text-slate-900">ClockWorth</span>
			</div>
			<a
				href="/signup"
				class="flex items-center gap-2 rounded-lg bg-gradient-to-r from-emerald-500 to-emerald-600 px-5 py-2 font-bold text-white shadow-lg shadow-emerald-500/30 transition-all hover:scale-105"
			>
				Get Started
				<ArrowRight size={18} />
			</a>
		</div>
	</div>
</nav>
