<script lang="ts">
	import { z, ZodNumber } from 'zod';
	import EditableForm from '../form/editable-form.svelte';
	import { NumberInput } from '../input';

	type Props = {
		prefix: string;
		name: string;
		value: number;
		schema: ZodNumber;
		onchange: (value: number) => void;
		description?: string;
	};

	let { name, prefix, value, schema, description, onchange }: Props = $props();
</script>

<div class="">
	<EditableForm
		schema={z.object({
			[name]: schema
		})}
		initialValues={{ [name]: value }}
		onsubmit={(values) => onchange(values[name])}
	>
		{#snippet placeholder()}
			<div class="flex flex-col items-center border">
				<div class="opacity-90">{name}</div>

				<div class="text-3xl font-extrabold">{value}</div>

				{#if description}
					<div class="opacity-90 text-xs">{description}</div>
				{/if}
			</div>
		{/snippet}

		{#snippet input({ fsm, values, errors })}
			<NumberInput
				bind:value={values[name]}
				name={prefix + '-' + name.toLocaleLowerCase()}
				errors={errors[name]}
				focusOnDiplay
				onblur={() => fsm.send('BLUR')}
				onescape={() => fsm.send('ESCAPE')}
			/>
		{/snippet}
	</EditableForm>
</div>
