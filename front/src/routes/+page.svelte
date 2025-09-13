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
	let loading = false;

	onMount(async () => {
		await loadEntries();
	});

	async function loadEntries() {
		if ($auth.isAuthenticated && $auth.token) {
			loading = true;
			try {
				const entriesList = await entriesApi.list($auth.token, {
					status:
						selectedStatus !== "all" ? selectedStatus : undefined,
				});
				entries.setEntries(entriesList);
			} catch (error) {
				console.error("Failed to load entries:", error);
			} finally {
				loading = false;
			}
		} else {
			// Load from local storage for guest mode
			const store = storage.getMediaStore();
			let selectedEntries = [];
			for (const entry of store.entries) {
				if (
					selectedStatus === "all" ||
					entry.status === selectedStatus
				) {
					selectedEntries.push(entry);
				}
			}
			entries.setEntries(selectedEntries);
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

	<!-- Entries List -->
	{#if loading}
		<div class="flex justify-center py-12">
			<div
				class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"
			/>
		</div>
	{:else if $entries.entries.length === 0}
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
			{#each $entries.entries as entry (entry.id)}
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
