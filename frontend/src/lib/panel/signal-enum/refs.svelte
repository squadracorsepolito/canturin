<script lang="ts">
	import type { SignalEnum } from '$lib/api/canturin';
	import ReferenceTree0 from '$lib/components/tree/reference-tree0.svelte';
	import { getSignalEnumState } from '$lib/panel/signal-enum/state.svelte';
	import type { PanelSectionProps } from '../types';

	let { entityId }: PanelSectionProps = $props();

	const ses = getSignalEnumState(entityId);

	$inspect(ses.entity.references);
</script>

{#snippet section(signalEnum: SignalEnum)}
	{#if signalEnum.references && signalEnum.references.length > 0}
		<h3 class="pb-5">References</h3>

		<ReferenceTree0 nodes={signalEnum.references} depth={4} />

		<!-- <ReferenceTree siblingNodes={getSignalReferenceTree(signalEnum.references)} depth={4} /> -->
	{/if}
{/snippet}

<section>
	{@render section(ses.entity)}
</section>
