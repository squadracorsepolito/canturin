<script lang="ts">
	import {
		MessageService,
		SignalTypeService,
		SignalUnitService,
		type StandardSignal
	} from '$lib/api/canturin';
	import { Attribute } from '$lib/components/attribute';
	import { LinkButton } from '$lib/components/button';
	import { Combobox } from '$lib/components/combobox';
	import Divider from '$lib/components/divider/divider.svelte';
	import layoutStateSvelte from '$lib/state/layout-state.svelte';
	import { signalTypeKindLables } from '../signal-type/utils';
	import { signalUnitKindLabels } from '../signal-unit/utils';
	import { getSignalState } from './state.svelte';

	type Props = {
		entityId: string;
		signal: StandardSignal;
	};

	let { entityId, signal }: Props = $props();

	const ss = getSignalState(entityId);

	function handleSignalType(sigTypeEntId: string) {
		ss.updateSignalType(sigTypeEntId);
	}

	function handleSignalUnit(sigUnitEntId: string) {
		ss.updateSignalUnit(sigUnitEntId);
	}
</script>

{#await Promise.all( [SignalTypeService.ListBrief(), MessageService.GetSpaceLeft(ss.entity.parentMessage.entityId)] ) then [signalTypes, spaceLeft]}
	{#if signalTypes}
		<Attribute
			label="Type"
			desc="The type of the signal. It specifies the size, min, max, scale, and offset of the signal"
		>
			<Combobox
				name="standard-signal-type"
				items={signalTypes.map((sigType) => ({
					...sigType,
					desc: `${signalTypeKindLables[sigType.kind]} of ${sigType.size} bits`
				}))}
				labelKey="name"
				valueKey="entityId"
				descKey="desc"
				bind:selected={signal.signalType.entityId}
				onselect={handleSignalType}
				filter={(item) => {
					return item.size > spaceLeft + signal.signalType.size;
				}}
			/>

			<div class="pt-2">
				<LinkButton
					label={`Go to ${signal.signalType.name}`}
					onclick={() => layoutStateSvelte.openPanel('signal_type', signal.signalType.entityId)}
				/>
			</div>
		</Attribute>
	{/if}
{/await}

<Divider />

{#await SignalUnitService.ListBrief() then signalUnits}
	{#if signalUnits}
		<Attribute label="Unit" desc="The unit of the signal">
			<Combobox
				name="standard-signal-unit"
				items={signalUnits.map((sigUnit) => ({
					...sigUnit,
					desc: signalUnitKindLabels[sigUnit.kind]
				}))}
				labelKey="name"
				valueKey="entityId"
				descKey="desc"
				bind:selected={signal.signalUnit.entityId}
				onselect={handleSignalUnit}
				onclear={() => handleSignalUnit('')}
			/>

			{#if signal.signalUnit.entityId}
				<div class="pt-2">
					<LinkButton
						label={`Go to ${signal.signalUnit.name}`}
						onclick={() => layoutStateSvelte.openPanel('signal_unit', signal.signalUnit.entityId)}
					/>
				</div>
			{/if}
		</Attribute>
	{/if}
{/await}
