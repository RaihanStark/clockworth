<script lang="ts">
	import { settings } from '$lib/stores/settings';
	import { api } from '$lib/api';
	import { DollarSign, Save } from '@lucide/svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
	
	let hourlyWorth = $state(data.hourlyWorth || $settings.hourlyWorth);
	let saved = $state(false);
	let saving = $state(false);
	let error = $state<string | null>(null);

	async function handleSave() {
		saving = true;
		error = null;
		
		try {
			// Save to API
			await api.updateHourlyWorth(hourlyWorth);
			
			// Update local store
			settings.update((s) => ({ ...s, hourlyWorth }));
			
			saved = true;
			setTimeout(() => {
				saved = false;
			}, 2000);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to save settings';
		} finally {
			saving = false;
		}
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 p-8">
	<!-- Header -->
	<div class="mb-8">
		<h1 class="text-4xl font-black text-slate-900">Settings</h1>
		<p class="mt-2 text-lg text-slate-600">Configure your time economics preferences</p>
	</div>

	<div class="mx-auto max-w-2xl">
		<div class="rounded-2xl bg-white p-8 shadow-lg">
			<div class="mb-6">
				<h2 class="text-2xl font-black text-slate-900">Hourly Worth</h2>
				<p class="mt-1 text-sm text-slate-600">
					Set your hourly rate to calculate time value and ROI
				</p>
			</div>

			<!-- Hourly Rate Input -->
			<div class="mb-8">
				<label
					for="hourlyWorth"
					class="mb-4 flex items-center gap-2 text-sm font-bold uppercase tracking-wide text-slate-700"
				>
					<DollarSign size={18} class="text-emerald-600" />
					Hourly Rate
				</label>

				<div class="relative flex items-center">
					<span class="absolute text-4xl font-black text-slate-400" style="left: 1.5rem;">$</span>
					<input
						id="hourlyWorth"
						type="number"
						bind:value={hourlyWorth}
						min="0"
						step="1"
						class="w-full rounded-2xl border-2 border-slate-200 bg-slate-50 px-20 py-6 text-center text-5xl font-black text-slate-900 transition-all focus:border-emerald-500 focus:bg-white focus:outline-none focus:ring-4 focus:ring-emerald-500/20"
					/>
					<span class="absolute right-6 text-2xl font-bold text-slate-400">/hr</span>
				</div>
			</div>

			<!-- Save Button -->
			<button
				onclick={handleSave}
				class="flex w-full items-center justify-center gap-2 rounded-xl bg-gradient-to-r from-emerald-500 to-emerald-600 px-6 py-4 font-bold text-white shadow-lg shadow-emerald-500/30 transition-all hover:scale-105 hover:shadow-xl active:scale-95 disabled:opacity-50 disabled:hover:scale-100"
				disabled={saving || hourlyWorth === $settings.hourlyWorth}
			>
				<Save size={20} />
				{saving ? 'Saving...' : saved ? 'Saved!' : 'Save Changes'}
			</button>

			{#if saved}
				<div class="mt-4 rounded-xl bg-emerald-50 p-4 text-center">
					<p class="font-semibold text-emerald-700">✓ Settings saved successfully!</p>
				</div>
			{/if}

			{#if error}
				<div class="mt-4 rounded-xl bg-red-50 p-4 text-center">
					<p class="font-semibold text-red-700">✗ {error}</p>
				</div>
			{/if}
		</div>
	</div>
</div>

<style>
	/* Remove number input spinners */
	input[type='number']::-webkit-inner-spin-button,
	input[type='number']::-webkit-outer-spin-button {
		-webkit-appearance: none;
		margin: 0;
	}

	input[type='number'] {
		-moz-appearance: textfield;
		appearance: textfield;
	}
</style>
