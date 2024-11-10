<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { Edge } from '@atlaskit/pragmatic-drag-and-drop-hitbox/dist/types/types';
	import DropIndicator from './drop-indicator.svelte';
	import { DragHandleIcon, SortIcon } from '../icon';
	import { type HighlightState } from './types';
	import { sortableItem } from '$lib/actions/sortable.svelte';

	type Props = {
		instanceId: string;
		id: string;
		highlightState: HighlightState;
		children: Snippet;
	};

	let { instanceId, id, highlightState, children }: Props = $props();

	let closestEdge = $state<Edge | null>(null);
	let isDragging = $state(false);

	let dragHandle = $state() as HTMLElement;
</script>

<div
	use:sortableItem={{
		id,
		instanceId,
		dragHandle,
		setClosestEdge(edge) {
			closestEdge = edge;
		},
		setIsDragging(isDragging) {
			isDragging = isDragging;
		}
	}}
	data-part="root"
	data-highlight-state={highlightState}
	data-dragging={isDragging ? 'true' : undefined}
>
	<div class="flex items-center gap-3 px-3">
		<div bind:this={dragHandle} data-part="handle">
			{#if highlightState === 'moving'}
				<SortIcon />
			{:else}
				<DragHandleIcon />
			{/if}
		</div>

		<div class="flex-1">
			{@render children()}
		</div>
	</div>

	{#if closestEdge}
		<DropIndicator edge={closestEdge} />
	{/if}
</div>

<style lang="postcss">
	[data-part='root'] {
		@apply relative rounded-btn border-2 transition-colors;

		&[data-dragging] {
			@apply opacity-20;
		}

		&[data-highlight-state='none'] {
			@apply border-transparent;
		}

		&[data-highlight-state='selecting'] {
			@apply border-secondary focus-ring-secondary text-secondary bg-secondary-ghost;
		}

		&[data-highlight-state='moving'] {
			@apply border-primary focus-ring-primary text-primary bg-primary-ghost;
		}
	}

	[data-part='handle'] {
		@apply rounded-btn cursor-grab;
	}
</style>
