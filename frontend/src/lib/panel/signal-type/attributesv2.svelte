<script lang="ts">
	import type { SignalType } from '$lib/api/canturin';
	import { SignalTypeKind } from '$lib/api/github.com/squadracorsepolito/acmelib/models';
	import { Attribute } from '$lib/components/attribute';
	import Divider from '$lib/components/divider/divider.svelte';
	import { NumberInput } from '$lib/components/input';
	import SegmentedControl from '$lib/components/segmented-control/segmented-control.svelte';
	import { getSignalTypeState } from '$lib/state/signal-type-state.svelte';
	import { z } from 'zod';
	import type { PanelSectionProps } from '../types';

	let { entityId }: PanelSectionProps = $props();

	const sts = getSignalTypeState(entityId);

	const kindOptions = [
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
	];

	function getKindValue() {
		switch (sts.entity.kind) {
			case SignalTypeKind.SignalTypeKindFlag:
				return 'flag';
			case SignalTypeKind.SignalTypeKindInteger:
				return 'integer';
			case SignalTypeKind.SignalTypeKindDecimal:
				return 'decimal';
			case SignalTypeKind.SignalTypeKindCustom:
				return 'custom';
			default:
				return 'custom';
		}
	}

	let kind = $state(getKindValue());
</script>

{#snippet section(signalType: SignalType)}
	{#if sts.isDraft}
		<Attribute label="Kind" desc="The kind of the type">
			<SegmentedControl name="signal-type-kind" options={kindOptions} bind:selectedValue={kind} />
		</Attribute>
	{:else}
		<Attribute label="Kind" desc="The kind of the type">
			<SegmentedControl
				name="signal-type-kind"
				options={kindOptions}
				bind:selectedValue={kind}
				readOnly
			/>
		</Attribute>
	{/if}

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

	<button onclick={() => (sts.isDraft = !sts.isDraft)}>draft</button>
{/snippet}

<section>
	{@render section(sts.entity)}
</section>
