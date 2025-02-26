<script lang="ts">
	import { SignalKind } from '$lib/api/canturin';
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';
	import Modal from './modal.svelte';
	import { SegmentedControl } from '../segmented-control';
	import { signalKindOptions } from '$lib/panel/signal/utils';

	type Props = {
		onsubmit: (signalKind: SignalKind) => void;
		trigger: Snippet<[{ getProps: () => HTMLButtonAttributes }]>;
	};

	let { onsubmit, trigger: triggerSnippet }: Props = $props();

	let signalKind = $state<SignalKind>(SignalKind.SignalKindStandard);
</script>

<Modal title="Add Signal" desc="Pick the kind of the signal you want to add">
	{#snippet trigger({ getProps })}
		{@render triggerSnippet({ getProps })}
	{/snippet}

	{#snippet content()}
		<SegmentedControl
			name="signal-kind"
			bind:selectedValue={signalKind}
			options={signalKindOptions}
		/>
	{/snippet}

	{#snippet actions({ close })}
		<button
			onclick={() => {
				onsubmit(signalKind);
				close();
			}}
			class="btn btn-primary">Add</button
		>
	{/snippet}
</Modal>
