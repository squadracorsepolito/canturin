<script lang="ts">
	import SortableList from '$lib/components/sortable/sortable-list.svelte';
	import { monitorForElements } from '@atlaskit/pragmatic-drag-and-drop/element/adapter';
	import { onMount } from 'svelte';
	import { reorder } from '@atlaskit/pragmatic-drag-and-drop/reorder';
	import { extractClosestEdge } from '@atlaskit/pragmatic-drag-and-drop-hitbox/closest-edge';
	import { getReorderDestinationIndex } from '@atlaskit/pragmatic-drag-and-drop-hitbox/util/get-reorder-destination-index';

	let items = $state([
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
	]);

	onMount(() => {
		return monitorForElements({
			canMonitor(args) {
				return args.source.data.instanceId === 'instance';
			},
			onDrop({ source, location }) {
				if (location.current.dropTargets.length === 0) return;

				const itemId = source.data.id;
				// const listId = location.initial.dropTargets[1].data.listId;

				const itemIdx = items.findIndex((item) => item.id === itemId);
				if (itemIdx === -1) return;

				// console.log('onDrop', itemId, itemIdx);

				console.log(location.current.dropTargets);

				// if (location.current.dropTargets.length === 1) {
				// 	console.log(
				// 		'dropTargets1',
				// 		location.current.dropTargets,
				// 		location.current.dropTargets.length
				// 	);
				// }

				// if (location.current.dropTargets.length === 1) {
				// Destructure and extract the destination card and column data from the drop targets
				const [destItemRecord] = location.current.dropTargets;

				// Find the index of the target card within the destination column's cards
				const indexOfTarget = items.findIndex((item) => item.id === destItemRecord.data.id);

				// Determine the closest edge of the target card: top or bottom
				const closestEdgeOfTarget = extractClosestEdge(destItemRecord.data);

				// Calculate the destination index for the card to be reordered within the same column
				const destinationIndex = getReorderDestinationIndex({
					startIndex: itemIdx,
					indexOfTarget,
					closestEdgeOfTarget,
					axis: 'vertical'
				});

				// Perform the card reordering within the same column
				reorderItems(itemIdx, destinationIndex);

				return;
				// }
			}
		});
	});

	function reorderItems(startIndex: number, finishIndex: number) {
		const updatedItems = reorder({
			list: items,
			startIndex,
			finishIndex
		});
		items = updatedItems;
	}
</script>

<SortableList {items}>
	{#snippet itemBody({ label })}
		<div class="border p-3">{label}</div>
	{/snippet}
</SortableList>
