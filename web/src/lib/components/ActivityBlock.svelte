<script lang="ts">
	import type { Activity } from '$lib/types';
	import { api } from '$lib/api';

	interface Props {
		activity: Activity;
		containerElement: HTMLElement | null;
		onActivityClick: (activity: Activity) => void;
		onActivityUpdate: (activity: Activity) => void;
	}

	let { activity, containerElement, onActivityClick, onActivityUpdate }: Props = $props();

	// Drag state
	let isDragging = $state(false);
	let dragStartY = $state(0);
	let dragCurrentY = $state(0);
	let originalTop = $state(0);

	// Resize state
	let isResizing = $state(false);
	let resizeMode = $state<'top' | 'bottom' | null>(null);
	let resizeStartY = $state(0);
	let resizeCurrentY = $state(0);
	let originalStartTime = $state('');
	let originalEndTime = $state('');

	// Click detection
	let clickTimeout: number;
	let wasDragged = $state(false);

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

	// Drag handlers
	function handleDragStart(e: MouseEvent) {
		if (!containerElement || isResizing) return;

		e.stopPropagation();
		e.preventDefault();

		wasDragged = false;
		isDragging = true;

		const rect = containerElement.getBoundingClientRect();
		dragStartY = e.clientY - rect.top + containerElement.scrollTop;
		dragCurrentY = dragStartY;

		const startMinutes = timeToMinutes(activity.startTime);
		originalTop = (startMinutes / (24 * 60)) * 2400;

		clickTimeout = setTimeout(() => {
			wasDragged = true;
		}, 150) as unknown as number;
	}

	function handleDragMove(e: MouseEvent) {
		if (!isDragging || !containerElement) return;

		const rect = containerElement.getBoundingClientRect();
		dragCurrentY = e.clientY - rect.top + containerElement.scrollTop;
	}

	function handleDragEnd() {
		if (!isDragging) return;

		isDragging = false;
		clearTimeout(clickTimeout);

		if (!wasDragged) {
			onActivityClick(activity);
			return;
		}

		const deltaY = dragCurrentY - dragStartY;
		const newTop = originalTop + deltaY;

		const newStartMinutes = Math.round((newTop / 2400) * 24 * 60);
		let roundedStart = Math.round(newStartMinutes / 15) * 15;

		const originalStart = timeToMinutes(activity.startTime);
		const originalEnd = timeToMinutes(activity.endTime);
		const duration = originalEnd - originalStart;

		// Enforce boundaries (0:00 to 23:59)
		const maxMinutes = 24 * 60 - 1; // 1439 minutes = 23:59
		const maxStartMinutes = maxMinutes - duration; // Max start time to keep end within 23:59
		roundedStart = Math.max(0, Math.min(roundedStart, maxStartMinutes));

		const newEndMinutes = roundedStart + duration;

		const startHour = Math.floor(roundedStart / 60);
		const startMin = roundedStart % 60;
		const endHour = Math.floor(newEndMinutes / 60);
		const endMin = newEndMinutes % 60;

		const newStartTime = `${String(startHour).padStart(2, '0')}:${String(startMin).padStart(2, '0')}`;
		const newEndTime = `${String(endHour).padStart(2, '0')}:${String(endMin).padStart(2, '0')}`;

		// Update activity via API
		api
			.updateActivity(activity.id, {
				title: activity.title,
				start_time: newStartTime,
				end_time: newEndTime,
				category: activity.category,
				roi: activity.roi
			})
			.then((apiActivity) => {
				const updatedActivity: Activity = {
					id: apiActivity.id,
					title: apiActivity.title,
					startTime: apiActivity.start_time,
					endTime: apiActivity.end_time,
					category: apiActivity.category,
					roi: apiActivity.roi as 'High' | 'Neutral' | 'Waste'
				};
				onActivityUpdate(updatedActivity);
			})
			.catch((error) => {
				console.error('Failed to update activity:', error);
			});
	}

	// Resize handlers
	function handleResizeStart(e: MouseEvent, mode: 'top' | 'bottom') {
		if (!containerElement) return;

		e.stopPropagation();
		e.preventDefault();

		isResizing = true;
		resizeMode = mode;
		originalStartTime = activity.startTime;
		originalEndTime = activity.endTime;

		const rect = containerElement.getBoundingClientRect();
		resizeStartY = e.clientY - rect.top + containerElement.scrollTop;
		resizeCurrentY = resizeStartY;
	}

	function handleResizeMove(e: MouseEvent) {
		if (!isResizing || !containerElement) return;

		const rect = containerElement.getBoundingClientRect();
		resizeCurrentY = e.clientY - rect.top + containerElement.scrollTop;
	}

	function handleResizeEnd() {
		if (!isResizing || !resizeMode) return;

		isResizing = false;

		const deltaY = resizeCurrentY - resizeStartY;
		const deltaMinutes = Math.round((deltaY / 2400) * 24 * 60);

		let newStartTime = originalStartTime;
		let newEndTime = originalEndTime;

		if (resizeMode === 'top') {
			const startMinutes = timeToMinutes(originalStartTime);
			let newStartMinutes = Math.round((startMinutes + deltaMinutes) / 15) * 15;
			const endMinutes = timeToMinutes(originalEndTime);

			// Enforce boundaries (0:00 minimum, maintain 15-minute minimum duration)
			newStartMinutes = Math.max(0, Math.min(newStartMinutes, endMinutes - 15));

			const startHour = Math.floor(newStartMinutes / 60);
			const startMin = newStartMinutes % 60;
			newStartTime = `${String(startHour).padStart(2, '0')}:${String(startMin).padStart(2, '0')}`;
		} else {
			const endMinutes = timeToMinutes(originalEndTime);
			let newEndMinutes = Math.round((endMinutes + deltaMinutes) / 15) * 15;
			const startMinutes = timeToMinutes(originalStartTime);

			// Enforce boundaries (23:59 maximum, maintain 15-minute minimum duration)
			const maxEndMinutes = 24 * 60 - 1; // 1439 minutes = 23:59
			newEndMinutes = Math.max(startMinutes + 15, Math.min(newEndMinutes, maxEndMinutes));

			const endHour = Math.floor(newEndMinutes / 60);
			const endMin = newEndMinutes % 60;
			newEndTime = `${String(endHour).padStart(2, '0')}:${String(endMin).padStart(2, '0')}`;
		}

		// Update activity via API
		api
			.updateActivity(activity.id, {
				title: activity.title,
				start_time: newStartTime,
				end_time: newEndTime,
				category: activity.category,
				roi: activity.roi
			})
			.then((apiActivity) => {
				const updatedActivity: Activity = {
					id: apiActivity.id,
					title: apiActivity.title,
					startTime: apiActivity.start_time,
					endTime: apiActivity.end_time,
					category: apiActivity.category,
					roi: apiActivity.roi as 'High' | 'Neutral' | 'Waste'
				};
				onActivityUpdate(updatedActivity);
			})
			.catch((error) => {
				console.error('Failed to update activity:', error);
			});

		resizeMode = null;
	}

	// Computed position
	const position = $derived.by(() => {
		const basePosition = getActivityPosition(activity);

		if (isDragging) {
			const deltaY = dragCurrentY - dragStartY;
			const newTop = ((originalTop + deltaY) / 2400) * 100;
			return {
				...basePosition,
				top: `${newTop}%`
			};
		}

		if (isResizing && resizeMode) {
			const deltaY = resizeCurrentY - resizeStartY;
			const deltaPercent = (deltaY / 2400) * 100;

			if (resizeMode === 'top') {
				const currentTop = parseFloat(basePosition.top);
				const currentHeight = parseFloat(basePosition.height);
				const newTop = currentTop + deltaPercent;
				const newHeight = currentHeight - deltaPercent;

				if (newHeight >= 1.04) {
					return {
						top: `${newTop}%`,
						height: `${newHeight}%`
					};
				}
			} else {
				const currentHeight = parseFloat(basePosition.height);
				const newHeight = currentHeight + deltaPercent;

				if (newHeight >= 1.04) {
					return {
						...basePosition,
						height: `${newHeight}%`
					};
				}
			}
		}

		return basePosition;
	});

	const roiColor = $derived(getROIColor(activity.roi));
	const icon = $derived(getROIIcon(activity.roi));
	const duration = $derived(getDuration(activity));

	// Global event listeners
	$effect(() => {
		if (isDragging) {
			window.addEventListener('mousemove', handleDragMove);
			window.addEventListener('mouseup', handleDragEnd);
			return () => {
				window.removeEventListener('mousemove', handleDragMove);
				window.removeEventListener('mouseup', handleDragEnd);
			};
		}
	});

	$effect(() => {
		if (isResizing) {
			window.addEventListener('mousemove', handleResizeMove);
			window.addEventListener('mouseup', handleResizeEnd);
			return () => {
				window.removeEventListener('mousemove', handleResizeMove);
				window.removeEventListener('mouseup', handleResizeEnd);
			};
		}
	});
