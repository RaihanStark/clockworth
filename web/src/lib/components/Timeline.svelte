<script lang="ts">
	import { activities, type Activity } from '$lib/stores/activities';
	import { Clock } from '@lucide/svelte';
	import { onMount, onDestroy } from 'svelte';

	interface Props {
		onActivityClick: (activity: Activity) => void;
		onAddActivity: (hour: number, startTime?: string, endTime?: string) => void;
	}

	let { onActivityClick, onAddActivity }: Props = $props();

	const hours = Array.from({ length: 24 }, (_, i) => i);
	let currentTimePosition = $state(0);
	let intervalId: number;

	// Drag selection state
	let isDragging = $state(false);
	let dragStartY = $state(0);
	let dragCurrentY = $state(0);
	let timelineContainer: HTMLDivElement;

	// Activity drag state
	let draggedActivity = $state<Activity | null>(null);
	let activityDragStartY = $state(0);
	let activityDragCurrentY = $state(0);
	let activityOriginalTop = $state(0);

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

	// Drag selection functions
	function handleMouseDown(e: MouseEvent) {
		// Only start drag on the timeline grid, not on activities or buttons
		const target = e.target as HTMLElement;
		if (target.closest('button') || target.closest('[role="button"]')) {
			return;
		}

		e.preventDefault(); // Prevent text selection
		isDragging = true;
		const rect = timelineContainer.getBoundingClientRect();
		dragStartY = e.clientY - rect.top + timelineContainer.scrollTop;
		dragCurrentY = dragStartY;
	}

	function handleMouseMove(e: MouseEvent) {
		if (isDragging) {
			const rect = timelineContainer.getBoundingClientRect();
			dragCurrentY = e.clientY - rect.top + timelineContainer.scrollTop;
		}

		if (draggedActivity) {
			const rect = timelineContainer.getBoundingClientRect();
			activityDragCurrentY = e.clientY - rect.top + timelineContainer.scrollTop;
		}
	}

	function handleMouseUp() {
		if (isDragging) {
			isDragging = false;

			const startY = Math.min(dragStartY, dragCurrentY);
			const endY = Math.max(dragStartY, dragCurrentY);

			// Convert Y positions to time
			const startMinutes = Math.round((startY / 2400) * 24 * 60);
			const endMinutes = Math.round((endY / 2400) * 24 * 60);

			// Round to nearest 15 minutes
			const roundedStart = Math.round(startMinutes / 15) * 15;
			const roundedEnd = Math.round(endMinutes / 15) * 15;

			// Ensure minimum 15 minutes duration
			const finalEnd = roundedEnd <= roundedStart ? roundedStart + 15 : roundedEnd;

			// Format times
			const startHour = Math.floor(roundedStart / 60);
			const startMin = roundedStart % 60;
			const endHour = Math.floor(finalEnd / 60);
			const endMin = finalEnd % 60;

			const startTime = `${String(startHour).padStart(2, '0')}:${String(startMin).padStart(2, '0')}`;
			const endTime = `${String(endHour).padStart(2, '0')}:${String(endMin).padStart(2, '0')}`;

			// Create activity with selected time range
			onAddActivity(startHour, startTime, endTime);
		}

		if (draggedActivity) {
			// Calculate new time based on drag position
			const deltaY = activityDragCurrentY - activityDragStartY;
			const newTop = activityOriginalTop + deltaY;

			// Convert to minutes
			const newStartMinutes = Math.round((newTop / 2400) * 24 * 60);
			const roundedStart = Math.round(newStartMinutes / 15) * 15;

			// Get original duration
			const originalStart = timeToMinutes(draggedActivity.startTime);
			const originalEnd = timeToMinutes(draggedActivity.endTime);
			const duration = originalEnd - originalStart;

			// Calculate new end time
			const newEndMinutes = roundedStart + duration;

			// Format times
			const startHour = Math.floor(roundedStart / 60);
			const startMin = roundedStart % 60;
			const endHour = Math.floor(newEndMinutes / 60);
			const endMin = newEndMinutes % 60;

			const newStartTime = `${String(startHour).padStart(2, '0')}:${String(startMin).padStart(2, '0')}`;
			const newEndTime = `${String(endHour).padStart(2, '0')}:${String(endMin).padStart(2, '0')}`;

			// Update activity
			activities.update(draggedActivity.id, {
				startTime: newStartTime,
				endTime: newEndTime
			});

			draggedActivity = null;
		}
	}

	// Activity drag and click handlers
	let clickTimeout: number;
	let isDraggingActivity = $state(false);

	function handleActivityMouseDown(e: MouseEvent, activity: Activity) {
		e.stopPropagation();
		e.preventDefault();

		isDraggingActivity = false;
		draggedActivity = activity;
		const rect = timelineContainer.getBoundingClientRect();
		activityDragStartY = e.clientY - rect.top + timelineContainer.scrollTop;
		activityDragCurrentY = activityDragStartY;

		// Store original position
		const startMinutes = timeToMinutes(activity.startTime);
		activityOriginalTop = (startMinutes / (24 * 60)) * 2400;

		// Set flag after short delay to distinguish drag from click
		clickTimeout = setTimeout(() => {
			isDraggingActivity = true;
		}, 150) as unknown as number;
	}

	function handleActivityClick(activity: Activity) {
		clearTimeout(clickTimeout);

		// Only trigger click if not dragging
		if (!isDraggingActivity && draggedActivity?.id !== activity.id) {
			onActivityClick(activity);
		}

		isDraggingActivity = false;
	}

	// Computed selection position for visual feedback
	const selectionStyle = $derived.by(() => {
		if (!isDragging) return { display: 'none' };

		const startY = Math.min(dragStartY, dragCurrentY);
		const endY = Math.max(dragStartY, dragCurrentY);
		const height = endY - startY;

		return {
			top: `${startY}px`,
			height: `${height}px`,
			display: 'block'
		};
	});

	// Get activity position with drag offset
	function getActivityDragPosition(activity: Activity) {
		const basePosition = getActivityPosition(activity);

		if (draggedActivity?.id === activity.id) {
			const deltaY = activityDragCurrentY - activityDragStartY;
			const newTop = ((activityOriginalTop + deltaY) / 2400) * 100;
			return {
				...basePosition,
				top: `${newTop}%`
			};
		}

		return basePosition;
	}
