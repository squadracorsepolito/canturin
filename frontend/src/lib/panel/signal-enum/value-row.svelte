<script lang="ts">
	import type { SignalEnumValue } from '$lib/api/canturin';
	import NumberEditablev2 from '$lib/components/editable/number-editablev2.svelte';
	import TextEditablev2 from '$lib/components/editable/text-editablev2.svelte';
	import { TableField } from '$lib/components/table';
	import { getSignalEnumState } from '$lib/state/signal-enum-state.svelte';
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
	<TextEditablev2
		bind:value={value.name}
		oncommit={handleName}
		name="signal-enum-value-name"
		fontWeight="medium"
		errors={nameErrors}
	/>
</TableField>

<TableField>
	<NumberEditablev2
		bind:value={value.index}
		oncommit={handleIndex}
		name="signal-enum-value-index"
		errors={indexErrors}
	/>
</TableField>

<TableField>
	<TextEditablev2
		bind:value={value.desc}
		oncommit={handleDesc}
		name="signal-enum-value-desc"
		placeholder="Add Description"
	/>
</TableField>
