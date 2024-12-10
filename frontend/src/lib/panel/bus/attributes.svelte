<script lang="ts">
	import { BusType, type Bus } from '$lib/api/canturin';
	import { Attribute } from '$lib/components/attribute';
	import Divider from '$lib/components/divider/divider.svelte';
	import { SegmentedControl } from '$lib/components/segmented-control';
	import type { SegmentedControlOption } from '$lib/components/segmented-control/types';
	import { Select } from '$lib/components/select';
	import type { PanelSectionProps } from '../types';
	import { getBusState } from './state.svelte';
	import { baudrateItems, getSelectItemFromBaudrate, type BaudrateSelectItem } from './utils';

	let { entityId }: PanelSectionProps = $props();

	const bs = getBusState(entityId);

	const busTypeOptions: SegmentedControlOption[] = [
		{
			label: 'CAN 2.0 A',
			value: BusType.BusTypeCAN2A,
			desc: '11 bit identifier, max payload 8 bytes, max speed 1 Mbit/s'
		}
	];

	let selectedBaudrate = $state(getSelectItemFromBaudrate(bs.entity.baudrate));

	function handleBaudrate(item: BaudrateSelectItem) {
		bs.updateBaudrate(item.valueAsNumber);
	}
</script>

{#snippet section(bus: Bus)}
	<Attribute label="Bus Type" desc="The type of the bus">
		<SegmentedControl
			name="bus-type"
			options={busTypeOptions}
			bind:selectedValue={bus.type}
			readOnly
		/>
	</Attribute>

	<Divider />

	<Attribute label="Baudrate" desc="The baudrate of the bus">
		<Select
			items={baudrateItems}
			labelKey="label"
			valueKey="value"
			bind:selected={selectedBaudrate}
			onselect={handleBaudrate}
		/>
	</Attribute>
{/snippet}

<section>
	<h3 class="pb-5">Attributes</h3>

	{@render section(bs.entity)}
</section>
