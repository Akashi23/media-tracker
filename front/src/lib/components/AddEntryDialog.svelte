<script lang="ts">
	import { createEventDispatcher } from "svelte";
	import { v4 as uuidv4 } from "uuid";
	import { auth } from "$stores/auth";
	import { entriesApi, mediaApi } from "$utils/api";
	import { storage } from "$utils/storage";
	import type {
		MediaType,
		Status,
		CreateEntryRequest,
		CreateMediaRequest,
	} from "$types";

	export let open = false;

	const dispatch = createEventDispatcher();

	let title = "";
	let mediaType: MediaType = "movie";
	let status: Status = "planned";
	let rating: number | undefined;
	let reviewMd = "";
	let year: number | undefined;
	let genres: string[] = [];
	let genreInput = "";
	let loading = false;
	let searchResults: any[] = [];
	let showSearchResults = false;

	const mediaTypes: Array<{ value: MediaType; label: string; icon: string }> =
		[
			{ value: "movie", label: "Movie", icon: "🎬" },
			{ value: "tv", label: "TV Show", icon: "📺" },
			{ value: "anime", label: "Anime", icon: "🌸" },
			{ value: "book", label: "Book", icon: "📚" },
			{ value: "game", label: "Game", icon: "🎮" },
			{ value: "video", label: "Video", icon: "📹" },
		];

	const statuses: Array<{ value: Status; label: string }> = [
		{ value: "planned", label: "Planned" },
		{ value: "in_progress", label: "In Progress" },
		{ value: "completed", label: "Completed" },
		{ value: "on_hold", label: "On Hold" },
		{ value: "dropped", label: "Dropped" },
	];

	async function searchMedia() {
		if (!title.trim()) return;

		try {
			const results = await mediaApi.search(title, mediaType);
			searchResults = results;
			showSearchResults = true;
		} catch (error) {
			console.error("Search failed:", error);
		}
	}

	async function selectMedia(media: any) {
		title = media.title;
		year = media.year;
		genres = media.genres || [];
		showSearchResults = false;
		searchResults = [];
	}

	function addGenre() {
		const genre = genreInput.trim();
		if (genre && !genres.includes(genre)) {
			genres = [...genres, genre];
			genreInput = "";
		}
	}

	function removeGenre(genreToRemove: string) {
		genres = genres.filter((genre) => genre !== genreToRemove);
	}

	function handleGenreKeydown(event: KeyboardEvent) {
		if (event.key === "Enter") {
			event.preventDefault();
			addGenre();
		}
	}

	async function handleSubmit() {
		if (!title.trim()) return;

		loading = true;
		try {
			let mediaId: string;
			let mediaData: CreateMediaRequest;

			mediaData = {
				type: mediaType,
				title,
				year,
				genres: genres,
			};

			// Check if media exists, if not create it
			if (searchResults && searchResults.length > 0) {
				mediaId = searchResults[0].id;
			} else {
				if ($auth.isAuthenticated && $auth.token) {
					const media = await mediaApi.create(mediaData, $auth.token);
					mediaId = media.id;
				} else {
					// Guest mode - create locally
					const media = {
						id: uuidv4(),
						...mediaData,
						created_at: new Date().toISOString(),
					};
					mediaId = media.id;
				}
			}

			// Create entry
			const entryData: CreateEntryRequest = {
				media_id: mediaId,
				status,
				rating,
				review_md: reviewMd || undefined,
			};

			if ($auth.isAuthenticated && $auth.token) {
				const entry = await entriesApi.create(entryData, $auth.token);
				dispatch("entry-added", entry);
			} else {
				// Guest mode - create locally
				const entry = {
					id: uuidv4(),
					user_id: storage.getGuestId(),
					...entryData,
					updated_at: new Date().toISOString(),
				};
				const media = {
					id: uuidv4(),
					...mediaData,
					created_at: new Date().toISOString(),
				};
				mediaId = media.id;
				storage.addEntry(entry, media);
				dispatch("entry-added", entry);
			}

			// Reset form
			title = "";
			rating = undefined;
			reviewMd = "";
			year = undefined;
			genres = [];
			genreInput = "";
			open = false;
		} catch (error) {
			console.error("Failed to create entry:", error);
			alert("Failed to create entry. Please try again.");
		} finally {
			loading = false;
		}
	}

	function handleClose() {
		open = false;
		title = "";
		rating = undefined;
		reviewMd = "";
		year = undefined;
		genres = [];
		genreInput = "";
		showSearchResults = false;
		searchResults = [];
	}

	$: if (title && title.length > 2) {
		searchMedia();
	}
</script>

