<script lang="ts">
	import { activities, type Activity } from '$lib/stores/activities';
	import { Plus, Clock } from '@lucide/svelte';
	import { onMount, onDestroy } from 'svelte';

	interface Props {
		onActivityClick: (activity: Activity) => void;
		onAddActivity: (hour: number) => void;
	}

	let { onActivityClick, onAddActivity }: Props = $props();

	const hours = Array.from({ length: 24 }, (_, i) => i);
	let currentTimePosition = $state(0);
	let intervalId: number;

	function updateCurrentTime() {
		const now = new Date();
		const hours = now.getHours();
		const minutes = now.getMinutes();
		const totalMinutes = hours * 60 + minutes;
		// Calculate position as percentage of the day (2400px total height)
		currentTimePosition = (totalMinutes / (24 * 60)) * 2400;
	}

	onMount(() => {
		updateCurrentTime();
		// Update every minute
		intervalId = setInterval(updateCurrentTime, 60000) as unknown as number;
	});

	onDestroy(() => {
		if (intervalId) {
			clearInterval(intervalId);
		}
	});

	function timeToMinutes(time: string): number {
		const [h, m] = time.split(':').map(Number);
		return h * 60 + m;
	}

	function getActivityPosition(activity: Activity) {
		const startMinutes = timeToMinutes(activity.startTime);
		const endMinutes = timeToMinutes(activity.endTime);
		const top = (startMinutes / (24 * 60)) * 100;
		const height = ((endMinutes - startMinutes) / (24 * 60)) * 100;
		return { top: `${top}%`, height: `${height}%` };
	}

	function getROIColor(roi: Activity['roi']) {
		switch (roi) {
			case 'High':
				return 'from-emerald-400 to-emerald-600 border-emerald-300';
			case 'Neutral':
				return 'from-slate-300 to-slate-500 border-slate-200';
			case 'Waste':
				return 'from-red-400 to-red-600 border-red-300';
		}
	}

	function getROIIcon(roi: Activity['roi']) {
		switch (roi) {
			case 'High':
				return '🎯';
			case 'Neutral':
				return '⚪';
			case 'Waste':
				return '⚠️';
		}
	}

	function formatHour(hour: number): string {
		const period = hour >= 12 ? 'PM' : 'AM';
		const displayHour = hour === 0 ? 12 : hour > 12 ? hour - 12 : hour;
		return `${displayHour} ${period}`;
	}

	function getDuration(activity: Activity): string {
		const start = timeToMinutes(activity.startTime);
		const end = timeToMinutes(activity.endTime);
		const minutes = end - start;
		const hours = Math.floor(minutes / 60);
		const mins = minutes % 60;
		if (hours > 0 && mins > 0) return `${hours}h ${mins}m`;
		if (hours > 0) return `${hours}h`;
		return `${mins}m`;
	}
</script>

<div class="relative flex h-full flex-col overflow-hidden rounded-2xl border-2 border-slate-200 bg-white shadow-xl">
	<div class="sticky top-0 z-10 border-b-2 border-slate-200 bg-gradient-to-r from-white to-slate-50 px-6 py-5">
		<div class="flex items-center gap-3">
			<div class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-blue-500 to-blue-600 text-white">
				<Clock size={20} />
			</div>
			<h2 class="text-2xl font-black text-slate-900">Today's Timeline</h2>
		</div>
	</div>

	<div class="relative flex-1 overflow-y-auto bg-gradient-to-br from-slate-50 to-white">
		<div class="relative min-h-[2400px]">
			<!-- Current Time Indicator -->
			<div
				class="absolute left-0 right-0 z-30 flex items-center"
				style="top: {currentTimePosition}px"
			>
				<div class="flex h-3 w-3 items-center justify-center rounded-full bg-red-500 shadow-lg">
					<div class="h-1.5 w-1.5 rounded-full bg-white"></div>
				</div>
				<div class="h-0.5 flex-1 bg-red-500 shadow-sm"></div>
			</div>

			<!-- Hour grid -->
			{#each hours as hour}
				<div
					class="group absolute left-0 right-0 flex h-[100px] cursor-pointer border-b border-slate-200/70 transition-all hover:bg-emerald-50/50"
					style="top: {(hour / 24) * 2400}px"
					role="button"
					tabindex="0"
					onclick={() => onAddActivity(hour)}
					onkeydown={(e) => {
						if (e.key === 'Enter' || e.key === ' ') {
							e.preventDefault();
							onAddActivity(hour);
						}
					}}
				>
					<div class="w-24 flex-shrink-0 border-r-2 border-slate-200 px-4 py-3">
						<div class="text-center">
							<span class="block text-lg font-bold text-slate-700">{hour === 0 ? 12 : hour > 12 ? hour - 12 : hour}</span>
							<span class="block text-xs font-semibold uppercase text-slate-400">{hour >= 12 ? 'PM' : 'AM'}</span>
						</div>
					</div>
					<div class="relative flex-1 px-2">
						<button
							class="absolute right-4 top-3 hidden rounded-xl bg-gradient-to-r from-emerald-500 to-emerald-600 p-2 text-white shadow-lg shadow-emerald-500/30 transition-all hover:scale-110 group-hover:block"
							onclick={(e) => {
								e.stopPropagation();
								onAddActivity(hour);
							}}
							aria-label="Add activity at {formatHour(hour)}"
						>
							<Plus size={18} />
						</button>
					</div>
				</div>
			{/each}

			<!-- Activities -->
			{#each $activities as activity (activity.id)}
				{@const position = getActivityPosition(activity)}
				{@const roiColor = getROIColor(activity.roi)}
				{@const icon = getROIIcon(activity.roi)}
				{@const duration = getDuration(activity)}
				<button
					class="absolute left-24 right-4 z-20 overflow-hidden rounded-2xl border-2 bg-gradient-to-r px-5 py-3 text-left shadow-lg transition-all hover:scale-[1.01] hover:shadow-2xl focus:outline-none focus:ring-4 focus:ring-emerald-500/30 {roiColor}"
					style="top: {position.top}; height: {position.height}; min-height: 60px"
					onclick={() => onActivityClick(activity)}
					aria-label="Edit {activity.title}"
				>
					<div class="flex h-full flex-col justify-center">
						<div class="flex items-start justify-between gap-2">
							<div class="flex-1">
								<div class="mb-1 flex items-center gap-2">
									<span class="text-lg">{icon}</span>
									<p class="truncate text-base font-bold text-white drop-shadow-sm">{activity.title}</p>
								</div>
								<div class="flex items-center gap-3 text-xs text-white/90">
									<span class="font-semibold">{activity.startTime} - {activity.endTime}</span>
									<span class="rounded-full bg-white/20 px-2 py-0.5 font-bold">{duration}</span>
								</div>
							</div>
							<div class="flex flex-col items-end gap-1">
								<span class="rounded-lg bg-white/25 px-2 py-1 text-xs font-bold uppercase tracking-wide text-white">{activity.roi}</span>
								<span class="rounded-lg bg-white/20 px-2 py-0.5 text-xs font-medium text-white/80">{activity.category}</span>
							</div>
						</div>
					</div>
				</button>
			{/each}
		</div>
	</div>
</div>
