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

	type Props = {
		id: string;
		children: Snippet;
	};

	let { id, children }: Props = $props();

	let closestEdge = $state<Edge | null>(null);

	const itemAction: Action<HTMLElement> = (el) => {
		const cleanup = combine(
			draggable({
				element: el,
				getInitialData() {
					return {
						instanceId: 'instance',
						type: 'item',
						id: id
					};
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

<div use:itemAction class="relative">
	{@render children()}

	{#if closestEdge}
		<DropIndicator edge={closestEdge} />
	{/if}
</div>
