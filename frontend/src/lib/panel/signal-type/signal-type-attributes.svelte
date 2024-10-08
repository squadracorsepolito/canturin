<script lang="ts">
	import type { SignalType } from '$lib/api/canturin';
	import { Attribute } from '$lib/components/attribute';
	import Divider from '$lib/components/divider/divider.svelte';
	import { NumberInput } from '$lib/components/input';
	import SegmentedControl from '$lib/components/segmented-control/segmented-control.svelte';
	import { getSignalTypeState } from '$lib/state/signal-type-state.svelte';
	import type { PanelSectionProps } from '../types';
	import { data } from './signal-type-data';

	let { entityId }: PanelSectionProps = $props();

	const sts = getSignalTypeState(entityId);
</script>

{#snippet section(signalType: SignalType)}
	<Attribute label="Kind" desc="The kind of the type">
		<SegmentedControl
			name="signal-type-kind"
			options={data.kind.options}
			bind:selectedValue={sts.entity.kind}
			readOnly
		/>
	</Attribute>

	<Divider />

	<Attribute label="Size" desc="The size in bits">
		<NumberInput name="signal-type-size" bind:value={signalType.size} />
	</Attribute>

	<Divider />

	<div class="grid grid-cols-2 gap-5">
		<Attribute label="Min" desc="The minimum value">
			<NumberInput name="signal-type-min" bind:value={signalType.min} />
		</Attribute>

		<Attribute label="Max" desc="The maximum value">
			<NumberInput name="signal-type-max" bind:value={signalType.max} />
		</Attribute>
	</div>

	<Divider />

	<div class="grid grid-cols-2 gap-5">
		<Attribute label="Scale" desc="The scale value">
			<NumberInput name="signal-type-scale" bind:value={signalType.scale} />
		</Attribute>

		<Attribute label="Offset" desc="The offset value">
			<NumberInput name="signal-type-offset" bind:value={signalType.offset} />
		</Attribute>
	</div>
{/snippet}

<section>
	{@render section(sts.entity)}
</section>