{#if open}
	<div
		class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
	>
		<div
			class="bg-white rounded-lg p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto"
		>
			<div class="flex justify-between items-center mb-6">
				<h2 class="text-xl font-bold">Add New Entry</h2>
				<button
					class="text-gray-400 hover:text-gray-600"
					on:click={handleClose}
				>
					✕
				</button>
			</div>

			<form on:submit|preventDefault={handleSubmit} class="space-y-4">
				<!-- Media Type -->
				<fieldset>
					<legend
						class="block text-sm font-medium text-gray-700 mb-2"
					>
						Media Type
					</legend>

					<div class="grid grid-cols-3 gap-2">
						{#each mediaTypes as type}
							<button
								type="button"
								class="p-3 border rounded-lg text-center {mediaType ===
								type.value
									? 'border-primary-500 bg-primary-50'
									: 'border-gray-300 hover:border-gray-400'}"
								on:click={() => (mediaType = type.value)}
							>
								<div class="text-2xl mb-1">{type.icon}</div>
								<div class="text-sm">{type.label}</div>
							</button>
						{/each}
					</div>
				</fieldset>

				<!-- Title -->
				<div>
					<label
						for="title"
						class="block text-sm font-medium text-gray-700 mb-1"
					>
						Title
					</label>
					<input
						id="title"
						type="text"
						bind:value={title}
						class="input"
						placeholder="Enter title..."
						required
					/>
				</div>

				<!-- Search Results -->
				{#if showSearchResults && searchResults && searchResults.length > 0}
					<div class="border rounded-lg p-2 max-h-32 overflow-y-auto">
						<p class="text-sm text-gray-600 mb-2">
							Search results:
						</p>
						{#each searchResults as result}
							<button
								type="button"
								class="w-full text-left p-2 hover:bg-gray-100 rounded"
								on:click={() => selectMedia(result)}
							>
								<div class="font-medium">{result.title}</div>
								{#if result.year}
									<div class="text-sm text-gray-600">
										({result.year})
									</div>
								{/if}
							</button>
						{/each}
					</div>
				{/if}

				<!-- Year -->
				<div>
					<label
						for="year"
						class="block text-sm font-medium text-gray-700 mb-1"
					>
						Year (optional)
					</label>
					<input
						id="year"
						type="number"
						bind:value={year}
						class="input"
						placeholder="2024"
						min="1900"
						max="2030"
					/>
				</div>

				<!-- Status -->
				<div>
					<label
						for="status"
						class="block text-sm font-medium text-gray-700 mb-1"
					>
						Status
					</label>
					<select id="status" bind:value={status} class="input">
						{#each statuses as statusOption}
							<option value={statusOption.value}
								>{statusOption.label}</option
							>
						{/each}
					</select>
				</div>

				<!-- Rating -->
				<div>
					<label
						for="rating"
						class="block text-sm font-medium text-gray-700 mb-1"
					>
						Rating (optional)
					</label>
					<input
						id="rating"
						type="number"
						bind:value={rating}
						class="input"
						placeholder="8.5"
						min="0"
						max="10"
						step="0.5"
					/>
				</div>

				<!-- Review -->
				<div>
					<label
						for="review"
						class="block text-sm font-medium text-gray-700 mb-1"
					>
						Review (optional)
					</label>
					<textarea
						id="review"
						bind:value={reviewMd}
						class="input"
						rows="4"
						placeholder="Write your review in Markdown..."
					/>
				</div>

				<!-- Genres -->
				<div>
					<label
						for="genres"
						class="block text-sm font-medium text-gray-700 mb-1"
					>
						Genres (optional)
					</label>
					<div class="space-y-2">
						<div class="flex space-x-2">
							<input
								id="genres"
								type="text"
								bind:value={genreInput}
								on:keydown={handleGenreKeydown}
								class="input flex-1"
								placeholder="Enter genre and press Enter..."
							/>
							<button
								type="button"
								on:click={addGenre}
								class="btn btn-secondary px-3"
							>
								Add
							</button>
						</div>
						{#if genres.length > 0}
							<div class="flex flex-wrap gap-2">
								{#each genres as genre}
									<span
										class="inline-flex items-center px-2 py-1 bg-blue-100 text-blue-800 text-sm rounded-full"
									>
										{genre}
										<button
											type="button"
											on:click={() => removeGenre(genre)}
											class="ml-1 text-blue-600 hover:text-blue-800"
										>
											×
										</button>
									</span>
								{/each}
							</div>
						{/if}
					</div>
				</div>

				<!-- Actions -->
				<div class="flex space-x-3 pt-4">
					<button
						type="submit"
						class="btn btn-primary flex-1"
						disabled={loading}
					>
						{loading ? "Adding..." : "Add Entry"}
					</button>
					<button
						type="button"
						class="btn btn-secondary"
						on:click={handleClose}
					>
						Cancel
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
