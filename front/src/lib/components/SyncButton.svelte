<script lang="ts">
    import { createEventDispatcher, onDestroy } from "svelte";
    import { auth } from "$stores/auth";
    import { entries } from "$stores/entries";
    import { entriesApi } from "$utils/api";
    import { storage } from "$utils/storage";
    import type { CreateEntryRequest, Entry } from "$types";

    const dispatch = createEventDispatcher();

    let loading = false;
    let lastSyncTime: string | null = null;
    let showErrorDetails = false;
    let lastSyncErrors: string[] = [];
    let syncTimeout: number | null = null;
    let lastSyncAttempt: number = 0;
    let showSyncErrorModal = false;
    let syncErrorMessage = "";
    let showSyncSuccessModal = false;
    let syncSuccessMessage = "";

    // Load last sync time from localStorage
    if (typeof window !== "undefined") {
        lastSyncTime = localStorage.getItem("lastSyncTime");
    }

    async function handleSync() {
        if (!$auth.isAuthenticated || !$auth.token) {
            syncErrorMessage = "Please login to sync your data";
            showSyncErrorModal = true;
            loading = false; // Ensure loading is reset
            return;
        }

        if (loading) {
            console.log("Sync already in progress, ignoring request");
            return;
        }

        // Prevent rapid successive clicks (debounce)
        const now = Date.now();
        if (now - lastSyncAttempt < 2000) {
            // 2 second debounce
            console.log("Sync debounced - too soon after last attempt");
            return;
        }
        lastSyncAttempt = now;

        loading = true;

        // Set a timeout to prevent endless loading
        syncTimeout = window.setTimeout(() => {
            if (loading) {
                console.error("Sync timeout - forcing stop");
                loading = false;
                lastSyncErrors = ["Sync timeout - operation took too long"];
                syncErrorMessage =
                    "Sync operation timed out. Please try again.";
                showSyncErrorModal = true;
                dispatch("sync-error", {
                    message: "Sync operation timed out. Please try again.",
                });
            }
        }, 30000); // 30 second timeout

        try {
            // Get local entries from storage
            const store = storage.getMediaStore();
            const localEntries = store.entries;

            if (localEntries.length === 0) {
                syncErrorMessage = "No local entries to sync";
                showSyncErrorModal = true;
                loading = false;
                if (syncTimeout) {
                    clearTimeout(syncTimeout);
                    syncTimeout = null;
                }
                return;
            }

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
            const response = (await Promise.race([
                entriesApi.sync(entriesToSync, $auth.token),
                new Promise<never>((_, reject) =>
                    setTimeout(
                        () => reject(new Error("Request timeout")),
                        25000
                    )
                ),
            ])) as {
                entries: Entry[];
                count: number;
                message: string;
                errors?: string[];
            };

            // Update local entries store with synced data
            entries.setEntries(response.entries || []);

            // Update last sync time
            lastSyncTime = new Date().toISOString();
            if (typeof window !== "undefined") {
                localStorage.setItem("lastSyncTime", lastSyncTime);
            }

            // Show success message
            let message = `Successfully synced ${response.count} entries`;
            if (response.errors && response.errors.length > 0) {
                message += ` (${response.errors.length} errors occurred)`;
                lastSyncErrors = response.errors;
                console.warn("Sync errors:", response.errors);
            } else {
                lastSyncErrors = [];
            }

            syncSuccessMessage = message;
            showSyncSuccessModal = true;

            dispatch("sync-success", {
                message,
                count: response.count,
                errors: response.errors || [],
            });
        } catch (error) {
            console.error("Sync failed:", error);
            const errorMessage =
                error instanceof Error ? error.message : "Unknown error";
            lastSyncErrors = [`Network error: ${errorMessage}`];
            syncErrorMessage = "Failed to sync entries. Please try again.";
            showSyncErrorModal = true;
            dispatch("sync-error", {
                message: "Failed to sync entries. Please try again.",
            });
        } finally {
            loading = false;
            if (syncTimeout) {
                clearTimeout(syncTimeout);
                syncTimeout = null;
            }
        }
    }

    function formatLastSync(): string {
        if (!lastSyncTime) return "Never";
        return new Date(lastSyncTime).toLocaleString();
    }

    function getLocalEntryCount(): number {
        const store = storage.getMediaStore();
        return store.entries.length;
    }

    function toggleErrorDetails() {
        showErrorDetails = !showErrorDetails;
    }

    function closeSyncErrorModal() {
        showSyncErrorModal = false;
        syncErrorMessage = "";
    }

    function closeSyncSuccessModal() {
        showSyncSuccessModal = false;
        syncSuccessMessage = "";
    }

    // Cleanup timeout on component destroy
    onDestroy(() => {
        if (syncTimeout) {
            clearTimeout(syncTimeout);
        }
    });
