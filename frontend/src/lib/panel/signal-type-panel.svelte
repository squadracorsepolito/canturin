<script lang="ts">
	import { type SignalType } from '$lib/api/canturin';
	import ReferenceTree from '$lib/components/tree/reference-tree.svelte';
	import { useSignalType } from '$lib/state/signal-type-state.svelte';
	import { getSignalReferenceTree } from '$lib/utils';
	import Panel from './panel.svelte';

	type Props = {
		entityId: string;
	};

	let { entityId }: Props = $props();

	let state = useSignalType(entityId);

	$effect(() => {
		state.reload(entityId);
	});
</script>

{#snippet sigTypePanel(sigType: SignalType)}
	<section>
		<h3>{sigType.name}</h3>
		<p>{sigType.desc}</p>
	</section>

	{#if sigType.references}
		<section>
			<h4>References</h4>

			<ReferenceTree siblingNodes={getSignalReferenceTree(sigType.references)} depth={4} />
		</section>
	{/if}
{/snippet}

<Panel>
	{#if !state.isLoading && state.signalType}
		{@render sigTypePanel(state.signalType)}
	{/if}
</Panel>
