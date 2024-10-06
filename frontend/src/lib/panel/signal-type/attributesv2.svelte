<script lang="ts">
	import type { SignalType } from '$lib/api/canturin';
	import SegmentedControl from '$lib/components/segmented-control/segmented-control.svelte';
	import { getSignalTypeState } from '$lib/state/signal-type-state.svelte';
	import type { PanelSectionProps } from '../types';

	let { entityId }: PanelSectionProps = $props();

	const sts = getSignalTypeState(entityId);

	let pippo = $state('flag');

	let pippo2 = $state('flag');
</script>

{#snippet section(signalType: SignalType)}
	<div class="grid grid-cols-2">
		<div>
			<h4 class="font-medium">Kind</h4>
			<p class="text-sm">The kind of the type</p>
		</div>

		<div>
			<SegmentedControl
				name="signal-type-kind"
				options={[
					{
						label: 'Flag',
						value: 'flag',
						desc: 'True or False'
					},
					{
						label: 'Integer',
						value: 'integer',
						desc: 'Any integer value'
					},
					{
						label: 'Decimal',
						value: 'decimal',
						desc: 'Any decimal value'
					},
					{
						label: 'Custom',
						value: 'custom',
						desc: 'Any value'
					}
				]}
				initialValue="flag"
				bind:selectedValue={pippo}
				disabled
			/>
		</div>
	</div>

	<div class="divider"></div>

	<div class="grid grid-cols-2">
		<div>
			<h4 class="font-medium">Kind</h4>
			<p class="text-sm">The kind of the type</p>
		</div>

		<div>
			<SegmentedControl
				name="signal-type-kind"
				options={[
					{
						label: 'Flag',
						value: 'flag',
						desc: 'True or False'
					},
					{
						label: 'Integer',
						value: 'integer',
						desc: 'Any integer value'
					},
					{
						label: 'Decimal',
						value: 'decimal',
						desc: 'Any decimal value'
					},
					{
						label: 'Custom',
						value: 'custom',
						desc: 'Any value'
					}
				]}
				initialValue="flag"
				bind:selectedValue={pippo2}
			/>
		</div>
	</div>
{/snippet}

<section>
	{@render section(sts.entity)}
</section>
