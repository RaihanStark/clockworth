<script lang="ts">
	import Timeline from '$lib/components/Timeline.svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import ActivityModal from '$lib/components/ActivityModal.svelte';
	import { Plus } from '@lucide/svelte';
	import type { Activity } from '$lib/types';
	import type { PageData } from './$types';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	let activities = $state<Activity[]>(data.activities);
	let selectedActivity = $state<Activity | null>(null);
	let showModal = $state(false);
	let initialHour = $state(9);
	let initialStartTime = $state<string | undefined>(undefined);
	let initialEndTime = $state<string | undefined>(undefined);

	function handleActivityClick(activity: Activity) {
		selectedActivity = activity;
		showModal = true;
	}

	function handleAddActivity(hour: number, startTime?: string, endTime?: string) {
		selectedActivity = null;
		initialHour = hour;
		initialStartTime = startTime;
		initialEndTime = endTime;
		showModal = true;
	}

	function handleCloseModal() {
		showModal = false;
		selectedActivity = null;
	}

	function handleActivityAdded(activity: Activity) {
		activities = [...activities, activity];
	}

	function handleActivityUpdated(activity: Activity) {
		activities = activities.map((a) => (a.id === activity.id ? activity : a));
	}

	function handleActivityDeleted(id: string) {
		activities = activities.filter((a) => a.id !== id);
	}

	const currentDate = new Date();
	const dayOfWeek = currentDate.toLocaleDateString('en-US', { weekday: 'long' });
	const monthDay = currentDate.toLocaleDateString('en-US', { month: 'long', day: 'numeric' });
</script>

<div class="bg-background h-screen overflow-hidden">
	<!-- Header -->
	<header
		class="border-border border-b bg-gradient-to-r from-slate-50 to-slate-100 px-4 py-4 shadow-sm sm:px-6 sm:py-5"
	>
		<div class="flex items-center justify-between gap-3">
			<!-- Left: Date & Title -->
			<div class="flex items-center gap-3 sm:gap-4">
				<div
					class="flex h-12 w-12 flex-col items-center justify-center rounded-xl bg-gradient-to-br from-emerald-500 to-emerald-600 text-white shadow-lg shadow-emerald-500/30 sm:h-14 sm:w-14 sm:rounded-2xl"
				>
					<span class="text-[10px] font-bold uppercase sm:text-xs"
						>{currentDate.toLocaleDateString('en-US', { month: 'short' })}</span
					>
					<span class="text-lg font-black sm:text-xl">{currentDate.getDate()}</span>
				</div>
				<div>
					<h1 class="text-xl font-black tracking-tight text-slate-900 sm:text-2xl lg:text-3xl">
						{dayOfWeek}
					</h1>
					<p class="text-xs font-medium text-slate-500 sm:text-sm">
						{monthDay}, {currentDate.getFullYear()}
					</p>
				</div>
			</div>

			<!-- Right: Actions -->
			<button
				onclick={() => handleAddActivity(9)}
				class="flex items-center gap-2 rounded-xl bg-gradient-to-r from-emerald-500 to-emerald-600 px-4 py-2.5 font-bold text-white shadow-lg shadow-emerald-500/30 transition-all hover:scale-105 hover:shadow-xl sm:px-5 lg:hidden"
			>
				<Plus size={20} />
				<span class="hidden sm:inline">New Activity</span>
			</button>
			<!-- Desktop only - larger button -->
			<button
				onclick={() => handleAddActivity(9)}
				class="hidden items-center gap-2 rounded-xl bg-gradient-to-r from-emerald-500 to-emerald-600 px-6 py-3 font-bold text-white shadow-lg shadow-emerald-500/30 transition-all hover:scale-105 hover:shadow-xl lg:flex"
			>
				<Plus size={20} />
				<span>New Activity</span>
			</button>
		</div>
	</header>

	<!-- Content -->
	<div class="h-[calc(100vh-150px)] lg:h-[calc(100vh-96px)]">
		<!-- Large Desktop Layout: Side by Side (1536px+) -->
		<div class="hidden h-full xl:block">
			<div class="flex h-full gap-6 px-6 py-6">
				<div class="w-full overflow-hidden">
					<Timeline
						{activities}
						onActivityClick={handleActivityClick}
						onAddActivity={handleAddActivity}
						onActivityUpdate={handleActivityUpdated}
					/>
				</div>
				<div class="w-full max-w-[400px] overflow-hidden">
					<Sidebar {activities} heatmapData={data.heatmapData} />
				</div>
			</div>
		</div>
		<!-- Mobile Layout: Full Height Timeline, No Sidebar (<1024px) -->
		<div class="h-full px-4 py-4 xl:hidden">
			<Timeline
				{activities}
				onActivityClick={handleActivityClick}
				onAddActivity={handleAddActivity}
				onActivityUpdate={handleActivityUpdated}
			/>
		</div>
	</div>
</div>

<!-- FAB for mobile - only show on very small screens -->
<button
	class="fixed bottom-6 right-6 z-50 rounded-full bg-gradient-to-br from-emerald-500 to-emerald-600 p-4 text-white shadow-xl shadow-emerald-500/40 transition-transform hover:scale-110 active:scale-95 sm:hidden"
	onclick={() => handleAddActivity(9)}
	aria-label="Add activity"
>
	<Plus size={24} strokeWidth={2.5} />
</button>

<!-- Modal -->
{#if showModal}
	<ActivityModal
		activity={selectedActivity}
		{initialHour}
		{initialStartTime}
		{initialEndTime}
		onClose={handleCloseModal}
		onActivityAdded={handleActivityAdded}
		onActivityUpdated={handleActivityUpdated}
		onActivityDeleted={handleActivityDeleted}
	/>
{/if}
