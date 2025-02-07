<script lang="ts">
	import { colorByName } from '$lib/actions/color-name.svelte';
	import type { BusLoadMessage } from '$lib/api/canturin';
	import { LinkButton } from '$lib/components/button';
	import layoutStateSvelte from '$lib/state/layout-state.svelte';
	import { getColorByName } from '$lib/utils';
	import type { PanelSectionProps } from '../types';
	import { getBusState } from './state.svelte';

	const loadMessageCount = 15;
	const othersMessageEntityId = 'others';

	let { entityId }: PanelSectionProps = $props();

	const bs = getBusState(entityId);

	function getOthersLoadMessage(loadMsgs: BusLoadMessage[]): BusLoadMessage {
		let percentage = 0;
		let bitsPerSec = 0;
		for (const loadMsg of loadMsgs) {
			percentage += loadMsg.percentage;
			bitsPerSec += loadMsg.bitsPerSec;
		}

		return {
			entityId: othersMessageEntityId,
			name: 'Others',
			desc: '',
			createTime: new Date(),
			percentage: percentage,
			bitsPerSec: bitsPerSec
		};
	}
</script>

{#snippet loadMessage(loadMsg: BusLoadMessage)}
	{@const color = getColorByName(loadMsg.name)}
	{@const isOthers = loadMsg.entityId === othersMessageEntityId}

	<div class="col-span-3 grid grid-cols-subgrid items-center">
		<div>
			{#if isOthers}
				<span class="font-medium px-2">{loadMsg.name}</span>
			{:else}
				<LinkButton
					label={loadMsg.name}
					onclick={() => layoutStateSvelte.openPanel('message', loadMsg.entityId)}
				/>
			{/if}
		</div>

		<div class="col-span-2 flex items-center gap-3">
			<div class="h-8" style:width="{loadMsg.percentage}%">
				{#if isOthers}
					<div class="h-full w-full bg-neutral rounded-btn"></div>
				{:else}
					<div class="h-full w-full rounded-btn" style:background-color={color.bgColor}></div>
				{/if}
			</div>

			<span class="text-sm">
				{loadMsg.percentage.toFixed(2)}%
			</span>
		</div>
	</div>
{/snippet}

<section>
	<h3 class="pb-5">Load</h3>

	{#await bs.getLoad() then load}
		<div
			class="inline-block p-5 rounded-box mb-5 {load.percentage < 33
				? 'bg-success text-success-content'
				: load.percentage < 66
					? 'bg-warning text-warning-content'
					: 'bg-error text-error-content'}"
		>
			<span class="text-lg font-semibold">
				{load.percentage.toFixed(2)}%
			</span>
		</div>

		{#if load.messages}
			<div class="grid grid-cols-3 gap-x-3 gap-y-2">
				{#each load.messages.slice(0, loadMessageCount) as loadMsg}
					{@render loadMessage(loadMsg)}
				{/each}

				{#if load.messages.length > loadMessageCount}
					{@render loadMessage(getOthersLoadMessage(load.messages.slice(loadMessageCount)))}
				{/if}
			</div>
		{/if}
	{/await}
</section>
