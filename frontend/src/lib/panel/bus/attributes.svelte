<script lang="ts">
	import { type Bus } from '$lib/api/canturin';
	import { Attribute } from '$lib/components/attribute';
	import Divider from '$lib/components/divider/divider.svelte';
	import { SegmentedControl } from '$lib/components/segmented-control';
	import { Select } from '$lib/components/select';
	import type { PanelSectionProps } from '../types';
	import { getBusState } from './state.svelte';
	import { baudrateItems, busTypeOptions, type BaudrateSelectItem } from './utils';

	let { entityId }: PanelSectionProps = $props();

	const bs = getBusState(entityId);

	let selectedBaudrate = $state(`${bs.entity.baudrate}`);

	function handleBaudrate(item: BaudrateSelectItem) {
		bs.updateBaudrate(item.valueAsNumber);
	}
</script>

{#snippet section(bus: Bus)}
	<Attribute label="Bus Type" desc="The type of the bus">
		<SegmentedControl name="bus-type" options={busTypeOptions} selectedValue={bus.type} readOnly />
	</Attribute>

	<Divider />

	<Attribute label="Baudrate" desc="The baudrate of the bus">
		<Select
			items={baudrateItems}
			name="bus-baudrate"
			labelKey="label"
			valueKey="value"
			bind:selected={selectedBaudrate}
			onitemselect={handleBaudrate}
		/>
	</Attribute>
{/snippet}

<section>
	<h3 class="pb-5">Attributes</h3>

	{@render section(bs.entity)}
</section>
