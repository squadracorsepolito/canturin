<script lang="ts">
	import type { SignalEnumValue } from '$lib/api/canturin';
	import { TextareaEditable, TextEditable } from '$lib/components/editable';
	import NumberEditable from '$lib/components/editable/number-editable.svelte';
	import TextEditablev2 from '$lib/components/editable/text-editablev2.svelte';
	import { AltArrowIcon, CloseIcon } from '$lib/components/icon';
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

	let expanded = $state(false);
</script>

<div class="border-t-2 border-base-content/10 py-5 col-span-8 grid grid-cols-subgrid items-center">
	<div>
		<NumberEditable
			validator={validateIndex}
			initialValue={value.index}
			name="signal-enum-value-index"
			onsubmit={handleIndex}
		/>
	</div>

	<div class="col-span-3">
		<TextEditable
			validator={validateName}
			name="signal-enum-value-name"
			initialValue={value.name}
			onsubmit={handleName}
			size="md"
		/>

		<TextEditablev2 bind:value={value.name} errors={nameErrors} oncommit={handleName} />
	</div>

	{#if !expanded}
		<button
			onclick={() => (expanded = true)}
			class="col-span-4 p-3 flex items-center justify-between gap-2 hover:bg-base-content/10 h-full transition-colors rounded-btn"
		>
			<div class="text-dimmed truncate">
				{#if value.desc}
					{value.desc}
				{:else}
					Click Here To Add Description
				{/if}
			</div>

			<div>
				<AltArrowIcon />
			</div>
		</button>
	{:else}
		<div class="col-start-8 flex justify-end">
			<button
				onclick={() => (expanded = false)}
				class="hover:bg-base-content/10 h-full transition-colors rounded-btn p-3"
			>
				<CloseIcon />
			</button>
		</div>

		<div class="col-start-2 col-span-7 pt-3">
			<TextareaEditable
				initialValue={value.desc}
				name="signal-enum-value-desc"
				onsubmit={handleDesc}
				triggerLabel="Add Description"
				rows={3}
			/>
		</div>
	{/if}
</div>
