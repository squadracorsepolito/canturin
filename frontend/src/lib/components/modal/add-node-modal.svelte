<script lang="ts">
	import { Validator } from '$lib/utils/validator.svelte';
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';
	import * as v from 'valibot';
	import Modal from './modal.svelte';
	import { NumberInput } from '../input';

	type Props = {
		onsubmit: (interfaceCount: number) => void;
		trigger: Snippet<[{ getProps: () => HTMLButtonAttributes }]>;
	};

	let { onsubmit, trigger: triggerSnippet }: Props = $props();

	let interfaceCount = $state(1);

	const intCountValidator = new Validator(
		v.pipe(v.number(), v.integer(), v.minValue(1)),
		() => interfaceCount
	);
</script>

<Modal title="Add Node" desc="Pick the number of interfaces of the node you want to add">
	{#snippet trigger({ getProps })}
		{@render triggerSnippet({ getProps })}
	{/snippet}

	{#snippet content()}
		<div class="pt-5">
			<h4 class="pb-1">Size</h4>
			<NumberInput
				name="signal-type-size"
				bind:value={interfaceCount}
				errors={intCountValidator.errors}
			/>
		</div>
	{/snippet}

	{#snippet actions({ close })}
		<button
			onclick={() => {
				onsubmit(interfaceCount);
				close();
			}}
			disabled={intCountValidator.errors ? true : false}
			class="btn btn-primary">Add</button
		>
	{/snippet}
</Modal>
