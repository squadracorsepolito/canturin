<script lang="ts">
	import type { SignalEnum } from '$lib/api/canturin';
	import ReferenceTree from '$lib/components/tree/reference-tree.svelte';
	import { getSignalEnumState } from '$lib/state/signal-enum-state.svelte';
	import { getSignalReferenceTree } from '$lib/utils';
	import { text } from '../signal-type/signal-type-text';
	import type { PanelSectionProps } from '../types';

	let { entityId }: PanelSectionProps = $props();

	const ses = getSignalEnumState(entityId);
</script>

{#snippet section(signalEnum: SignalEnum)}
	{#if signalEnum.references && signalEnum.references.length > 0}
		<h3 class="pb-5">{text.headings.refs}</h3>

		<ReferenceTree siblingNodes={getSignalReferenceTree(signalEnum.references)} depth={4} />
	{/if}
{/snippet}

<section>
	{@render section(ses.entity)}
</section>
