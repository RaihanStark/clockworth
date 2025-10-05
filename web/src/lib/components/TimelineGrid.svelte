<script lang="ts">
	const hours = Array.from({ length: 24 }, (_, i) => i);
</script>

<!-- Hour grid -->
{#each hours as hour}
	<div
		class="absolute left-0 right-0 flex h-[100px] border-b border-slate-200/70"
		style="top: {(hour / 24) * 2400}px"
	>
		<div class="w-24 flex-shrink-0 border-r-2 border-slate-200 px-4 py-3">
			<div class="text-center">
				<span class="block text-lg font-bold text-slate-700"
					>{hour === 0 ? 12 : hour > 12 ? hour - 12 : hour}</span
				>
				<span class="block text-xs font-semibold uppercase text-slate-400"
					>{hour >= 12 ? 'PM' : 'AM'}</span
				>
			</div>
		</div>
		<div class="relative flex-1 px-2"></div>
	</div>
{/each}
