<script lang="ts">
	import { type SignalUnit } from '$lib/api/canturin';
	import Summary from '$lib/components/summary/summary.svelte';
	import { ReferenceTree } from '$lib/components/tree';
	import { useSignalUnit } from '$lib/state/signal-unit-state.svelte';
	import { getSignalReferenceTree } from '$lib/utils';
	import Panel from './panel.svelte';

	type Props = {
		entityId: string;
	};

	let { entityId }: Props = $props();

	let state = useSignalUnit(entityId);

	$effect(() => {
		state.reload(entityId);
	});
</script>

{#snippet sigUnitPanel(sigUnit: SignalUnit)}
	<section>
		<h3>{sigUnit.name}</h3>
		<p>{sigUnit.desc}</p>
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
