<script lang="ts">
	import { Clock, TrendingUp, Target } from '@lucide/svelte';
</script>

<section class="px-6 py-24">
	<div class="mx-auto max-w-6xl">
		<!-- Section Header -->
		<div class="mb-16 text-center">
			<h2 class="mb-4 text-4xl font-black text-slate-900 sm:text-5xl">
				Everything You Need to <span class="bg-gradient-to-r from-emerald-500 to-emerald-600 bg-clip-text text-transparent">Master Your Time</span>
			</h2>
			<p class="text-xl text-slate-600">Simple tools. Powerful insights. Real results.</p>
		</div>

		<!-- Features Grid -->
		<div class="grid gap-8 lg:grid-cols-3">
			<!-- Feature 1 -->
			<div class="group relative overflow-hidden rounded-3xl bg-gradient-to-br from-emerald-500 to-emerald-600 p-8 shadow-2xl transition-all hover:scale-105 hover:shadow-emerald-500/50">
				<div class="absolute -right-8 -top-8 h-32 w-32 rounded-full bg-white/10 blur-2xl"></div>
				<div class="relative">
					<div class="mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-white/20 backdrop-blur-sm">
						<Clock class="text-white" size={32} />
					</div>
					<h3 class="mb-4 text-2xl font-black text-white">Timeline View</h3>
					<p class="text-lg leading-relaxed text-emerald-50">
						Visual hour-by-hour timeline that makes tracking your day effortless and intuitive.
					</p>
				</div>
			</div>

			<!-- Feature 2 -->
			<div class="group relative overflow-hidden rounded-3xl bg-gradient-to-br from-slate-800 to-slate-900 p-8 shadow-2xl transition-all hover:scale-105 hover:shadow-slate-500/50">
				<div class="absolute -right-8 -top-8 h-32 w-32 rounded-full bg-emerald-500/20 blur-2xl"></div>
				<div class="relative">
					<div class="mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-emerald-500/20 backdrop-blur-sm">
						<TrendingUp class="text-emerald-400" size={32} />
					</div>
					<h3 class="mb-4 text-2xl font-black text-white">Worth Calculator</h3>
					<p class="text-lg leading-relaxed text-slate-300">
						Set your hourly rate and instantly see the monetary value of every activity you do.
					</p>
				</div>
			</div>

			<!-- Feature 3 -->
			<div class="group relative overflow-hidden rounded-3xl bg-gradient-to-br from-white to-slate-50 p-8 shadow-2xl transition-all hover:scale-105 hover:shadow-xl">
				<div class="absolute -right-8 -top-8 h-32 w-32 rounded-full bg-emerald-500/10 blur-2xl"></div>
				<div class="relative">
					<div class="mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-gradient-to-br from-emerald-500 to-emerald-600">
						<Target class="text-white" size={32} />
					</div>
					<h3 class="mb-4 text-2xl font-black text-slate-900">Smart Analytics</h3>
					<p class="text-lg leading-relaxed text-slate-600">
						Get actionable insights on productivity patterns and time allocation across activities.
					</p>
				</div>
			</div>
		</div>
	</div>
</section>
