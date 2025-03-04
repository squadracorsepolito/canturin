<script lang="ts">
	import type { BusBase, Network } from '$lib/api/canturin';
	import { IconButton, LinkButton } from '$lib/components/button';
	import { HoverPreview } from '$lib/components/hover-preview';
	import { AddIcon, DeleteIcon } from '$lib/components/icon';
	import { Table, TableField, TableTitle } from '$lib/components/table';
	import layout from '$lib/state/layout-state.svelte';
	import { getSelectItemFromBaudrate } from '../bus/utils';
	import { getNetworkState } from './state.svelte';

	const ns = getNetworkState();

	function handleAdd() {
		ns.addBus();
	}

	function handleBulkDelete(buses: BusBase[]) {
		ns.deleteBuses(buses.map((b) => b.entityId));
	}

	function handleDelete(bus: BusBase) {
		ns.deleteBuses([bus.entityId]);
	}
</script>

{#snippet preview(sig: BusBase)}
	<div>
		<span class="font-medium text-sm pr-1">{sig.name}</span>
	</div>

	{#if sig.desc}
		<div class="text-xs text-dimmed pt-1">{sig.desc}</div>
	{/if}
{/snippet}

{#snippet section(net: Network)}
	{#if net.buses}
		<Table items={net.buses} idKey="entityId">
			{#snippet bulkActions({ selectedCount, selectedItems, deselectAll })}
				<div class="flex justify-end gap-5">
					<IconButton onclick={handleAdd} label="Add Bus" themeColor="primary">
						<AddIcon />
					</IconButton>

					<IconButton
						onclick={() => {
							handleBulkDelete(selectedItems);
							deselectAll();
						}}
						label={`Delete Buses ${selectedCount > 0 ? ` (${selectedCount})` : ''}`}
						disabled={selectedCount === 0}
						themeColor="error"
					>
						<DeleteIcon />
					</IconButton>
				</div>
			{/snippet}

			{#snippet header()}
				<TableTitle>Name</TableTitle>

				<TableTitle>Baudrate</TableTitle>
			{/snippet}

			{#snippet row(bus)}
				<TableField>
					<HoverPreview placement="right">
						{#snippet trigger()}
							<LinkButton label={bus.name} onclick={() => layout.openPanel('bus', bus.entityId)} />
						{/snippet}

						{#snippet content()}
							{@render preview(bus)}
						{/snippet}
					</HoverPreview>
				</TableField>

				<TableField>
					<span class="px-2">
						{getSelectItemFromBaudrate(bus.baudrate).label}
					</span>
				</TableField>
			{/snippet}

			{#snippet rowActions(bus)}
				<IconButton onclick={() => handleDelete(bus)} themeColor="error">
					<DeleteIcon />
				</IconButton>
			{/snippet}
		</Table>
	{/if}
{/snippet}

<section>
	<h3 class="pb-5">Buses</h3>

	{@render section(ns.entity)}
</section>
