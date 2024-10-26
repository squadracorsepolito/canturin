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

	let { entityId }: PanelSectionProps = $props();

	const ses = getSignalEnumState(entityId);

	function getValueInvalidNames(values: SignalEnumValue[], entId: string) {
		return values.filter((val) => val.entityId !== entId).map((val) => val.name);
	}

	function getValueInvalidIndexes(values: SignalEnumValue[], entId: string) {
		return values.filter((val) => val.entityId !== entId).map((val) => val.index);
	}

	let sortToggled = $state(false);
</script>

{#snippet section(signalEnum: SignalEnum)}
	{#if signalEnum.values && signalEnum.values.length > 0}
		<Toggle bind:toggled={sortToggled} name="sort">
			<SortIcon />
		</Toggle>

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
{/snippet}

<section>
	<h3 class="pb-5">{text.headings.values}</h3>

	{@render section(ses.entity)}
</section>
