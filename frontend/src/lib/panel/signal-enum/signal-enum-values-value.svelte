<script lang="ts">
	import type { SignalEnumValue } from '$lib/api/canturin';
	import { TextareaEditable, TextEditable } from '$lib/components/editable';
	import NumberEditable from '$lib/components/editable/number-editable.svelte';
	import { getSignalEnumState } from '$lib/state/signal-enum-state.svelte';
	import { z } from 'zod';

	type Props = {
		entityId: string;
		value: SignalEnumValue;
		invalidNames: string[];
		invalidIndexes: number[];
	};

	let { entityId, value, invalidNames, invalidIndexes }: Props = $props();

	const ses = getSignalEnumState(entityId);

	const nameSchema = z.object({
		name: z
			.string()
			.min(1)
			.refine((n) => !invalidNames.includes(n), { message: 'Duplicated' })
	});

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

	const indexSchema = z.object({
		index: z
			.number()
			.min(0)
			.refine((n) => !invalidIndexes.includes(n), { message: 'Duplicated' })
	});

	function validateIndex(index: number) {
		const res = indexSchema.safeParse({ index });
		if (res.success) {
			return undefined;
		}
		return res.error.flatten().fieldErrors.index;
	}

	function handleIndex(index: number) {
		ses.updateValueIndex(value.entityId, index);
	}

	function handleDesc(desc: string) {
		ses.updateValueDesc(value.entityId, desc);
	}
</script>

<div class="py-3 col-span-8 grid grid-cols-subgrid">
	<div class="self-center">
		<NumberEditable
			validator={validateIndex}
			initialValue={value.index}
			name="signal-enum-value-index"
			onsubmit={handleIndex}
		/>
	</div>

	<div class="col-span-3 self-center">
		<TextEditable
			validator={validateName}
			name="signal-enum-value-name"
			initialValue={value.name}
			onsubmit={handleName}
			size="md"
		/>
	</div>

	<div class="col-span-4">
		<TextareaEditable
			initialValue={value.desc}
			name="signal-enum-value-desc"
			onsubmit={handleDesc}
			triggerLabel="Set Description"
			rows={1}
		/>
	</div>
</div>
