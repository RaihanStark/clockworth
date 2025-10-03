<script lang="ts">
	import { ArrowRight, Clock, TrendingUp, Target } from '@lucide/svelte';
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-50 via-emerald-50 to-slate-50">
	<!-- Navigation -->
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

	<!-- Hero Section -->
	<section class="px-6 py-20 sm:py-32">
		<div class="mx-auto max-w-4xl text-center">
			<h1 class="mb-6 text-5xl font-black tracking-tight text-slate-900 sm:text-7xl">
				Track Your Time.<br />
				<span class="bg-gradient-to-r from-emerald-500 to-emerald-600 bg-clip-text text-transparent">
					Know Your Worth.
				</span>
			</h1>
			<p class="mb-10 text-xl text-slate-600 sm:text-2xl">
				Understand how you spend your time and make every hour count with ClockWorth.
			</p>
			<a
				href="/signup"
				class="inline-flex items-center gap-3 rounded-xl bg-gradient-to-r from-emerald-500 to-emerald-600 px-8 py-4 text-lg font-bold text-white shadow-xl shadow-emerald-500/30 transition-all hover:scale-105"
			>
				Start Tracking Free
				<ArrowRight size={24} />
			</a>
		</div>
	</section>

	<!-- Value of Time Section -->
	<section class="relative overflow-hidden bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 px-6 py-24">
		<!-- Background decoration -->
		<div class="absolute inset-0 opacity-10">
			<div class="absolute left-1/4 top-1/4 h-96 w-96 rounded-full bg-emerald-500 blur-3xl"></div>
			<div class="absolute bottom-1/4 right-1/4 h-96 w-96 rounded-full bg-emerald-400 blur-3xl"></div>
		</div>

		<div class="relative mx-auto max-w-6xl">
			<!-- Header -->
			<div class="mb-16 text-center">
				<h2 class="mb-4 text-5xl font-black text-white sm:text-6xl">
					Time = <span class="bg-gradient-to-r from-emerald-400 to-emerald-300 bg-clip-text text-transparent">Money</span>
				</h2>
				<p class="text-2xl font-bold text-emerald-400">
					You can always make more money, but you can never make more time.
				</p>
			</div>

			<!-- Stats Grid -->
			<div class="mb-12 grid gap-6 md:grid-cols-3">
				<div class="group relative overflow-hidden rounded-2xl bg-gradient-to-br from-emerald-500/20 to-emerald-600/10 p-8 backdrop-blur-sm transition-all hover:scale-105">
					<div class="absolute right-4 top-4 text-6xl opacity-10">⏰</div>
					<div class="relative">
						<div class="mb-2 text-5xl font-black text-emerald-400">24</div>
						<div class="text-lg font-bold text-white">Hours Every Day</div>
						<div class="mt-2 text-sm text-slate-300">Non-renewable. Use them wisely.</div>
					</div>
				</div>

				<div class="group relative overflow-hidden rounded-2xl bg-gradient-to-br from-amber-500/20 to-amber-600/10 p-8 backdrop-blur-sm transition-all hover:scale-105">
					<div class="absolute right-4 top-4 text-6xl opacity-10">💰</div>
					<div class="relative">
						<div class="mb-2 text-5xl font-black text-amber-400">$270</div>
						<div class="text-lg font-bold text-white">Lost Per Task</div>
						<div class="mt-2 text-sm text-slate-300">When you don't value your time.</div>
					</div>
				</div>

				<div class="group relative overflow-hidden rounded-2xl bg-gradient-to-br from-rose-500/20 to-rose-600/10 p-8 backdrop-blur-sm transition-all hover:scale-105">
					<div class="absolute right-4 top-4 text-6xl opacity-10">🚀</div>
					<div class="relative">
						<div class="mb-2 text-5xl font-black text-rose-400">10x</div>
						<div class="text-lg font-bold text-white">Productivity Gain</div>
						<div class="mt-2 text-sm text-slate-300">Focus on high-value work.</div>
					</div>
				</div>
			</div>

			<!-- Key Message -->
			<div class="rounded-3xl border border-emerald-500/30 bg-gradient-to-br from-slate-800/80 to-slate-900/80 p-12 text-center backdrop-blur-sm">
				<p class="mb-6 text-3xl font-bold leading-relaxed text-white">
					Successful people don't trade time for money—<br />
					<span class="bg-gradient-to-r from-emerald-400 to-emerald-300 bg-clip-text text-transparent">
						they leverage their time to create value.
					</span>
				</p>
				<p class="text-xl text-slate-300">
					Track it. Measure it. Optimize it. Your future self will thank you.
				</p>
			</div>
		</div>
	</section>

	<!-- Features Section -->
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

	<!-- CTA Section -->
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

	<!-- Footer -->
	<footer class="border-t border-slate-200 bg-white px-6 py-8">
		<div class="mx-auto max-w-7xl text-center text-sm text-slate-600">
			<p>&copy; 2025 ClockWorth. Make every hour count.</p>
		</div>
	</footer>
</div>
