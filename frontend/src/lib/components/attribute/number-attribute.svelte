<script lang="ts">
	import { z } from 'zod';
	import EditableForm from '../form/editable-form.svelte';
	import { NumberInput } from '../input';
	import AttributePlaceholder from './attribute-placeholder.svelte';

	type Props = {
		prefix: string;
		name: string;
		value: number;
		schema: z.ZodNumber;
		onchange: (value: number) => void;
		desc?: string;
	};

	let { name, prefix, value, schema, desc, onchange }: Props = $props();

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
		submitOnClickOutside
	>
		{#snippet placeholder(fsm)}
			<AttributePlaceholder {value} title={name} {desc} isEditing={fsm.current === 'editing'} />
		{/snippet}

		{#snippet input({ fsm, values, errors })}
			<div class="px-4 py-2 border-2 border-primary rounded-b-box">
				<NumberInput
					bind:value={values[name]}
					name={prefix + '-' + name.toLocaleLowerCase()}
					errors={errors[name]}
					focusOnDisplay
					onescape={() => fsm.send('ESCAPE')}
					label={name}
					{min}
					{max}
				/>
			</div>
		{/snippet}
	</EditableForm>
</div>
