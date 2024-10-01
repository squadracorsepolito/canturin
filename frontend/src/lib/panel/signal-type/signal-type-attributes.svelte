<script lang="ts">
	import NumberAttribute from '$lib/components/attribute/number-attribute.svelte';
	import type { PanelSectionProps } from '../types';
	import { getSignalTypeContext } from './signal-type-context.svelte';
	import {
		maxSchema,
		minSchema,
		offsetSchema,
		scaleSchema,
		sizeSchema
	} from './signal-type-schema';

	let { entityId }: PanelSectionProps = $props();

	const { signalType } = getSignalTypeContext(entityId);

	function handleSize(size: number) {}

	function handleMin(min: number) {}

	function handleMax(max: number) {}

	function handleScale(scale: number) {}

	function handleOffset(offset: number) {}
</script>

<section class="grid grid-cols-4">
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
</section>
