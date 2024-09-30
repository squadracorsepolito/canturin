<script lang="ts">
	import { type SignalUnit } from '$lib/api/canturin';
	import SignalUnitIcon from '$lib/components/icon/signal-type-icon.svelte';
	import NameInput from '$lib/components/input/name-input.svelte';
	import Summary from '$lib/components/summary/summary.svelte';
	import DescTextarea from '$lib/components/textarea/desc-textarea.svelte';
	import { ReferenceTree } from '$lib/components/tree';
	import { useSignalUnit } from '$lib/state/signal-unit-state.svelte';
	import { getSignalReferenceTree } from '$lib/utils';
	import Panel from './panel.svelte';
	import { onMount } from 'svelte';

	type Props = {
		entityId: string;
	};

	let { entityId }: Props = $props();
	let state = useSignalUnit(entityId);
	let invalidNames: string[] = [];

	$effect(() => {
		state.reload(entityId);
	});

	function handleDesc(desc: string) {
		state.updateDesc(desc);
	}

	onMount(async () => {
		invalidNames = await state.getInvalidNames();
	});
</script>

{#snippet sigUnitPanel(sigUnit: SignalUnit)}
	<section>
		<NameInput
			label="Signal Unit Name"
			prefixName="signal_unit"
			initialValue={sigUnit.name}
			onSubmit={(n) => {
				state.reload(n);
			}}
			invalidNames={invalidNames}
		>
			<SignalUnitIcon height="48" width="48" />
		</NameInput>

		<div class="pt-5">
			<DescTextarea initialDesc={sigUnit.desc} onSubmit={handleDesc} />
		</div>
	</section>

	<section>
		<Summary
			infos={[
				{
					title: 'Symbol',
					value: sigUnit.symbol,
					desc: 'The symbol of the unit'
				},
				{
					title: 'Reference Count',
					value: sigUnit.referenceCount,
					desc: 'The number of signals using the unit'
				}
			]}
		/>
	</section>

	{#if sigUnit.references}
		<section>
			<h4>References</h4>

			<ReferenceTree siblingNodes={getSignalReferenceTree(sigUnit.references)} depth={4} />
		</section>
	{/if}
{/snippet}

<Panel>
	{#if !state.isLoading && state.signalUnit}
		{@render sigUnitPanel(state.signalUnit)}
	{/if}
</Panel>
