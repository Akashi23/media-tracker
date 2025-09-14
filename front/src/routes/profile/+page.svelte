<script lang="ts">
    import { onMount } from "svelte";
    import { auth } from "$stores/auth";
    import { entries } from "$stores/entries";
    import { collections } from "$stores/collections";
    import { authApi, entriesApi, collectionsApi } from "$utils/api";
    import { storage } from "$utils/storage";
    import type { User, Entry, Collection } from "$types";

    let user: User | null = null;
    let userStats = {
        totalEntries: 0,
        totalCollections: 0,
        completedEntries: 0,
        inProgressEntries: 0,
        plannedEntries: 0,
        averageRating: 0,
        favoriteGenres: [] as string[],
        recentEntries: [] as Entry[],
        recentCollections: [] as Collection[],
    };
    let loading = true;
    let error: string | null = null;

    onMount(async () => {
        // Initialize auth state from localStorage
        auth.init();

        if ($auth.isAuthenticated && $auth.token) {
            await loadProfile();
            await loadUserStats();
        } else {
            // For guest users, show guest profile
            user = {
                id: storage.getGuestId(),
                email: "guest@example.com",
                name: "Guest User",
                created_at: new Date().toISOString(),
            };
            await loadGuestStats();
        }

        loading = false;
    });

    async function loadProfile() {
        if (!$auth.token) return;

        try {
            user = await authApi.getProfile($auth.token);
        } catch (error) {
            console.error("Failed to load profile:", error);
            error = "Failed to load profile";
        }
    }

    async function loadUserStats() {
        if (!$auth.token) return;

        try {
            // Load entries and collections
            const [entriesList, collectionsList] = await Promise.all([
                entriesApi.list($auth.token),
                collectionsApi.list($auth.token),
            ]);

            // Calculate statistics
            userStats.totalEntries = entriesList.length;
            userStats.totalCollections = collectionsList.length;
            userStats.completedEntries = entriesList.filter(
                (e) => e.status === "completed"
            ).length;
            userStats.inProgressEntries = entriesList.filter(
                (e) => e.status === "in_progress"
            ).length;
            userStats.plannedEntries = entriesList.filter(
                (e) => e.status === "planned"
            ).length;

            // Calculate average rating
            const ratedEntries = entriesList.filter(
                (e) => e.rating && e.rating > 0
            );
            if (ratedEntries.length > 0) {
                userStats.averageRating =
                    ratedEntries.reduce((sum, e) => sum + (e.rating || 0), 0) /
                    ratedEntries.length;
            }

            // Get favorite genres
            const genreCount: Record<string, number> = {};
            entriesList.forEach((entry) => {
                if (entry.media?.genres) {
                    entry.media.genres.forEach((genre) => {
                        genreCount[genre] = (genreCount[genre] || 0) + 1;
                    });
                }
            });
            userStats.favoriteGenres = Object.entries(genreCount)
                .sort(([, a], [, b]) => b - a)
                .slice(0, 5)
                .map(([genre]) => genre);

            // Get recent entries (last 5)
            userStats.recentEntries = entriesList
                .sort(
                    (a, b) =>
                        new Date(b.updated_at || "").getTime() -
                        new Date(a.updated_at || "").getTime()
                )
                .slice(0, 5);

            // Get recent collections (last 3)
            userStats.recentCollections = collectionsList
                .sort(
                    (a, b) =>
                        new Date(b.created_at).getTime() -
                        new Date(a.created_at).getTime()
                )
                .slice(0, 3);
        } catch (error) {
            console.error("Failed to load user stats:", error);
        }
    }

    async function loadGuestStats() {
        try {
            const store = storage.getMediaStore();
            const entriesList = store.entries;
            const collectionsList = store.collections;

            // Calculate statistics for guest
            userStats.totalEntries = entriesList.length;
            userStats.totalCollections = collectionsList.length;
            userStats.completedEntries = entriesList.filter(
                (e) => e.status === "completed"
            ).length;
            userStats.inProgressEntries = entriesList.filter(
                (e) => e.status === "in_progress"
            ).length;
            userStats.plannedEntries = entriesList.filter(
                (e) => e.status === "planned"
            ).length;

            // Calculate average rating
            const ratedEntries = entriesList.filter(
                (e) => e.rating && e.rating > 0
            );
            if (ratedEntries.length > 0) {
                userStats.averageRating =
                    ratedEntries.reduce((sum, e) => sum + (e.rating || 0), 0) /
                    ratedEntries.length;
            }

            // Get favorite genres
            const genreCount: Record<string, number> = {};
            entriesList.forEach((entry) => {
                if (entry.media?.genres) {
                    entry.media.genres.forEach((genre) => {
                        genreCount[genre] = (genreCount[genre] || 0) + 1;
                    });
                }
            });
            userStats.favoriteGenres = Object.entries(genreCount)
                .sort(([, a], [, b]) => b - a)
                .slice(0, 5)
                .map(([genre]) => genre);

            // Get recent entries (last 5)
            userStats.recentEntries = entriesList
                .sort(
                    (a, b) =>
                        new Date(b.updated_at || "").getTime() -
                        new Date(a.updated_at || "").getTime()
                )
                .slice(0, 5);

            // Get recent collections (last 3)
            userStats.recentCollections = collectionsList
                .sort(
                    (a, b) =>
                        new Date(b.created_at).getTime() -
                        new Date(a.created_at).getTime()
                )
                .slice(0, 3);
        } catch (error) {
            console.error("Failed to load guest stats:", error);
        }
    }

    function formatDate(dateString: string): string {
        return new Date(dateString).toLocaleDateString();
    }

    function formatRating(rating: number): string {
        return rating.toFixed(1);
    }
