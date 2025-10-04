<script lang="ts">
	import { ArrowRight, Clock } from '@lucide/svelte';
</script>

<section class="relative overflow-hidden bg-gradient-to-br from-slate-900 via-emerald-900 to-slate-900 px-6 py-24">
	<!-- Background Pattern -->
	<div class="absolute inset-0 opacity-5">
		<div class="absolute left-0 top-0 h-full w-full" style="background-image: radial-gradient(circle, white 1px, transparent 1px); background-size: 50px 50px;"></div>
	</div>

	<div class="relative mx-auto max-w-4xl text-center">
		<div class="mb-6 inline-flex items-center gap-2 rounded-full border border-emerald-500/30 bg-emerald-500/10 px-4 py-2 backdrop-blur-sm">
			<Clock class="text-emerald-400" size={18} />
			<span class="text-sm font-bold text-emerald-400">Limited Time Offer</span>
		</div>

		<h2 class="mb-4 text-4xl font-black text-white sm:text-5xl">
			Ready to Take Control of Your Time?
		</h2>
		<p class="mb-10 text-xl text-slate-300">
			Join hundreds of people making every hour count. Start tracking for free today.
		</p>

		<a
			href="/signup"
			class="inline-flex items-center gap-2 rounded-xl bg-gradient-to-r from-emerald-500 to-emerald-600 px-10 py-4 text-lg font-bold text-white shadow-xl shadow-emerald-500/30 transition-all hover:scale-105"
		>
			Get Started Free
			<ArrowRight size={20} />
		</a>

		<p class="mt-8 text-sm text-slate-400">No credit card required • Free forever</p>
	</div>
</section>
