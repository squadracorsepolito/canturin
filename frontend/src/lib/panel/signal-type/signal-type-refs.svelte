<script lang="ts">
	import { ReferenceTree } from '$lib/components/tree';
	import { getSignalReferenceTree } from '$lib/utils';
	import type { PanelSectionProps } from '../types';
	import { getSignalTypeState } from '../../state/signal-type-state.svelte';
	import { type SignalType } from '$lib/api/canturin';

	let { entityId }: PanelSectionProps = $props();

	const sts = getSignalTypeState(entityId);
</script>

{#snippet section(signalType: SignalType)}
	{#if signalType.references}
		<section>
			<h4>References</h4>

			<ReferenceTree siblingNodes={getSignalReferenceTree(signalType.references)} depth={4} />
		</section>
	{/if}
{/snippet}

{@render section(sts.entity)}
