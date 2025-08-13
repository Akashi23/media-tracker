<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { auth } from '$stores/auth';
	import { page } from '$app/stores';
	import Header from '$components/Header.svelte';
	import { authApi } from '$utils/api';

	onMount(() => {
		auth.init();
	});

	$: isAuthenticated = $auth.isAuthenticated;
	$: isGuest = $auth.isGuest;
	$: currentPath = $page.url.pathname;
</script>

<svelte:head>
	<title>Media Tracker</title>
	<meta name="description" content="Track your media consumption - movies, books, anime, games, and more" />
</svelte:head>

<div class="min-h-screen bg-gray-50">
	<Header />
	
	<main class="container mx-auto px-4 py-8">
		<slot />
	</main>
</div>
