<script lang="ts" generics="T extends { [K in keyof T]: any }">
	import type { Snippet } from 'svelte';
	import { Checkbox } from '../checkbox';
	import { flip } from 'svelte/animate';
	import TableField from './table-field.svelte';
	import TableTitle from './table-title.svelte';
	import { Sortable } from '$lib/actions/sortable.svelte';
	import { DragHandleIcon, SortIcon } from '../icon';
	import { Toggle } from '../toggle';
	import './styles.css';
	import { uniqueId, type KeyOfString } from '$lib/utils';

	type Props = {
		items: T[];
		idKey: KeyOfString<T>;
		reorder?: (id: string, from: number, to: number) => void;
		bulkActions?: Snippet<[{ selectedCount: number; selectedItems: T[] }]>;
		header: Snippet;
		row: Snippet<[T]>;
		rowActions?: Snippet<[T]>;
	};

	let { items, idKey, reorder, bulkActions, header, row, rowActions }: Props = $props();

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

	$effect(() => {
		if (items.length === 0) {
			handleSelectAll(false);
		}
	});

	function handleSelectAll(checked: boolean) {
		if (checked) {
			selectedCount = itemSelector.items.length;
		} else {
			selectedCount = 0;
		}

		for (const item of itemSelector.items) {
			item.selected = checked;
		}

		if (allItemsSelected !== checked) {
			allItemsSelected = checked;
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

	const sortable = new Sortable({
		instanceId: uniqueId(),
		enabled: false,
		itemsGetter: () => items.map((item) => ({ id: item[idKey] })),
		reorder: (id, from, to) => {
			reorder?.(id, from, to);
		}
	});
</script>

<div class="flex items-center gap-5 justify-end pb-5">
	{#if reorder}
		<Toggle bind:toggled={sortable.enabled} name="sortable-list-enable">
			<SortIcon />
		</Toggle>
	{/if}

	{#if bulkActions}
		<div>
			{@render bulkActions({ selectedCount, selectedItems })}
		</div>
	{/if}
</div>

<table class="table">
	<thead>
		<tr>
			<th>
				<div class="min-h-6">
					{#if !sortable.enabled}
						<Checkbox bind:checked={allItemsSelected} oncheckchange={handleSelectAll} />
					{/if}
				</div>
			</th>

			{@render header()}

			{#if rowActions}
				<TableTitle>Actions</TableTitle>
			{/if}
		</tr>
	</thead>

	<tbody use:sortable.root class="table-body">
		{#if items.length === itemSelector.items.length}
			{#each items as item, idx (item[idKey])}
				<tr animate:flip={{ duration: 150 }} use:sortable.item={{ id: item[idKey] }}>
					<td>
						<div class="flex">
							{#if sortable.enabled}
								<div use:sortable.dragHandle={{ id: item[idKey] }}>
									{#if sortable.isItemMoving(item[idKey])}
										<SortIcon height="20" width="20" />
									{:else}
										<DragHandleIcon height="20" width="20" />
									{/if}
								</div>
							{:else}
								<Checkbox
									bind:checked={itemSelector.items[idx].selected}
									oncheckchange={handleRowSelect}
								/>
							{/if}
						</div>
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
