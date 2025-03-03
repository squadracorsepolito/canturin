<script lang="ts">
	import type { SignalType } from '$lib/api/canturin';
	import { Attribute } from '$lib/components/attribute';
	import Divider from '$lib/components/divider/divider.svelte';
	import { NumberEditable } from '$lib/components/editable';
	import { Readonly } from '$lib/components/readonly';
	import SegmentedControl from '$lib/components/segmented-control/segmented-control.svelte';
	import { getSignalTypeState } from '$lib/panel/signal-type/state.svelte';
	import type { PanelSectionProps } from '../types';
	import { Switch } from '$lib/components/switch';
	import { signalTypeKindOptions } from './utils';
	import { Validator } from '$lib/utils/validator.svelte';
	import * as v from 'valibot';

	let { entityId }: PanelSectionProps = $props();

	const sts = getSignalTypeState(entityId);

	const minValidator = new Validator(v.number(), () => sts.entity.min);

	const maxValidator = new Validator(v.number(), () => sts.entity.max);

	const scaleValidator = new Validator(v.number(), () => sts.entity.scale);

	const offsetValidator = new Validator(v.number(), () => sts.entity.offset);

	function handleSigned(signed: boolean) {
		sts.updateSigned(signed);
	}

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
	<Attribute label="Kind" desc="The kind of the type">
		<SegmentedControl
			name="signal-type-kind"
			options={signalTypeKindOptions}
			bind:selectedValue={sts.entity.kind}
			readOnly
		/>
	</Attribute>

	<Divider />

	<div class="grid grid-cols-2 gap-5">
		<Attribute label="Size" desc="The size in bits">
			<Readonly>
				<span class="font-medium">
					{signalType.size}
				</span>
			</Readonly>
		</Attribute>

		<Attribute label="Signed" desc="Whether the value is signed">
			<Switch checked={signalType.signed} oncheckedchange={handleSigned} />
		</Attribute>
	</div>

	<Divider />

	<div class="grid grid-cols-2 gap-5">
		<Attribute label="Min" desc="The minimum value">
			<NumberEditable
				bind:value={signalType.min}
				name="signal-type-min"
				oncommit={handleMin}
				errors={minValidator.errors}
			/>
		</Attribute>

		<Attribute label="Max" desc="The maximum value">
			<NumberEditable
				bind:value={signalType.max}
				name="signal-type-max"
				oncommit={handleMax}
				errors={maxValidator.errors}
			/>
		</Attribute>
	</div>

	<Divider />

	<div class="grid grid-cols-2 gap-5">
		<Attribute label="Scale" desc="The scale factor">
			<NumberEditable
				bind:value={signalType.scale}
				name="signal-type-scale"
				oncommit={handleScale}
				errors={scaleValidator.errors}
			/>
		</Attribute>

		<Attribute label="Offset" desc="The offset value">
			<NumberEditable
				bind:value={signalType.offset}
				name="signal-type-offset"
				oncommit={handleOffset}
				errors={offsetValidator.errors}
			/>
		</Attribute>
	</div>
{/snippet}

<section>
	<h3 class="pb-5">Attributes</h3>

	{@render section(sts.entity)}
</section>
