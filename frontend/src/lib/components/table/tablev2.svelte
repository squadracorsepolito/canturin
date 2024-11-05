<script lang="ts" generics="T">
	import type { Snippet } from 'svelte';
	import { Checkbox } from '../checkbox';

	type RowItem = {
		selected: boolean;
	} & T;

	type Props = {
		items: T[];
		bulkActions?: Snippet<[{ selectedCount: number; selectedItems: RowItem[] }]>;
		header: Snippet;
		row: Snippet<[T]>;
		rowActions?: Snippet<[RowItem]>;
	};

	let { items, bulkActions, header, row, rowActions }: Props = $props();

	class RowItems {
		#items = $state<RowItem[]>([]);

		constructor(getter: () => T[]) {
			$effect(() => {
				this.#items = getter().map((item) => ({ selected: false, ...item }));
			});
		}

		get items() {
			return this.#items;
		}
	}

	let rowItems = new RowItems(() => items);

	let allItemsSelected = $state(false);
	let selectedCount = $state(0);

	function handleSelectAll(checked: boolean) {
		if (checked) {
			selectedCount = rowItems.items.length;
		} else {
			selectedCount = 0;
		}

		for (const item of rowItems.items) {
			item.selected = checked;
		}
	}

	function handleRowSelect(checked: boolean) {
		if (checked) {
			selectedCount += 1;
		} else {
			allItemsSelected = false;
			selectedCount -= 1;
		}

		if (selectedCount === rowItems.items.length) {
			allItemsSelected = true;
		}
	}
</script>

{#if bulkActions}
	<div class="pb-3">
		{@render bulkActions({
			selectedCount,
			selectedItems: rowItems.items.filter((i) => i.selected)
		})}
	</div>
{/if}

<table class="table">
	<thead>
		<tr>
			<th>
				<Checkbox bind:checked={allItemsSelected} oncheckchange={handleSelectAll} />
			</th>

			{@render header()}

			{#if rowActions}
				<th>Actions</th>
			{/if}
		</tr>
	</thead>

	<tbody>
		{#each rowItems.items as item}
			<tr>
				<td>
					<Checkbox bind:checked={item.selected} oncheckchange={handleRowSelect} />
				</td>

				{@render row(item)}

				{#if rowActions}
					<td>
						{@render rowActions(item)}
					</td>
				{/if}
			</tr>
		{/each}
	</tbody>
</table>
