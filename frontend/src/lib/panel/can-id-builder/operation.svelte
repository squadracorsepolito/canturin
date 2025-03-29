<script lang="ts">
	import { CANIDBuilderOpKind, type CANIDBuilderOp } from '$lib/api/canturin';
	import { IconButton } from '$lib/components/button';
	import Divider from '$lib/components/divider/divider.svelte';
	import { NumberEditable } from '$lib/components/editable';
	import { AddIcon, DeleteIcon } from '$lib/components/icon';
	import { AddCanIdBuilderOpModal } from '$lib/components/modal';
	import { Validator } from '$lib/utils/validator.svelte';
	import { canIdBuilderOpKindDescs, canIdBuilderOpKindLabels } from './utils';
	import * as v from 'valibot';

	type Props = {
		baseSize: number;
		op: CANIDBuilderOp;
		partialRes: number;
		isLast?: boolean;
		insert: (opIndex: number, opKind: CANIDBuilderOpKind, opFrom: number, opLen: number) => void;
		delete: (opIndex: number) => void;
		updateFrom: (opIndex: number, opFrom: number) => void;
		updateLen: (opIndex: number, opLen: number) => void;
	};

	let {
		baseSize,
		op,
		partialRes,
		isLast,
		insert,
		delete: delFunc,
		updateFrom,
		updateLen
	}: Props = $props();

	function getColor(kind: CANIDBuilderOpKind) {
		switch (kind) {
			case CANIDBuilderOpKind.CANIDBuilderOpKindMessagePriority:
				return 'bg-primary text-primary-content';
			case CANIDBuilderOpKind.CANIDBuilderOpKindMessageID:
				return 'bg-secondary text-secondary-content';
			case CANIDBuilderOpKind.CANIDBuilderOpKindNodeID:
				return 'bg-info text-info-content';
			case CANIDBuilderOpKind.CANIDBuilderOpKindBitMask:
				return 'bg-error text-error-content';
			default:
				return '';
		}
	}

	let partialBits = $derived.by(() => {
		return partialRes.toString(2).padStart(baseSize, '0').split('').map(Number);
	});

	const fromValidator = new Validator(
		v.pipe(v.number(), v.integer(), v.minValue(0), v.maxValue(baseSize - 1)),
		() => op.from
	);

	function getLenMaxValue() {
		return baseSize - op.from;
	}

	const lenValidator = new Validator(
		v.pipe(v.number(), v.integer(), v.minValue(1), v.maxValue(getLenMaxValue())),
		() => op.len
	);
</script>

{#snippet bitSlice()}
	{#each { length: baseSize } as _, idx}
		{@const bit = baseSize - idx - 1}
		{@const isInRange = bit >= op.from && bit < op.from + op.len}

		<div
			class={[
				'join-item bg-base-200 aspect-square border-l border-y border-neutral justify-center justify-items-center m-auto w-full flex items-center',
				{ [getColor(op.kind)]: isInRange },
				{ 'rounded-l-btn': idx === 0 },
				{ 'border-r rounded-r-btn': idx === baseSize - 1 }
			]}
		>
			<span>
				{bit}
			</span>
		</div>
	{/each}
{/snippet}

{#snippet actions(incIndex?: boolean)}
	<div class="grid grid-cols-subgrid {op.index === 0 ? 'col-span-3' : 'col-span-2'}">
		<Divider />

		<div class="flex justify-center items-center">
			<AddCanIdBuilderOpModal
				onsubmit={({ opKind, opFrom, opLen }) => {
					insert(incIndex ? op.index + 1 : op.index, opKind, opFrom, opLen);
				}}
				maxSize={baseSize}
			>
				{#snippet trigger({ getProps })}
					<IconButton {...getProps()} themeColor="primary">
						<AddIcon />
					</IconButton>
				{/snippet}
			</AddCanIdBuilderOpModal>
		</div>

		{#if op.index === 0}
			<Divider />
		{/if}
	</div>
{/snippet}

{@render actions()}

<div class="flex gap-3 justify-between">
	<div>
		<h4>{canIdBuilderOpKindLabels[op.kind]}</h4>

		<p class="opacity-85 text-sm">{canIdBuilderOpKindDescs[op.kind]}</p>

		<div class="flex gap-3 pt-2">
			<NumberEditable
				bind:value={op.from}
				errors={fromValidator.errors}
				name="can-id-builder-op-from"
				oncommit={() => updateFrom(op.index, op.from)}
				label="From"
			/>

			<NumberEditable
				bind:value={op.len}
				errors={lenValidator.errors}
				name="can-id-builder-op-len"
				oncommit={() => updateLen(op.index, op.len)}
				label="Len"
			/>
		</div>
	</div>

	<IconButton onclick={() => delFunc(op.index)} themeColor="error">
		<DeleteIcon />
	</IconButton>
</div>

<div class="flex flex-col items-center">
	<div class="flex-1 block bg-neutral w-2 rounded-t-box"></div>

	<div class="w-10 h-10 rounded-box bg-neutral flex items-center justify-center">
		<span class="text-neutral-content">
			{op.index + 1}
		</span>
	</div>

	<div class="flex-1 block bg-neutral w-2 rounded-b-box"></div>
</div>

<div
	class="grid row-span-2 grid-rows-subgrid"
	style:grid-template-columns="repeat({baseSize}, minmax(0, 1fr))"
>
	{@render bitSlice()}

	{#each { length: baseSize } as _, idx}
		{@const bit = baseSize - idx - 1}
		{@const isInRange = bit >= op.from && bit < op.from + op.len}

		<div class="h-full flex justify-center join-item">
			{#if isInRange}
				<div
					class={['h-full border border-neutral rounded-box flex items-center', getColor(op.kind)]}
				>
					<div class="px-2">
						{partialBits[idx]}
					</div>
				</div>
			{:else}
				<div class="h-full bg-base-200 w-4 border border-base-300 rounded-box"></div>
			{/if}
		</div>
	{/each}
</div>

{#if isLast}
	{@render actions(true)}
{/if}
