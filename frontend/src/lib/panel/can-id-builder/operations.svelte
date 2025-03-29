<script lang="ts">
	import { type CANIDBuilderOp, CANIDBuilderOpKind, MessagePriority } from '$lib/api/canturin';
	import { Attribute } from '$lib/components/attribute';
	import Divider from '$lib/components/divider/divider.svelte';
	import { NumberInput } from '$lib/components/input';
	import { Readonly } from '$lib/components/readonly';
	import { Select } from '$lib/components/select';
	import { getBinaryNumber, getHexNumber } from '$lib/utils';
	import { Validator } from '$lib/utils/validator.svelte';
	import { prioritySelectItems, priorityToNumber } from '../message/utils';
	import Operation from './operation.svelte';
	import type { CanIdBuilderState } from './state.svete';
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

	let partialresults = $state<number[]>([]);

	async function handleCalculate(msgPriority: MessagePriority, msgId: number, nodeId: number) {
		if (msgIdValidator.errors || nodeIdValidator.errors) return;

		const canIds = await s.calculateCANIds(msgPriority, msgId, nodeId);

		const pRes: number[] = [];
		for (const canId of canIds) {
			pRes.push(canId);
		}
		partialresults = pRes;

		if (canIds.length > 1) {
			resultCanId = canIds[canIds.length - 1];
		}
	}

	$effect(() => {
		handleCalculate(msgPriority, msgId, nodeId);
	});

	let resultCanId = $state(0);

	function handleInsertOperation(
		opIndex: number,
		opKind: CANIDBuilderOpKind,
		opFrom: number,
		opLen: number
	) {
		s.insertOperation(opIndex, opKind, opFrom, opLen);
	}

	function handleDeleteOperation(opIndex: number) {
		s.deleteOperation(opIndex);
	}

	function handleUpdateFrom(opIndex: number, opFrom: number) {
		s.updateOperationFrom(opIndex, opFrom);
	}

	function handleUpdateLen(opIndex: number, opLen: number) {
		s.updateOperationLen(opIndex, opLen);
	}
</script>

{#snippet operations(operations: CANIDBuilderOp[])}
	<div class="mt-8 grid gap-x-8 gap-y-2" style:grid-template-columns="1fr auto 1fr">
		{#each operations as op, idx}
			<Operation
				{op}
				baseSize={canIdSize}
				partialRes={partialresults.length > idx ? partialresults[idx] : 0}
				isLast={idx === operations.length - 1}
				insert={handleInsertOperation}
				delete={handleDeleteOperation}
				updateFrom={handleUpdateFrom}
				updateLen={handleUpdateLen}
			/>
		{/each}

		<div class="flex items-end text-base-content bg-base-200 rounded-box">
			<div class="py-3 px-5 text-4xl">Result</div>
		</div>

		<div></div>

		<div>
			<Divider />

			<div class="flex items-center justify-between text-4xl font-bold gap-3">
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
