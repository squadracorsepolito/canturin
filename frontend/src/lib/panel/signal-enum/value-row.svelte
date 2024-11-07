<script lang="ts">
	import type { SignalEnumValue } from '$lib/api/canturin';
	import TextEditablev2 from '$lib/components/editable/text-editablev2.svelte';
	import { ExpandableTableField, TableField } from '$lib/components/table';
	import { getSignalEnumState } from '$lib/state/signal-enum-state.svelte';
	import { z } from 'zod';

	type Props = {
		enumEntityId: string;
		value: SignalEnumValue;
		invalidNames: string[];
	};

	let { enumEntityId, value, invalidNames }: Props = $props();

	const ses = getSignalEnumState(enumEntityId);

	const nameSchema = z.object({
		name: z
			.string()
			.min(1)
			.refine((n) => !invalidNames.includes(n), { message: 'Duplicated' })
	});

	let nameErrors = $derived.by(() => validateName(value.name));

	function validateName(name: string) {
		const res = nameSchema.safeParse({ name });
		if (res.success) {
			return undefined;
		}
		return res.error.flatten().fieldErrors.name;
	}

	function handleName(name: string) {
		ses.updateValueName(value.entityId, name);
	}
</script>

<TableField>
	<TextEditablev2
		bind:value={value.name}
		oncommit={handleName}
		name="signal-enum-value-name"
		errors={nameErrors}
	/>
</TableField>

<TableField>
	{value.index}
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
