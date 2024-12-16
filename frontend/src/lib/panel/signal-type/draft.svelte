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
	import Divider from '$lib/components/divider/divider.svelte';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { zod } from 'sveltekit-superforms/adapters';
	import { SignalTypeIcon } from '$lib/components/icon';
	import ToggleInput from '$lib/components/input/toggle-input.svelte';
	import layout from '$lib/state/layout-state.svelte';
	import { signalTypeKindOptions } from './utils';

	let invalidNames = $state<string[]>([]);

	onMount(async () => {
		const res = await SignalTypeService.GetInvalidNames('');
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

				const tmpSigType = await SignalTypeService.Create({
					kind,
					name,
					desc,
					size,
					signed,
					min,
					max,
					scale,
					offset
				});

				layout.openPanel('signal_type', tmpSigType.entityId);
			}
		}
	});
</script>

<form use:enhance method="POST">
	<div class="flex gap-3 items-center pb-8">
		<SignalTypeIcon height="48" width="48" />

		<h2>Create new Signal Type</h2>
	</div>

	<FormField label="Name" cols={4}>
		<TextInput name="signal-type-name" bind:value={$form.name} errors={$errors.name} />
	</FormField>

	<div class="py-5">
		<FormField label="Description (optional)" cols={4}>
			<Textarea name="signal-type-desc" bind:value={$form.desc} />
		</FormField>
	</div>

	<Divider />

	<FormField label="Kind" desc="The kind of the type">
		<SegmentedControl
			name="signal-type-kind"
			options={signalTypeKindOptions}
			bind:selectedValue={$form.kind}
		/>
	</FormField>

	<Divider />

	<div class="grid grid-cols-2 gap-5">
		<FormField label="Size" desc="The size in bits">
			<NumberInput name="signal-type-size" bind:value={$form.size} errors={$errors.size} />
		</FormField>

		<FormField label="Signed" desc="Whether the value is signed">
			<ToggleInput name="signal-type-signed" bind:checked={$form.signed} />
		</FormField>
	</div>

	<Divider />

	<div class="grid grid-cols-2 gap-5">
		<FormField label="Min" desc="The minimum value">
			<NumberInput name="signal-type-min" bind:value={$form.min} errors={$errors.min} />
		</FormField>

		<FormField label="Max" desc="The maximum value">
			<NumberInput name="signal-type-max" bind:value={$form.max} errors={$errors.max} />
		</FormField>
	</div>

	<Divider />

	<div class="grid grid-cols-2 gap-5">
		<FormField label="Scale" desc="The scale factor">
			<NumberInput name="signal-type-scale" bind:value={$form.scale} errors={$errors.scale} />
		</FormField>

		<FormField label="Offset" desc="The offset">
			<NumberInput name="signal-type-offset" bind:value={$form.offset} errors={$errors.offset} />
		</FormField>
	</div>

	<div class="flex justify-end pt-5">
		<SubmitButton label="Create Signal Type" />
	</div>
</form>
