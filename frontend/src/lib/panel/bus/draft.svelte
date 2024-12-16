<script lang="ts">
	import { BusService, BusType } from '$lib/api/canturin';
	import FormField from '$lib/components/form-field/form-field.svelte';
	import { BusIcon } from '$lib/components/icon';
	import { TextInput } from '$lib/components/input';
	import { SegmentedControl } from '$lib/components/segmented-control';
	import Textarea from '$lib/components/textarea/textarea.svelte';
	import { onMount } from 'svelte';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import { baudrateItems, busTypeOptions, getSelectItemFromBaudrate } from './utils';
	import Divider from '$lib/components/divider/divider.svelte';
	import { Select } from '$lib/components/select';
	import { SubmitButton } from '$lib/components/button';
	import layout from '$lib/state/layout-state.svelte';

	let invalidNames = $state<string[]>([]);

	onMount(async () => {
		const res = await BusService.GetInvalidNames('');
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
		busType: z.nativeEnum(BusType).default(BusType.BusTypeCAN2A)
	});

	let selectedBaudrate = $state(getSelectItemFromBaudrate(1_000_000));

	const { enhance, errors, form } = superForm(defaults(zod(schema)), {
		SPA: true,
		validators: zod(schema),
		onUpdate: async ({ form }) => {
			if (form.valid) {
				const { name, desc, busType } = form.data;

				const tmpBus = await BusService.Create({
					name,
					desc,
					busType,
					baudrate: selectedBaudrate.valueAsNumber
				});

				layout.openPanel('bus', tmpBus.entityId);
			}
		}
	});
</script>

<form use:enhance method="POST">
	<div class="flex gap-3 items-center pb-8">
		<BusIcon height="48" width="48" />

		<h2>Create new Bus</h2>
	</div>

	<FormField label="Name" cols={4}>
		<TextInput name="bus-name" bind:value={$form.name} errors={$errors.name} />
	</FormField>

	<div class="py-5">
		<FormField label="Description (optional)" cols={4}>
			<Textarea name="bus-desc" bind:value={$form.desc} />
		</FormField>
	</div>

	<Divider />

	<FormField label="Bus Type" desc="The type of the bus">
		<SegmentedControl name="bus-type" options={busTypeOptions} bind:selectedValue={$form.busType} />
	</FormField>

	<Divider />

	<FormField label="Baudrate" desc="The baudrate of the bus">
		<Select
			name="bus-baudrate"
			items={baudrateItems}
			labelKey="label"
			valueKey="value"
			bind:selected={selectedBaudrate}
		/>
	</FormField>

	<div class="flex justify-end pt-5">
		<SubmitButton label="Create Bus" />
	</div>
</form>
