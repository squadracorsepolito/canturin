<script lang="ts">
	import { onMount } from 'svelte';
	import * as v from 'valibot';
	import { nameSchema } from '$lib/utils/validator.svelte';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import { MessageIcon } from '$lib/components/icon';
	import FormField from '$lib/components/form-field/form-field.svelte';
	import { NumberInput, TextInput } from '$lib/components/input';
	import { Textarea } from '$lib/components/textarea';
	import Divider from '$lib/components/divider/divider.svelte';
	import { SubmitButton } from '$lib/components/button';
	import { Switch } from '$lib/components/switch';
	import { Readonly } from '$lib/components/readonly';
	import { getHexNumber } from '$lib/utils';
	import { MessageByteOrder, MessageSendType, MessageService } from '$lib/api/canturin';
	import { SegmentedControl } from '$lib/components/segmented-control';
	import { byteOrderOptions, sendTypeSelectItems } from './utils';
	import { Select } from '$lib/components/select';

	let invalidNames = $state<string[]>([]);

	onMount(async () => {
		const res = await MessageService.GetInvalidNames('');
		if (res) {
			invalidNames = res;
		}
	});

	let hasStaticCANID = $state<boolean>(false);

	const schema = v.object({
		name: nameSchema(() => invalidNames),
		desc: v.string(),
		id: v.pipe(v.number(), v.integer(), v.minValue(0)),
		canId: v.pipe(v.number(), v.integer(), v.minValue(0)),
		size: v.optional(v.pipe(v.number(), v.integer(), v.minValue(1), v.maxValue(8)), 1),
		byteOrder: v.optional(v.enum(MessageByteOrder), MessageByteOrder.MessageByteOrderLittleEndian),
		cycleTime: v.pipe(v.number(), v.integer(), v.minValue(0)),
		sendType: v.optional(v.enum(MessageSendType), MessageSendType.MessageSendTypeCyclic),
		delayTime: v.pipe(v.number(), v.integer(), v.minValue(0)),
		startDelayTime: v.pipe(v.number(), v.integer(), v.minValue(0))
	});

	const { enhance, form, errors } = superForm(defaults(valibot(schema)), {
		SPA: true,
		validators: valibot(schema),
		onUpdate: async ({ form }) => {
			if (form.valid) {
				console.log(form.data);
			}
		}
	});
</script>

<form method="POST" use:enhance>
	<div class="flex gap-3 items-center pb-8">
		<MessageIcon height="48" width="48" />

		<h2>Create new Message</h2>
	</div>

	<FormField label="Name" cols={4}>
		<TextInput name="message-name" bind:value={$form.name} errors={$errors.name} />
	</FormField>

	<div class="py-5">
		<FormField label="Description (optional)" cols={4}>
			<Textarea name="message-desc" bind:value={$form.desc} />
		</FormField>
	</div>

	<Divider />

	<div class="pb-5">
		<FormField label="Static CAN ID" desc="Whether the CAN ID is static">
			<Switch bind:checked={hasStaticCANID} />
		</FormField>
	</div>

	<div class="grid grid-cols-2 gap-5">
		<FormField label="ID" desc="The ID of the message">
			<NumberInput name="message-id" bind:value={$form.id} errors={$errors.id} min={0} />
		</FormField>

		<FormField label="Hex ID" desc="The ID of the message in hex">
			<Readonly>{getHexNumber($form.id)}</Readonly>
		</FormField>

		<FormField
			label="CAN ID"
			desc="The CAN ID of the message. It is the actual value of the id field of the CAN bus header"
		>
			<NumberInput name="message-id" bind:value={$form.canId} errors={$errors.canId} min={0} />
		</FormField>

		<FormField label="Hex CAN ID" desc="The CAN ID of the message in hex">
			<Readonly>{getHexNumber($form.canId)}</Readonly>
		</FormField>
	</div>

	<Divider />

	<FormField label="Size" desc="The size of the message in bytes">
		<NumberInput
			name="message-size"
			bind:value={$form.size}
			errors={$errors.size}
			min={1}
			max={8}
		/>
	</FormField>

	<FormField label="Byte Order" desc="The byte order of the message">
		<SegmentedControl
			name="message-byte-order"
			options={byteOrderOptions}
			bind:selectedValue={$form.byteOrder}
		/>
	</FormField>

	<Divider />

	<div class="grid grid-cols-2 gap-5">
		<FormField label="Cycle Time" desc="The cycle time of the message in ms">
			<NumberInput
				name="message-cycle-time"
				bind:value={$form.cycleTime}
				errors={$errors.cycleTime}
			/>
		</FormField>

		<FormField label="Send Type" desc="How the message is sent">
			<Select
				items={sendTypeSelectItems}
				valueKey="value"
				labelKey="label"
				bind:selected={$form.sendType}
				name="message-send-type"
			/>
		</FormField>

		<FormField label="Delay Time" desc="The delay time of the message in ms">
			<NumberInput
				name="message-delay-time"
				bind:value={$form.delayTime}
				errors={$errors.delayTime}
			/>
		</FormField>

		<FormField label="Start Delay Time" desc="The start delay time of the message in ms">
			<NumberInput
				name="message-start-delay-time"
				bind:value={$form.startDelayTime}
				errors={$errors.startDelayTime}
			/>
		</FormField>
	</div>

	<div class="flex justify-end pt-5">
		<SubmitButton label="Create Message" />
	</div>
</form>
