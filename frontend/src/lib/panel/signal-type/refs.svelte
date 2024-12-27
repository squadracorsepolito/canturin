<script lang="ts">
	import type { PanelSectionProps } from '../types';
	import { getSignalTypeState } from './state.svelte';
	import { type SignalType } from '$lib/api/canturin';
	import ReferenceTree0 from '$lib/components/tree/reference-tree0.svelte';

	let { entityId }: PanelSectionProps = $props();

	const sts = getSignalTypeState(entityId);
</script>

{#snippet section(signalType: SignalType)}
	{#if signalType.references && signalType.references.length > 0}
		<h3 class="pb-5">References</h3>

		<ReferenceTree0 nodes={signalType.references} depth={4} />
		<!-- <ReferenceTree siblingNodes={getSignalReferenceTree(signalType.references)} depth={4} /> -->
	{/if}
{/snippet}

<section>
	{@render section(sts.entity)}
</section>
