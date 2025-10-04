<script lang="ts">
	import { Clock } from '@lucide/svelte';
</script>

<div class="flex min-h-screen items-center justify-center bg-gradient-to-br from-slate-50 to-slate-100 p-4">
	<div class="w-full max-w-md">
		<!-- Logo -->
		<div class="mb-8 text-center">
			<div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-gradient-to-br from-emerald-400 to-emerald-600 shadow-lg shadow-emerald-500/50">
				<Clock size={32} class="text-white" />
			</div>
			<h1 class="text-3xl font-black text-slate-900">Clockworth</h1>
			<p class="mt-1 text-sm font-semibold text-emerald-600">Time Economics</p>
		</div>

		<!-- Card -->
		<div class="rounded-2xl bg-white p-8 shadow-xl">
			<slot />
		</div>
	</div>
</div>
