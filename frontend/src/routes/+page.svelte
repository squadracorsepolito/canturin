<script lang="ts">
	import Sortable from 'sortablejs';

	let items = [
		{ id: 1, name: 'one' },
		{ id: 2, name: 'two' },
		{ id: 3, name: 'three' },
		{ id: 4, name: 'four' }
	];

	function initSortable(list: HTMLElement) {
		console.log(list);

		const sortable = Sortable.create(list, {
			animation: 150,
			group: 'pippo'
		});

		console.log(sortable);

		return sortable.destroy();
	}
</script>

<ul use:initSortable>
	{#each items as item (item.id)}
		<li>{item.name}</li>
	{/each}
</ul>
