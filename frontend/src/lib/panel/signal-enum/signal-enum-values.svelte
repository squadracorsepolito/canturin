<script lang="ts">
	import type { SignalEnum, SignalEnumValue } from '$lib/api/canturin';
	import { SortIcon } from '$lib/components/icon';
	import { Toggle } from '$lib/components/toggle';
	import { getSignalEnumState } from '$lib/state/signal-enum-state.svelte';
	import type { PanelSectionProps } from '../types';
	import { text } from './signal-enum-text';
	import SignalEnumValuesValue from './signal-enum-values-value.svelte';
	import { SortableList } from '$lib/components/sortable';
	import { reorder } from '@atlaskit/pragmatic-drag-and-drop/reorder';

	let { entityId }: PanelSectionProps = $props();

	const ses = getSignalEnumState(entityId);

	function getValueInvalidNames(values: SignalEnumValue[], entId: string) {
		return values.filter((val) => val.entityId !== entId).map((val) => val.name);
	}

	function getValueInvalidIndexes(values: SignalEnumValue[], entId: string) {
		return values.filter((val) => val.entityId !== entId).map((val) => val.index);
	}

	let sortToggled = $state(false);

	function reorderItems(startIndex: number, finishIndex: number) {
		if (!ses.entity.values) return;

		const updatedItems = reorder({
			list: ses.entity.values,
			startIndex,
			finishIndex
		});
		ses.entity.values = updatedItems;
	}
</script>

{#snippet section(signalEnum: SignalEnum)}
	{#if signalEnum.values && signalEnum.values.length > 0}
		<div class="px-2 pb-3">
			<Toggle bind:toggled={sortToggled} name="signal-enum-values-sort">
				<SortIcon />
			</Toggle>
		</div>

		{#if sortToggled}
			<SortableList
				items={signalEnum.values.map((val) => {
					return { id: val.entityId, value: val };
				})}
				instanceId={entityId + 'signal-enum-values'}
				reorder={reorderItems}
			>
				{#snippet itemBody({ item: { value } })}
					<div class="p-3 grid grid-cols-8 gap-5 items-center">
						<div class="justify-self-center"><span class="text-h3">{value.index}</span></div>

						<div class="col-span-7 text-h4">
							{value.name}
						</div>
					</div>
				{/snippet}
			</SortableList>
		{:else}
			<div class="grid grid-cols-8 gap-2">
				{#each signalEnum.values as sigEnumValue}
					<SignalEnumValuesValue
						{entityId}
						value={sigEnumValue}
						invalidNames={getValueInvalidNames(signalEnum.values, sigEnumValue.entityId)}
						invalidIndexes={getValueInvalidIndexes(signalEnum.values, sigEnumValue.entityId)}
					/>
				{/each}
			</div>
		{/if}
	{/if}
{/snippet}

<section>
	<h3 class="pb-5">{text.headings.values}</h3>

	{@render section(ses.entity)}
</section>
