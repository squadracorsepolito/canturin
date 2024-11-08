<script lang="ts">
	import type { SignalEnumValue } from '$lib/api/canturin';
	import NumberEditablev2 from '$lib/components/editable/number-editablev2.svelte';
	import TextEditablev2 from '$lib/components/editable/text-editablev2.svelte';
	import { ExpandableTableField, TableField } from '$lib/components/table';
	import { getSignalEnumState } from '$lib/state/signal-enum-state.svelte';
	import { z } from 'zod';

	type Props = {
		enumEntityId: string;
		value: SignalEnumValue;
	};

	let { enumEntityId, value }: Props = $props();

	const ses = getSignalEnumState(enumEntityId);

	let tmpName = $state(value.name);

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
		const res = nameSchema.safeParse({ name: tmpName });
		if (res.success) {
			return undefined;
		}
		return res.error.flatten().fieldErrors.name;
	});

	function handleName(name: string) {
		ses.updateValueName(value.entityId, name);
	}

	let tmpIndex = $state(value.index);

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
		const res = indexSchema.safeParse({ index: tmpIndex });
		if (res.success) {
			return undefined;
		}
		return res.error.flatten().fieldErrors.index;
	});

	function handleIndex(index: number) {
		ses.updateValueIndex(value.entityId, index);
	}
</script>

<TableField>
	<TextEditablev2
		bind:value={tmpName}
		oncommit={handleName}
		name="signal-enum-value-name"
		errors={nameErrors}
	/>
</TableField>

<TableField>
	<NumberEditablev2
		bind:value={tmpIndex}
		oncommit={handleIndex}
		name="signal-enum-value-index"
		errors={indexErrors}
	/>
</TableField>

{#if value.desc}
	<ExpandableTableField>
		{value.desc}
	</ExpandableTableField>
{:else}
	<TableField>
		<i class="text-dimmed">No Description</i>
	</TableField>
{/if}
