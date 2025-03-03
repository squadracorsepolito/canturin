<script lang="ts">
	import { type SignalUnit } from '$lib/api/canturin';
	import { Attribute } from '$lib/components/attribute';
	import Divider from '$lib/components/divider/divider.svelte';
	import { TextEditable } from '$lib/components/editable';
	import { SegmentedControl } from '$lib/components/segmented-control';
	import type { PanelSectionProps } from '../types';
	import { getSignalUnitState } from './state.svelte';
	import { signalUnitKindOptions } from './utils';

	let { entityId }: PanelSectionProps = $props();

	let sus = getSignalUnitState(entityId);

	function handleSymbol(sym: string) {
		sus.updateSymbol(sym);
	}

	function handleKind(kind: string) {
		sus.updateKind(kind);
	}
</script>

{#snippet section(sigUnit: SignalUnit)}
	<Attribute label="Kind" desc="The kind of the type">
		<SegmentedControl
			name="signal-type-kind"
			options={signalUnitKindOptions}
			bind:selectedValue={sus.entity.kind}
			onchange={handleKind}
		/>
	</Attribute>

	<Divider />

	<Attribute label="Symbol" desc="Symbol of the Unit">
		<TextEditable name="signal-unit-symbol" bind:value={sigUnit.symbol} oncommit={handleSymbol} />
	</Attribute>
{/snippet}

<section>
	<h3 class="pb-5">Attributes</h3>

	{@render section(sus.entity)}
</section>
