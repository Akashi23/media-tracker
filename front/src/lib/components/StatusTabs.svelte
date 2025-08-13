<script lang="ts">
	import { createEventDispatcher } from "svelte";
	import type { Status } from "$types";

	export let selectedStatus: Status | "all" = "all";

	const dispatch = createEventDispatcher();

	const statuses: Array<{
		value: Status | "all";
		label: string;
		count?: number;
	}> = [
		{ value: "all", label: "All" },
		{ value: "planned", label: "Planned" },
		{ value: "in_progress", label: "In Progress" },
		{ value: "completed", label: "Completed" },
		{ value: "on_hold", label: "On Hold" },
		{ value: "dropped", label: "Dropped" },
	];

	function handleStatusChange(status: Status | "all") {
		selectedStatus = status;
		dispatch("change", status);
	}
</script>

<div class="border-b border-gray-200">
	<nav class="-mb-px flex space-x-8">
		{#each statuses as status}
			<button
				class="py-2 px-1 border-b-2 font-medium text-sm {selectedStatus ===
				status.value
					? 'border-primary-500 text-primary-600'
					: 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
				on:click={() => handleStatusChange(status.value)}
			>
				{status.label}
				{#if status.count !== undefined}
					<span
						class="ml-2 bg-gray-100 text-gray-900 py-0.5 px-2.5 rounded-full text-xs"
					>
						{status.count}
					</span>
				{/if}
			</button>
		{/each}
	</nav>
</div>
