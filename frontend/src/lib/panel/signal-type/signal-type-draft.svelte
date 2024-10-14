<script lang="ts">
	import { SignalTypeKind, SignalTypeService } from '$lib/api/canturin';
	import TextInput from '$lib/components/input/text-input.svelte';
	import Textarea from '$lib/components/textarea/textarea.svelte';
	import { onMount } from 'svelte';
	import { z } from 'zod';
	import { SegmentedControl } from '$lib/components/segmented-control';
	import NumberInput from '$lib/components/input/number-input.svelte';
	import SubmitButton from '$lib/components/button/submit-button.svelte';
	import FormField from '$lib/components/form-field/form-field.svelte';
	import { text } from './signal-type-text';
	import Divider from '$lib/components/divider/divider.svelte';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { zod } from 'sveltekit-superforms/adapters';
	import { SignalTypeIcon } from '$lib/components/icon';
	import ToggleInput from '$lib/components/input/toggle-input.svelte';
	import layout from '$lib/state/layout-state.svelte';

	let invalidNames = $state<string[]>([]);

	onMount(async () => {
		const res = await SignalTypeService.GetNames();
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
		kind: z.nativeEnum(SignalTypeKind).default(SignalTypeKind.SignalTypeKindCustom),
		size: z.number().min(1).max(64).default(1),
		signed: z.boolean().default(false),
		min: z.number().default(0),
		max: z.number().default(1),
		scale: z.number().default(1),
		offset: z.number().default(0)
	});

	const { enhance, errors, form } = superForm(defaults(zod(schema)), {
		SPA: true,
		validators: zod(schema),
		onUpdate: async ({ form }) => {
			if (form.valid) {
				const { name, desc, kind, size, signed, min, max, scale, offset } = form.data;

				const tmpSigType = await SignalTypeService.Create(
					kind,
					name,
					desc,
					size,
					signed,
					min,
					max,
					scale,
					offset
				);

				layout.openPanel('signal_type', tmpSigType.entityId);
			}
		}
	});
</script>

<form use:enhance method="POST">
	<div class="flex gap-3 items-center pb-8">
		<SignalTypeIcon height="48" width="48" />

		<h2>{text.headings.draft}</h2>
	</div>

	<FormField {...text.name} cols={4}>
		<TextInput name="signal-type-name" bind:value={$form.name} errors={$errors.name} />
	</FormField>

	<div class="py-5">
		<FormField {...text.desc} cols={4}>
			<Textarea name="signal-type-desc" bind:value={$form.desc} />
		</FormField>
	</div>

	<FormField {...text.kind}>
		<SegmentedControl
			name="signal-type-kind"
			options={text.kind.options}
			bind:selectedValue={$form.kind}
		/>
	</FormField>

	<Divider />

	<div class="grid grid-cols-2 gap-5">
		<FormField {...text.size}>
			<NumberInput name="signal-type-size" bind:value={$form.size} errors={$errors.size} />
		</FormField>

		<FormField {...text.signed}>
			<ToggleInput name="signal-type-signed" bind:checked={$form.signed} />
		</FormField>
	</div>

	<Divider />

	<div class="grid grid-cols-2 gap-5">
		<FormField {...text.min}>
			<NumberInput name="signal-type-min" bind:value={$form.min} errors={$errors.min} />
		</FormField>

		<FormField {...text.max}>
			<NumberInput name="signal-type-max" bind:value={$form.max} errors={$errors.max} />
		</FormField>
	</div>

	<Divider />

	<div class="grid grid-cols-2 gap-5">
		<FormField {...text.scale}>
			<NumberInput name="signal-type-scale" bind:value={$form.scale} errors={$errors.scale} />
		</FormField>

		<FormField {...text.offset}>
			<NumberInput name="signal-type-offset" bind:value={$form.offset} errors={$errors.offset} />
		</FormField>
	</div>

	<div class="flex justify-end pt-5">
		<SubmitButton label={text.buttons.draft.submit} />
	</div>
</form>
