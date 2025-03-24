<script lang="ts">
	import { type CANIDBuilderOp, CANIDBuilderOpKind, MessagePriority } from '$lib/api/canturin';
	import { Attribute, AttributeGroup } from '$lib/components/attribute';
	import { IconButton } from '$lib/components/button';
	import Divider from '$lib/components/divider/divider.svelte';
	import { NumberEditable } from '$lib/components/editable';
	import { AddIcon } from '$lib/components/icon';
	import { NumberInput } from '$lib/components/input';
	import NumberInput0 from '$lib/components/input/number-input0.svelte';
	import { Readonly } from '$lib/components/readonly';
	import { Select } from '$lib/components/select';
	import { getBinaryNumber, getHexNumber } from '$lib/utils';
	import { Validator } from '$lib/utils/validator.svelte';
	import { prioritySelectItems, priorityToNumber } from '../message/utils';
	import type { CanIdBuilderState } from './state.svete';
	import { canIdBuilderOpKindDescs, canIdBuilderOpKindLabels } from './utils';
	import * as v from 'valibot';

	const canIdSize = 11;

	type Props = {
		s: CanIdBuilderState;
	};

	let { s }: Props = $props();

	let msgPriority = $state(MessagePriority.MessagePriorityMedium);

	let msgId = $state(1);

	const msgIdValidator = new Validator(v.pipe(v.number(), v.integer(), v.minValue(0)), () => msgId);

	let nodeId = $state(1);

	const nodeIdValidator = new Validator(
		v.pipe(v.number(), v.integer(), v.minValue(0)),
		() => nodeId
	);

	let partialBits = $state<number[][]>([]);

	async function handleCalculate(msgPriority: MessagePriority, msgId: number, nodeId: number) {
		if (msgIdValidator.errors || nodeIdValidator.errors) return;

		const canIds = await s.calculateCANIds(msgPriority, msgId, nodeId);

		const partials: number[][] = [];
		for (const canId of canIds) {
			partials.push(getBinaryArray(canId));
		}

		partialBits = partials;

		if (canIds.length > 1) {
			resultCanId = canIds[canIds.length - 1];
		}
	}

	$effect(() => {
		handleCalculate(msgPriority, msgId, nodeId);
	});

	let resultCanId = $state(0);

	function getBinaryArray(num: number) {
		return num.toString(2).padStart(canIdSize, '0').split('').map(Number);
	}

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
</script>

