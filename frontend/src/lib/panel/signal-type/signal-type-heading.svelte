<script lang="ts">
	import { EditableForm } from '$lib/components/form';
	import { SignalTypeIcon } from '$lib/components/icon';
	import { ResizeableTextInput } from '$lib/components/input';
	import type { PanelSectionProps } from '../types';
	import { getSignalTypeContext } from './signal-type-context.svelte';
	import { nameSchema } from './signal-type-schema';

	let { entityId, invalidNames }: PanelSectionProps & { invalidNames: string[] } = $props();

	const { signalType } = getSignalTypeContext(entityId);

	function handleName(name: string) {}
</script>

<section>
	<div class="flex h-24 gap-3 items-center">
		<SignalTypeIcon width="48" height="48" />

		<div class="flex-1">
			<EditableForm
				schema={nameSchema(invalidNames)}
				initialValues={{ name: signalType.name }}
				hidePlaceholder
				onsubmit={({ name }) => handleName(name)}
			>
				{#snippet placeholder()}
					<div class="text-xl font-bold break-all">{signalType.name}</div>
				{/snippet}

				{#snippet input({ fsm, values, errors })}
					<ResizeableTextInput
						size="sm"
						label="Name"
						name="signal-type-name"
						bind:value={values.name}
						errors={errors.name}
						focusOnDiplay
						onescape={() => fsm.send('ESCAPE')}
						onblur={() => fsm.send('BLUR')}
					/>
				{/snippet}
			</EditableForm>
		</div>
	</div>
</section>
