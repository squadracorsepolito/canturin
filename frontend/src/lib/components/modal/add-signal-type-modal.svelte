<script lang="ts">
	import type { HTMLButtonAttributes } from 'svelte/elements';
	import Modal from './modal.svelte';
	import type { Snippet } from 'svelte';
	import { SignalTypeKind } from '$lib/api/canturin';
	import { Validator } from '$lib/utils/validator.svelte';
	import * as v from 'valibot';
	import { SegmentedControl } from '../segmented-control';
	import { signalTypeKindOptions } from '$lib/panel/signal-type/utils';
	import NumberInput from '../input/number-input.svelte';

	type Props = {
		onsubmit: (signalTypeKind: SignalTypeKind, size: number) => void;
		trigger: Snippet<[{ getProps: () => HTMLButtonAttributes }]>;
	};

	let { onsubmit, trigger: triggerSnippet }: Props = $props();

	let signalTypeKind = $state<SignalTypeKind>(SignalTypeKind.SignalTypeKindFlag);

	let size = $state(1);

	const sizeValidator = new Validator(
		v.pipe(v.number(), v.integer(), v.minValue(1), v.maxValue(64)),
		() => size
	);
</script>

<Modal
	title="Add Signal Type"
	desc="Pick the kind and set the size of the signal type you want to add"
>
	{#snippet trigger({ getProps })}
		{@render triggerSnippet({ getProps })}
	{/snippet}

	{#snippet content()}
		<SegmentedControl
			name="signal-type-kind"
			bind:selectedValue={signalTypeKind}
			options={signalTypeKindOptions}
		/>

		<div class="pt-5">
			<h4 class="pb-1">Size</h4>
			<NumberInput name="signal-type-size" bind:value={size} errors={sizeValidator.errors} />
		</div>
	{/snippet}

	{#snippet actions({ close })}
		<button
			onclick={() => {
				onsubmit(signalTypeKind, size);
				close();
			}}
			disabled={sizeValidator.errors ? true : false}
			class="btn btn-primary">Add</button
		>
	{/snippet}
</Modal>
