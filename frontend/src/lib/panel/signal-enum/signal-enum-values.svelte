<script lang="ts">
	import type { SignalEnum, SignalEnumValue } from '$lib/api/canturin';
	import { SortIcon } from '$lib/components/icon';
	import { Toggle } from '$lib/components/toggle';
	import { getSignalEnumState } from '$lib/state/signal-enum-state.svelte';
	import { onMount } from 'svelte';
	import type { PanelSectionProps } from '../types';
	import { text } from './signal-enum-text';
	import SignalEnumVal from './signal-enum-val.svelte';
	import Sortable from 'sortablejs';
	import type { Action } from 'svelte/action';
	import SortableList from '$lib/components/sortable/sortable-list.svelte';
	import { monitorForElements } from '@atlaskit/pragmatic-drag-and-drop/element/adapter';

	import { reorder } from '@atlaskit/pragmatic-drag-and-drop/reorder';
	import { extractClosestEdge } from '@atlaskit/pragmatic-drag-and-drop-hitbox/closest-edge';
	import { getReorderDestinationIndex } from '@atlaskit/pragmatic-drag-and-drop-hitbox/util/get-reorder-destination-index';

	let { entityId }: PanelSectionProps = $props();

	const ses = getSignalEnumState(entityId);

	function getValueInvalidNames(values: SignalEnumValue[], entId: string) {
		return values.filter((val) => val.entityId !== entId).map((val) => val.name);
	}

	function getValueInvalidIndexes(values: SignalEnumValue[], entId: string) {
		return values.filter((val) => val.entityId !== entId).map((val) => val.index);
	}

	let sortToggled = $state(false);

	onMount(() => {
		return monitorForElements({
			canMonitor(args) {
				return args.source.data.instanceId === 'instance';
			},
			onDrop({ source, location }) {
				if (location.current.dropTargets.length === 0) return;

				const itemId = source.data.id;
				// const listId = location.initial.dropTargets[1].data.listId;

				if (!ses.entity.values) return;

				const itemIdx = ses.entity.values.findIndex((item) => item.entityId === itemId);
				if (itemIdx === -1) return;

				// console.log('onDrop', itemId, itemIdx);

				// console.log(location.current.dropTargets);

				// if (location.current.dropTargets.length === 1) {
				// 	console.log(
				// 		'dropTargets1',
				// 		location.current.dropTargets,
				// 		location.current.dropTargets.length
				// 	);
				// }

				// if (location.current.dropTargets.length === 1) {
				// Destructure and extract the destination card and column data from the drop targets
				const [destItemRecord] = location.current.dropTargets;

				// Find the index of the target card within the destination column's cards
				const indexOfTarget = ses.entity.values.findIndex(
					(item) => item.entityId === destItemRecord.data.id
				);

				// Determine the closest edge of the target card: top or bottom
				const closestEdgeOfTarget = extractClosestEdge(destItemRecord.data);

				// Calculate the destination index for the card to be reordered within the same column
				const destinationIndex = getReorderDestinationIndex({
					startIndex: itemIdx,
					indexOfTarget,
					closestEdgeOfTarget,
					axis: 'vertical'
				});

				// Perform the card reordering within the same column
				reorderItems(itemIdx, destinationIndex);

				return;
				// }
			}
		});
	});

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
		<Toggle bind:toggled={sortToggled} name="sort">
			<SortIcon />
		</Toggle>

		{#if sortToggled}
			<SortableList
				items={signalEnum.values.map((val) => {
					return { id: val.entityId, value: val };
				})}
			>
				{#snippet itemBody({ value })}
					<div class="p-3">{value.name}</div>
				{/snippet}
			</SortableList>
		{:else}
			<div class="grid grid-cols-8 gap-2">
				{#each signalEnum.values as sigEnumValue}
					<SignalEnumVal
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