</script>

<div class="flex items-center space-x-3">
    {#if $auth.isAuthenticated}
        <div class="text-sm text-gray-600">
            <div>Local entries: {getLocalEntryCount()}</div>
            <div>Last sync: {formatLastSync()}</div>
            {#if lastSyncErrors.length > 0}
                <div class="text-red-600">
                    {lastSyncErrors.length} error(s) in last sync
                    <button
                        class="text-blue-600 underline ml-1"
                        on:click={toggleErrorDetails}
                    >
                        {showErrorDetails ? "hide" : "show"} details
                    </button>
                </div>
            {/if}
        </div>

        <button
            class="btn btn-secondary flex items-center space-x-2 {loading
                ? 'opacity-75 cursor-not-allowed'
                : ''}"
            on:click={handleSync}
            disabled={loading || getLocalEntryCount() === 0}
            title={loading
                ? "Sync in progress..."
                : "Sync local entries with server"}
        >
            {#if loading}
                <div
                    class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"
                />
                <span>Syncing...</span>
            {:else}
                <span>🔄</span>
                <span>Sync</span>
            {/if}
        </button>
    {:else}
        <div class="text-sm text-gray-500">Login to sync your data</div>
    {/if}
</div>

<!-- Sync Success Modal -->
{#if showSyncSuccessModal}
    <div
        class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
        <div class="bg-white rounded-lg p-6 w-full max-w-md mx-4">
            <div class="flex justify-between items-center mb-4">
                <h3 class="text-lg font-bold text-green-600">
                    Sync Successful
                </h3>
                <button
                    class="text-gray-400 hover:text-gray-600"
                    on:click={closeSyncSuccessModal}
                >
                    ✕
                </button>
            </div>

            <div class="mb-4">
                <p class="text-gray-700">{syncSuccessMessage}</p>
                {#if lastSyncErrors.length > 0}
                    <p class="text-sm text-orange-600 mt-2">
                        Some entries had errors. Check the sync button for
                        details.
                    </p>
                {/if}
            </div>

            <div class="flex justify-end space-x-3">
                <button
                    class="btn btn-primary"
                    on:click={closeSyncSuccessModal}
                >
                    OK
                </button>
            </div>
        </div>
    </div>
{/if}

<!-- Sync Error Modal -->
{#if showSyncErrorModal}
    <div
        class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
        <div class="bg-white rounded-lg p-6 w-full max-w-md mx-4">
            <div class="flex justify-between items-center mb-4">
                <h3 class="text-lg font-bold text-red-600">Sync Error</h3>
                <button
                    class="text-gray-400 hover:text-gray-600"
                    on:click={closeSyncErrorModal}
                >
                    ✕
                </button>
            </div>

            <div class="mb-4">
                <p class="text-gray-700">{syncErrorMessage}</p>
            </div>

            <div class="flex justify-end space-x-3">
                <button
                    class="btn btn-secondary"
                    on:click={closeSyncErrorModal}
                >
                    Close
                </button>
            </div>
        </div>
    </div>
{/if}

<!-- Error Details Modal -->
{#if showErrorDetails && lastSyncErrors.length > 0}
    <div
        class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
        <div
            class="bg-white rounded-lg p-6 w-full max-w-2xl max-h-[80vh] overflow-y-auto"
        >
            <div class="flex justify-between items-center mb-4">
                <h3 class="text-lg font-bold text-red-600">Sync Errors</h3>
                <button
                    class="text-gray-400 hover:text-gray-600"
                    on:click={toggleErrorDetails}
                >
                    ✕
                </button>
            </div>

            <div class="space-y-2">
                {#each lastSyncErrors as error, index}
                    <div
                        class="p-3 bg-red-50 border border-red-200 rounded text-sm"
                    >
                        <div class="font-medium text-red-800">
                            Error {index + 1}:
                        </div>
                        <div class="text-red-700 mt-1">{error}</div>
                    </div>
                {/each}
            </div>

            <div class="mt-4 flex justify-end">
                <button class="btn btn-secondary" on:click={toggleErrorDetails}>
                    Close
                </button>
            </div>
        </div>
    </div>
{/if}