</script>

<div
	class="relative flex h-full flex-col overflow-hidden rounded-2xl border-2 border-slate-200 bg-white shadow-xl"
>
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

	<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
	<div
		class="relative flex-1 select-none overflow-y-auto bg-gradient-to-br from-slate-50 to-white"
		bind:this={timelineContainer}
		onmousedown={handleMouseDown}
		onmousemove={handleMouseMove}
		onmouseup={handleMouseUp}
		onmouseleave={() => {
			isDragging = false;
			draggedActivity = null;
		}}
		style="cursor: {isDragging
			? 'grabbing'
			: 'default'}; user-select: none; -webkit-user-select: none;"
		role="application"
		aria-label="Timeline drag to create activity"
	>
		<div class="relative min-h-[2400px]">
			<!-- Drag Selection Overlay -->
			{#if isDragging}
				<div
					class="z-25 pointer-events-none absolute left-24 right-4 rounded-lg border-2 border-emerald-500 bg-emerald-500/20"
					style="top: {selectionStyle.top}; height: {selectionStyle.height};"
				></div>
			{/if}
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

			<!-- Activities -->
			{#each $activities as activity (activity.id)}
				{@const position = getActivityDragPosition(activity)}
				{@const roiColor = getROIColor(activity.roi)}
				{@const icon = getROIIcon(activity.roi)}
				{@const duration = getDuration(activity)}
				{@const isBeingDragged = draggedActivity?.id === activity.id}
				<button
					class="absolute left-24 right-4 z-20 overflow-hidden rounded-2xl border-2 bg-gradient-to-r px-5 py-3 text-left shadow-lg focus:outline-none focus:ring-4 focus:ring-emerald-500/30 {roiColor} {isBeingDragged ? 'scale-105 shadow-2xl cursor-grabbing opacity-80 ring-4 ring-blue-500/50' : 'transition-all hover:scale-[1.01] hover:shadow-2xl cursor-grab'}"
					style="top: {position.top}; height: {position.height}; min-height: 60px; {isBeingDragged ? 'transition: none;' : ''}"
					onmousedown={(e) => handleActivityMouseDown(e, activity)}
					onclick={() => handleActivityClick(activity)}
					aria-label="Edit {activity.title}"
				>
					<div class="flex h-full flex-col justify-center">
						<div class="flex items-start justify-between gap-2">
							<div class="flex-1">
								<div class="mb-1 flex items-center gap-2">
									<span class="text-lg">{icon}</span>
									<p class="truncate text-base font-bold text-white drop-shadow-sm">
										{activity.title}
									</p>
								</div>
								<div class="flex items-center gap-3 text-xs text-white/90">
									<span class="font-semibold">{activity.startTime} - {activity.endTime}</span>
									<span class="rounded-full bg-white/20 px-2 py-0.5 font-bold">{duration}</span>
								</div>
							</div>
							<div class="flex flex-col items-end gap-1">
								<span
									class="rounded-lg bg-white/25 px-2 py-1 text-xs font-bold uppercase tracking-wide text-white"
									>{activity.roi}</span
								>
								<span class="rounded-lg bg-white/20 px-2 py-0.5 text-xs font-medium text-white/80"
									>{activity.category}</span
								>
							</div>
						</div>
					</div>
				</button>
			{/each}
		</div>
	</div>
</div>
