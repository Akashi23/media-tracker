<script lang="ts">
	import { onMount } from "svelte";
	import { auth } from "$stores/auth";
	import { entries } from "$stores/entries";
	import { entriesApi, authApi } from "$utils/api";
	import { storage } from "$utils/storage";
	import EntryCard from "$components/EntryCard.svelte";
	import AddEntryDialog from "$components/AddEntryDialog.svelte";
	import EditEntryDialog from "$components/EditEntryDialog.svelte";
	import StatusTabs from "$components/StatusTabs.svelte";
	import type { Status, Entry } from "$types";

	let showAddDialog = false;
	let showEditDialog = false;
	let selectedEntry: Entry | null = null;
	let selectedStatus: Status | "all" = "all";
	let selectedGenre: string | "all" = "all";
	let availableGenres: string[] = [];
	let loading = false;
	let filteredEntries: Entry[] = [];

	// Extract available genres from entries
	$: {
		const genreSet = new Set<string>();
		($entries.entries || []).forEach((entry) => {
			if (entry.media?.genres) {
				entry.media.genres.forEach((genre) => {
					if (genre.trim()) {
						genreSet.add(genre.trim());
					}
				});
			}
		});
		availableGenres = Array.from(genreSet).sort();
	}

	// Reactive filtered entries based on selected status and genre
	$: {
		console.log(
			"Entries store updated:",
			$entries.entries.length,
			"entries"
		);
		console.log("All entries:", $entries.entries);
		filteredEntries = ($entries.entries || []).filter((entry) => {
			const statusMatch =
				selectedStatus === "all" || entry.status === selectedStatus;
			const genreMatch =
				selectedGenre === "all" ||
				(entry.media?.genres &&
					entry.media.genres.includes(selectedGenre));
			return statusMatch && genreMatch;
		});
		console.log(
			"Filtered entries:",
			filteredEntries.length,
			"for status:",
			selectedStatus,
			"and genre:",
			selectedGenre
		);
	}

	onMount(async () => {
		// Initialize auth state from localStorage
		console.log("Initializing auth from localStorage...");
		auth.init();
		console.log(
			"Auth initialized:",
			$auth.isAuthenticated,
			"Token:",
			!!$auth.token
		);

		// If user is authenticated, sync local entries with backend
		if ($auth.isAuthenticated && $auth.token) {
			await syncLocalEntries();
		}

		await loadEntries();
	});

	// Reload entries when status changes (for backend filtering if needed)
	$: if (selectedStatus) {
		// Note: We now do client-side filtering, so this is just for consistency
		// The filteredEntries reactive statement above handles the actual filtering
	}

	async function loadEntries() {
		console.log(
			"Loading entries...",
			"Authenticated:",
			$auth.isAuthenticated,
			"Token:",
			!!$auth.token
		);

		if ($auth.isAuthenticated && $auth.token) {
			loading = true;
			try {
				// Always load all entries from backend, filtering will be done in the UI
				const entriesList = await entriesApi.list($auth.token);
				console.log("Loaded entries from backend:", entriesList.length);
				entries.setEntries(entriesList);
			} catch (error) {
				console.error("Failed to load entries:", error);
			} finally {
				loading = false;
			}
		} else {
			// Load from local storage for guest mode
			const store = storage.getMediaStore();
			console.log(
				"Loaded entries from local storage:",
				store.entries.length
			);
			entries.setEntries(store.entries);
			console.log(
				"Set entries in store, current store entries:",
				$entries.entries.length
			);
		}
	}

	async function syncLocalEntries() {
		try {
			console.log("Syncing local entries with backend...");

			// Get local entries from storage
			const store = storage.getMediaStore();
			const localEntries = store.entries;

			if (localEntries.length === 0) {
				console.log("No local entries to sync");
				return;
			}

			console.log("Found", localEntries.length, "local entries to sync");

			// Convert local entries to SyncEntryRequest format
			const entriesToSync = localEntries.map((entry) => ({
				media: {
					type: entry.media?.type || "video",
					title: entry.media?.title || "Unknown Title",
					original_title: entry.media?.original_title,
					year: entry.media?.year,
					cover_url: entry.media?.cover_url,
					creators: entry.media?.creators,
					genres: entry.media?.genres || [],
					duration: entry.media?.duration,
					metadata: entry.media?.metadata,
				},
				status: entry.status,
				rating: entry.rating,
				review_md: entry.review_md,
				progress: entry.progress,
				started_at: entry.started_at,
				finished_at: entry.finished_at,
			}));

			// Sync with server
			const response = await entriesApi.sync(entriesToSync, $auth.token!);
			console.log("Sync completed:", response.count, "entries synced");

			if (response.errors && response.errors.length > 0) {
				console.warn("Sync errors:", response.errors);
			}

			// Clear local storage after successful sync
			storage.clearEntries();
			console.log("Local entries cleared after sync");
		} catch (error) {
			console.error("Failed to sync local entries:", error);
		}
	}

	function handleEditEntry(entry: Entry) {
		selectedEntry = entry;
		showEditDialog = true;
	}

	function handleDeleteEntry(id: string) {
		storage.removeEntry(id); // remove from localStorage
		loadEntries();
	}

	$: if (selectedStatus) {
		loadEntries();
	}