</script>

<svelte:head>
    <title>Profile - Media Tracker</title>
</svelte:head>

<div class="container mx-auto px-4 py-8">
    {#if loading}
        <div class="flex justify-center items-center h-64">
            <div
                class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"
            />
        </div>
    {:else if error}
        <div
            class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded"
        >
            {error}
        </div>
    {:else}
        <!-- Profile Header -->
        <div class="bg-white rounded-lg shadow-md p-6 mb-6">
            <div class="flex items-center space-x-4">
                <div
                    class="w-16 h-16 bg-blue-500 rounded-full flex items-center justify-center text-white text-2xl font-bold"
                >
                    {user?.name?.charAt(0).toUpperCase() || "G"}
                </div>
                <div>
                    <h1 class="text-2xl font-bold text-gray-900">
                        {user?.name || "Guest User"}
                    </h1>
                    <p class="text-gray-600">{user?.email}</p>
                    <p class="text-sm text-gray-500">
                        Member since {user?.created_at
                            ? formatDate(user.created_at)
                            : "Recently"}
                    </p>
                </div>
            </div>
        </div>

        <!-- Statistics Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
            <div class="bg-white rounded-lg shadow-md p-6">
                <div class="flex items-center">
                    <div class="p-2 bg-blue-100 rounded-lg">
                        <svg
                            class="w-6 h-6 text-blue-600"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                            />
                        </svg>
                    </div>
                    <div class="ml-4">
                        <p class="text-sm font-medium text-gray-600">
                            Total Entries
                        </p>
                        <p class="text-2xl font-bold text-gray-900">
                            {userStats.totalEntries}
                        </p>
                    </div>
                </div>
            </div>

            <div class="bg-white rounded-lg shadow-md p-6">
                <div class="flex items-center">
                    <div class="p-2 bg-green-100 rounded-lg">
                        <svg
                            class="w-6 h-6 text-green-600"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                            />
                        </svg>
                    </div>
                    <div class="ml-4">
                        <p class="text-sm font-medium text-gray-600">
                            Completed
                        </p>
                        <p class="text-2xl font-bold text-gray-900">
                            {userStats.completedEntries}
                        </p>
                    </div>
                </div>
            </div>

            <div class="bg-white rounded-lg shadow-md p-6">
                <div class="flex items-center">
                    <div class="p-2 bg-purple-100 rounded-lg">
                        <svg
                            class="w-6 h-6 text-purple-600"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
                            />
                        </svg>
                    </div>
                    <div class="ml-4">
                        <p class="text-sm font-medium text-gray-600">
                            Collections
                        </p>
                        <p class="text-2xl font-bold text-gray-900">
                            {userStats.totalCollections}
                        </p>
                    </div>
                </div>
            </div>

            <div class="bg-white rounded-lg shadow-md p-6">
                <div class="flex items-center">
                    <div class="p-2 bg-yellow-100 rounded-lg">
                        <svg
                            class="w-6 h-6 text-yellow-600"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"
                            />
                        </svg>
                    </div>
                    <div class="ml-4">
                        <p class="text-sm font-medium text-gray-600">
                            Avg Rating
                        </p>
                        <p class="text-2xl font-bold text-gray-900">
                            {userStats.averageRating > 0
                                ? formatRating(userStats.averageRating)
                                : "N/A"}
                        </p>
                    </div>
                </div>
            </div>
        </div>

        <!-- Status Breakdown -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
            <!-- Status Chart -->
            <div class="bg-white rounded-lg shadow-md p-6">
                <h3 class="text-lg font-semibold text-gray-900 mb-4">
                    Status Breakdown
                </h3>
                <div class="space-y-3">
                    <div class="flex items-center justify-between">
                        <span class="text-sm text-gray-600">Completed</span>
                        <div class="flex items-center space-x-2">
                            <div class="w-24 bg-gray-200 rounded-full h-2">
                                <div
                                    class="bg-green-500 h-2 rounded-full"
                                    style="width: {userStats.totalEntries > 0
                                        ? (userStats.completedEntries /
                                              userStats.totalEntries) *
                                          100
                                        : 0}%"
                                />
                            </div>
                            <span class="text-sm font-medium text-gray-900"
                                >{userStats.completedEntries}</span
                            >
                        </div>
                    </div>
                    <div class="flex items-center justify-between">
                        <span class="text-sm text-gray-600">In Progress</span>
                        <div class="flex items-center space-x-2">
                            <div class="w-24 bg-gray-200 rounded-full h-2">
                                <div
                                    class="bg-blue-500 h-2 rounded-full"
                                    style="width: {userStats.totalEntries > 0
                                        ? (userStats.inProgressEntries /
                                              userStats.totalEntries) *
                                          100
                                        : 0}%"
                                />
                            </div>
                            <span class="text-sm font-medium text-gray-900"
                                >{userStats.inProgressEntries}</span
                            >
                        </div>
                    </div>
                    <div class="flex items-center justify-between">
                        <span class="text-sm text-gray-600">Planned</span>
                        <div class="flex items-center space-x-2">
                            <div class="w-24 bg-gray-200 rounded-full h-2">
                                <div
                                    class="bg-yellow-500 h-2 rounded-full"
                                    style="width: {userStats.totalEntries > 0
                                        ? (userStats.plannedEntries /
                                              userStats.totalEntries) *
                                          100
                                        : 0}%"
                                />
                            </div>
                            <span class="text-sm font-medium text-gray-900"
                                >{userStats.plannedEntries}</span
                            >
                        </div>
                    </div>
                </div>
            </div>

            <!-- Favorite Genres -->
            <div class="bg-white rounded-lg shadow-md p-6">
                <h3 class="text-lg font-semibold text-gray-900 mb-4">
                    Favorite Genres
                </h3>
                {#if userStats.favoriteGenres.length > 0}
                    <div class="flex flex-wrap gap-2">
                        {#each userStats.favoriteGenres as genre}
                            <span
                                class="px-3 py-1 bg-blue-100 text-blue-800 text-sm rounded-full"
                            >
                                {genre}
                            </span>
                        {/each}
                    </div>
                {:else}
                    <p class="text-gray-500 text-sm">
                        No genres data available
                    </p>
                {/if}
            </div>
        </div>

        <!-- Recent Activity -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <!-- Recent Entries -->
            <div class="bg-white rounded-lg shadow-md p-6">
                <h3 class="text-lg font-semibold text-gray-900 mb-4">
                    Recent Entries
                </h3>
                {#if userStats.recentEntries.length > 0}
                    <div class="space-y-3">
                        {#each userStats.recentEntries as entry}
                            <div
                                class="flex items-center space-x-3 p-3 bg-gray-50 rounded-lg"
                            >
                                <div
                                    class="w-10 h-10 bg-gray-200 rounded flex items-center justify-center"
                                >
                                    <span class="text-xs text-gray-600">
                                        {entry.media?.type
                                            ?.charAt(0)
                                            .toUpperCase() || "M"}
                                    </span>
                                </div>
                                <div class="flex-1 min-w-0">
                                    <p
                                        class="text-sm font-medium text-gray-900 truncate"
                                    >
                                        {entry.media?.title || "Unknown Title"}
                                    </p>
                                    <p class="text-xs text-gray-500">
                                        {entry.status}
                                        {entry.rating
                                            ? `• ${entry.rating}/10`
                                            : ""}
                                    </p>
                                </div>
                            </div>
                        {/each}
                    </div>
                {:else}
                    <p class="text-gray-500 text-sm">No recent entries</p>
                {/if}
            </div>

            <!-- Recent Collections -->
            <div class="bg-white rounded-lg shadow-md p-6">
                <h3 class="text-lg font-semibold text-gray-900 mb-4">
                    Recent Collections
                </h3>
                {#if userStats.recentCollections.length > 0}
                    <div class="space-y-3">
                        {#each userStats.recentCollections as collection}
                            <div
                                class="flex items-center space-x-3 p-3 bg-gray-50 rounded-lg"
                            >
                                <div
                                    class="w-10 h-10 bg-purple-100 rounded flex items-center justify-center"
                                >
                                    <svg
                                        class="w-5 h-5 text-purple-600"
                                        fill="none"
                                        stroke="currentColor"
                                        viewBox="0 0 24 24"
                                    >
                                        <path
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                            stroke-width="2"
                                            d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
                                        />
                                    </svg>
                                </div>
                                <div class="flex-1 min-w-0">
                                    <p
                                        class="text-sm font-medium text-gray-900 truncate"
                                    >
                                        {collection.title}
                                    </p>
                                    <p class="text-xs text-gray-500">
                                        {collection.entries?.length || 0} entries
                                    </p>
                                </div>
                            </div>
                        {/each}
                    </div>
                {:else}
                    <p class="text-gray-500 text-sm">No recent collections</p>
                {/if}
            </div>
        </div>
    {/if}
</div>