</script>

<div
	class="group absolute left-24 right-4 z-20 overflow-visible rounded-2xl"
	style="top: {position.top}; height: {position.height}; min-height: 60px; {isDragging || isResizing
		? 'transition: none;'
		: ''}"
>
	<!-- Top Resize Handle -->
	<div
		class="absolute -top-1 left-0 right-0 z-30 h-2 cursor-ns-resize opacity-0 transition-opacity group-hover:opacity-100"
		onmousedown={(e) => handleResizeStart(e, 'top')}
		role="button"
		tabindex="-1"
		aria-label="Resize activity start time"
	>
		<div class="h-full w-full"></div>
	</div>

	<!-- Activity Content -->
	<button
		class="relative h-full w-full overflow-hidden rounded-2xl border-2 bg-gradient-to-r px-5 py-3 text-left shadow-lg focus:outline-none focus:ring-4 focus:ring-emerald-500/30 {roiColor} {isDragging
			? 'scale-105 cursor-grabbing opacity-80 shadow-2xl ring-4 ring-blue-500/50'
			: isResizing
				? 'shadow-2xl ring-4 ring-purple-500/50'
				: 'cursor-grab transition-all hover:scale-[1.01] hover:shadow-2xl'}"
		onmousedown={handleDragStart}
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

	<!-- Bottom Resize Handle -->
	<div
		class="absolute -bottom-1 left-0 right-0 z-30 h-2 cursor-ns-resize opacity-0 transition-opacity group-hover:opacity-100"
		onmousedown={(e) => handleResizeStart(e, 'bottom')}
		role="button"
		tabindex="-1"
		aria-label="Resize activity end time"
	>
		<div class="h-full w-full"></div>
	</div>
</div>
