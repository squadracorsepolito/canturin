<script lang="ts">
	import { SignalTypeIcon } from '$lib/components/icon';
	import type { PanelSectionProps } from '../types';
	import { TextEditable } from '$lib/components/editable';
	import { type SignalUnit } from '$lib/api/canturin';
	import { ReferenceTree } from '$lib/components/tree';
	import { getSignalUnitState } from '$lib/state/signal-unit-state.svelte';
	import { getSignalReferenceTree } from '$lib/utils';
	import TextareaEditable from '$lib/components/editable/textarea-editable.svelte';
	import { text } from './signal-unit-text';
	import Attribute from '$lib/components/attribute/attribute.svelte';
	import Readonly from '$lib/components/readonly/readonly.svelte';

	let { entityId }: PanelSectionProps = $props();

	let state = getSignalUnitState(entityId);

	$effect(() => {
		state.reload(entityId);
	});

	function handleName(name: string) {
		state.updateName(name);
	}

	function handleDesc(desc: string) {
		state.updateDesc(desc);
	}

	function handleSymbol(sym: string) {
		state.updateSymbol(sym);
	}

</script>

{#snippet sigUnitPanel(sigUnit: SignalUnit)}
	<!-- name box -->
	<div class="flex gap-2 items-center">
		<SignalTypeIcon width="48" height="48" />
		
		<TextEditable
			validator={() => undefined}
			name="signal-type-name"
			initialValue={sigUnit.name}
			onsubmit={handleName}
			placeholder="Name"
		/>
	</div>

	<!-- description box -->
	<div class="pt-8">
		<TextareaEditable
			initialValue={sigUnit.desc}
			name="signal-type-desc"
			triggerLabel={text.buttons.heading.descTriggerLabel}
			onsubmit={handleDesc}
		/>
	</div>

	<!-- attribute grid -->
	<div class="grid grid-cols-2 gap-5 pt-8">
		<Attribute {...text.symbol}>
			<TextEditable
				validator={() => undefined}
				name="Symbol"
				initialValue={sigUnit.symbol}
				onsubmit={handleSymbol}
				placeholder={sigUnit.symbol}
			/>
		</Attribute>
	
		<Attribute {...text.refCount}>
			<Readonly>
				<span class="font-medium">
					{sigUnit.referenceCount}
				</span>
			</Readonly>
		</Attribute>
	</div>

	{#if sigUnit.references}
		<div class="pt-8">
			<section>
				<h3 class="pb-5">{text.headings.refs}</h3>

				<ReferenceTree siblingNodes={getSignalReferenceTree(sigUnit.references)} depth={4} />
			</section>
		</div>
	{/if}
{/snippet}

<section>
	{@render sigUnitPanel(state.entity)}
</section>
