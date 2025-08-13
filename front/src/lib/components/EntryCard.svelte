<script lang="ts">
	import { createEventDispatcher } from "svelte";
	import type { Entry, Status } from "$types";
	import { marked } from "marked";
	import { storage } from "$lib/utils/storage";

	export let entry: Entry;

	const dispatch = createEventDispatcher();

	const statusColors = {
		planned: "bg-blue-100 text-blue-800",
		in_progress: "bg-yellow-100 text-yellow-800",
		completed: "bg-green-100 text-green-800",
		on_hold: "bg-orange-100 text-orange-800",
		dropped: "bg-red-100 text-red-800",
	};

	const typeIcons = {
		movie: "🎬",
		tv: "📺",
		anime: "🌸",
		book: "📚",
		game: "🎮",
		video: "📹",
	};

	function getStatusLabel(status: Status): string {
		const labels = {
			planned: "Planned",
			in_progress: "In Progress",
			completed: "Completed",
			on_hold: "On Hold",
			dropped: "Dropped",
		};
		return labels[status];
	}

	function formatRating(rating?: number): string {
		if (!rating) return "No rating";
		return `${rating}/10`;
	}

	function formatDate(dateStr?: string): string {
		if (!dateStr) return "";
		return new Date(dateStr).toLocaleDateString();
	}

	function handleEdit() {
		dispatch("edit", entry);
	}

	function handleDelete() {
		if (confirm("Are you sure you want to delete this entry?")) {
			dispatch("delete", entry.id);
		}
	}
</script>

<div class="card hover:shadow-lg transition-shadow">
	<!-- Header -->
	<div class="flex justify-between items-start mb-4">
		<div class="flex items-center space-x-2">
			<span class="text-2xl"
				>{typeIcons[entry.media?.type || "video"]}</span
			>
			<div>
				<h3 class="font-semibold text-lg text-gray-900">
					{entry.media?.title || "Unknown Title"}
				</h3>
				{#if entry.media?.year}
					<p class="text-sm text-gray-600">({entry.media.year})</p>
				{/if}
			</div>
		</div>

		<div class="flex space-x-2">
			<button
				class="text-gray-400 hover:text-gray-600"
				on:click={handleEdit}
			>
				✏️
			</button>
			<button
				class="text-gray-400 hover:text-red-600"
				on:click={handleDelete}
			>
				🗑️
			</button>
		</div>
	</div>

	<!-- Status and Rating -->
	<div class="flex justify-between items-center mb-3">
		<span
			class="px-2 py-1 rounded-full text-xs font-medium {statusColors[
				entry.status
			]}"
		>
			{getStatusLabel(entry.status)}
		</span>

		{#if entry.rating}
			<div class="flex items-center space-x-1">
				<span class="text-yellow-500">⭐</span>
				<span class="text-sm font-medium"
					>{formatRating(entry.rating)}</span
				>
			</div>
		{/if}
	</div>

	<!-- Progress -->
	{#if entry.progress}
		<div class="mb-3">
			{#if entry.media?.type === "book" && entry.progress.pagesRead && entry.progress.pagesTotal}
				<div class="text-sm text-gray-600">
					{entry.progress.pagesRead} / {entry.progress.pagesTotal} pages
				</div>
			{:else if (entry.media?.type === "anime" || entry.media?.type === "tv") && entry.progress.episodesSeen && entry.progress.episodesTotal}
				<div class="text-sm text-gray-600">
					{entry.progress.episodesSeen} / {entry.progress
						.episodesTotal} episodes
				</div>
			{:else if entry.media?.type === "game" && entry.progress.gamePercent}
				<div class="text-sm text-gray-600">
					{entry.progress.gamePercent}% complete
				</div>
			{/if}
		</div>
	{/if}

	<!-- Dates -->
	<div class="flex justify-between text-xs text-gray-500 mb-3">
		{#if entry.started_at}
			<span>Started: {formatDate(entry.started_at)}</span>
		{/if}
		{#if entry.finished_at}
			<span>Finished: {formatDate(entry.finished_at)}</span>
		{/if}
	</div>

	<!-- Review Preview -->
	{#if entry.review_md}
		<div class="border-t pt-3">
			<div class="text-sm text-gray-700 line-clamp-3">
				{@html marked(entry.review_md.substring(0, 150))}
				{#if entry.review_md.length > 150}...{/if}
			</div>
		</div>
	{/if}

	<!-- Genres -->
	{#if entry.media?.genres && entry.media.genres.length > 0}
		<div class="mt-3 flex flex-wrap gap-1">
			{#each entry.media.genres.slice(0, 3) as genre}
				<span
					class="px-2 py-1 bg-gray-100 text-gray-700 text-xs rounded"
				>
					{genre}
				</span>
			{/each}
			{#if entry.media.genres.length > 3}
				<span
					class="px-2 py-1 bg-gray-100 text-gray-700 text-xs rounded"
				>
					+{entry.media.genres.length - 3}
				</span>
			{/if}
		</div>
	{/if}
</div>

<style>
	.line-clamp-3 {
		display: -webkit-box;
		-webkit-line-clamp: 3;
		line-clamp: 3;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>
