<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/stores";
    import { publicApi } from "$utils/api";
    import type { Collection, Entry } from "$types";

    let shareData: Collection | Entry[] | null = null;
    let loading = true;
    let error: string | null = null;
    let shareType: "collection" | "profile" | null = null;
    let collectionData: Collection | null = null;
    let profileData: Entry[] | null = null;

    onMount(async () => {
        const token = $page.params.token;
        if (!token) {
            error = "Invalid share link";
            loading = false;
            return;
        }

        try {
            const data = await publicApi.getShare(token);
            console.log("Share data received:", data);
            shareData = data;

            // Determine share type based on data structure
            if (
                data &&
                typeof data === "object" &&
                "title" in data &&
                "entries" in data
            ) {
                shareType = "collection";
                collectionData = data as Collection;
                console.log("Collection data:", collectionData);
                console.log("Collection entries:", collectionData?.entries);
            } else if (Array.isArray(data)) {
                shareType = "profile";
                profileData = data as Entry[];
                console.log("Profile data:", profileData);
            }
        } catch (err) {
            console.error("Failed to load share:", err);
            error = "Share not found or expired";
        } finally {
            loading = false;
        }
    });

    function formatDate(dateStr: string): string {
        return new Date(dateStr).toLocaleDateString();
    }

    function getStatusColor(status: string): string {
        switch (status) {
            case "completed":
                return "bg-green-100 text-green-800";
            case "in_progress":
                return "bg-blue-100 text-blue-800";
            case "planned":
                return "bg-yellow-100 text-yellow-800";
            case "on_hold":
                return "bg-orange-100 text-orange-800";
            case "dropped":
                return "bg-red-100 text-red-800";
            default:
                return "bg-gray-100 text-gray-800";
        }
    }
</script>

<svelte:head>
    <title>Shared Content - Media Tracker</title>
</svelte:head>

