<script lang="ts">
	import type { SignalEnum, SignalEnumValue } from '$lib/api/canturin';
	import { AddIcon, DeleteIcon, SortIcon } from '$lib/components/icon';
	import { Toggle } from '$lib/components/toggle';
	import { getSignalEnumState } from '$lib/state/signal-enum-state.svelte';
	import type { PanelSectionProps } from '../types';
	import { text } from './signal-enum-text';
	import Tablev2 from '$lib/components/table/tablev2.svelte';
	import { IconButton } from '$lib/components/button';
	import ValueRow from './value-row.svelte';
	import { TableTitle } from '$lib/components/table';

	let { entityId }: PanelSectionProps = $props();

	const ses = getSignalEnumState(entityId);

	function getValueInvalidNames(values: SignalEnumValue[], entId: string) {
		return values.filter((val) => val.entityId !== entId).map((val) => val.name);
	}

	function getValueInvalidIndexes(values: SignalEnumValue[], entId: string) {
		return values.filter((val) => val.entityId !== entId).map((val) => val.index);
	}

	let reorderToggled = $state(false);

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
		<Tablev2 items={signalEnum.values} idKey="entityId" reorder={handleReorderValue}>
			{#snippet bulkActions({ selectedCount, selectedItems })}
				<div class="flex justify-end gap-5">
					<IconButton onclick={() => handleAdd()} label="Add Value" color="primary">
						<AddIcon />
					</IconButton>

					<IconButton
						onclick={() => handleBulkDelete(selectedItems)}
						label={`Delete Values ${selectedCount > 0 ? ` (${selectedCount})` : ''}`}
						disabled={selectedCount === 0}
						color="error"
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
				<IconButton onclick={() => handleDelete(value)} color="error">
					<DeleteIcon />
				</IconButton>
			{/snippet}
		</Tablev2>
	{/if}
{/snippet}

<section>
	<h3 class="pb-5">{text.headings.values}</h3>

	{@render section(ses.entity)}
</section>
