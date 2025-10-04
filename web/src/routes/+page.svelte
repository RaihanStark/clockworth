<script lang="ts">
	import Navigation from '$lib/components/landing/Navigation.svelte';
	import Hero from '$lib/components/landing/Hero.svelte';
	import ValueOfTime from '$lib/components/landing/ValueOfTime.svelte';
	import Features from '$lib/components/landing/Features.svelte';
	import CTA from '$lib/components/landing/CTA.svelte';
	import Footer from '$lib/components/landing/Footer.svelte';
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-50 via-emerald-50 to-slate-50">
	<Navigation />
	<Hero />
	<ValueOfTime />
	<Features />
	<CTA />
	<Footer />
</div>
