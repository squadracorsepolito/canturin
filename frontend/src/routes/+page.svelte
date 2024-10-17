<script lang="ts">
	import SortableItem from '$lib/components/sortable/sortable-item.svelte';

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
</script>

<div>
	{#each items as item}
		<SortableItem id={item.id}>
			<div class="h-10 w-32 border-2">
				{item.label}
			</div>
		</SortableItem>
	{/each}
</div>
