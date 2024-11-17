<script lang="ts">
	import { SignalTypeIcon } from '$lib/components/icon';
	import { TextEditable } from '$lib/components/editable';
	import { type SignalUnit } from '$lib/api/canturin';
	import Summary from '$lib/components/summary/summary.svelte';
	import Divider from '$lib/components/divider/divider.svelte';
	import { ReferenceTree } from '$lib/components/tree';
	import { useSignalUnit } from '$lib/state/signal-unit-state.svelte';
	import { getSignalReferenceTree } from '$lib/utils';
	import Panel from '../panel.svelte';
	import TextareaEditable from '$lib/components/editable/textarea-editable.svelte';
	import { text } from './signal-unit-text';
	import Attribute from '$lib/components/attribute/attribute.svelte';
	import Readonly from '$lib/components/readonly/readonly.svelte';

	type Props = {
		entityId: string;
	};

	let { entityId }: Props = $props();

	let state = useSignalUnit(entityId);

	$effect(() => {
		state.reload(entityId);
	});

	function handleName(name: string) {
		state.updateName(name);
	}

	function handleDesc(desc: string) {
		state.updateDesc(desc);
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

	<div class="grid grid-cols-2 gap-5">
		<Attribute {...text.symbol}>
			<Readonly>
				<span class="font-medium">
					{sigUnit.symbol}
				</span>
			</Readonly>
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
		<section>
			<h4>References</h4>

			<ReferenceTree siblingNodes={getSignalReferenceTree(sigUnit.references)} depth={4} />
		</section>
	{/if}
{/snippet}

<Panel>
	{#if !state.isLoading && state.signalUnit}
		{@render sigUnitPanel(state.signalUnit)}
	{/if}
</Panel>
