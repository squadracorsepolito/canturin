<script lang="ts">
	import { NodeService } from '$lib/api/canturin';
	import { onMount } from 'svelte';
	import { z } from 'zod';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { zod } from 'sveltekit-superforms/adapters';
	import { NodeIcon } from '$lib/components/icon';
	import FormField from '$lib/components/form-field/form-field.svelte';
	import { NumberInput, TextInput } from '$lib/components/input';
	import { Textarea } from '$lib/components/textarea';
	import Divider from '$lib/components/divider/divider.svelte';
	import { SubmitButton } from '$lib/components/button';
	import layout from '$lib/state/layout-state.svelte';

	let invalidNames = $state<string[]>([]);
	let invalidNodeIds = $state<number[]>([]);

	onMount(async () => {
		const names = await NodeService.GetInvalidNames('');
		if (names) {
			invalidNames = names;
		}

		const nodeIds = await NodeService.GetInvalidNodeIDs('');
		if (nodeIds) {
			invalidNodeIds = nodeIds;
		}
	});

	const schema = z.object({
		name: z
			.string()
			.min(1)
			.refine((n) => !invalidNames.includes(n), { message: 'Duplicated' })
			.default(''),
		desc: z.string().optional().default(''),
		nodeId: z
			.number()
			.min(0)
			.refine((n) => !invalidNodeIds.includes(n), { message: 'Duplicated' })
			.default(0),
		interfaceCount: z.number().min(1).default(1)
	});

	const { enhance, errors, form } = superForm(defaults(zod(schema)), {
		SPA: true,
		validators: zod(schema),
		onUpdate: async ({ form }) => {
			if (form.valid) {
				const tmpNode = await NodeService.Create(form.data);

				layout.openPanel('node', tmpNode.entityId);

				console.log(tmpNode);
			}
		}
	});
</script>

<form use:enhance method="POST">
	<div class="flex gap-3 items-center pb-8">
		<NodeIcon height="48" width="48" />

		<h2>Create new Node</h2>
	</div>

	<FormField label="Name" cols={4}>
		<TextInput name="node-name" bind:value={$form.name} errors={$errors.name} />
	</FormField>

	<div class="py-5">
		<FormField label="Description (optional)" cols={4}>
			<Textarea name="node-desc" bind:value={$form.desc} />
		</FormField>
	</div>

	<Divider />

	<FormField label="Node ID" desc="The unique ID of the node expressed as a decimal number">
		<NumberInput name="node-id" bind:value={$form.nodeId} errors={$errors.nodeId} min={0} />
	</FormField>

	<Divider />

	<FormField label="Interface Count" desc="The number of interfaces that are present on the node">
		<NumberInput
			name="node-interface-count"
			bind:value={$form.interfaceCount}
			errors={$errors.interfaceCount}
			min={1}
		/>
	</FormField>

	<div class="flex justify-end pt-5">
		<SubmitButton label="Create Node" />
	</div>
</form>
