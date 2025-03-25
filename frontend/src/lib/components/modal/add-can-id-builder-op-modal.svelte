<script lang="ts">
	import Modal from './modal.svelte';
	import { CANIDBuilderOpKind } from '$lib/api/canturin';
	import type { ModalProps } from './types';
	import { SegmentedControl } from '../segmented-control';
	import { canIdBuilderOpKindOptions } from '$lib/panel/can-id-builder/utils';
	import { Validator } from '$lib/utils/validator.svelte';
	import * as v from 'valibot';
	import { NumberInput } from '../input';

	type Props = ModalProps<{ opKind: CANIDBuilderOpKind; opFrom: number; opLen: number }> & {
		maxSize: number;
	};

	let { onsubmit, trigger: triggerSnippet, maxSize }: Props = $props();

	let opKind = $state(CANIDBuilderOpKind.CANIDBuilderOpKindMessagePriority);

	let opFrom = $state(0);
	const opFromValidator = new Validator(
		v.pipe(v.number(), v.integer(), v.minValue(0), v.maxValue(maxSize - 1)),
		() => opFrom
	);

	function getOpLenMaxValue() {
		return maxSize - opFrom;
	}

	let opLen = $state(1);
	const opLenValidator = new Validator(
		v.pipe(
			v.number(),
			v.integer(),
			v.minValue(1),
			v.check(
				(opLen) => opLen <= getOpLenMaxValue(),
				() => `Invalid op len: Max len of ${getOpLenMaxValue()} exceeded`
			)
		),
		() => opLen
	);

	let submitDisabled = $derived.by(() => {
		if (!opFromValidator.errors && !opLenValidator.errors) return false;

		return true;
	});
</script>

<Modal
	title="Add CAN ID Builder Operation"
	desc="Set the kind, from, and length of the operation you want to add"
>
	{#snippet trigger({ getProps })}
		{@render triggerSnippet({ getProps })}
	{/snippet}

	{#snippet content()}
		<div class="flex flex-col gap-6">
			<SegmentedControl
				name="can-id-builder-op-kind"
				bind:selectedValue={opKind}
				options={canIdBuilderOpKindOptions}
			/>

			<NumberInput
				name="can-id-builder-op-from"
				bind:value={opFrom}
				errors={opFromValidator.errors}
				min={0}
				label="From"
				desc="The bit from which the operation starts"
			/>

			<NumberInput
				name="can-id-builder-op-len"
				bind:value={opLen}
				errors={opLenValidator.errors}
				min={1}
				label="Len"
				desc="The length of the operation in bits"
			/>
		</div>
	{/snippet}

	{#snippet actions({ close })}
		<button
			onclick={() => {
				onsubmit({ opKind, opFrom, opLen });
				close();
			}}
			disabled={submitDisabled}
			class="btn btn-primary">Add</button
		>
	{/snippet}
</Modal>
