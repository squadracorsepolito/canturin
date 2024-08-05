<script lang="ts" generics="T">
	import type { Snippet } from 'svelte';

	type Props = {
		items: T[];
		header: Snippet;
		row: Snippet<[T]>;
	};

	let { items, header, row }: Props = $props();
</script>

<table class="table">
	<thead>
		<tr>
			{@render header()}
		</tr>
	</thead>

	<tbody>
		{#each items as item}
			<tr>
				{@render row(item)}
			</tr>
		{/each}
	</tbody>
</table>