</script>

<svelte:head>
	<title>Dashboard - Media Tracker</title>
</svelte:head>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex justify-between items-center">
		<div>
			<h1 class="text-3xl font-bold text-gray-900">
				{#if $auth.isAuthenticated}
					Welcome back, {$auth.user?.name || "User"}!
				{:else}
					Your Media Tracker
				{/if}
			</h1>
			<p class="text-gray-600 mt-1">
				{#if $auth.isGuest}
					Guest mode - your data is stored locally
				{:else}
					Track your media consumption
				{/if}
			</p>
		</div>

		<button class="btn btn-primary" on:click={() => (showAddDialog = true)}>
			+ Add Entry
		</button>
	</div>

	<!-- Status Tabs -->
	<StatusTabs bind:selectedStatus />

	<!-- Genre Filter -->
	{#if availableGenres.length > 0}
		<div class="mb-6">
			<div class="flex flex-wrap gap-2">
				<button
					on:click={() => (selectedGenre = "all")}
					class="px-3 py-1 text-sm font-medium rounded-full transition-colors {selectedGenre ===
					'all'
						? 'bg-blue-500 text-white'
						: 'bg-gray-200 text-gray-700 hover:bg-gray-300'}"
				>
					All Genres
				</button>
				{#each availableGenres as genre}
					<button
						on:click={() => (selectedGenre = genre)}
						class="px-3 py-1 text-sm font-medium rounded-full transition-colors {selectedGenre ===
						genre
							? 'bg-blue-500 text-white'
							: 'bg-gray-200 text-gray-700 hover:bg-gray-300'}"
					>
						{genre}
					</button>
				{/each}
			</div>
		</div>
	{/if}

	<!-- Entries List -->
	{#if loading}
		<div class="flex justify-center py-12">
			<div
				class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"
			/>
		</div>
	{:else if !filteredEntries || filteredEntries.length === 0}
		<div class="text-center py-12">
			<div class="text-gray-400 text-6xl mb-4">📚</div>
			<h3 class="text-lg font-medium text-gray-900 mb-2">
				No entries yet
			</h3>
			<p class="text-gray-600 mb-4">
				Start tracking your media consumption by adding your first
				entry.
			</p>
			<button
				class="btn btn-primary"
				on:click={() => (showAddDialog = true)}
			>
				Add Your First Entry
			</button>
		</div>
	{:else}
		<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
			{#each filteredEntries as entry (entry.id)}
				<EntryCard
					{entry}
					on:edit={(e) => handleEditEntry(e.detail)}
					on:delete={(e) => handleDeleteEntry(e.detail)}
				/>
			{/each}
		</div>
	{/if}
</div>

{#if showAddDialog}
	<AddEntryDialog
		bind:open={showAddDialog}
		on:entry-added={() => {
			showAddDialog = false;
			loadEntries();
		}}
	/>
{/if}

{#if showEditDialog}
	<EditEntryDialog
		bind:open={showEditDialog}
		entry={selectedEntry}
		on:entry-updated={() => {
			showEditDialog = false;
			selectedEntry = null;
			loadEntries();
		}}
	/>
{/if}