{#snippet bitSlice(baseLen: number, op: CANIDBuilderOp)}
	<div class="join join-horizontal">
		{#each { length: baseLen } as _, idx}
			{@const bit = baseLen - idx - 1}
			{@const isInRange = bit >= op.from && bit < op.from + op.len}

			<div
				class={[
					'join-item p-4 w-14 bg-base-200 border border-neutral text-center',
					{ [getColor(op.kind)]: isInRange }
				]}
			>
				{bit}
			</div>
		{/each}
	</div>
{/snippet}

{#snippet decoration(onTop?: boolean)}
	<div class="grid grid-cols-subgrid {onTop ? 'col-span-3' : 'col-span-2'}">
		<Divider />

		<div class="flex justify-center items-center">
			<IconButton themeColor="primary">
				<AddIcon />
			</IconButton>
		</div>

		{#if onTop}
			<Divider />
		{/if}
	</div>
{/snippet}

{#snippet operations(operations: CANIDBuilderOp[])}
	<div class="mt-8 grid gap-x-8 gap-y-2" style:grid-template-columns="1fr auto 1fr">
		{#each operations as op, idx}
			{@render decoration(idx === 0)}

			<div>
				<h4>{canIdBuilderOpKindLabels[op.kind]}</h4>

				<p class="opacity-85 text-sm">{canIdBuilderOpKindDescs[op.kind]}</p>

				<div class="flex gap-3 pt-2">
					<NumberEditable bind:value={op.from} name="can-id-builder-op-from" label="From" />

					<NumberEditable bind:value={op.len} name="can-id-builder-op-len" label="Len" />
				</div>
			</div>

			<div class="flex flex-col items-center">
				<div class="flex-1 block bg-neutral w-2 rounded-t-box"></div>

				<div class="w-10 h-10 rounded-box bg-neutral flex items-center justify-center">
					<span class="text-neutral-content">
						{idx + 1}
					</span>
				</div>

				<div class="flex-1 block bg-neutral w-2 rounded-b-box"></div>
			</div>

			<div class="grid row-span-2 grid-rows-subgrid">
				<div class="flex items-center">
					{@render bitSlice(canIdSize, op)}
				</div>

				<div class="join">
					{#each { length: canIdSize }}
						<div class="h-full w-14 flex justify-center join-item">
							<div class="h-full bg-base-200 w-4 border border-base-300 rounded-box"></div>
						</div>
					{/each}
				</div>

				{@render partial(canIdSize, op, idx)}
			</div>

			{#if idx === operations.length - 1}
				{@render decoration()}
			{/if}
		{/each}

		{@render result()}
	</div>
{/snippet}

{#snippet partial(baseLen: number, op: CANIDBuilderOp, opIdx: number)}
	<div class="join">
		{#each { length: baseLen } as _, idx}
			{@const bit = canIdSize - idx - 1}
			{@const isInRange = bit >= op.from && bit < op.from + op.len}

			<div class="h-full w-14 flex justify-center join-item">
				{#if isInRange}
					<div
						class={[
							'h-full border border-neutral rounded-box flex items-center',
							getColor(op.kind)
						]}
					>
						<div class="px-2">
							{#if partialBits.length !== 0}
								{partialBits[opIdx][idx] ? '1' : '0'}
							{/if}
						</div>
					</div>
				{:else}
					<div class="h-full bg-base-200 w-4 border border-base-300 rounded-box"></div>
				{/if}
			</div>
		{/each}
	</div>
{/snippet}

{#snippet result()}
	<div class="flex items-end text-base-content bg-neutral rounded-box">
		<div class="py-3 px-5 text-4xl text-neutral-content">Result</div>
	</div>

	<div></div>

	<div>
		<Divider />

		<div class="flex items-center justify-between text-4xl font-bold">
			<div class="bg-base-200 rounded-box px-5 py-3">
				<div>{resultCanId}</div>
			</div>

			<div class="bg-base-200 rounded-box px-5 py-3">
				<div>{getHexNumber(resultCanId)}</div>
			</div>

			<div class="bg-base-200 rounded-box px-5 py-3">
				<div>{getBinaryNumber(resultCanId)}</div>
			</div>
		</div>
	</div>
{/snippet}

{#snippet testData()}
	<Attribute
		label="Test Message Priority"
		desc="The message priority used for testing the CAN ID builder"
	>
		<div class="flex gap-5">
			<Select
				items={prioritySelectItems}
				bind:selected={msgPriority}
				name="can-id-builder-message-priority"
				valueKey="value"
				labelKey="label"
			/>

			<Readonly>{getHexNumber(priorityToNumber(msgPriority))}</Readonly>

			<Readonly>{getBinaryNumber(priorityToNumber(msgPriority))}</Readonly>
		</div>
	</Attribute>

	<Divider />

	<Attribute label="Test Message ID" desc="The message ID used for testing the CAN ID builder">
		<div class="flex gap-5">
			<NumberInput
				name="can-id-builder-message-id"
				bind:value={msgId}
				errors={msgIdValidator.errors}
				min={0}
			/>

			<Readonly>{getHexNumber(msgId)}</Readonly>

			<Readonly>{getBinaryNumber(msgId)}</Readonly>
		</div>
	</Attribute>

	<Divider />

	<Attribute label="Test Node ID" desc="The node ID used for testing the CAN ID builder">
		<div class="flex gap-5">
			<NumberInput
				name="can-id-builder-node-id"
				bind:value={nodeId}
				errors={nodeIdValidator.errors}
				min={0}
			/>

			<Readonly>{getHexNumber(nodeId)}</Readonly>

			<Readonly>{getBinaryNumber(nodeId)}</Readonly>
		</div>
	</Attribute>
{/snippet}

<section>
	<h3 class="pb-5">Operations</h3>

	{@render testData()}

	{@render operations(s.entity.operations || [])}
</section>
