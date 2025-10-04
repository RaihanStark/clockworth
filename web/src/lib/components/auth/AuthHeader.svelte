<script lang="ts">
	let { title, description } = $props<{
		title: string;
		description: string;
	}>();
</script>

<div class="mb-8 text-center">
	<h2 class="text-2xl font-black text-slate-900">{title}</h2>
	<p class="mt-2 text-sm text-slate-600">{description}</p>
</div>
