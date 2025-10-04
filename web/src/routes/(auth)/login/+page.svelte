<script lang="ts">
	import AuthLayout from '$lib/components/auth/AuthLayout.svelte';
	import AuthHeader from '$lib/components/auth/AuthHeader.svelte';
	import GoogleAuthButton from '$lib/components/auth/GoogleAuthButton.svelte';
	import AuthFooterLink from '$lib/components/auth/AuthFooterLink.svelte';
	import { authStore } from '$lib/stores/auth';

	let loading = false;

	async function handleGoogleSignIn() {
		loading = true;
		try {
			await authStore.login();
		} catch (error) {
			console.error('Login failed:', error);
			loading = false;
		}
	}
</script>

<AuthLayout>
	<AuthHeader
		title="Welcome Back"
		description="Sign in to continue to Clockworth"
	/>

	<GoogleAuthButton
		{loading}
		text="Sign in with Google"
		onClick={handleGoogleSignIn}
	/>

	<AuthFooterLink
		text="Don't have an account?"
		linkText="Sign Up"
		href="/signup"
	/>
</AuthLayout>
