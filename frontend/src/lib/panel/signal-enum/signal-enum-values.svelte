<script lang="ts">
	import type { SignalEnum, SignalEnumValue } from '$lib/api/canturin';
	import { AddIcon, DeleteIcon, SortIcon } from '$lib/components/icon';
	import { Toggle } from '$lib/components/toggle';
	import { getSignalEnumState } from '$lib/state/signal-enum-state.svelte';
	import type { PanelSectionProps } from '../types';
	import { text } from './signal-enum-text';
	import SignalEnumValuesValue from './signal-enum-values-value.svelte';
	import { SortableList } from '$lib/components/sortable';
	import { reorder } from '@atlaskit/pragmatic-drag-and-drop/reorder';
	import Readonly from '$lib/components/readonly/readonly.svelte';
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
		if (!ses.entity.values) return;

		const updatedItems = reorder({
			list: ses.entity.values,
			startIndex: from,
			finishIndex: to
		});
		ses.entity.values = updatedItems;

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
		<Tablev2 items={signalEnum.values} idKey="entityId">
			{#snippet bulkActions({ selectedCount, selectedItems })}
				<div class="flex justify-end gap-3">
					<Toggle bind:toggled={reorderToggled} name="signal-enum-values-reorder">
						<SortIcon />
					</Toggle>

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

		<!-- <div class="grid grid-cols-8 gap-x-3">
			<div class="col-span-8 grid grid-cols-subgrid font-bold text-sm text-dimmed py-5">
				<div class="pl-2">INDEX</div>

				<div class="pl-2 col-span-3">NAME</div>

				<div class="pl-2 col-span-4">DESCRIPTION</div>
			</div>

			{#if reorderToggled}
				<div
					style:grid-row="span {signalEnum.values.length} / span {signalEnum.values.length}"
					class="grid grid-rows-subgrid items-center"
				>
					{#each ses.indexes as idx}
						<Readonly>
							<span>{idx}</span>
						</Readonly>
					{/each}
				</div>

				<div
					style:grid-row="span {signalEnum.values.length} / span {signalEnum.values.length}"
					class="col-span-7"
				>
					<SortableList
						items={signalEnum.values.map((val) => {
							return { id: val.entityId, value: val };
						})}
						instanceId={entityId + 'signal-enum-values'}
						reorder={handleReorderValue}
					>
						{#snippet itemBody({ item: { value } })}
							<div class="p-3 grid grid-cols-8 gap-5 items-center">
								<div class="col-span-7 text-h4">
									{value.name}
								</div>
							</div>
						{/snippet}
					</SortableList>
				</div>
			{:else}
				{#each signalEnum.values as sigEnumValue}
					<SignalEnumValuesValue
						{entityId}
						value={sigEnumValue}
						invalidNames={getValueInvalidNames(signalEnum.values, sigEnumValue.entityId)}
						invalidIndexes={getValueInvalidIndexes(signalEnum.values, sigEnumValue.entityId)}
					/>
				{/each}
			{/if}
		</div> -->
	{/if}
{/snippet}

<section>
	<h3 class="pb-5"># {text.headings.values}</h3>

	{@render section(ses.entity)}
</section>
