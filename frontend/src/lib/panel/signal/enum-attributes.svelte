<script lang="ts">
	import { MessageService, SignalEnumService, type EnumSignal } from '$lib/api/canturin';
	import { Attribute } from '$lib/components/attribute';
	import { LinkButton } from '$lib/components/button';
	import { Combobox } from '$lib/components/combobox';
	import layout from '$lib/state/layout-state.svelte';
	import { getSignalState } from './state.svelte';

	type Props = {
		entityId: string;
		signal: EnumSignal;
	};

	let { entityId, signal }: Props = $props();

	const ss = getSignalState(entityId);

	function handleSignalEnum(sigEnumEntId: string) {
		ss.updateSignalEnum(sigEnumEntId);
	}
</script>

{#await Promise.all( [SignalEnumService.ListBrief(), MessageService.GetSpaceLeft(ss.entity.parentMessage.entityId)] ) then [signalEnums, spaceLeft]}
	{#if signalEnums}
		<Attribute
			label="Enum"
			desc="The enum of the signal. It specifies the values that the signal can have"
		>
			<Combobox
				name="enum-signal-enum"
				items={signalEnums.map((sigEnum) => ({
					...sigEnum,
					desc: `Size of ${sigEnum.size} bits`
				}))}
				labelKey="name"
				valueKey="entityId"
				descKey="desc"
				bind:selected={signal.signalEnum.entityId}
				onselect={handleSignalEnum}
				filter={(item) => {
					return item.size - signal.signalEnum.size > spaceLeft;
				}}
			/>

			<div class="pt-2">
				<LinkButton
					label={`Go to ${signal.signalEnum.name}`}
					onclick={() => layout.openPanel('signal_enum', signal.signalEnum.entityId)}
				/>
			</div>
		</Attribute>
	{/if}
{/await}
