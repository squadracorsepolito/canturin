<script lang="ts">
	import {
		MessageByteOrder,
		MessageService,
		type Message,
		type MessageSendType
	} from '$lib/api/canturin';
	import { Attribute, AttributeGroup } from '$lib/components/attribute';
	import { NumberEditable } from '$lib/components/editable';
	import { Select } from '$lib/components/select';
	import { checkUnused, Validator } from '$lib/utils/validator.svelte';
	import type { PanelSectionProps } from '../types';
	import { getMessageState } from './state.svelte';
	import * as v from 'valibot';
	import { byteOrderOptions, sendTypeSelectItems } from './utils';
	import Divider from '$lib/components/divider/divider.svelte';
	import Readonly from '$lib/components/readonly/readonly.svelte';
	import { getHexNumber } from '$lib/utils';
	import { Switch } from '$lib/components/switch';
	import { onMount } from 'svelte';
	import { SegmentedControl } from '$lib/components/segmented-control';

	let { entityId }: PanelSectionProps = $props();

	const ms = getMessageState(entityId);

	let invalidMsgIds = $state<number[]>([]);
	let invalidCanIds = $state<number[]>([]);

	onMount(async () => {
		const msgIds = await ms.getInvalidMessageIds();
		const canIds = await ms.getInvalidCanIds();

		invalidMsgIds = msgIds;
		invalidCanIds = canIds;
	});

	const messageIdValidator = new Validator(
		v.pipe(
			v.number(),
			v.integer(),
			v.minValue(0),
			checkUnused(() => invalidMsgIds)
		),
		() => ms.entity.id
	);

	const canIdValidator = new Validator(
		v.pipe(
			v.number(),
			v.integer(),
			v.minValue(0),
			checkUnused(() => invalidCanIds)
		),
		() => ms.entity.canId
	);

	const cycleTimeValidator = new Validator(
		v.pipe(v.number(), v.integer(), v.minValue(0)),
		() => ms.entity.cycleTime
	);

	const delayTimeValidator = new Validator(
		v.pipe(v.number(), v.integer(), v.minValue(0)),
		() => ms.entity.delayTime
	);

	const startDelayTimeValidator = new Validator(
		v.pipe(v.number(), v.integer(), v.minValue(0)),
		() => ms.entity.startDelayTime
	);

	function handleMessageId(msgId: number) {
		ms.updateMessageId(msgId);
	}

	function handleCanId(canId: number) {
		ms.updateStaticCanId(canId);
	}

	function handleByteOrder(byteOrder: string) {
		ms.updateByteOrder(byteOrder as MessageByteOrder);
	}

	function handleCycleTime(cycleTime: number) {
		ms.updateCycleTime(cycleTime);
	}

	function handleSendType(sendType: MessageSendType) {
		ms.updateSendType(sendType);
	}

	function handleDelayTime(delayTime: number) {
		ms.updateDelayTime(delayTime);
	}

	function handleStartDelayTime(startDelayTime: number) {
		ms.updateStartDelayTime(startDelayTime);
	}
</script>

{#snippet section(msg: Message)}
	<AttributeGroup>
		<Attribute label="Static CAN ID" desc="Whether the CAN ID is static" fullSpan>
			<Switch bind:checked={msg.hasStaticCANID} />
		</Attribute>

		<Attribute label="ID" desc="The ID of the message">
			<NumberEditable
				bind:value={msg.id}
				name="message-id"
				errors={messageIdValidator.errors}
				oncommit={handleMessageId}
				readonly={msg.hasStaticCANID}
			/>
		</Attribute>

		<Attribute label="Hex ID" desc="The ID of the message in hex">
			<Readonly>{getHexNumber(msg.id)}</Readonly>
		</Attribute>

		<Attribute
			label="CAN ID"
			desc="The CAN ID of the message. It is the actual value of the id field of the CAN bus header"
		>
			<NumberEditable
				bind:value={msg.canId}
				name="message-can-id"
				errors={canIdValidator.errors}
				oncommit={handleCanId}
				readonly={!msg.hasStaticCANID}
			/>
		</Attribute>

		<Attribute label="Hex CAN ID" desc="The CAN ID of the message in hex">
			<Readonly>{getHexNumber(msg.canId)}</Readonly>
		</Attribute>
	</AttributeGroup>

	<Divider />

	<AttributeGroup>
		<Attribute label="Size" desc="The size of the message in bytes" fullSpan>
			<Readonly>
				{msg.sizeByte}
			</Readonly>
		</Attribute>

		<Attribute label="Byte Order" desc="The byte order of the message payload" fullSpan>
			<SegmentedControl
				name="message-byte-order"
				options={byteOrderOptions}
				bind:selectedValue={msg.byteOrder}
				onchange={handleByteOrder}
			/>
		</Attribute>
	</AttributeGroup>

	<Divider />

	<AttributeGroup>
		<Attribute label="Cycle Time" desc="The cycle time of the message in ms">
			<NumberEditable
				bind:value={msg.cycleTime}
				name="message-cycle-time"
				errors={cycleTimeValidator.errors}
				oncommit={handleCycleTime}
			/>
		</Attribute>

		<Attribute label="Send Type" desc="How the message is sent">
			<Select
				items={sendTypeSelectItems}
				valueKey="value"
				labelKey="label"
				bind:selected={msg.sendType}
				name="message-send-type"
				onselect={handleSendType}
			/>
		</Attribute>

		<Attribute label="Delay Time" desc="The delay time of the message in ms">
			<NumberEditable
				bind:value={msg.delayTime}
				name="message-delay-time"
				errors={delayTimeValidator.errors}
				oncommit={handleDelayTime}
			/>
		</Attribute>

		<Attribute label="Start Delay Time" desc="The start delay time of the message in ms">
			<NumberEditable
				bind:value={msg.startDelayTime}
				name="message-start-delay-time"
				errors={startDelayTimeValidator.errors}
				oncommit={handleStartDelayTime}
			/>
		</Attribute>
	</AttributeGroup>
{/snippet}

<section>
	<h3 class="pb-5">Attributes</h3>

	{@render section(ms.entity)}
</section>
