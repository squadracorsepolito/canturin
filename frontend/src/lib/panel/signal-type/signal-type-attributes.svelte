<script lang="ts">
	import { SignalTypeKind } from '$lib/api/github.com/squadracorsepolito/acmelib';

	import AttributeGroup from '$lib/components/attribute/attribute-group.svelte';
	import NumberAttribute from '$lib/components/attribute/number-attribute.svelte';
	import type { PanelSectionProps } from '../types';
	import {
		kindSchema,
		maxSchema,
		minSchema,
		offsetSchema,
		scaleSchema,
		sizeSchema
	} from './signal-type-schema';
	import { getSignalTypeState } from '../../state/signal-type-state.svelte';
	import { type SignalType } from '$lib/api/canturin';
	import { RadioAttribute } from '$lib/components/attribute';
	import type { RadioInputOption } from '$lib/components/input/types';

	let { entityId }: PanelSectionProps = $props();

	const sts = getSignalTypeState(entityId);

	function handleKind(kindId: number) {
		sts.updateKind(kindId);
	}

	function handleSize(size: number) {
		sts.updateSize(size);
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

	function getSelectedKindId(kind: SignalTypeKind) {
		switch (kind) {
			case SignalTypeKind.SignalTypeKindFlag:
				return 1;
			case SignalTypeKind.SignalTypeKindInteger:
				return 2;
			case SignalTypeKind.SignalTypeKindDecimal:
				return 3;
			case SignalTypeKind.SignalTypeKindCustom:
				return 4;
			default:
				return 4;
		}
	}

	let kindOptions = $derived.by(() => {
		const opts: RadioInputOption[] = [
			{
				id: getSelectedKindId(SignalTypeKind.SignalTypeKindFlag),
				label: 'Flag',
				name: 'flag'
			},
			{
				id: getSelectedKindId(SignalTypeKind.SignalTypeKindInteger),
				label: 'Integer',
				name: 'integer'
			},
			{
				id: getSelectedKindId(SignalTypeKind.SignalTypeKindDecimal),
				label: 'Decimal',
				name: 'decimal'
			},
			{
				id: getSelectedKindId(SignalTypeKind.SignalTypeKindCustom),
				label: 'Custom',
				name: 'custom'
			}
		];

		const size = sts.entity.size;
		if (size > 1) {
			opts[0].disabled = true;
		}

		const min = sts.entity.min;
		const max = sts.entity.min;
		const scale = sts.entity.scale;
		const offset = sts.entity.offset;
		if (
			!Number.isInteger(min) ||
			!Number.isInteger(max) ||
			!Number.isInteger(scale) ||
			!Number.isInteger(offset)
		) {
			opts[1].disabled = true;
		}
		return opts;
	});
</script>

{#snippet section(signalType: SignalType)}
	<AttributeGroup>
		{#snippet attributes()}
			<RadioAttribute
				prefix="signal-type"
				name="Kind"
				selected={getSelectedKindId(signalType.kind)}
				schema={kindSchema}
				options={kindOptions}
				onchange={handleKind}
				desc="The kind of the type"
			/>

			<NumberAttribute
				prefix="signal-type"
				name="Size"
				value={signalType.size}
				schema={sizeSchema}
				onchange={handleSize}
				desc="The size in bits"
			/>

			<NumberAttribute
				prefix="signal-type"
				name="Min"
				value={signalType.min}
				schema={minSchema}
				onchange={handleMin}
				desc="The minimum value"
			/>

			<NumberAttribute
				prefix="signal-type"
				name="Max"
				value={signalType.max}
				schema={maxSchema}
				onchange={handleMax}
				desc="The maximum value"
			/>

			<NumberAttribute
				prefix="signal-type"
				name="Scale"
				value={signalType.scale}
				schema={scaleSchema}
				onchange={handleScale}
				desc="The scale value"
			/>

			<NumberAttribute
				prefix="signal-type"
				name="Offset"
				value={signalType.offset}
				schema={offsetSchema}
				onchange={handleOffset}
				desc="The offset value"
			/>
		{/snippet}
	</AttributeGroup>
{/snippet}

<section>
	{@render section(sts.entity)}
</section>
