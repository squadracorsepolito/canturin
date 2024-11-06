<script lang="ts">
	import type { SignalType } from '$lib/api/canturin';
	import { Attribute } from '$lib/components/attribute';
	import Divider from '$lib/components/divider/divider.svelte';
	import NumberEditable from '$lib/components/editable/number-editable.svelte';
	import { Readonly, ToggleReadonly } from '$lib/components/readonly';
	import SegmentedControl from '$lib/components/segmented-control/segmented-control.svelte';
	import { getSignalTypeState } from '$lib/state/signal-type-state.svelte';
	import { z } from 'zod';
	import type { PanelSectionProps } from '../types';
	import { text } from './signal-type-text';
	import { Switch } from '$lib/components/switch';

	let { entityId }: PanelSectionProps = $props();

	const sts = getSignalTypeState(entityId);

	const minSchema = z.object({
		min: z.number()
	});

	function validateMin(min: number) {
		const res = minSchema.safeParse({ min });
		if (res.success) {
			return undefined;
		}
		return res.error.flatten().fieldErrors.min;
	}

	function handleMin(min: number) {
		sts.updateMin(min);
	}

	const maxSchema = z.object({
		max: z.number()
	});

	function validateMax(max: number) {
		const res = maxSchema.safeParse({ max });
		if (res.success) {
			return undefined;
		}
		return res.error.flatten().fieldErrors.max;
	}

	function handleMax(max: number) {
		sts.updateMax(max);
	}

	const scaleSchema = z.object({
		scale: z.number()
	});

	function validateScale(scale: number) {
		const res = scaleSchema.safeParse({ scale });
		if (res.success) {
			return undefined;
		}
		return res.error.flatten().fieldErrors.scale;
	}

	function handleScale(scale: number) {
		sts.updateScale(scale);
	}

	const offsetSchema = z.object({
		offset: z.number()
	});

	function validateOffset(offset: number) {
		const res = offsetSchema.safeParse({ offset });
		if (res.success) {
			return undefined;
		}
		return res.error.flatten().fieldErrors.offset;
	}

	function handleOffset(offset: number) {
		sts.updateOffset(offset);
	}
</script>

{#snippet section(signalType: SignalType)}
	<Attribute label="Kind" desc="The kind of the type">
		<SegmentedControl
			name="signal-type-kind"
			options={text.kind.options}
			bind:selectedValue={sts.entity.kind}
			readOnly
		/>
	</Attribute>

	<Divider />

	<div class="grid grid-cols-2 gap-5">
		<Attribute {...text.size}>
			<Readonly>
				<span class="font-medium">
					{signalType.size}
				</span>
			</Readonly>
		</Attribute>

		<Attribute {...text.signed}>
			<Switch checked={signalType.signed} readOnly />
		</Attribute>
	</div>

	<Divider />

	<div class="grid grid-cols-2 gap-5">
		<Attribute {...text.min}>
			<NumberEditable
				validator={validateMin}
				initialValue={signalType.min}
				name="signal-type-min"
				onsubmit={handleMin}
				placeholder={text.min.label}
			/>
		</Attribute>

		<Attribute {...text.max}>
			<NumberEditable
				validator={validateMax}
				initialValue={signalType.max}
				name="signal-type-max"
				onsubmit={handleMax}
				placeholder={text.max.label}
			/>
		</Attribute>
	</div>

	<Divider />

	<div class="grid grid-cols-2 gap-5">
		<Attribute {...text.scale}>
			<NumberEditable
				validator={validateScale}
				initialValue={signalType.scale}
				name="signal-type-scale"
				onsubmit={handleScale}
				placeholder={text.scale.label}
			/>
		</Attribute>

		<Attribute {...text.offset}>
			<NumberEditable
				validator={validateOffset}
				initialValue={signalType.offset}
				name="signal-type-offset"
				onsubmit={handleOffset}
				placeholder={text.offset.label}
			/>
		</Attribute>
	</div>
{/snippet}

<section>
	<h3 class="pb-5">{text.headings.attributes}</h3>

	{@render section(sts.entity)}
</section>
