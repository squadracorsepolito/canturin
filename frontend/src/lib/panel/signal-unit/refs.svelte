<script lang="ts">
	import { ReferenceTree } from '$lib/components/tree';
	import { getSignalReferenceTree } from '$lib/utils';
	import type { PanelSectionProps } from '../types';
	import type { SignalUnit } from '$lib/api/canturin';
	import { getSignalUnitState } from './signal-unit-state.svelte';

	let { entityId }: PanelSectionProps = $props();

	const sus = getSignalUnitState(entityId);
</script>

{#snippet section(sigUnit: SignalUnit)}
	{#if sigUnit.references && sigUnit.references.length > 0}
		<h3 class="pb-5">References</h3>

		<ReferenceTree siblingNodes={getSignalReferenceTree(sigUnit.references)} depth={4} />
	{/if}
{/snippet}

<section>
	{@render section(sus.entity)}
</section>
