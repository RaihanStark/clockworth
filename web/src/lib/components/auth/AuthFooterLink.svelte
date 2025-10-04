<script lang="ts">
	let { text, linkText, href } = $props<{
		text: string;
		linkText: string;
		href: string;
	}>();
</script>

<div class="mt-6 text-center">
	<p class="text-sm text-slate-600">
		{text}
		<a href={href} class="font-bold text-emerald-600 hover:text-emerald-700 hover:underline">
			{linkText}
		</a>
	</p>
</div>
