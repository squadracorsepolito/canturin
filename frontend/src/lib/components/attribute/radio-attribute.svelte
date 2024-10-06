<script lang="ts">
	import { z } from 'zod';
	import { EditableForm } from '../form';
	import type { RadioInputOption } from '../input/types';
	import { RadioInput } from '../input';
	import { Button, SubmitButton } from '../button';
	import AttributePlaceholder from './attribute-placeholder.svelte';

	type Props = {
		prefix: string;
		name: string;
		options: RadioInputOption[];
		selected: number;
		schema: z.ZodNumber;
		desc?: string;
		onchange: (selected: number) => void;
	};

	let { prefix, name, options, selected, schema, desc, onchange }: Props = $props();

	let selectedLabel = $derived(options[selected - 1].label);
</script>

<div>
	<EditableForm
		schema={z.object({
			[name]: schema
		})}
		initialValues={{ [name]: selected }}
		onsubmit={(values) => onchange(values[name])}
		blurOnSubmit
		submitOnClickOutside
	>
		{#snippet placeholder(fsm)}
			<AttributePlaceholder
				value={selectedLabel}
				title={name}
				{desc}
				isEditing={fsm.current === 'editing'}
			/>
		{/snippet}

		{#snippet input({ fsm, values })}
			<div class="px-4 py-2 border-2 border-primary rounded-b-box">
				<RadioInput
					bind:selected={values[name]}
					name={prefix + '-' + name.toLocaleLowerCase()}
					{options}
				/>

				<div class="pt-3 flex gap-2 justify-end">
					<Button label="Close" onclick={() => fsm.send('ESCAPE')} />
					<SubmitButton label="Save" />
				</div>
			</div>
		{/snippet}
	</EditableForm>
</div>