<div class="container mx-auto px-4 py-8">
    {#if loading}
        <div class="flex justify-center items-center h-64">
            <div
                class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"
            />
        </div>
    {:else if error}
        <div class="text-center py-12">
            <div class="text-red-600 text-xl font-semibold mb-4">{error}</div>
            <p class="text-gray-600 mb-6">
                This share link may be invalid or expired.
            </p>
            <a href="/" class="btn btn-primary">Go to Home</a>
        </div>
    {:else if shareData}
        {#if shareType === "collection"}
            <!-- Collection Share -->
            <div class="max-w-4xl mx-auto">
                <div class="bg-white rounded-lg shadow-md p-6 mb-6">
                    <div class="flex items-center justify-between mb-4">
                        <div>
                            <h1 class="text-3xl font-bold text-gray-900 mb-2">
                                {collectionData?.title}
                            </h1>
                            <div
                                class="flex items-center space-x-4 text-sm text-gray-600"
                            >
                                <span
                                    >Created: {formatDate(
                                        collectionData?.created_at || ""
                                    )}</span
                                >
                                <span>•</span>
                                <span
                                    >{collectionData?.entries?.length || 0} entries</span
                                >
                                <span>•</span>
                                <span
                                    class="px-2 py-1 bg-green-100 text-green-800 text-xs rounded-full"
                                >
                                    Public Collection
                                </span>
                            </div>
                        </div>
                    </div>

                    {#if collectionData?.entries && collectionData.entries.length > 0}
                        <div
                            class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"
                        >
                            {#each collectionData.entries as entry}
                                <div
                                    class="border rounded-lg p-4 hover:shadow-md transition-shadow"
                                >
                                    <div class="flex items-start space-x-3">
                                        <div
                                            class="w-12 h-12 bg-gray-200 rounded flex items-center justify-center flex-shrink-0"
                                        >
                                            <span
                                                class="text-sm font-medium text-gray-600"
                                            >
                                                {entry.media?.type
                                                    ?.charAt(0)
                                                    .toUpperCase() || "M"}
                                            </span>
                                        </div>

                                        <div class="flex-1 min-w-0">
                                            <h3
                                                class="font-medium text-gray-900 truncate"
                                            >
                                                {entry.media?.title ||
                                                    "Unknown Title"}
                                            </h3>
                                            <p class="text-sm text-gray-500">
                                                {entry.media?.year
                                                    ? entry.media.year
                                                    : "Unknown Year"} • {entry
                                                    .media?.type ||
                                                    "Unknown Type"}
                                            </p>
                                            <div
                                                class="flex items-center space-x-2 mt-2"
                                            >
                                                <span
                                                    class="px-2 py-1 text-xs font-medium rounded-full {getStatusColor(
                                                        entry.status
                                                    )}"
                                                >
                                                    {entry.status.replace(
                                                        "_",
                                                        " "
                                                    )}
                                                </span>
                                                {#if entry.rating}
                                                    <span
                                                        class="text-sm text-gray-600"
                                                    >
                                                        {entry.rating}/10
                                                    </span>
                                                {/if}
                                            </div>

                                            {#if entry.review_md && entry.review_md.trim()}
                                                <div
                                                    class="mt-3 p-3 bg-gray-50 rounded-lg"
                                                >
                                                    <h4
                                                        class="text-sm font-medium text-gray-700 mb-2"
                                                    >
                                                        Review:
                                                    </h4>
                                                    <div
                                                        class="text-sm text-gray-600 line-clamp-3"
                                                    >
                                                        {entry.review_md}
                                                    </div>
                                                </div>
                                            {/if}
                                        </div>
                                    </div>
                                </div>
                            {/each}
                        </div>
                    {:else}
                        <p class="text-gray-500 text-center py-8">
                            This collection is empty.
                        </p>
                    {/if}
                </div>
            </div>
        {:else if shareType === "profile"}
            <!-- Profile Share -->
            <div class="max-w-4xl mx-auto">
                <div class="bg-white rounded-lg shadow-md p-6 mb-6">
                    <h1 class="text-3xl font-bold text-gray-900 mb-4">
                        Public Profile
                    </h1>
                    <p class="text-gray-600 mb-6">
                        Viewing {profileData?.length || 0} entries from this user's
                        collection.
                    </p>

                    {#if profileData && profileData.length > 0}
                        <div
                            class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"
                        >
                            {#each profileData as entry}
                                <div
                                    class="border rounded-lg p-4 hover:shadow-md transition-shadow"
                                >
                                    <div class="flex items-start space-x-3">
                                        <div
                                            class="w-12 h-12 bg-gray-200 rounded flex items-center justify-center flex-shrink-0"
                                        >
                                            <span
                                                class="text-sm font-medium text-gray-600"
                                            >
                                                {entry.media?.type
                                                    ?.charAt(0)
                                                    .toUpperCase() || "M"}
                                            </span>
                                        </div>

                                        <div class="flex-1 min-w-0">
                                            <h3
                                                class="font-medium text-gray-900 truncate"
                                            >
                                                {entry.media?.title ||
                                                    "Unknown Title"}
                                            </h3>
                                            <p class="text-sm text-gray-500">
                                                {entry.media?.year
                                                    ? entry.media.year
                                                    : "Unknown Year"} • {entry
                                                    .media?.type ||
                                                    "Unknown Type"}
                                            </p>
                                            <div
                                                class="flex items-center space-x-2 mt-2"
                                            >
                                                <span
                                                    class="px-2 py-1 text-xs font-medium rounded-full {getStatusColor(
                                                        entry.status
                                                    )}"
                                                >
                                                    {entry.status.replace(
                                                        "_",
                                                        " "
                                                    )}
                                                </span>
                                                {#if entry.rating}
                                                    <span
                                                        class="text-sm text-gray-600"
                                                    >
                                                        {entry.rating}/10
                                                    </span>
                                                {/if}
                                            </div>

                                            {#if entry.review_md && entry.review_md.trim()}
                                                <div
                                                    class="mt-3 p-3 bg-gray-50 rounded-lg"
                                                >
                                                    <h4
                                                        class="text-sm font-medium text-gray-700 mb-2"
                                                    >
                                                        Review:
                                                    </h4>
                                                    <div
                                                        class="text-sm text-gray-600 line-clamp-3"
                                                    >
                                                        {entry.review_md}
                                                    </div>
                                                </div>
                                            {/if}
                                        </div>
                                    </div>
                                </div>
                            {/each}
                        </div>
                    {:else}
                        <p class="text-gray-500 text-center py-8">
                            No entries to display.
                        </p>
                    {/if}
                </div>
            </div>
        {/if}

        <!-- Footer -->
        <div class="text-center mt-8">
            <p class="text-gray-500 text-sm">
                This is a public share from <a
                    href="/"
                    class="text-blue-600 hover:underline">Media Tracker</a
                >
            </p>
        </div>
    {/if}
</div>
