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

<div class="h-screen overflow-hidden bg-background">
	<!-- Header -->
	<header class="border-b border-border bg-gradient-to-r from-slate-50 to-slate-100 px-6 py-5 shadow-sm">
		<div class="flex items-center justify-between">
			<!-- Left: Date & Title -->
			<div class="flex items-center gap-4">
				<div class="flex h-14 w-14 flex-col items-center justify-center rounded-2xl bg-gradient-to-br from-emerald-500 to-emerald-600 text-white shadow-lg shadow-emerald-500/30">
					<span class="text-xs font-bold uppercase">{currentDate.toLocaleDateString('en-US', { month: 'short' })}</span>
					<span class="text-xl font-black">{currentDate.getDate()}</span>
				</div>
				<div>
					<h1 class="text-3xl font-black tracking-tight text-slate-900">{dayOfWeek}</h1>
					<p class="text-sm font-medium text-slate-500">{monthDay}, {currentDate.getFullYear()}</p>
				</div>
			</div>

			<!-- Right: Actions -->
			<button
				onclick={() => handleAddActivity(9)}
				class="flex items-center gap-2 rounded-xl bg-gradient-to-r from-emerald-500 to-emerald-600 px-5 py-2.5 font-bold text-white shadow-lg shadow-emerald-500/30 transition-all hover:scale-105 hover:shadow-xl"
			>
				<Plus size={20} />
				<span class="hidden sm:inline">New Activity</span>
			</button>
		</div>
	</header>

	<!-- Content -->
	<div class="h-[calc(100vh-96px)] p-6">
		<!-- Desktop Layout: Split View -->
		<div class="hidden h-full gap-6 lg:grid lg:grid-cols-3">
			<div class="col-span-2 overflow-hidden">
				<Timeline
					{activities}
					onActivityClick={handleActivityClick}
					onAddActivity={handleAddActivity}
					onActivityUpdate={handleActivityUpdated}
				/>
			</div>
			<div class="col-span-1 overflow-hidden">
				<Sidebar {activities} heatmapData={data.heatmapData} />
			</div>
		</div>

		<!-- Mobile Layout: Stacked -->
		<div class="flex h-full flex-col gap-6 lg:hidden">
			<div class="flex-1 overflow-hidden">
				<Timeline
					{activities}
					onActivityClick={handleActivityClick}
					onAddActivity={handleAddActivity}
					onActivityUpdate={handleActivityUpdated}
				/>
			</div>
			<div class="h-[300px] overflow-hidden">
				<Sidebar {activities} heatmapData={data.heatmapData} />
			</div>
		</div>
	</div>
</div>

<!-- FAB for mobile -->
<button
	class="fixed bottom-6 right-6 rounded-full bg-primary p-4 text-primary-foreground shadow-lg transition-transform hover:scale-110 lg:hidden"
	onclick={() => handleAddActivity(9)}
	aria-label="Add activity"
>
	<Plus size={24} />
</button>

<!-- Modal -->
{#if showModal}
	<ActivityModal
		activity={selectedActivity}
		initialHour={initialHour}
		initialStartTime={initialStartTime}
		initialEndTime={initialEndTime}
		onClose={handleCloseModal}
		onActivityAdded={handleActivityAdded}
		onActivityUpdated={handleActivityUpdated}
		onActivityDeleted={handleActivityDeleted}
	/>
{/if}
