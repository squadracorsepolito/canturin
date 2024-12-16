<script lang="ts">
	import type { SignalUnit } from '$lib/api/canturin';
	import { Attribute } from '$lib/components/attribute';
	import { TextEditable } from '$lib/components/editable';
	import { Readonly } from '$lib/components/readonly';
	import type { PanelSectionProps } from '../types';
	import { getSignalUnitState } from './signal-unit-state.svelte';

	let { entityId }: PanelSectionProps = $props();

	let sus = getSignalUnitState(entityId);

	function handleSymbol(sym: string) {
		sus.updateSymbol(sym);
	}
</script>

{#snippet section(sigUnit: SignalUnit)}
	<div class="grid grid-cols-2 gap-5 pt-8">
		<Attribute label="Symbol" desc="Symbol of the Unit">
			<TextEditable name="signal-unit-symbol" bind:value={sigUnit.symbol} oncommit={handleSymbol} />
		</Attribute>

		<Attribute label="Reference Count">
			<Readonly>
				<span class="font-medium">
					{sigUnit.referenceCount}
				</span>
			</Readonly>
		</Attribute>
	</div>
{/snippet}

<section>
	<h3 class="pb-5">Attributes</h3>

	{@render section(sus.entity)}
</section>
