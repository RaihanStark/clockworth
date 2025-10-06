<script lang="ts">
	import type { Activity, ROICategory } from '$lib/types';
	import { api } from '$lib/api';
	import { X, Clock, Tag, Target, Calendar } from '@lucide/svelte';
	import { onMount } from 'svelte';

	interface Props {
		activity?: Activity | null;
		initialHour?: number;
		initialStartTime?: string;
		initialEndTime?: string;
		onClose: () => void;
		onActivityAdded?: (activity: Activity) => void;
		onActivityUpdated?: (activity: Activity) => void;
		onActivityDeleted?: (id: string) => void;
	}

	let {
		activity = null,
		initialHour = 9,
		initialStartTime,
		initialEndTime,
		onClose,
		onActivityAdded,
		onActivityUpdated,
		onActivityDeleted
	}: Props = $props();

	let modalElement: HTMLDivElement;
	let firstFocusable: HTMLElement;

	let formData = $state({
		title: activity?.title || '',
		startTime:
			activity?.startTime || initialStartTime || `${String(initialHour).padStart(2, '0')}:00`,
		endTime:
			activity?.endTime || initialEndTime || `${String(initialHour + 1).padStart(2, '0')}:00`,
		category: activity?.category || 'Development',
		roi: (activity?.roi || 'Neutral') as ROICategory
	});

	let errors = $state({
		title: '',
		time: '',
		overlap: ''
	});

	const categories = [
		{ name: 'Development', icon: '💻' },
		{ name: 'Meetings', icon: '👥' },
		{ name: 'Personal', icon: '🏠' },
		{ name: 'Health', icon: '💪' },
		{ name: 'Learning', icon: '📚' },
		{ name: 'Other', icon: '📌' }
	];

	const roiOptions: { value: ROICategory; label: string; icon: string; description: string }[] = [
		{ value: 'High', label: 'High ROI', icon: '🎯', description: 'Productive & valuable' },
		{ value: 'Neutral', label: 'Neutral', icon: '⚪', description: 'Necessary tasks' },
		{ value: 'Waste', label: 'Waste', icon: '⚠️', description: 'Low-value time' }
	];

	function validateForm(): boolean {
		errors.title = '';
		errors.time = '';
		errors.overlap = '';

		if (!formData.title.trim()) {
			errors.title = 'Title is required';
			return false;
		}

		const startMinutes = timeToMinutes(formData.startTime);
		const endMinutes = timeToMinutes(formData.endTime);

		if (endMinutes <= startMinutes) {
			errors.time = 'End time must be after start time';
			return false;
		}

		return true;
	}

	function timeToMinutes(time: string): number {
		const [h, m] = time.split(':').map(Number);
		return h * 60 + m;
	}

	async function handleSubmit() {
		if (!validateForm()) return;

		try {
			if (activity) {
				// Update existing activity
				const apiActivity = await api.updateActivity(activity.id, {
					title: formData.title,
					start_time: formData.startTime,
					end_time: formData.endTime,
					category: formData.category,
					roi: formData.roi
				});

				const updatedActivity: Activity = {
					id: apiActivity.id,
					title: apiActivity.title,
					startTime: apiActivity.start_time,
					endTime: apiActivity.end_time,
					category: apiActivity.category,
					roi: apiActivity.roi as ROICategory
				};

				onActivityUpdated?.(updatedActivity);
			} else {
				// Create new activity
				const today = new Date().toISOString().split('T')[0];
				const apiActivity = await api.createActivity({
					title: formData.title,
					start_time: formData.startTime,
					end_time: formData.endTime,
					category: formData.category,
					roi: formData.roi,
					activity_date: today
				});

				const newActivity: Activity = {
					id: apiActivity.id,
					title: apiActivity.title,
					startTime: apiActivity.start_time,
					endTime: apiActivity.end_time,
					category: apiActivity.category,
					roi: apiActivity.roi as ROICategory
				};

				onActivityAdded?.(newActivity);
			}
			closeModal();
		} catch (error: any) {
			console.error('Failed to save activity:', error);
			// Check if it's an overlap error (409 Conflict)
			if (error.message && error.message.includes('overlap')) {
				errors.overlap = 'This time slot overlaps with an existing activity';
			} else {
				errors.overlap = 'Failed to save activity. Please try again.';
			}
		}
	}

	async function handleDelete() {
		if (!activity) return;

		const confirmed = confirm('Are you sure you want to delete this activity?');
		if (!confirmed) return;

		try {
			await api.deleteActivity(activity.id);
			onActivityDeleted?.(activity.id);
			closeModal();
		} catch (error) {
			console.error('Failed to delete activity:', error);
			// You could show an error message to the user here
		}
	}

	function closeModal() {
		onClose();
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			closeModal();
		}

		// Focus trap
		if (e.key === 'Tab') {
			const focusableElements = modalElement.querySelectorAll(
				'button:not([disabled]), input:not([disabled]), select:not([disabled]), textarea:not([disabled])'
			);
			const firstElement = focusableElements[0] as HTMLElement;
			const lastElement = focusableElements[focusableElements.length - 1] as HTMLElement;

			if (e.shiftKey && document.activeElement === firstElement) {
				e.preventDefault();
				lastElement.focus();
			} else if (!e.shiftKey && document.activeElement === lastElement) {
				e.preventDefault();
				firstElement.focus();
			}
		}
	}

	function handleBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget) {
			closeModal();
		}
	}

	onMount(() => {
		firstFocusable = modalElement.querySelector('input, select, button') as HTMLElement;
		firstFocusable?.focus();
	});
