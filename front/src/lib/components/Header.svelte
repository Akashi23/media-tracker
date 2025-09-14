<script lang="ts">
	import { auth } from "$stores/auth";
	import { authApi } from "$utils/api";
	import { goto } from "$app/navigation";
	import SyncButton from "./SyncButton.svelte";

	let showLoginDialog = false;
	let email = "";
	let loading = false;

	async function handleLogin() {
		if (!email.trim()) return;

		loading = true;
		try {
			const response = await authApi.login({ email });
			const user = await authApi.getProfile(response.token);
			auth.login(user, response.token);
			showLoginDialog = false;
			email = "";
		} catch (error) {
			console.error("Login failed:", error);
			alert("Login failed. Please try again.");
		} finally {
			loading = false;
		}
	}

	function handleLogout() {
		auth.logout();
		goto("/");
	}

	function switchToGuest() {
		auth.setGuest(true);
		goto("/");
	}
</script>

<header class="bg-white shadow-sm border-b">
	<div class="container mx-auto px-4">
		<div class="flex justify-between items-center h-16">
			<!-- Logo -->
			<div class="flex items-center">
				<a href="/" class="text-xl font-bold text-primary-600">
					Media Tracker
				</a>
			</div>

			<!-- Navigation -->
			<nav class="hidden md:flex space-x-8">
				<a href="/" class="text-gray-700 hover:text-primary-600"
					>Dashboard</a
				>
				<a
					href="/collections"
					class="text-gray-700 hover:text-primary-600">Collections</a
				>
				<a href="/genres" class="text-gray-700 hover:text-primary-600"
					>Genres</a
				>
				<a href="/profile" class="text-gray-700 hover:text-primary-600"
					>Profile</a
				>
			</nav>

			<!-- Auth Section -->
			<div class="flex items-center space-x-4">
				{#if $auth.isAuthenticated}
					<SyncButton
						on:sync-success={(e) => {
							console.log("Sync success:", e.detail);
							// Success messages are now handled by the SyncButton component itself
						}}
						on:sync-error={(e) => {
							console.log("Sync error:", e.detail);
							// Error messages are now handled by the SyncButton component itself
						}}
					/>
					<span class="text-sm text-gray-600">
						{$auth.user?.email}
					</span>
					<button class="btn btn-secondary" on:click={handleLogout}>
						Logout
					</button>
				{:else if $auth.isGuest}
					<button
						class="btn btn-secondary"
						on:click={() => (showLoginDialog = true)}
					>
						Login
					</button>
				{:else}
					<button
						class="btn btn-primary"
						on:click={() => (showLoginDialog = true)}
					>
						Login
					</button>
					<button class="btn btn-secondary" on:click={switchToGuest}>
						Guest Mode
					</button>
				{/if}
			</div>
		</div>
	</div>
</header>

<!-- Login Dialog -->
{#if showLoginDialog}
	<div
		class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
	>
		<div class="bg-white rounded-lg p-6 w-full max-w-md">
			<h2 class="text-xl font-bold mb-4">Login</h2>
			<p class="text-gray-600 mb-4">
				Enter your email to login or create an account.
			</p>

			<form on:submit|preventDefault={handleLogin} class="space-y-4">
				<div>
					<label
						for="email"
						class="block text-sm font-medium text-gray-700 mb-1"
					>
						Email
					</label>
					<input
						id="email"
						type="email"
						bind:value={email}
						class="input"
						placeholder="your@email.com"
						required
					/>
				</div>

				<div class="flex space-x-3">
					<button
						type="submit"
						class="btn btn-primary flex-1"
						disabled={loading}
					>
						{loading ? "Logging in..." : "Login"}
					</button>
					<button
						type="button"
						class="btn btn-secondary"
						on:click={() => (showLoginDialog = false)}
					>
						Cancel
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
