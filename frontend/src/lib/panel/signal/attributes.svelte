<script lang="ts">
	import { type Signal, SignalKind } from '$lib/api/canturin';
	import { Attribute, AttributeGroup } from '$lib/components/attribute';
	import Divider from '$lib/components/divider/divider.svelte';
	import { Readonly } from '$lib/components/readonly';
	import { SegmentedControl } from '$lib/components/segmented-control';
	import type { PanelSectionProps } from '../types';
	import StandardAttributes from './standard-attributes.svelte';
	import { getSignalState } from './state.svelte';
	import { signalKindOptions } from './utils';

	let { entityId }: PanelSectionProps = $props();

	const ss = getSignalState(entityId);
</script>

{#snippet section(sig: Signal)}
	<Attribute label="Kind" desc="The kind of the signal">
		<SegmentedControl
			name="signal-kind"
			options={signalKindOptions}
			bind:selectedValue={ss.entity.kind}
			readOnly
		/>
	</Attribute>

	<Divider />

	<AttributeGroup>
		<Attribute label="Size" desc="The size of the signal in bits">
			<Readonly>{sig.size}</Readonly>
		</Attribute>

		<Attribute
			label="Start Pos"
			desc="The start position of the signal within a payload of a message or multiplexed signal"
		>
			<Readonly>{sig.startPos}</Readonly>
		</Attribute>
	</AttributeGroup>

	<Divider />

	{#if sig.kind === SignalKind.SignalKindStandard}
		<StandardAttributes {entityId} signal={sig.standard} />
	{/if}
{/snippet}

<section>
	<h3 class="pb-5">Attributes</h3>

	{@render section(ss.entity)}
</section>
