<script lang="ts">
	import type { SignalEnumValue } from '$lib/api/canturin';
	import { TextEditable, NumberEditable } from '$lib/components/editable';
	import { TableField } from '$lib/components/table';
	import { getSignalEnumState } from '$lib/panel/signal-enum/state.svelte';
	import { z } from 'zod';

	type Props = {
		enumEntityId: string;
		value: SignalEnumValue;
	};

	let { enumEntityId, value }: Props = $props();

	const ses = getSignalEnumState(enumEntityId);

	let invalidNames = $derived.by(() => {
		if (!ses.entity.values) return [];

		return ses.entity.values
			.filter((val) => val.entityId !== value.entityId)
			.map((val) => val.name);
	});

	const nameSchema = z.object({
		name: z
			.string()
			.min(1)
			.refine((n) => !invalidNames.includes(n), { message: 'Duplicated' })
	});

	let nameErrors = $derived.by(() => {
		const res = nameSchema.safeParse({ name: value.name });
		if (res.success) {
			return undefined;
		}
		return res.error.flatten().fieldErrors.name;
	});

	function handleName(name: string) {
		ses.updateValueName(value.entityId, name);
	}

	let invalidIndexes = $derived.by(() => {
		if (!ses.entity.values) return [];

		return ses.entity.values
			.filter((val) => val.entityId !== value.entityId)
			.map((val) => val.index);
	});

	const indexSchema = z.object({
		index: z
			.number()
			.min(0)
			.max(4_294_967_295)
			.refine((n) => !invalidIndexes.includes(n), { message: 'Duplicated' })
	});

	let indexErrors = $derived.by(() => {
		const res = indexSchema.safeParse({ index: value.index });
		if (res.success) {
			return undefined;
		}
		return res.error.flatten().fieldErrors.index;
	});

	function handleIndex(index: number) {
		ses.updateValueIndex(value.entityId, index);
	}

	function handleDesc(desc: string) {
		ses.updateValueDesc(value.entityId, desc);
	}
</script>

<TableField>
	<TextEditable
		bind:value={value.name}
		oncommit={handleName}
		name="signal-enum-value-name"
		fontWeight="medium"
		errors={nameErrors}
		border="transparent"
	/>
</TableField>

<TableField>
	<NumberEditable
		bind:value={value.index}
		oncommit={handleIndex}
		name="signal-enum-value-index"
		errors={indexErrors}
		border="transparent"
	/>
</TableField>

<TableField>
	<TextEditable
		bind:value={value.desc}
		oncommit={handleDesc}
		name="signal-enum-value-desc"
		placeholder="Add Description"
		border="transparent"
	/>
</TableField>