</script>

<svelte:window onkeydown={handleKeydown} />

<div
	class="fixed inset-0 z-50 flex items-center justify-center bg-black/80 p-4"
	style="will-change: contents;"
	onclick={handleBackdropClick}
	onkeydown={(e) => e.key === 'Enter' && closeModal()}
	role="dialog"
	aria-modal="true"
	aria-labelledby="modal-title"
	tabindex="-1"
>
	<div
		bind:this={modalElement}
		class="w-full max-w-lg overflow-hidden rounded-3xl bg-white shadow-2xl"
		style="will-change: contents;"
	>
		<!-- Header -->
		<div
			class="relative overflow-hidden bg-gradient-to-r from-emerald-500 to-emerald-600 px-8 py-6 text-white"
		>
			<div class="relative z-10">
				<div class="mb-2 flex items-center gap-2">
					<Calendar size={24} />
					<h3 id="modal-title" class="text-2xl font-black">
						{activity ? 'Edit Activity' : 'New Activity'}
					</h3>
				</div>
				<p class="text-sm text-emerald-100">Track your time and maximize your ROI</p>
			</div>
			<button
				onclick={closeModal}
				class="absolute right-4 top-4 z-20 rounded-xl bg-white/20 p-2 text-white backdrop-blur-sm transition-all hover:scale-110 hover:bg-white/30"
				aria-label="Close modal"
			>
				<X size={20} />
			</button>
			<!-- Decorative circles -->
			<div class="absolute -right-8 -top-8 h-32 w-32 rounded-full bg-white/10"></div>
			<div class="absolute -bottom-6 -left-6 h-24 w-24 rounded-full bg-white/5"></div>
		</div>

		<!-- Form -->
		<form
			onsubmit={(e) => {
				e.preventDefault();
				handleSubmit();
			}}
			class="space-y-6 p-8"
		>
			<!-- Title -->
			<div>
				<label
					for="title"
					class="mb-2 flex items-center gap-2 text-sm font-bold uppercase tracking-wide text-slate-700"
				>
					<Tag size={16} />
					Activity Title <span class="text-red-500">*</span>
				</label>
				<input
					id="title"
					type="text"
					bind:value={formData.title}
					class="w-full rounded-xl border-2 border-slate-200 bg-slate-50 px-4 py-3 text-slate-900 transition-all placeholder:text-slate-400 focus:border-emerald-500 focus:bg-white focus:outline-none focus:ring-4 focus:ring-emerald-500/20"
					placeholder="e.g., Deep Work Session"
				/>
				{#if errors.title}
					<p class="mt-2 flex items-center gap-1 text-sm font-medium text-red-500">
						⚠️ {errors.title}
					</p>
				{/if}
			</div>

			<!-- Time -->
			<div>
				<label
					class="mb-2 flex items-center gap-2 text-sm font-bold uppercase tracking-wide text-slate-700"
				>
					<Clock size={16} />
					Time Range
				</label>
				<div class="grid grid-cols-2 gap-4">
					<div>
						<label for="startTime" class="mb-1 block text-xs font-semibold text-slate-500"
							>Start Time</label
						>
						<input
							id="startTime"
							type="time"
							bind:value={formData.startTime}
							class="w-full rounded-xl border-2 border-slate-200 bg-slate-50 px-4 py-3 text-slate-900 transition-all focus:border-emerald-500 focus:bg-white focus:outline-none focus:ring-4 focus:ring-emerald-500/20"
						/>
					</div>
					<div>
						<label for="endTime" class="mb-1 block text-xs font-semibold text-slate-500"
							>End Time</label
						>
						<input
							id="endTime"
							type="time"
							bind:value={formData.endTime}
							class="w-full rounded-xl border-2 border-slate-200 bg-slate-50 px-4 py-3 text-slate-900 transition-all focus:border-emerald-500 focus:bg-white focus:outline-none focus:ring-4 focus:ring-emerald-500/20"
						/>
					</div>
				</div>
				{#if errors.time}
					<p class="mt-2 flex items-center gap-1 text-sm font-medium text-red-500">
						⚠️ {errors.time}
					</p>
				{/if}
				{#if errors.overlap}
					<p class="mt-2 flex items-center gap-1 text-sm font-medium text-red-500">
						⚠️ {errors.overlap}
					</p>
				{/if}
			</div>

			<!-- Category -->
			<div>
				<label
					class="mb-3 flex items-center gap-2 text-sm font-bold uppercase tracking-wide text-slate-700"
				>
					<Target size={16} />
					Category
				</label>
				<div class="grid grid-cols-3 gap-2">
					{#each categories as cat}
						<button
							type="button"
							onclick={() => (formData.category = cat.name)}
							class="flex flex-col items-center gap-2 rounded-xl border-2 p-3 transition-all {formData.category ===
							cat.name
								? 'border-emerald-500 bg-emerald-50 shadow-lg shadow-emerald-500/20'
								: 'border-slate-200 bg-white hover:border-slate-300 hover:bg-slate-50'}"
						>
							<span class="text-2xl">{cat.icon}</span>
							<span
								class="text-xs font-bold {formData.category === cat.name
									? 'text-emerald-700'
									: 'text-slate-600'}">{cat.name}</span
							>
						</button>
					{/each}
				</div>
			</div>

			<!-- ROI -->
			<div>
				<label
					class="mb-3 flex items-center gap-2 text-sm font-bold uppercase tracking-wide text-slate-700"
				>
					<Target size={16} />
					ROI Impact
				</label>
				<div class="space-y-2">
					{#each roiOptions as roi}
						<button
							type="button"
							onclick={() => (formData.roi = roi.value)}
							class="group relative flex w-full items-center gap-4 overflow-hidden rounded-xl border-2 p-4 text-left transition-all {formData.roi ===
							roi.value
								? roi.value === 'High'
									? 'border-emerald-500 bg-gradient-to-r from-emerald-50 to-emerald-100 shadow-lg shadow-emerald-500/20'
									: roi.value === 'Neutral'
										? 'border-slate-400 bg-gradient-to-r from-slate-50 to-slate-100 shadow-lg shadow-slate-400/20'
										: 'border-red-500 bg-gradient-to-r from-red-50 to-red-100 shadow-lg shadow-red-500/20'
								: 'border-slate-200 bg-white hover:border-slate-300 hover:bg-slate-50'}"
						>
							<span class="text-2xl">{roi.icon}</span>
							<div class="flex-1">
								<p
									class="font-bold {formData.roi === roi.value
										? roi.value === 'High'
											? 'text-emerald-700'
											: roi.value === 'Neutral'
												? 'text-slate-700'
												: 'text-red-700'
										: 'text-slate-700'}"
								>
									{roi.label}
								</p>
								<p class="text-xs text-slate-500">{roi.description}</p>
							</div>
							{#if formData.roi === roi.value}
								<div
									class="flex h-6 w-6 items-center justify-center rounded-full {roi.value === 'High'
										? 'bg-emerald-500'
										: roi.value === 'Neutral'
											? 'bg-slate-500'
											: 'bg-red-500'} text-white"
								>
									✓
								</div>
							{/if}
						</button>
					{/each}
				</div>
			</div>

			<!-- Delete Button (for editing) -->
			{#if activity}
				<div class="pt-4">
					<button
						type="button"
						onclick={handleDelete}
						class="w-full text-center text-sm font-medium text-red-600 transition-colors hover:text-red-700 hover:underline"
					>
						Delete Activity
					</button>
				</div>
			{/if}

			<!-- Actions -->
			<div class="flex gap-3 pt-4">
				<button
					type="button"
					onclick={closeModal}
					class="flex-1 rounded-xl border-2 border-slate-200 bg-white px-6 py-3 font-bold text-slate-700 transition-all hover:bg-slate-50"
				>
					Cancel
				</button>
				<button
					type="submit"
					class="flex-1 rounded-xl bg-gradient-to-r from-emerald-500 to-emerald-600 px-6 py-3 font-bold text-white shadow-lg shadow-emerald-500/30 transition-all hover:scale-105 hover:shadow-xl"
				>
					{activity ? '💾 Save Changes' : '✨ Create Activity'}
				</button>
			</div>
		</form>
	</div>
</div>
