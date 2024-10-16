<script lang="ts">
	import { SignalEnumService } from '$lib/api/canturin';
	import { SignalEnumIcon } from '$lib/components/icon';
	import layout from '$lib/state/layout-state.svelte';
	import { onMount } from 'svelte';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import { text } from './signal-enum-text';
	import FormField from '$lib/components/form-field/form-field.svelte';
	import { NumberInput, TextInput } from '$lib/components/input';
	import { Textarea } from '$lib/components/textarea';
	import Divider from '$lib/components/divider/divider.svelte';
	import { SubmitButton } from '$lib/components/button';

	let invalidNames = $state<string[]>([]);

	onMount(async () => {
		const res = await SignalEnumService.GetNames();
		if (res) {
			invalidNames = res;
		}
	});

	const schema = z.object({
		name: z
			.string()
			.min(1)
			.refine((n) => !invalidNames.includes(n), { message: 'Duplicated' })
			.default(''),
		desc: z.string().optional().default(''),
		minSize: z.number().min(1).max(64).default(1)
	});

	const { enhance, errors, form } = superForm(defaults(zod(schema)), {
		SPA: true,
		validators: zod(schema),
		onUpdate: async ({ form }) => {
			if (form.valid) {
				const { name, desc, minSize } = form.data;

				const tmpSignalEnum = await SignalEnumService.Create(name, desc, minSize);

				layout.openPanel('signal_enum', tmpSignalEnum.entityId);
			}
		}
	});
</script>

<form use:enhance method="POST">
	<div class="flex gap-3 items-center pb-8">
		<SignalEnumIcon height="48" width="48" />

		<h2>{text.headings.draft}</h2>
	</div>

	<FormField {...text.name} cols={4}>
		<TextInput name="signal-enum-name" bind:value={$form.name} errors={$errors.name} />
	</FormField>

	<div class="py-5">
		<FormField {...text.desc} cols={4}>
			<Textarea name="signal-enum-desc" bind:value={$form.desc} />
		</FormField>
	</div>

	<Divider />

	<FormField {...text.minSize}>
		<NumberInput name="signal-enum-min-size" bind:value={$form.minSize} errors={$errors.minSize} />
	</FormField>

	<div class="flex justify-end pt-5">
		<SubmitButton label={text.buttons.draft.submit} />
	</div>
</form>
