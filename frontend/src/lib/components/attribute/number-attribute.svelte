<script lang="ts">
	import { z } from 'zod';
	import EditableForm from '../form/editable-form.svelte';
	import { NumberInput } from '../input';

	type Props = {
		prefix: string;
		name: string;
		value: number;
		schema: z.ZodNumber;
		onchange: (value: number) => void;
		description?: string;
	};

	let { name, prefix, value, schema, description, onchange }: Props = $props();

	let min = $derived.by<number | undefined>(() => {
		const match = schema._def.checks.find(({ kind }) => kind === 'min');
		if (match !== undefined) {
			if (match.kind === 'min') {
				return match.value;
			}
		}
	});

	let max = $derived.by<number | undefined>(() => {
		const match = schema._def.checks.find(({ kind }) => kind === 'max');
		if (match !== undefined) {
			if (match.kind === 'max') {
				return match.value;
			}
		}
	});
</script>

<div>
	<EditableForm
		schema={z.object({
			[name]: schema
		})}
		initialValues={{ [name]: value }}
		onsubmit={(values) => onchange(values[name])}
	>
		{#snippet placeholder(fsm)}
			<div
				class="flex flex-col items-center p-3 rounded-t-box {fsm.current === 'editing'
					? 'bg-primary text-primary-content'
					: 'hover:bg-base-200 rounded-b-box'}"
			>
				<div class="opacity-90">{name}</div>

				<div class="text-3xl font-extrabold">{value}</div>

				{#if description}
					<div class="opacity-90 text-xs">{description}</div>
				{/if}
			</div>
		{/snippet}

		{#snippet input({ fsm, values, errors })}
			<div class="px-4 py-2 border-2 border-primary rounded-b-box">
				<NumberInput
					bind:value={values[name]}
					name={prefix + '-' + name.toLocaleLowerCase()}
					errors={errors[name]}
					focusOnDisplay
					onblur={() => fsm.send('BLUR')}
					onescape={() => fsm.send('ESCAPE')}
					label={name}
					{min}
					{max}
				/>
			</div>
		{/snippet}
	</EditableForm>
</div>
