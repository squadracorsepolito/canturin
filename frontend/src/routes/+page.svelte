<script lang="ts">
	import SortableItem from '$lib/components/sortable/sortable-item.svelte';
	import { dropTargetForElements } from '@atlaskit/pragmatic-drag-and-drop/element/adapter';
	import type { Action } from 'svelte/action';

	let items = [
		{
			id: 'item-1',
			label: 'Item 1'
		},
		{
			id: 'item-2',
			label: 'Item 2'
		},
		{
			id: 'item-3',
			label: 'Item 3'
		},
		{
			id: 'item-4',
			label: 'Item 4'
		},
		{
			id: 'item-5',
			label: 'Item 5'
		}
	];

	const dropTarget: Action<HTMLDivElement> = (el) => {
		const cleanup = dropTargetForElements({
			element: el,
			onDragEnter: () => {
				console.log('is-over:');
			},
			onDragLeave: () => {
				console.log('is leaving');
			},
			onDrop: () => console.log('drop')
		});

		return {
			destroy() {
				cleanup();
			}
		};
	};
</script>

<div use:dropTarget class="block border-2 h-96 w-96"></div>

<div class="w-60 border-yellow-300 border">
	{#each items as item}
		<SortableItem id={item.id}>
			<div class="w-32 border-2">
				{item.label}
			</div>
		</SortableItem>
	{/each}
</div>
