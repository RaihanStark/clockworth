<script lang="ts">
	import type { Activity } from '$lib/types';
	import { settings } from '$lib/stores/settings';
	import { TrendingUp, TrendingDown, Clock, Target, Calendar } from '@lucide/svelte';

	interface Props {
		activities: Activity[];
	}

	let { activities }: Props = $props();

	interface ROIBreakdown {
		High: number;
		Neutral: number;
		Waste: number;
	}

	const netTimeWorth = $derived.by(() => {
		let total = 0;
		const hourlyRate = $settings.hourlyWorth;

		for (const activity of activities) {
			const duration = getDuration(activity);
			const multiplier = getMultiplier(activity.roi);
			total += duration * hourlyRate * multiplier;
		}

		return total;
	});

	const roiBreakdown = $derived.by(() => {
		const breakdown: ROIBreakdown = { High: 0, Neutral: 0, Waste: 0 };

		for (const activity of activities) {
			const duration = getDuration(activity);
			breakdown[activity.roi] += duration;
		}

		return breakdown;
	});

	const totalHours = $derived(roiBreakdown.High + roiBreakdown.Neutral + roiBreakdown.Waste);

	// Calendar data - static mock data for visualization
	const currentDate = new Date();
	const currentMonth = currentDate.getMonth();
	const currentYear = currentDate.getFullYear();

	// Mock data: status for each day (green = positive, red = negative, gray = neutral/no data)
	const mockDayStatus: Record<string, 'positive' | 'negative' | 'neutral'> = {
		'2025-09-04': 'positive',
		'2025-09-05': 'positive',
		'2025-09-06': 'negative',
		'2025-09-07': 'positive',
		'2025-09-08': 'neutral',
		'2025-09-09': 'neutral',
		'2025-09-10': 'neutral',
		'2025-09-11': 'positive',
		'2025-09-12': 'positive',
		'2025-09-13': 'negative',
		'2025-09-14': 'positive',
		'2025-09-15': 'positive',
		'2025-09-16': 'neutral',
		'2025-09-17': 'neutral',
		'2025-09-18': 'positive',
		'2025-09-19': 'negative',
		'2025-09-20': 'positive',
		'2025-09-21': 'positive',
		'2025-09-22': 'positive',
		'2025-09-23': 'neutral',
		'2025-09-24': 'neutral',
		'2025-09-25': 'positive',
		'2025-09-26': 'negative',
		'2025-09-27': 'positive',
		'2025-09-28': 'positive',
		'2025-09-29': 'neutral',
		'2025-09-30': 'neutral',
		'2025-10-01': 'positive',
		'2025-10-02': 'positive',
		'2025-10-03': 'positive'
	};

	// Build last 30 days calendar grid
	const calendarDays: Array<{
		day: number;
		month: string;
		status: string | null;
		isToday?: boolean;
		dateKey: string;
	}> = [];

	// Calculate the date 30 days ago
	const thirtyDaysAgo = new Date(currentDate);
	thirtyDaysAgo.setDate(currentDate.getDate() - 29); // -29 to include today

	// Get the day of week for 30 days ago to align the grid
	const startDayOfWeek = thirtyDaysAgo.getDay();

	// Add empty cells to align with the correct day of week
	for (let i = 0; i < startDayOfWeek; i++) {
		const emptyDate = new Date(thirtyDaysAgo);
		emptyDate.setDate(thirtyDaysAgo.getDate() - (startDayOfWeek - i));
		const dateKey = `${emptyDate.getFullYear()}-${String(emptyDate.getMonth() + 1).padStart(2, '0')}-${String(emptyDate.getDate()).padStart(2, '0')}`;
		calendarDays.push({
			day: emptyDate.getDate(),
			month: emptyDate.toLocaleDateString('en-US', { month: 'short' }),
			status: null,
			dateKey
		});
	}

	// Add the last 30 days
	for (let i = 0; i < 30; i++) {
		const date = new Date(thirtyDaysAgo);
		date.setDate(thirtyDaysAgo.getDate() + i);
		const dateKey = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`;
		const isToday = date.toDateString() === currentDate.toDateString();

		calendarDays.push({
			day: date.getDate(),
			month: date.toLocaleDateString('en-US', { month: 'short' }),
			status: mockDayStatus[dateKey] || 'neutral',
			isToday,
			dateKey
		});
	}

	function getDuration(activity: Activity): number {
		const start = timeToMinutes(activity.startTime);
		const end = timeToMinutes(activity.endTime);
		return (end - start) / 60;
	}

	function timeToMinutes(time: string): number {
		const [h, m] = time.split(':').map(Number);
		return h * 60 + m;
	}

	function getMultiplier(roi: Activity['roi']): number {
		switch (roi) {
			case 'High':
				return 1.0;
			case 'Neutral':
				return 0.0;
			case 'Waste':
				return -1.0;
		}
	}

	function getPercentage(hours: number): number {
		return totalHours > 0 ? (hours / totalHours) * 100 : 0;
	}

	function getDayColor(status: string | null) {
		switch (status) {
			case 'positive':
				return 'bg-emerald-500 hover:bg-emerald-600';
			case 'negative':
				return 'bg-red-500 hover:bg-red-600';
			case 'neutral':
			default:
				return 'bg-slate-300 hover:bg-slate-400';
		}
	}
</script>

<div
	class="from-card to-accent/20 flex h-full flex-col gap-4 overflow-y-auto rounded-xl bg-gradient-to-br px-1"
>
	<!-- Net Time Worth - Hero Card -->
	<div
		class="relative overflow-hidden rounded-2xl bg-gradient-to-br from-emerald-500 to-emerald-600 p-6 text-white shadow-lg {netTimeWorth <
		0
			? 'from-red-500 to-red-600'
			: ''}"
	>
		<div class="relative z-10">
			<div class="mb-2 flex items-center gap-2">
				{#if netTimeWorth >= 0}
					<TrendingUp size={20} class="animate-bounce" />
				{:else}
					<TrendingDown size={20} class="animate-bounce" />
				{/if}
				<span class="text-sm font-medium uppercase tracking-wide opacity-90">Net Time Worth</span>
			</div>
			<p class="text-5xl font-black tracking-tight">
				${netTimeWorth >= 0 ? '+' : ''}{Math.abs(netTimeWorth).toFixed(0)}
			</p>
			<p class="mt-2 text-sm opacity-80">
				@ ${$settings.hourlyWorth}/hr · {totalHours.toFixed(1)}h tracked
			</p>
		</div>
		<!-- Decorative background element -->
		<div class="absolute -right-4 -top-4 h-32 w-32 rounded-full bg-white/10"></div>
		<div class="absolute -bottom-6 -left-6 h-40 w-40 rounded-full bg-white/5"></div>
	</div>

	<!-- ROI Breakdown Cards Grid -->
	<div class="grid grid-cols-3 gap-3">
		<!-- High ROI -->
		<div
			class="group relative overflow-hidden rounded-xl bg-emerald-500/10 p-4 transition-all hover:scale-105 hover:bg-emerald-500/20"
		>
			<div class="absolute right-2 top-2 h-8 w-8 rounded-full bg-emerald-500/20"></div>
			<Target size={18} class="mb-2 text-emerald-600" />
			<p class="text-2xl font-bold text-emerald-700">{roiBreakdown.High.toFixed(1)}</p>
			<p class="text-xs font-medium text-emerald-600">High ROI</p>
			<div class="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-emerald-200">
				<div
					class="h-full bg-emerald-600 transition-all"
					style="width: {getPercentage(roiBreakdown.High)}%"
				></div>
			</div>
		</div>

		<!-- Neutral -->
		<div
			class="group relative overflow-hidden rounded-xl bg-slate-500/10 p-4 transition-all hover:scale-105 hover:bg-slate-500/20"
		>
			<div class="absolute right-2 top-2 h-8 w-8 rounded-full bg-slate-500/20"></div>
			<Clock size={18} class="mb-2 text-slate-600" />
			<p class="text-2xl font-bold text-slate-700">{roiBreakdown.Neutral.toFixed(1)}</p>
			<p class="text-xs font-medium text-slate-600">Neutral</p>
			<div class="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-slate-200">
				<div
					class="h-full bg-slate-600 transition-all"
					style="width: {getPercentage(roiBreakdown.Neutral)}%"
				></div>
			</div>
		</div>

		<!-- Waste -->
		<div
			class="group relative overflow-hidden rounded-xl bg-red-500/10 p-4 transition-all hover:scale-105 hover:bg-red-500/20"
		>
			<div class="absolute right-2 top-2 h-8 w-8 rounded-full bg-red-500/20"></div>
			<TrendingDown size={18} class="mb-2 text-red-600" />
			<p class="text-2xl font-bold text-red-700">{roiBreakdown.Waste.toFixed(1)}</p>
			<p class="text-xs font-medium text-red-600">Waste</p>
			<div class="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-red-200">
				<div
					class="h-full bg-red-600 transition-all"
					style="width: {getPercentage(roiBreakdown.Waste)}%"
				></div>
			</div>
		</div>
	</div>

	<!-- Visual Bar -->
	<div class="bg-muted relative overflow-hidden rounded-xl p-4">
		<p class="text-foreground mb-3 text-sm font-bold uppercase tracking-wide">Distribution</p>
		<div class="h-12 overflow-hidden rounded-lg shadow-inner">
			{#if totalHours > 0}
				<div class="flex h-full">
					{#if roiBreakdown.High > 0}
						<div
							class="bg-gradient-to-b from-emerald-400 to-emerald-600 transition-all hover:opacity-80"
							style="width: {getPercentage(roiBreakdown.High)}%"
							title="High: {roiBreakdown.High.toFixed(1)}h"
						></div>
					{/if}
					{#if roiBreakdown.Neutral > 0}
						<div
							class="bg-gradient-to-b from-slate-300 to-slate-500 transition-all hover:opacity-80"
							style="width: {getPercentage(roiBreakdown.Neutral)}%"
							title="Neutral: {roiBreakdown.Neutral.toFixed(1)}h"
						></div>
					{/if}
					{#if roiBreakdown.Waste > 0}
						<div
							class="bg-gradient-to-b from-red-400 to-red-600 transition-all hover:opacity-80"
							style="width: {getPercentage(roiBreakdown.Waste)}%"
							title="Waste: {roiBreakdown.Waste.toFixed(1)}h"
						></div>
					{/if}
				</div>
			{:else}
				<div class="bg-muted flex h-full items-center justify-center">
					<p class="text-muted-foreground text-sm">No data</p>
				</div>
			{/if}
		</div>
	</div>

	<!-- Last 30 Days Calendar Heatmap -->
	<div class="rounded-xl bg-white p-5 shadow-md">
		<div class="mb-4 flex items-center justify-between">
			<div class="flex items-center gap-2">
				<Calendar size={18} class="text-slate-600" />
				<h3 class="text-sm font-bold uppercase tracking-wide text-slate-700">
					Last 30 Days
				</h3>
			</div>
		</div>

		<!-- Day labels -->
		<div class="mb-2 grid grid-cols-7 gap-1">
			{#each ['S', 'M', 'T', 'W', 'T', 'F', 'S'] as dayLabel}
				<div class="text-center text-xs font-bold text-slate-500">{dayLabel}</div>
			{/each}
		</div>

		<!-- Calendar grid -->
		<div class="grid grid-cols-7 gap-1">
			{#each calendarDays as { day, month, status, isToday }}
				{#if status === null}
					<div class="aspect-square opacity-30">
						<div class="flex h-full items-center justify-center rounded-lg bg-slate-100">
							<span class="text-xs font-medium text-slate-400">{day}</span>
						</div>
					</div>
				{:else}
					<button
						class="group relative aspect-square rounded-lg transition-all {getDayColor(
							status
						)} {isToday
							? 'ring-3 scale-105 shadow-md shadow-blue-500/40 ring-blue-500'
							: 'shadow-sm'}"
						title="{month} {day} - {status === 'positive'
							? 'Positive Net Worth'
							: status === 'negative'
								? 'Negative Net Worth'
								: 'No Data'}"
					>
						<span
							class="flex h-full items-center justify-center text-xs font-bold text-white {isToday
								? 'scale-110'
								: ''}"
						>
							{day}
						</span>
					</button>
				{/if}
			{/each}
		</div>

		<!-- Legend -->
		<div class="mt-4 flex items-center justify-center gap-4 text-xs">
			<div class="flex items-center gap-1.5">
				<div class="h-3 w-3 rounded bg-emerald-500"></div>
				<span class="font-medium text-slate-600">Positive</span>
			</div>
			<div class="flex items-center gap-1.5">
				<div class="h-3 w-3 rounded bg-slate-300"></div>
				<span class="font-medium text-slate-600">Neutral</span>
			</div>
			<div class="flex items-center gap-1.5">
				<div class="h-3 w-3 rounded bg-red-500"></div>
				<span class="font-medium text-slate-600">Negative</span>
			</div>
		</div>
	</div>

</div>
