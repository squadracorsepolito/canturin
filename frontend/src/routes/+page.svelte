<script lang="ts">
	import Sortable from 'sortablejs';
	import { SortableList } from '@jhubbardsf/svelte-sortablejs';

	let items = [
		{ id: 1, name: 'one' },
		{ id: 2, name: 'two' },
		{ id: 3, name: 'three' },
		{ id: 4, name: 'four' }
	];

	function initSortable(list: HTMLElement) {
		console.log(list);

		new Sortable(list, {
			group: 'pippo',
			animation: 150
		});
	}
</script>

<SortableList class="flex flex-col">
	{#each items as item}
		<div>
			{item.name}
		</div>
	{/each}
</SortableList>
