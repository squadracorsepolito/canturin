<script lang="ts" generics="T extends { [K in keyof T]: any }">
	import type { Snippet } from 'svelte';
	import { Checkbox } from '../checkbox';
	import { flip } from 'svelte/animate';
	import TableField from './table-field.svelte';
	import TableTitle from './table-title.svelte';

	type Props = {
		items: T[];
		idKey: {
			[K in keyof T]: T[K] extends string ? K : never;
		}[keyof T];
		bulkActions?: Snippet<[{ selectedCount: number; selectedItems: T[] }]>;
		header: Snippet;
		row: Snippet<[T]>;
		rowActions?: Snippet<[T]>;
	};

	let { items, idKey, bulkActions, header, row, rowActions }: Props = $props();

	class Selector {
		#items = $state<
			{
				id: string;
				selected: boolean;
			}[]
		>([]);

		constructor(getter: () => T[]) {
			$effect(() => {
				this.#items = getter().map((item) => ({ selected: false, id: item[idKey] }));
			});
		}

		get items() {
			return this.#items;
		}
	}

	let itemSelector = new Selector(() => items);

	let allItemsSelected = $state(false);
	let selectedCount = $state(0);

	function handleSelectAll(checked: boolean) {
		if (checked) {
			selectedCount = itemSelector.items.length;
		} else {
			selectedCount = 0;
		}

		for (const item of itemSelector.items) {
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

		if (selectedCount === itemSelector.items.length) {
			allItemsSelected = true;
		}
	}

	let selectedItems = $derived.by(() => {
		const selectedIds = itemSelector.items.filter((item) => item.selected).map((item) => item.id);
		return items.filter((item) => selectedIds.includes(item[idKey]));
	});
</script>

{#if bulkActions}
	<div class="pb-3">
		{@render bulkActions({ selectedCount, selectedItems })}
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
				<TableTitle>Actions</TableTitle>
			{/if}
		</tr>
	</thead>

	<tbody>
		{#if items.length === itemSelector.items.length}
			{#each items as item, idx (item[idKey])}
				<tr
					animate:flip={{
						duration: 150
					}}
				>
					<td>
						<Checkbox
							bind:checked={itemSelector.items[idx].selected}
							oncheckchange={handleRowSelect}
						/>
					</td>

					{@render row(item)}

					{#if rowActions}
						<TableField>
							{@render rowActions(item)}
						</TableField>
					{/if}
				</tr>
			{/each}
		{/if}
	</tbody>
</table>
