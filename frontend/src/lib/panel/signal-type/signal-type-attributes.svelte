<script lang="ts">
	import AttributeGroup from '$lib/components/attribute/attribute-group.svelte';
	import NumberAttribute from '$lib/components/attribute/number-attribute.svelte';
	import type { PanelSectionProps } from '../types';
	import {
		maxSchema,
		minSchema,
		offsetSchema,
		scaleSchema,
		sizeSchema
	} from './signal-type-schema';
	import { getSignalTypeState } from '../../state/signal-type-state.svelte';
	import { type SignalType } from '$lib/api/canturin';

	let { entityId }: PanelSectionProps = $props();

	const sts = getSignalTypeState(entityId);

	function handleSize(size: number) {}

	function handleMin(min: number) {
		sts.updateMin(min);
	}

	function handleMax(max: number) {
		sts.updateMax(max);
	}

	function handleScale(scale: number) {
		sts.updateScale(scale);
	}

	function handleOffset(offset: number) {
		sts.updateOffset(offset);
	}
</script>

{#snippet section(signalType: SignalType)}
	<AttributeGroup>
		{#snippet attributes()}
			<NumberAttribute
				prefix="signal-type"
				name="Size"
				value={signalType.size}
				schema={sizeSchema}
				onchange={handleSize}
				description="The size in bits"
			/>

			<NumberAttribute
				prefix="signal-type"
				name="Min"
				value={signalType.min}
				schema={minSchema}
				onchange={handleMin}
				description="The minimum value"
			/>

			<NumberAttribute
				prefix="signal-type"
				name="Max"
				value={signalType.max}
				schema={maxSchema}
				onchange={handleMax}
				description="The maximum value"
			/>

			<NumberAttribute
				prefix="signal-type"
				name="Scale"
				value={signalType.scale}
				schema={scaleSchema}
				onchange={handleScale}
				description="The scale value"
			/>

			<NumberAttribute
				prefix="signal-type"
				name="Offset"
				value={signalType.offset}
				schema={offsetSchema}
				onchange={handleOffset}
				description="The offset value"
			/>
		{/snippet}
	</AttributeGroup>
{/snippet}

<section>
	{@render section(sts.entity)}
</section>
