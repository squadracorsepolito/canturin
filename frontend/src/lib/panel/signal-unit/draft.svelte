<script lang="ts">
	import { SignalUnitKind, SignalUnitService } from '$lib/api/canturin';
	import Divider from '$lib/components/divider/divider.svelte';
	import FormField from '$lib/components/form-field/form-field.svelte';
	import { SignalUnitIcon } from '$lib/components/icon';
	import { TextInput } from '$lib/components/input';
	import { SegmentedControl } from '$lib/components/segmented-control';
	import { Textarea } from '$lib/components/textarea';
	import { onMount } from 'svelte';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import { signalUnitKindOptions } from './utils';
	import { SubmitButton } from '$lib/components/button';
	import layout from '$lib/state/layout-state.svelte';

	let invalidNames = $state<string[]>([]);

	onMount(async () => {
		const res = await SignalUnitService.GetInvalidNames('');
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
		kind: z.nativeEnum(SignalUnitKind).default(SignalUnitKind.SignalUnitKindCustom),
		symbol: z.string().default('')
	});

	const { enhance, errors, form } = superForm(defaults(zod(schema)), {
		SPA: true,
		validators: zod(schema),
		onUpdated: async ({ form }) => {
			if (form.valid) {
				const tmpSigUnit = await SignalUnitService.Create(form.data);

				layout.openPanel('signal_unit', tmpSigUnit.entityId);
			}
		}
	});
</script>

<form use:enhance method="POST">
	<div class="flex gap-3 items-center pb-8">
		<SignalUnitIcon height="48" width="48" />

		<h2>Create new Signal Unit</h2>
	</div>

	<FormField label="Name">
		<TextInput name="signal-unit-name" bind:value={$form.name} errors={$errors.name} />
	</FormField>

	<div class="py-5">
		<FormField label="Description (optional)" cols={4}>
			<Textarea name="signal-unit-desc" bind:value={$form.desc} />
		</FormField>
	</div>

	<Divider />

	<FormField label="Kind" desc="The kind of the unit" cols={4}>
		<SegmentedControl
			name="signal-unit-kind"
			options={signalUnitKindOptions}
			bind:selectedValue={$form.kind}
		/>
	</FormField>

	<Divider />

	<FormField label="Symbol" desc="The symbol of the unit">
		<TextInput name="signal-unit-symbol" bind:value={$form.symbol} />
	</FormField>

	<div class="flex justify-end pt-5">
		<SubmitButton label="Create Signal Unit" />
	</div>
</form>
