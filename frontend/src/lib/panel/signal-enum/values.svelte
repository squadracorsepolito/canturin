<script lang="ts">
	import type { SignalEnum, SignalEnumValue } from '$lib/api/canturin';
	import { AddIcon, DeleteIcon } from '$lib/components/icon';
	import { getSignalEnumState } from '$lib/panel/signal-enum/state.svelte';
	import type { PanelSectionProps } from '../types';
	import { Table } from '$lib/components/table';
	import { IconButton } from '$lib/components/button';
	import ValueRow from './value-row.svelte';
	import { TableTitle } from '$lib/components/table';

	let { entityId }: PanelSectionProps = $props();

	const ses = getSignalEnumState(entityId);

	function handleReorderValue(valueEntId: string, from: number, to: number) {
		ses.reorderValue(valueEntId, from, to);
	}

	function handleAdd() {
		ses.addValue();
	}

	function handleBulkDelete(values: SignalEnumValue[]) {
		ses.deleteValues(values.map((val) => val.entityId));
	}

	function handleDelete(value: SignalEnumValue) {
		ses.deleteValue(value.entityId);
	}
</script>

{#snippet section(signalEnum: SignalEnum)}
	{#if signalEnum.values && signalEnum.values.length > 0}
		<Table items={signalEnum.values} idKey="entityId" reorder={handleReorderValue}>
			{#snippet bulkActions({ selectedCount, selectedItems, deselectAll })}
				<div class="flex justify-end gap-5">
					<IconButton onclick={() => handleAdd()} label="Add Value" themeColor="primary">
						<AddIcon />
					</IconButton>

					<IconButton
						onclick={() => {
							handleBulkDelete(selectedItems);
							deselectAll();
						}}
						label={`Delete Values ${selectedCount > 0 ? ` (${selectedCount})` : ''}`}
						disabled={selectedCount === 0}
						themeColor="error"
					>
						<DeleteIcon />
					</IconButton>
				</div>
			{/snippet}

			{#snippet header()}
				<TableTitle>Name</TableTitle>

				<TableTitle>Index</TableTitle>

				<TableTitle>Description</TableTitle>
			{/snippet}

			{#snippet row(value)}
				<ValueRow enumEntityId={entityId} {value} />
			{/snippet}

			{#snippet rowActions(value)}
				<IconButton onclick={() => handleDelete(value)} themeColor="error">
					<DeleteIcon />
				</IconButton>
			{/snippet}
		</Table>
	{/if}
{/snippet}

<section>
	<h3 class="pb-5">Values</h3>

	{@render section(ses.entity)}
</section>
