<script lang="ts">
	import AuthLayout from '$lib/components/auth/AuthLayout.svelte';
	import AuthHeader from '$lib/components/auth/AuthHeader.svelte';
	import GoogleAuthButton from '$lib/components/auth/GoogleAuthButton.svelte';
	import AuthFooterLink from '$lib/components/auth/AuthFooterLink.svelte';
	import { authStore } from '$lib/stores/auth';

	let name = $state('');
	let email = $state('');
	let password = $state('');
	let loading = $state(false);
	let googleLoading = $state(false);
	let error = $state('');

	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';

		if (!email || !password) {
			error = 'Please enter both email and password.';
			return;
		}

		if (password.length < 6) {
			error = 'Password must be at least 6 characters.';
			return;
		}

		loading = true;
		try {
			await authStore.registerWithEmail(email, password, name);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Registration failed';
		} finally {
			loading = false;
		}
	}

	async function handleGoogleSignUp() {
		error = '';
		googleLoading = true;
		try {
			await authStore.login();
		} catch (err) {
			console.error('Sign up failed:', err);
			error = 'Failed to connect to Google Sign-In.';
			googleLoading = false;
		}
	}
</script>

<AuthLayout>
	<AuthHeader
		title="Create Account"
		description="Start tracking your time economics today"
	/>

	{#if error}
		<div class="mb-5 rounded-xl border border-red-200 bg-red-50 p-3.5 text-sm font-medium text-red-700">
			{error}
		</div>
	{/if}

	<form onsubmit={handleSubmit} class="space-y-4">
		<div>
			<label for="name" class="block text-xs font-semibold uppercase tracking-wider text-slate-700">
				Full Name
			</label>
			<input
				id="name"
				type="text"
				bind:value={name}
				autocomplete="name"
				placeholder="Your name"
				class="mt-1.5 block w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-900 transition-all placeholder:text-slate-400 focus:border-emerald-500 focus:bg-white focus:outline-none focus:ring-2 focus:ring-emerald-500/20"
			/>
		</div>

		<div>
			<label for="email" class="block text-xs font-semibold uppercase tracking-wider text-slate-700">
				Email Address
			</label>
			<input
				id="email"
				type="email"
				bind:value={email}
				required
				autocomplete="email"
				placeholder="you@example.com"
				class="mt-1.5 block w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-900 transition-all placeholder:text-slate-400 focus:border-emerald-500 focus:bg-white focus:outline-none focus:ring-2 focus:ring-emerald-500/20"
			/>
		</div>

		<div>
			<label for="password" class="block text-xs font-semibold uppercase tracking-wider text-slate-700">
				Password
			</label>
			<input
				id="password"
				type="password"
				bind:value={password}
				required
				minlength="6"
				autocomplete="new-password"
				placeholder="Min. 6 characters"
				class="mt-1.5 block w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-900 transition-all placeholder:text-slate-400 focus:border-emerald-500 focus:bg-white focus:outline-none focus:ring-2 focus:ring-emerald-500/20"
			/>
		</div>

		<button
			type="submit"
			disabled={loading}
			class="flex w-full items-center justify-center rounded-xl bg-emerald-600 px-6 py-3.5 font-semibold text-white shadow-sm transition-all hover:bg-emerald-700 hover:shadow-md active:scale-95 disabled:cursor-not-allowed disabled:opacity-50"
		>
			{#if loading}
				<div class="h-5 w-5 animate-spin rounded-full border-2 border-white/30 border-t-white"></div>
			{:else}
				Create Account
			{/if}
		</button>
	</form>

	<div class="relative my-6">
		<div class="absolute inset-0 flex items-center">
			<div class="w-full border-t border-slate-200"></div>
		</div>
		<div class="relative flex justify-center text-xs uppercase">
			<span class="bg-white px-3 font-semibold text-slate-400">Or continue with</span>
		</div>
	</div>

	<GoogleAuthButton
		loading={googleLoading}
		text="Sign up with Google"
		onClick={handleGoogleSignUp}
	/>

	<AuthFooterLink
		text="Already have an account?"
		linkText="Sign In"
		href="/login"
	/>
</AuthLayout>
