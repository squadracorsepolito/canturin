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
	import { DragHandleIcon } from '../icon';

	type Props = {
		id: string;
		children: Snippet;
	};

	let { id, children }: Props = $props();

	let closestEdge = $state<Edge | null>(null);
	let isDragging = $state(false);

	let dragHandle: HTMLElement;

	const itemAction: Action<HTMLElement> = (el) => {
		const cleanup = combine(
			draggable({
				element: el,
				dragHandle: dragHandle,
				getInitialData() {
					return {
						instanceId: 'instance',
						type: 'item',
						id: id
					};
				},
				onDragStart() {
					isDragging = true;
				},
				onDrop() {
					isDragging = false;

					el.animate([{ backgroundColor: '#37cdbe' }, {}], {
						duration: 500,
						easing: 'cubic-bezier(0.25, 0.1, 0.25, 1.0)',
						iterations: 1
					});
				}
			}),
			dropTargetForElements({
				element: el,
				canDrop({ source }) {
					return source.data.instanceId === 'instance';
				},
				getData(args) {
					return attachClosestEdge(
						{
							type: 'item',
							id: id
						},
						{
							element: args.element,
							input: args.input,
							allowedEdges: ['top', 'bottom']
						}
					);
				},
				onDragEnter(args) {
					if (args.source.data.id !== id) {
						closestEdge = extractClosestEdge(args.self.data);
					}
				},
				onDrag(args) {
					if (args.source.data.id !== id) {
						closestEdge = extractClosestEdge(args.self.data);
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

<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
<li use:itemAction tabindex="0" class="relative {isDragging && 'opacity-25'}">
	<div class="flex items-center gap-3 px-3">
		<div bind:this={dragHandle}>
			<DragHandleIcon />
		</div>

		<div class="flex-1">
			{@render children()}
		</div>
	</div>

	{#if closestEdge}
		<DropIndicator edge={closestEdge} />
	{/if}
</li>
