<script lang="ts">
	import type { Action } from 'svelte/action';
	import {
		draggable,
		dropTargetForElements
	} from '@atlaskit/pragmatic-drag-and-drop/element/adapter';
	import {
		attachClosestEdge,
		extractClosestEdge
	} from '@atlaskit/pragmatic-drag-and-drop-hitbox/closest-edge';
	import type { Snippet } from 'svelte';
	import type { Edge } from '@atlaskit/pragmatic-drag-and-drop-hitbox/dist/types/types';
	import DropIndicator from './drop-indicator.svelte';
	import { combine } from '@atlaskit/pragmatic-drag-and-drop/combine';
	import { DragHandleIcon, SortIcon } from '../icon';
	import { getItem, isItem, type HighlightState } from './types';

	type Props = {
		instanceId: string;
		id: string;
		highlightState: HighlightState;
		children: Snippet;
	};

	let { instanceId, id, highlightState, children }: Props = $props();

	let closestEdge = $state<Edge | null>(null);
	let isDragging = $state(false);

	let dragHandle: HTMLElement;

	const itemAction: Action<HTMLElement> = (el) => {
		const cleanup = combine(
			draggable({
				element: el,
				dragHandle: dragHandle,
				getInitialData() {
					return getItem({ instanceId, id });
				},
				onDragStart() {
					isDragging = true;
				},
				onDrop() {
					isDragging = false;

					el.animate([{ backgroundColor: '#37cdbe' }, {}], {
						duration: 600,
						delay: 150,
						easing: 'cubic-bezier(0.25, 0.1, 0.25, 1.0)',
						iterations: 1
					});
				}
			}),
			dropTargetForElements({
				element: el,
				canDrop({ source }) {
					return isItem(source.data) && source.data.instanceId === instanceId;
				},
				getData({ element, input }) {
					return attachClosestEdge(getItem({ instanceId, id }), {
						element: element,
						input: input,
						allowedEdges: ['top', 'bottom']
					});
				},
				onDragEnter({ source, self }) {
					if (isItem(source.data) && source.data.id !== id) {
						closestEdge = extractClosestEdge(self.data);
					}
				},
				onDrag({ source, self }) {
					if (isItem(source.data) && source.data.id !== id) {
						closestEdge = extractClosestEdge(self.data);
					}
				},
				onDragLeave() {
					closestEdge = null;
				},
				onDrop() {
					closestEdge = null;
				}
			})
		);

		return {
			destroy() {
				cleanup();
			}
		};
	};
</script>

<div
	use:itemAction
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
