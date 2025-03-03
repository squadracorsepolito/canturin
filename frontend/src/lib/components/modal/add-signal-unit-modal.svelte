<script lang="ts">
	import { SignalUnitKind } from '$lib/api/canturin';
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';
	import Modal from './modal.svelte';
	import { SegmentedControl } from '../segmented-control';
	import { signalUnitKindOptions } from '$lib/panel/signal-unit/utils';

	type Props = {
		onsubmit: (signalUnitKind: SignalUnitKind) => void;
		trigger: Snippet<[{ getProps: () => HTMLButtonAttributes }]>;
	};

	let { onsubmit, trigger: triggerSnippet }: Props = $props();

	let signalUnitKind = $state<SignalUnitKind>(SignalUnitKind.SignalUnitKindCustom);
</script>

<Modal title="Add Signal Unit" desc="Pick the kind of the signal unit you want to add">
	{#snippet trigger({ getProps })}
		{@render triggerSnippet({ getProps })}
	{/snippet}

	{#snippet content()}
		<SegmentedControl
			name="signal-unit-kind"
			bind:selectedValue={signalUnitKind}
			options={signalUnitKindOptions}
		/>
	{/snippet}

	{#snippet actions({ close })}
		<button
			onclick={() => {
				onsubmit(signalUnitKind);
				close();
			}}
			class="btn btn-primary">Add</button
		>
	{/snippet}
</Modal>
