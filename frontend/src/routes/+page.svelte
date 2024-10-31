<script lang="ts">
	import Divider from '$lib/components/divider/divider.svelte';
	import SortableList from '$lib/components/sortable/sortable-list.svelte';
	import { reorder } from '@atlaskit/pragmatic-drag-and-drop/reorder';

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

	function reorderItems(startIndex: number, finishIndex: number) {
		console.log({ startIndex, finishIndex });

		const updatedItems = reorder({
			list: items,
			startIndex,
			finishIndex
		});
		items = updatedItems;
	}
</script>

<Divider></Divider>

<SortableList {items} instanceId="items" reorder={reorderItems}>
	{#snippet itemBody({ label })}
		<div class="p-3">{label}</div>
	{/snippet}
</SortableList>
