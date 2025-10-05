<script lang="ts">
	import { activities, type Activity } from '$lib/stores/activities';
	import { Clock } from '@lucide/svelte';
	import { onMount, onDestroy } from 'svelte';
	import ActivityBlock from './ActivityBlock.svelte';
	import TimelineGrid from './TimelineGrid.svelte';

	interface Props {
		onActivityClick: (activity: Activity) => void;
		onAddActivity: (hour: number, startTime?: string, endTime?: string) => void;
	}

	let { onActivityClick, onAddActivity }: Props = $props();

	let timelineContainer: HTMLDivElement;
	let currentTimePosition = $state(0);
	let intervalId: number;

	// Drag-to-create state
	let isDraggingToCreate = $state(false);
	let dragStartY = $state(0);
	let dragCurrentY = $state(0);

	function updateCurrentTime() {
		const now = new Date();
		const hours = now.getHours();
		const minutes = now.getMinutes();
		const totalMinutes = hours * 60 + minutes;
		currentTimePosition = (totalMinutes / (24 * 60)) * 2400;
	}

	onMount(() => {
		updateCurrentTime();
		intervalId = setInterval(updateCurrentTime, 60000) as unknown as number;
	});

	onDestroy(() => {
		if (intervalId) {
			clearInterval(intervalId);
		}
	});

	// Drag-to-create handlers
	function handleMouseDown(e: MouseEvent) {
		const target = e.target as HTMLElement;
		if (target.closest('button') || target.closest('[role="button"]')) {
			return;
		}

		e.preventDefault();
		isDraggingToCreate = true;
		const rect = timelineContainer.getBoundingClientRect();
		dragStartY = e.clientY - rect.top + timelineContainer.scrollTop;
		dragCurrentY = dragStartY;
	}

	function handleMouseMove(e: MouseEvent) {
		if (!isDraggingToCreate) return;

		const rect = timelineContainer.getBoundingClientRect();
		dragCurrentY = e.clientY - rect.top + timelineContainer.scrollTop;
	}

	function handleMouseUp() {
		if (!isDraggingToCreate) return;

		isDraggingToCreate = false;

		const startY = Math.min(dragStartY, dragCurrentY);
		const endY = Math.max(dragStartY, dragCurrentY);

		const startMinutes = Math.round((startY / 2400) * 24 * 60);
		const endMinutes = Math.round((endY / 2400) * 24 * 60);

		let roundedStart = Math.round(startMinutes / 15) * 15;
		let roundedEnd = Math.round(endMinutes / 15) * 15;

		// Enforce boundaries (0:00 to 23:59)
		const maxMinutes = 24 * 60 - 1; // 1439 minutes = 23:59
		roundedStart = Math.max(0, Math.min(roundedStart, maxMinutes));
		roundedEnd = Math.max(0, Math.min(roundedEnd, maxMinutes));

		const finalEnd = roundedEnd <= roundedStart ? roundedStart + 15 : roundedEnd;

		const startHour = Math.floor(roundedStart / 60);
		const startMin = roundedStart % 60;
		const endHour = Math.floor(finalEnd / 60);
		const endMin = finalEnd % 60;

		const startTime = `${String(startHour).padStart(2, '0')}:${String(startMin).padStart(2, '0')}`;
		const endTime = `${String(endHour).padStart(2, '0')}:${String(endMin).padStart(2, '0')}`;

		onAddActivity(startHour, startTime, endTime);
	}

	const selectionStyle = $derived.by(() => {
		if (!isDraggingToCreate) return { display: 'none' };

		const startY = Math.min(dragStartY, dragCurrentY);
		const endY = Math.max(dragStartY, dragCurrentY);
		const height = endY - startY;

		return {
			top: `${startY}px`,
			height: `${height}px`,
			display: 'block'
		};
	});
</script>

<div
	class="relative flex h-full flex-col overflow-hidden rounded-2xl border-2 border-slate-200 bg-white shadow-xl"
>
	<!-- Header -->
	<div
		class="sticky top-0 z-10 border-b-2 border-slate-200 bg-gradient-to-r from-white to-slate-50 px-6 py-5"
	>
		<div class="flex items-center gap-3">
			<div
				class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-blue-500 to-blue-600 text-white"
			>
				<Clock size={20} />
			</div>
			<h2 class="text-2xl font-black text-slate-900">Today's Timeline</h2>
		</div>
	</div>

	<!-- Timeline Container -->
	<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
	<div
		class="relative flex-1 select-none overflow-y-auto bg-gradient-to-br from-slate-50 to-white"
		bind:this={timelineContainer}
		onmousedown={handleMouseDown}
		onmousemove={handleMouseMove}
		onmouseup={handleMouseUp}
		onmouseleave={() => {
			isDraggingToCreate = false;
		}}
		style="cursor: {isDraggingToCreate
			? 'grabbing'
			: 'default'}; user-select: none; -webkit-user-select: none;"
		role="application"
		aria-label="Timeline drag to create activity"
	>
		<div class="relative min-h-[2400px]">
			<!-- Drag Selection Overlay -->
			{#if isDraggingToCreate}
				<div
					class="z-25 pointer-events-none absolute left-24 right-4 rounded-lg border-2 border-emerald-500 bg-emerald-500/20"
					style="top: {selectionStyle.top}; height: {selectionStyle.height};"
				></div>
			{/if}

			<!-- Current Time Indicator -->
			<div class="absolute left-0 right-0 z-30 flex items-center" style="top: {currentTimePosition}px">
				<div class="flex h-3 w-3 items-center justify-center rounded-full bg-red-500 shadow-lg">
					<div class="h-1.5 w-1.5 rounded-full bg-white"></div>
				</div>
				<div class="h-0.5 flex-1 bg-red-500 shadow-sm"></div>
			</div>

			<!-- Hour Grid -->
			<TimelineGrid />

			<!-- Activities -->
			{#each $activities as activity (activity.id)}
				<ActivityBlock {activity} containerElement={timelineContainer} {onActivityClick} />
			{/each}
		</div>
	</div>
</div>
