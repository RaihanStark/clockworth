<script lang="ts">
	import AuthLayout from '$lib/components/auth/AuthLayout.svelte';
	import AuthHeader from '$lib/components/auth/AuthHeader.svelte';
	import GoogleAuthButton from '$lib/components/auth/GoogleAuthButton.svelte';
	import AuthFooterLink from '$lib/components/auth/AuthFooterLink.svelte';
	import { authStore } from '$lib/stores/auth';

	let loading = false;

	async function handleGoogleSignUp() {
		loading = true;
		try {
			await authStore.login();
		} catch (error) {
			console.error('Sign up failed:', error);
			loading = false;
		}
	}
</script>

<AuthLayout>
	<AuthHeader
		title="Create Account"
		description="Start tracking your time economics today"
	/>

	<GoogleAuthButton
		{loading}
		text="Sign up with Google"
		onClick={handleGoogleSignUp}
	/>

	<AuthFooterLink
		text="Already have an account?"
		linkText="Sign In"
		href="/login"
	/>
</AuthLayout>
