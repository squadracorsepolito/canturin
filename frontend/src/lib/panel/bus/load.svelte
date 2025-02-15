<script lang="ts">
	import { colorByName } from '$lib/actions/color-name.svelte';
	import type { BusLoadMessage } from '$lib/api/canturin';
	import { LinkButton } from '$lib/components/button';
	import { Collapsible } from '$lib/components/collapsible';
	import { AltArrowIcon, CloseIcon } from '$lib/components/icon';
	import layoutStateSvelte from '$lib/state/layout-state.svelte';
	import type { PanelSectionProps } from '../types';
	import { getBusState } from './state.svelte';

	const loadMessageCount = 15;

	let { entityId }: PanelSectionProps = $props();

	const bs = getBusState(entityId);

	function formatPercentage(percentage: number) {
		return `${percentage.toFixed(2)}%`;
	}

	function summarizeLoad(loadMsgs: BusLoadMessage[]) {
		let percentage = 0;

		for (let i = loadMessageCount; i < loadMsgs.length; i++) {
			percentage += loadMsgs[i].percentage;
		}

		return percentage;
	}
</script>

{#snippet percentage(percentage: number, name?: string)}
	<div class="col-span-2 flex items-center gap-3">
		<div class="h-8" style:width="{percentage}%">
			{#if name}
				<div use:colorByName={{ name }} class="h-full w-full rounded-btn"></div>
			{:else}
				<div class="h-full w-full bg-neutral rounded-btn"></div>
			{/if}
		</div>

		<span class="text-sm">
			{formatPercentage(percentage)}
		</span>
	</div>
{/snippet}

{#snippet loadMessage(loadMsg: BusLoadMessage)}
	<div class="col-span-3 grid grid-cols-subgrid items-center">
		<div class="h-full">
			<LinkButton
				label={loadMsg.name}
				onclick={() => layoutStateSvelte.openPanel('message', loadMsg.entityId)}
			/>
		</div>

		{@render percentage(loadMsg.percentage, loadMsg.name)}
	</div>
{/snippet}

<section>
	<h3 class="pb-5">Load</h3>

	{#await bs.getLoad() then load}
		<div
			class="inline-block p-3 rounded-box mb-5 {load.percentage < 33
				? 'bg-success text-success-content'
				: load.percentage < 66
					? 'bg-warning text-warning-content'
					: 'bg-error text-error-content'}"
		>
			<span class="text-lg font-semibold">
				{formatPercentage(load.percentage)}
			</span>
		</div>

		{#if load.messages}
			<div class="grid grid-cols-3 gap-x-3 gap-y-2">
				{#each load.messages.slice(0, loadMessageCount) as loadMsg}
					{@render loadMessage(loadMsg)}
				{/each}
			</div>

			{#if load.messages.length > loadMessageCount}
				<div class="pt-3">
					<Collapsible initialCollapsed raw>
						{#snippet trigger({ collapsed })}
							<div class="grid grid-cols-3 gap-3 items-center py-2">
								<div class="flex items-center gap-2 pl-2">
									{#if collapsed}
										<AltArrowIcon height="18" width="18" />
									{:else}
										<CloseIcon height="18" width="18" />
									{/if}

									<span>Others</span>
								</div>

								{#if load.messages}
									{@render percentage(summarizeLoad(load.messages))}
								{/if}
							</div>
						{/snippet}

						{#snippet content()}
							{#if load.messages}
								<div class="grid grid-cols-3 gap-x-3 gap-y-2 bg-base-200 py-2 rounded-btn mt-2">
									{#each load.messages.slice(loadMessageCount) as loadMsg}
										{@render loadMessage(loadMsg)}
									{/each}
								</div>
							{/if}
						{/snippet}
					</Collapsible>
				</div>
			{/if}
		{/if}
	{/await}
</section>
