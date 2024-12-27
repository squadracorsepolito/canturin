<script lang="ts">
	import type { PanelSectionProps } from '../types';
	import type { SignalUnit } from '$lib/api/canturin';
	import { getSignalUnitState } from './state.svelte';
	import ReferenceTree0 from '$lib/components/tree/reference-tree0.svelte';

	let { entityId }: PanelSectionProps = $props();

	const sus = getSignalUnitState(entityId);

	$inspect(sus.entity.references);
</script>

{#snippet section(sigUnit: SignalUnit)}
	{#if sigUnit.references && sigUnit.references.length > 0}
		<h3 class="pb-5">References</h3>

		<ReferenceTree0 nodes={sigUnit.references} depth={4} />
		<!-- <ReferenceTree siblingNodes={getSignalReferenceTree(sigUnit.references)} depth={4} /> -->
	{/if}
{/snippet}

<section>
	{@render section(sus.entity)}
</section>
